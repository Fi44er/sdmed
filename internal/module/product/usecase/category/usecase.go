package category_usecase

import (
	"context"

	"github.com/Fi44er/sdmed/internal/module/product/entity"
	product_constant "github.com/Fi44er/sdmed/internal/module/product/pkg"
	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/Fi44er/sdmed/pkg/postgres/uow"
)

type IFileUsecaseAdapter interface {
	UploadFile(ctx context.Context, file *product_entity.File) error
	GetFile(ctx context.Context, name string) (*product_entity.File, error)
}

type ICategoryUsecase interface {
	Create(ctx context.Context, category *product_entity.Category) error
}

type ICategoryRepository interface {
	Create(ctx context.Context, entity *product_entity.Category) error
	GetByID(ctx context.Context, id string) (*product_entity.Category, error)
	GetByName(ctx context.Context, name string) (*product_entity.Category, error)
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
) ICategoryUsecase {
	return &CategoryUsecase{
		logger:      logger,
		repository:  repository,
		fileUsecase: fileUsease,
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

		if existCategory == nil {
			return product_constant.ErrCategoryNotFound
		}

		if err := u.repository.Create(ctx, category); err != nil {
			return err
		}

		for _, image := range category.Images {
			if err := u.fileUsecase.UploadFile(ctx, &image); err != nil {
				return err
			}
		}

		needCleanup = false
		return nil
	})
}
