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
	GetByID(ctx context.Context, id string) (*product_entity.Product, error)

	GetFilters(ctx context.Context, categoryID string) ([]product_entity.Filter, error)
	CreateMany(ctx context.Context, products []*product_entity.Product) error
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

func (u *ProductUsecase) GetByID(ctx context.Context, id string) (*product_entity.Product, error) {
	u.logger.Debugf("Getting product by ID: %s", id)
	product, err := u.repository.GetByID(ctx, id)
	if err != nil {
		u.logger.Errorf("Failed to get product by ID: %v", err)
		return nil, err
	}
	return product, nil
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

func (u *ProductUsecase) CreateMany(ctx context.Context, products []*product_entity.Product) error {
	if len(products) == 0 {
		return nil
	}

	u.logger.Infof("Batch creating %d products", len(products))

	return u.uow.Do(ctx, func(ctx context.Context) error {
		repo, err := u.uow.GetRepository(ctx, ownerType)
		if err != nil {
			return err
		}
		productRepo := repo.(product_usecase_contracts.IProductRepository)

		// 1. Фильтруем продукты, которые уже существуют
		// В идеале в репозиторий нужно добавить метод GetByArticles([]string)
		// Если его нет, используем простую проверку (или пропускаем, если в репо стоит OnConflict DoNothing)
		toCreate := make([]*product_entity.Product, 0, len(products))
		for _, p := range products {
			exist, err := productRepo.GetByArticle(ctx, p.Article)
			if err != nil {
				return err
			}
			if exist == nil {
				p.Slogify() // Генерируем слаг
				toCreate = append(toCreate, p)
			} else {
				// Если товар уже существует, но TRUCodeID обновился или привязался впервые, сохраняем изменения
				shouldUpdate := false
				if exist.TRUCodeID == nil && p.TRUCodeID != nil {
					exist.TRUCodeID = p.TRUCodeID
					shouldUpdate = true
				} else if exist.TRUCodeID != nil && p.TRUCodeID != nil && *exist.TRUCodeID != *p.TRUCodeID {
					exist.TRUCodeID = p.TRUCodeID
					shouldUpdate = true
				}
				
				if shouldUpdate {
					u.logger.Infof("Updating TRUCodeID for existing product %s to %s", exist.Article, *p.TRUCodeID)
					if err := productRepo.Update(ctx, exist); err != nil {
						u.logger.Errorf("Failed to update existing product TRUCodeID: %v", err)
					}
				}
			}
		}

		if len(toCreate) == 0 {
			u.logger.Info("All products already exist, checking for links completed")
			return nil
		}

		// 2. Пакетное сохранение продуктов
		// После этого метода у объектов в toCreate должны заполниться ID (если это делает репо)
		if err := productRepo.CreateMany(ctx, toCreate); err != nil {
			return err
		}

		// 3. Обработка характеристик и файлов
		var allCharValues []*product_entity.ProductCharValue

		for _, p := range toCreate {
			// Подготавливаем файлы (картинки)
			if len(p.Images) > 0 {
				imageNames := make([]string, 0, len(p.Images))
				for _, img := range p.Images {
					imageNames = append(imageNames, img.Name)
				}
				// Делаем файлы постоянными
				if err := u.fileUsecase.MakeFilesPermanent(ctx, imageNames, p.ID, ownerType); err != nil {
					return err
				}
			}

			// Привязываем характеристики к ID созданного продукта
			for i := range p.CharValues {
				p.CharValues[i].ProductID = p.ID
				allCharValues = append(allCharValues, &p.CharValues[i])
			}
		}

		// 4. Пакетное сохранение всех характеристик всех продуктов одним запросом
		if len(allCharValues) > 0 {
			if err := u.charValueUsecase.CreateMany(ctx, allCharValues); err != nil {
				return err
			}
		}

		u.logger.Infof("Successfully created %d products and their characteristics", len(toCreate))
		return nil
	})
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

		charValues := make([]*product_entity.ProductCharValue, 0, len(product.CharValues))
		for i := range product.CharValues {
			charValues = append(charValues, &product.CharValues[i])
		}

		if len(charValues) > 0 {
			if err := u.charValueUsecase.CreateMany(ctx, charValues); err != nil {
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
