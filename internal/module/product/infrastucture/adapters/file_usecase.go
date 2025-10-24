package product_adapters

import (
	"context"

	file_entity "github.com/Fi44er/sdmed/internal/module/file/entity"
	file_usecase "github.com/Fi44er/sdmed/internal/module/file/usecase/file"
	"github.com/Fi44er/sdmed/internal/module/product/entity"
)

type FileUsecaseAdapter struct {
	fileUsecase *file_usecase.FileUsecase
}

func NewFileUsecaseAdapter(fileUsecase *file_usecase.FileUsecase) *FileUsecaseAdapter {
	return &FileUsecaseAdapter{
		fileUsecase: fileUsecase,
	}
}

func (a *FileUsecaseAdapter) UploadFile(ctx context.Context, file *product_entity.File) error {
	fileproduct_entity := toFileproduct_entity(file)
	return a.fileUsecase.Upload(ctx, fileproduct_entity)
}

func toFileproduct_entity(file *product_entity.File) *file_entity.File {
	return &file_entity.File{
		ID:        file.ID,
		Name:      file.Name,
		OwnerID:   file.OwnerID,
		OwnerType: file.OwnerType,
		Data:      file.Data,
	}
}
