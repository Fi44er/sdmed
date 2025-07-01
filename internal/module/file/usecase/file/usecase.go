package usecase

import (
	"context"

	"github.com/Fi44er/sdmed/internal/module/file/dto"
	"github.com/Fi44er/sdmed/internal/module/file/entity"
	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/Fi44er/sdmed/pkg/postgres/uow"
)

type IFileRepository interface {
	Create(ctx context.Context, file *entity.File) error
	GetByName(ctx context.Context, name string) (*entity.File, error)
}

type IFileStorage interface {
	Upload(name *string, data []byte) error
	Delete(name string) error
	Get(name string) ([]byte, error)
}

type FileUsecase struct {
	repository  IFileRepository
	uow         uow.Uow
	fileStorage IFileStorage
	logger      *logger.Logger
}

func NewFileUsecase(
	repository IFileRepository,
	uow uow.Uow,
	fileStorage IFileStorage,
	logger *logger.Logger,
) *FileUsecase {
	return &FileUsecase{
		repository:  repository,
		uow:         uow,
		fileStorage: fileStorage,
		logger:      logger,
	}
}

func (u *FileUsecase) Upload(ctx context.Context, dto *dto.UploadFiles) error {
	return u.uow.Do(ctx, func(ctx context.Context) error {
		if err := dto.File.GenerateName(); err != nil {
			u.logger.Errorf("failed to generate file name: %s", err)
			return err
		}

		needCleanup := true
		defer func() {
			if needCleanup {
				if err := u.fileStorage.Delete(dto.File.Name); err != nil {
					u.logger.Errorf("failed to cleanup file %s: %v", dto.File.Name, err)
				}
			}
		}()

		if err := u.fileStorage.Upload(&dto.File.Name, dto.Data); err != nil {
			return err
		}

		repo, err := u.uow.GetRepository(ctx, "file")
		if err != nil {

			return err
		}
		fileRepo := repo.(IFileRepository)

		if err := fileRepo.Create(ctx, &dto.File); err != nil {
			return err
		}

		needCleanup = false
		return nil
	})
}

func (u *FileUsecase) Get(ctx context.Context, name string) (*entity.File, error) {
	file, err := u.repository.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}

	data, err := u.fileStorage.Get(name)
	if err != nil {
		return nil, err
	}
	file.Data = data
	return file, nil
}
