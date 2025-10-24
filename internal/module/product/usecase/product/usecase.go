package usecase

import (
	"context"

	"github.com/Fi44er/sdmed/internal/module/product/entity"
	"github.com/Fi44er/sdmed/internal/module/product/usecase/product/contracts"
	"github.com/Fi44er/sdmed/pkg/logger"
)

type ProductUsecase struct {
	repository  contracts.IProductRepository
	logger      *logger.Logger
	fileUsecase contracts.IFileUsecase
}

func NewProductUsecase(
	repository contracts.IProductRepository,
	logger *logger.Logger,
	fileUsecase contracts.IFileUsecase,
) *ProductUsecase {
	return &ProductUsecase{
		repository:  repository,
		logger:      logger,
		fileUsecase: fileUsecase,
	}
}

func (u *ProductUsecase) GetAll(ctx context.Context, limit, offset int) ([]product_entity.Product, error) {
	return u.repository.GetAll(ctx, limit, offset)
}

func (u *ProductUsecase) GetByID(ctx context.Context, id string) (*product_entity.Product, error) {
	return u.repository.GetByID(ctx, id)
}

func (u *ProductUsecase) Create(ctx context.Context, product_entity *product_entity.Product) error {
	return u.repository.Create(ctx, product_entity)
}

func (u *ProductUsecase) UploadImages(ctx context.Context, file []product_entity.File) error {
	for _, f := range file {
		err := u.fileUsecase.UploadFile(ctx, &f)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *ProductUsecase) Update(ctx context.Context, product_entity *product_entity.Product) error {
	return u.repository.Update(ctx, product_entity)
}

func (u *ProductUsecase) Delete(ctx context.Context, id string) error {
	return u.repository.Delete(ctx, id)
}
