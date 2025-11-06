package file_usecase

import (
	"context"
	"time"

	"github.com/Fi44er/sdmed/internal/module/file/entity"
	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/Fi44er/sdmed/pkg/postgres/uow"
)

type IFileRepository interface {
	GetByName(ctx context.Context, name string) (*entity.File, error)
	GetExpiredTemporaryFiles(ctx context.Context, before time.Time) ([]*entity.File, error)

	Create(ctx context.Context, file *entity.File) error
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, file *entity.File) error
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

type IFileUsecase interface {
	Upload(ctx context.Context, file *entity.File) error
	Get(ctx context.Context, name string) (*entity.File, error)
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

func (u *FileUsecase) UploadTemporary(ctx context.Context, file *entity.File, ttl time.Duration) error {
	return u.uploadWithStatus(ctx, file, entity.FileStatusTemporary, ttl)
}

func (u *FileUsecase) UploadPermanent(ctx context.Context, file *entity.File, ownerID, ownerType string) error {
	if err := u.uploadWithStatus(ctx, file, entity.FileStatusPermanent, 0); err != nil {
		return err
	}

	// Обновляем файл как постоянный
	return u.uow.Do(ctx, func(ctx context.Context) error {
		repo, err := u.uow.GetRepository(ctx, "file")
		if err != nil {
			return err
		}
		fileRepo := repo.(IFileRepository)

		file.MarkAsPermanent(ownerID, ownerType)
		return fileRepo.Update(ctx, file)
	})
}

func (u *FileUsecase) uploadWithStatus(ctx context.Context, file *entity.File, status entity.FileStatus, ttl time.Duration) error {
	return u.uow.Do(ctx, func(ctx context.Context) error {
		if err := file.GenerateName(); err != nil {
			u.logger.Errorf("failed to generate file name: %s", err)
			return err
		}

		// Устанавливаем статус и TTL
		file.Status = status
		if status == entity.FileStatusTemporary && ttl > 0 {
			expiresAt := time.Now().Add(ttl)
			file.ExpiresAt = &expiresAt
		}

		needCleanup := true
		defer func() {
			if needCleanup {
				if err := u.fileStorage.Delete(file.Name); err != nil {
					u.logger.Errorf("failed to cleanup file %s: %v", file.Name, err)
				}
			}
		}()

		if err := u.fileStorage.Upload(&file.Name, file.Data); err != nil {
			return err
		}

		repo, err := u.uow.GetRepository(ctx, "file")
		if err != nil {
			return err
		}
		fileRepo := repo.(IFileRepository)

		if err := fileRepo.Create(ctx, file); err != nil {
			return err
		}

		needCleanup = false
		return nil
	})
}

func (u *FileUsecase) MakeFilesPermanent(ctx context.Context, fileIDs []string, ownerID, ownerType string) error {
	return u.uow.Do(ctx, func(ctx context.Context) error {
		repo, err := u.uow.GetRepository(ctx, "file")
		if err != nil {
			return err
		}
		fileRepo := repo.(IFileRepository)

		for _, fileID := range fileIDs {
			file, err := fileRepo.GetByName(ctx, fileID)
			if err != nil {
				return err
			}

			file.MarkAsPermanent(ownerID, ownerType)
			if err := fileRepo.Update(ctx, file); err != nil {
				return err
			}
		}
		return nil
	})
}
