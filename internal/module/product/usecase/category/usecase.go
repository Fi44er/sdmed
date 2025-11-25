package category_usecase

import (
	"context"
	"fmt"

	product_entity "github.com/Fi44er/sdmed/internal/module/product/entity"
	product_constant "github.com/Fi44er/sdmed/internal/module/product/pkg"
	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/Fi44er/sdmed/pkg/postgres/uow"
	"github.com/Fi44er/sdmed/pkg/utils"
)

const ownerType = "category"

type IFileUsecaseAdapter interface {
	MakeFilesPermanent(ctx context.Context, fileIDs []string, ownerID, ownerType string) error
	GetByOwner(ctx context.Context, ownerID, ownerType string) ([]product_entity.File, error)
	GetByOwners(ctx context.Context, ownerIDs []string, ownerType string) (map[string][]product_entity.File, error)

	DeleteByOwner(ctx context.Context, ownerID, ownerType string) error
	DeleteByID(ctx context.Context, id string) error
}

type ICategoryUsecase interface {
	Create(ctx context.Context, category *product_entity.Category) error
	GetByID(ctx context.Context, id string) (*product_entity.Category, error)
	GetAll(ctx context.Context, offset, limit int) ([]product_entity.Category, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, category *product_entity.Category) error
}

type ICategoryRepository interface {
	Create(ctx context.Context, entity *product_entity.Category) error
	Update(ctx context.Context, category *product_entity.Category) error
	GetByID(ctx context.Context, id string) (*product_entity.Category, error)
	GetByName(ctx context.Context, name string) (*product_entity.Category, error)
	GetAll(ctx context.Context, offset, limit int) ([]product_entity.Category, error)
	Delete(ctx context.Context, id string) error
}

type CategoryUsecase struct {
	repository  ICategoryRepository
	uow         uow.Uow
	logger      *logger.Logger
	fileUsecase IFileUsecaseAdapter
}

func NewCategoryUsecase(
	logger *logger.Logger,
	repository ICategoryRepository,
	fileUsease IFileUsecaseAdapter,
	uow uow.Uow,
) ICategoryUsecase {
	return &CategoryUsecase{
		logger:      logger,
		repository:  repository,
		fileUsecase: fileUsease,
		uow:         uow,
	}
}

func (u *CategoryUsecase) Update(ctx context.Context, category *product_entity.Category) error {
	u.logger.Infof("Updating category: %s", category.Name)

	return u.uow.Do(ctx, func(ctx context.Context) error {
		repo, err := u.uow.GetRepository(ctx, ownerType)
		if err != nil {
			u.logger.Errorf("Failed to get repository: %v", err)
			return err
		}
		categoryRepo := repo.(ICategoryRepository)

		existCategory, err := categoryRepo.GetByID(ctx, category.ID)
		if err != nil {
			u.logger.Errorf("Failed to get category from repository: %v", err)
			return err
		}

		if err := categoryRepo.Update(ctx, category); err != nil {
			u.logger.Errorf("Failed to update category in repository: %v", err)
			return err
		}

		files, err := u.fileUsecase.GetByOwner(ctx, category.ID, ownerType)
		if err != nil {
			u.logger.Warnf("Failed to get files for category %s: %v", category.ID, err)
		} else {
			u.logger.Debugf("Found %d files for category %s", len(files), category.ID)
		}

		existCategory.Images = files

		deletedImg, addedImg := utils.FindDifferences(existCategory.Images, category.Images, func(f product_entity.File) (string, string) { return f.Name, f.ID })

		// TODO: оптимизировать удаление файлов
		for _, fileID := range deletedImg {
			if err := u.fileUsecase.DeleteByID(ctx, fileID); err != nil {
				u.logger.Warnf("Failed to delete file %s: %v", fileID, err)
			}
		}

		if err := u.fileUsecase.MakeFilesPermanent(ctx, addedImg, category.ID, ownerType); err != nil {
			u.logger.Errorf("Failed to make files permanent for category %s: %v", category.ID, err)
			return err
		}

		return nil
	})
}

