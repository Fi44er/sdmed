package category_usecase

import (
	"context"

	product_entity "github.com/Fi44er/sdmed/internal/module/product/entity"
	product_constant "github.com/Fi44er/sdmed/internal/module/product/pkg"
	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/Fi44er/sdmed/pkg/postgres/uow"
)

type IFileUsecaseAdapter interface {
	MakeFilesPermanent(ctx context.Context, fileIDs []string, ownerID, ownerType string) error
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

		if err := u.fileUsecase.MakeFilesPermanent(ctx, imagesNames, category.ID, "category"); err != nil {
			return err
		}

		needCleanup = false
		return nil
	})
}
