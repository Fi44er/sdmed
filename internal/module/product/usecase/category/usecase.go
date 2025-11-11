package category_usecase

import (
	"context"
	"fmt"

	product_entity "github.com/Fi44er/sdmed/internal/module/product/entity"
	product_constant "github.com/Fi44er/sdmed/internal/module/product/pkg"
	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/Fi44er/sdmed/pkg/postgres/uow"
)

const ownerType = "category"

type IFileUsecaseAdapter interface {
	MakeFilesPermanent(ctx context.Context, fileIDs []string, ownerID, ownerType string) error
	GetByOwner(ctx context.Context, ownerID, ownerType string) ([]product_entity.File, error)
	GetByOwners(ctx context.Context, ownerIDs []string, ownerType string) (map[string][]product_entity.File, error)
}

type ICategoryUsecase interface {
	Create(ctx context.Context, category *product_entity.Category) error
	GetByID(ctx context.Context, id string) (*product_entity.Category, error)
	GetAll(ctx context.Context, offset, limit int) ([]product_entity.Category, error)
}

type ICategoryRepository interface {
	Create(ctx context.Context, entity *product_entity.Category) error
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

func (u *CategoryUsecase) Create(ctx context.Context, category *product_entity.Category) error {
	return u.uow.Do(ctx, func(ctx context.Context) error {
		needCleanup := true
		defer func() {
			if needCleanup {
				if err := u.repository.Delete(ctx, category.ID); err != nil {
					u.logger.Errorf("failed to cleanup file %s: %v", category.ID, err)
				}
			}
		}()

		existCategory, err := u.repository.GetByName(ctx, category.Name)
		if err != nil {
			return err
		}

		if existCategory != nil {
			return product_constant.ErrCategoryAlreadyExist
		}

		if err := u.repository.Create(ctx, category); err != nil {
			return err
		}

		imagesNames := make([]string, 0)
		for _, image := range category.Images {
			imagesNames = append(imagesNames, image.Name)
		}

		if err := u.fileUsecase.MakeFilesPermanent(ctx, imagesNames, category.ID, ownerType); err != nil {
			return err
		}

		needCleanup = false
		return nil
	})
}

func (u *CategoryUsecase) GetByID(ctx context.Context, id string) (*product_entity.Category, error) {
	category, err := u.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if category == nil {
		return nil, product_constant.ErrCategoryNotFound
	}

	files, err := u.fileUsecase.GetByOwner(ctx, id, ownerType)
	if err != nil {
		u.logger.Warnf("failed to enrich category with images category_id: %v; error: %v", id, err)
	}

	category.Images = files

	return category, nil
}

func (u *CategoryUsecase) GetAll(ctx context.Context, offset, limit int) ([]product_entity.Category, error) {
	categories, err := u.repository.GetAll(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	if len(categories) == 0 {
		return categories, nil
	}

	if err := u.enrichWithBatch(ctx, categories); err != nil {
		u.logger.Warnf("failed to enrich categories with images category_count: %v; error: %v", len(categories), err)
	}

	return categories, nil
}

func (u *CategoryUsecase) Delete(ctx context.Context, id string) error {
	return nil
}

func (u *CategoryUsecase) enrichWithBatch(
	ctx context.Context,
	categories []product_entity.Category,
) error {
	categoryIDs := make([]string, len(categories))
	for i, category := range categories {
		categoryIDs[i] = category.ID
	}

	filesByOwner, err := u.fileUsecase.GetByOwners(ctx, categoryIDs, ownerType)
	if err != nil {
		return fmt.Errorf("batch get files by owners: %w", err)
	}

	for i := range categories {
		if files, exists := filesByOwner[categories[i].ID]; exists {
			categories[i].Images = files
		}
	}

	return nil
}
