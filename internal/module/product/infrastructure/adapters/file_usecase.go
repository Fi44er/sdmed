package product_adapters

import (
	"context"
	"fmt"

	file_entity "github.com/Fi44er/sdmed/internal/module/file/entity"
	file_usecase "github.com/Fi44er/sdmed/internal/module/file/usecase/file"
	product_entity "github.com/Fi44er/sdmed/internal/module/product/entity"
)

type IFileUsecaseAdapter interface {
	MakeFilesPermanent(ctx context.Context, fileIDs []string, ownerID, ownerType string) error
	GetByOwner(ctx context.Context, ownerID, ownerType string) ([]product_entity.File, error)
	GetByOwners(ctx context.Context, ownerIDs []string, ownerType string) (map[string][]product_entity.File, error)
	DeleteByOwner(ctx context.Context, ownerID, ownerType string) error
	DeleteByID(ctx context.Context, id string) error
}

type FileUsecaseAdapter struct {
	fileUsecase file_usecase.IFileUsecase
}

func NewFileUsecaseAdapter(fileUsecase file_usecase.IFileUsecase) IFileUsecaseAdapter {
	return &FileUsecaseAdapter{
		fileUsecase: fileUsecase,
	}
}

func (a *FileUsecaseAdapter) DeleteByID(ctx context.Context, id string) error {
	return a.fileUsecase.DeleteByID(ctx, id)
}

func (a *FileUsecaseAdapter) DeleteByOwner(ctx context.Context, ownerID, ownerType string) error {
	return a.fileUsecase.DeleteByOwner(ctx, ownerID, ownerType)
}

func (a *FileUsecaseAdapter) MakeFilesPermanent(ctx context.Context, fileIDs []string, ownerID, ownerType string) error {
	return a.fileUsecase.MakeFilesPermanent(ctx, fileIDs, ownerID, ownerType)
}

func (a *FileUsecaseAdapter) GetByOwner(ctx context.Context, ownerID, ownerType string) ([]product_entity.File, error) {
	files, err := a.fileUsecase.GetByOwner(ctx, ownerID, ownerType)
	if err != nil {
		return nil, fmt.Errorf("file usecase get by owner: %w", err)
	}

	return a.convertFiles(files), nil
}

func (a *FileUsecaseAdapter) GetByOwners(ctx context.Context, ownerIDs []string, ownerType string) (map[string][]product_entity.File, error) {
	filesByOwner, err := a.fileUsecase.GetByOwners(ctx, ownerIDs, ownerType)
	if err != nil {
		return nil, fmt.Errorf("file usecase get by owners: %w", err)
	}

	return a.convertFilesByOwner(filesByOwner), nil
}

func (a *FileUsecaseAdapter) convertFiles(files []file_entity.File) []product_entity.File {
	if files == nil {
		return nil
	}

	result := make([]product_entity.File, len(files))
	for i, file := range files {
		result[i] = *a.toProductFile(&file)
	}
	return result
}

func (a *FileUsecaseAdapter) convertFilesByOwner(filesByOwner map[string][]file_entity.File) map[string][]product_entity.File {
	if filesByOwner == nil {
		return nil
	}

	result := make(map[string][]product_entity.File, len(filesByOwner))
	for ownerID, files := range filesByOwner {
		result[ownerID] = a.convertFiles(files)
	}
	return result
}

func (a *FileUsecaseAdapter) toFile(file *product_entity.File) *file_entity.File {
	return &file_entity.File{
		ID:        file.ID,
		Name:      file.Name,
		OwnerID:   file.OwnerID,
		OwnerType: file.OwnerType,
		Data:      file.Data,
	}
}

func (a *FileUsecaseAdapter) toProductFile(file *file_entity.File) *product_entity.File {
	return &product_entity.File{
		ID:        file.ID,
		Name:      file.Name,
		OwnerID:   file.OwnerID,
		OwnerType: file.OwnerType,
		Data:      file.Data,
	}
}