func (u *CategoryUsecase) Create(ctx context.Context, category *product_entity.Category) error {
	u.logger.Infof("Creating category: %s", category.Name)

	return u.uow.Do(ctx, func(ctx context.Context) error {
		repo, err := u.uow.GetRepository(ctx, ownerType)
		if err != nil {
			u.logger.Errorf("Failed to get repository: %v", err)
			return err
		}

		categoryRepo := repo.(ICategoryRepository)

		needCleanup := true
		defer func() {
			if needCleanup {
				u.logger.Warnf("Cleaning up category due to failed creation: %s", category.ID)
				if err := categoryRepo.Delete(ctx, category.ID); err != nil {
					u.logger.Errorf("Failed to cleanup category %s: %v", category.ID, err)
				}
			}
		}()

		existCategory, err := categoryRepo.GetByName(ctx, category.Name)
		if err != nil {
			u.logger.Errorf("Failed to check category existence by name %s: %v", category.Name, err)
			return err
		}

		if existCategory != nil {
			u.logger.Warnf("Category already exists: %s", category.Name)
			return product_constant.ErrCategoryAlreadyExist
		}

		if err := categoryRepo.Create(ctx, category); err != nil {
			u.logger.Errorf("Failed to create category in repository: %v", err)
			return err
		}

		imagesNames := make([]string, 0)
		for _, image := range category.Images {
			imagesNames = append(imagesNames, image.Name)
		}

		if len(imagesNames) > 0 {
			u.logger.Infof("Making %d files permanent for category %s", len(imagesNames), category.ID)
			if err := u.fileUsecase.MakeFilesPermanent(ctx, imagesNames, category.ID, ownerType); err != nil {
				u.logger.Errorf("Failed to make files permanent for category %s: %v", category.ID, err)
				return err
			}
		}
		needCleanup = false
		u.logger.Infof("Category created successfully: %s (ID: %s)", category.Name, category.ID)
		return nil
	})
}

func (u *CategoryUsecase) GetByID(ctx context.Context, id string) (*product_entity.Category, error) {
	u.logger.Debugf("Getting category by ID: %s", id)

	category, err := u.repository.GetByID(ctx, id)
	if err != nil {
		u.logger.Errorf("Failed to get category by ID %s: %v", id, err)
		return nil, err
	}

	if category == nil {
		u.logger.Debugf("Category not found: %s", id)
		return nil, product_constant.ErrCategoryNotFound
	}

	files, err := u.fileUsecase.GetByOwner(ctx, id, ownerType)
	if err != nil {
		u.logger.Warnf("Failed to get files for category %s: %v", id, err)
	} else {
		u.logger.Debugf("Found %d files for category %s", len(files), id)
	}

	category.Images = files
	u.logger.Debugf("Category retrieved successfully: %s", id)
	return category, nil
}

func (u *CategoryUsecase) GetAll(ctx context.Context, offset, limit int) ([]product_entity.Category, error) {
	u.logger.Debugf("Getting all categories (offset: %d, limit: %d)", offset, limit)

	categories, err := u.repository.GetAll(ctx, offset, limit)
	if err != nil {
		u.logger.Errorf("Failed to get categories: %v", err)
		return nil, err
	}

	u.logger.Debugf("Found %d categories", len(categories))

	if len(categories) == 0 {
		return categories, nil
	}

	if err := u.enrichWithBatch(ctx, categories); err != nil {
		u.logger.Warnf("Failed to enrich categories with images (count: %d): %v", len(categories), err)
	} else {
		u.logger.Debugf("Successfully enriched %d categories with images", len(categories))
	}

	return categories, nil
}

func (u *CategoryUsecase) Delete(ctx context.Context, id string) error {
	u.logger.Infof("Deleting category: %s", id)

	return u.uow.Do(ctx, func(ctx context.Context) error {
		repo, err := u.uow.GetRepository(ctx, ownerType)
		if err != nil {
			u.logger.Errorf("Failed to get repository for deletion: %v", err)
			return err
		}
		categoryRepo := repo.(ICategoryRepository)

		if err := categoryRepo.Delete(ctx, id); err != nil {
			u.logger.Errorf("Failed to delete category %s: %v", id, err)
			return err
		}

		if err := u.fileUsecase.DeleteByOwner(ctx, id, ownerType); err != nil {
			u.logger.Errorf("Failed to delete files for category %s: %v", id, err)
			return err
		}

		u.logger.Infof("Category deleted successfully: %s", id)
		return nil
	})
}

func (u *CategoryUsecase) enrichWithBatch(ctx context.Context, categories []product_entity.Category) error {
	u.logger.Debugf("Enriching %d categories with files", len(categories))

	categoryIDs := make([]string, len(categories))
	for i, category := range categories {
		categoryIDs[i] = category.ID
	}

	filesByOwner, err := u.fileUsecase.GetByOwners(ctx, categoryIDs, ownerType)
	if err != nil {
		return fmt.Errorf("batch get files by owners: %w", err)
	}

	enrichedCount := 0
	for i := range categories {
		if files, exists := filesByOwner[categories[i].ID]; exists {
			categories[i].Images = files
			enrichedCount++
		}
	}

	u.logger.Debugf("Enriched %d out of %d categories with files", enrichedCount, len(categories))
	return nil
}
