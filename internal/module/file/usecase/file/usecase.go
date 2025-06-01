package usecase

import (
	"context"

	"github.com/Fi44er/sdmed/internal/module/file/dto"
	"github.com/Fi44er/sdmed/pkg/logger"
)

type IFileRepository interface{}

type IFileStorage interface {
	UploadFiles(name string, data []byte) error
}

type FileUsecase struct {
	repository  IFileRepository
	fileStorage IFileStorage
	logger      *logger.Logger
}

func NewFileUsecase(
	repository IFileRepository,
	fileStorage IFileStorage,
	logger *logger.Logger,
) *FileUsecase {
	return &FileUsecase{
		repository:  repository,
		fileStorage: fileStorage,
		logger:      logger,
	}
}

func (u *FileUsecase) UploadFiles(ctx *context.Context, dto dto.UploadFiles) error {
	if err := u.fileStorage.UploadFiles(dto.Name, dto.Data); err != nil {
		return err
	}
	return nil
}
