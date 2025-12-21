package product_usecase

import (
	"context"
	"fmt"

	product_entity "github.com/Fi44er/sdmed/internal/module/product/entity"
	product_constant "github.com/Fi44er/sdmed/internal/module/product/pkg"
	product_usecase_contracts "github.com/Fi44er/sdmed/internal/module/product/usecase/product/contracts"
	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/Fi44er/sdmed/pkg/postgres/uow"
)

const ownerType = "product"

type IProductUsecase interface {
	Create(ctx context.Context, product *product_entity.Product) error
	GetBySlug(ctx context.Context, slug string) (*product_entity.Product, error)
	GetAll(ctx context.Context, params *product_entity.ProductFilterParams) ([]product_entity.Product, int64, error)

	GetFilters(ctx context.Context, categoryID string) ([]product_entity.Filter, error)
}

type ProductUsecase struct {
	repository product_usecase_contracts.IProductRepository
	logger     *logger.Logger
	cache      product_usecase_contracts.ICache

	uow              uow.Uow
	fileUsecase      product_usecase_contracts.IFileUsecaseAdapter
	charValueUsecase product_usecase_contracts.ICharValueUsecase
}

func NewProductUsecase(
	repository product_usecase_contracts.IProductRepository,
	logger *logger.Logger,
	uow uow.Uow,
	cache product_usecase_contracts.ICache,
	fileUsecase product_usecase_contracts.IFileUsecaseAdapter,
	charValueUsecase product_usecase_contracts.ICharValueUsecase,
) IProductUsecase {
	return &ProductUsecase{
		repository:       repository,
		logger:           logger,
		uow:              uow,
		cache:            cache,
		fileUsecase:      fileUsecase,
		charValueUsecase: charValueUsecase,
	}
}

func (u *ProductUsecase) GetAll(ctx context.Context, params *product_entity.ProductFilterParams) ([]product_entity.Product, int64, error) {
	u.logger.Debugf("Getting all products (page: %d, pageSize: %d)", params.Page, params.PageSize)

	u.logger.Debugf("Filter params: %+v", params)
	products, total, err := u.repository.GetAll(ctx, *params)
	if err != nil {
		u.logger.Errorf("Failed to get all products: %v", err)
		return nil, 0, err
	}

	u.logger.Debugf("Found %d products", len(products))

	if len(products) == 0 {
		return products, 0, nil
	}

	if err := u.enrichWithBatch(ctx, products); err != nil {
		u.logger.Warnf("Failed to enrich products with images (count: %d): %v", len(products), err)
	} else {
		u.logger.Debugf("Successfully enriched %d products with images", len(products))
	}

	return products, total, nil
}

func (u *ProductUsecase) GetBySlug(ctx context.Context, slug string) (*product_entity.Product, error) {
	u.logger.Debugf("Getting product by slug: %s", slug)

	product, err := u.repository.GetBySlug(ctx, slug)
	if err != nil {
		u.logger.Errorf("Failed to get product by slug %s: %v", slug, err)
		return nil, err
	}

	if product == nil {
		u.logger.Debugf("Product not found: %s", slug)
		return nil, product_constant.ErrProductNotFound
	}

	files, err := u.fileUsecase.GetByOwner(ctx, product.ID, ownerType)
	if err != nil {
		u.logger.Warnf("Failed to get files for product %s: %v", product.ID, err)
	} else {
		u.logger.Debugf("Found %d files for product %s", len(files), product.ID)
	}

	product.Images = files
	u.logger.Debugf("Product %s retrieved successfully", product.ID)
	return product, nil
}

