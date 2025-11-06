package product_adapters

import (
	"context"

	file_entity "github.com/Fi44er/sdmed/internal/module/file/entity"
	file_usecase "github.com/Fi44er/sdmed/internal/module/file/usecase/file"
	product_entity "github.com/Fi44er/sdmed/internal/module/product/entity"
)

type IFileUsecaseAdapter interface {
	UploadFile(ctx context.Context, file *product_entity.File) error
	GetFile(ctx context.Context, name string) (*product_entity.File, error)
}

type FileUsecaseAdapter struct {
	fileUsecase file_usecase.IFileUsecase
}

func NewFileUsecaseAdapter(fileUsecase file_usecase.IFileUsecase) IFileUsecaseAdapter {
	return &FileUsecaseAdapter{
		fileUsecase: fileUsecase,
	}
}

func (a *FileUsecaseAdapter) UploadFile(ctx context.Context, file *product_entity.File) error {
	fileproduct_entity := toFile(file)
	return a.fileUsecase.Upload(ctx, fileproduct_entity)
}

func (a *FileUsecaseAdapter) GetFile(ctx context.Context, name string) (*product_entity.File, error) {
	file, err := a.fileUsecase.Get(ctx, name)
	return toProductFile(file), err
}

func toFile(file *product_entity.File) *file_entity.File {
	return &file_entity.File{
		ID:        file.ID,
		Name:      file.Name,
		OwnerID:   file.OwnerID,
		OwnerType: file.OwnerType,
		Data:      file.Data,
	}
}

func toProductFile(file *file_entity.File) *product_entity.File {
	return &product_entity.File{
		ID:        file.ID,
		Name:      file.Name,
		OwnerID:   file.OwnerID,
		OwnerType: file.OwnerType,
		Data:      file.Data,
	}
}