func (u *ProductUsecase) Create(ctx context.Context, product *product_entity.Product) error {
	u.logger.Infof("Creating product: %s", product.Name)

	return u.uow.Do(ctx, func(ctx context.Context) error {
		repo, err := u.uow.GetRepository(ctx, ownerType)
		if err != nil {
			u.logger.Errorf("Failed to get repository: %v", err)
			return err
		}
		productRepo := repo.(product_usecase_contracts.IProductRepository)

		existProduct, err := productRepo.GetByArticle(ctx, product.Article)
		if err != nil {
			u.logger.Errorf("Failed to check if product with article %s exists: %v", product.Article, err)
			return err
		}

		if existProduct != nil {
			u.logger.Errorf("Product with article %s already exists", product.Article)
			return product_constant.ErrProductAlreadyExists
		}

		product.Slogify()
		if err := productRepo.Create(ctx, product); err != nil {
			u.logger.Errorf("Failed to create product: %v", err)
			return err
		}

		imagesNames := make([]string, 0, len(product.Images))
		for _, image := range product.Images {
			imagesNames = append(imagesNames, image.Name)
		}

		if len(imagesNames) > 0 {
			u.logger.Infof("Making %d files permanent for category %s", len(imagesNames), product.ID)
			if err := u.fileUsecase.MakeFilesPermanent(ctx, imagesNames, product.ID, ownerType); err != nil {
				u.logger.Errorf("Failed to make files permanent for category %s: %v", product.ID, err)
				return err
			}
		}

		for i := range product.CharValues {
			product.CharValues[i].ProductID = product.ID
		}

		if len(product.CharValues) > 0 {
			if err := u.charValueUsecase.CreateMany(ctx, product.CharValues); err != nil {
				u.logger.Errorf("Failed to create char values for product %s: %v", product.ID, err)
				return err
			}
		}

		u.logger.Infof("Product created successfully: %s (ID: %s)", product.Name, product.ID)
		return nil
	})
}

func (u *ProductUsecase) Update(ctx context.Context, product_entity *product_entity.Product) error {
	return u.repository.Update(ctx, product_entity)
}

func (u *ProductUsecase) Delete(ctx context.Context, id string) error {
	u.logger.Infof("Deleting product with ID: %s", id)

	return u.uow.Do(ctx, func(ctx context.Context) error {
		repo, err := u.uow.GetRepository(ctx, ownerType)
		if err != nil {
			u.logger.Errorf("Failed to get repository for product deletion: %v", err)
			return err
		}
		productRepo := repo.(product_usecase_contracts.IProductRepository)

		if err := productRepo.Delete(ctx, id); err != nil {
			u.logger.Errorf("Failed to delete product with ID %s: %v", id, err)
			return err
		}

		if err := u.fileUsecase.DeleteByOwner(ctx, id, ownerType); err != nil {
			u.logger.Errorf("Failed to delete files for category %s: %v", id, err)
			return err
		}

		u.logger.Infof("Product with ID %s deleted successfully", id)
		return nil
	})
}

func (u *ProductUsecase) GetFilters(ctx context.Context, categoryID string) ([]product_entity.Filter, error) {
	cachedFilters := new([]product_entity.Filter)
	key := product_constant.CategoryFiltersKeyPrefix + categoryID
	if err := u.cache.Get(ctx, key, cachedFilters); err == nil {
		u.logger.Debugf("Get filters for category %s from cache successfully", categoryID)
		return *cachedFilters, nil
	}

	filters, err := u.repository.GetFiltersByCategory(ctx, categoryID)
	if err != nil {
		u.logger.Errorf("Failed to get filters for category %s: %v", categoryID, err)
		return nil, err
	}

	if err := u.cache.Set(ctx, key, filters, product_constant.FilterExpered); err != nil {
		u.logger.Errorf("Failed to set filters for category %s to cache: %v", categoryID, err)
	}

	return filters, nil
}

func (u *ProductUsecase) enrichWithBatch(ctx context.Context, products []product_entity.Product) error {
	u.logger.Debugf("Enriching %d products with files", len(products))

	productIDs := make([]string, len(products))
	for i, product := range products {
		productIDs[i] = product.ID
	}

	filesByOwner, err := u.fileUsecase.GetByOwners(ctx, productIDs, ownerType)
	if err != nil {
		return fmt.Errorf("batch get files by owners: %w", err)
	}

	enrichedCount := 0
	for i := range products {
		if files, exists := filesByOwner[products[i].ID]; exists {
			products[i].Images = files
			enrichedCount++
		}
	}

	u.logger.Debugf("Enriched %d out of %d categories with files", enrichedCount, len(products))
	return nil
}
