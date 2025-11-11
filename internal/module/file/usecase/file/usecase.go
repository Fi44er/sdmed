package file_usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/Fi44er/sdmed/internal/config"
	file_entity "github.com/Fi44er/sdmed/internal/module/file/entity"
	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/Fi44er/sdmed/pkg/postgres/uow"
)

type IFileRepository interface {
	GetByName(ctx context.Context, name string) (*file_entity.File, error)
	GetExpiredTemporaryFiles(ctx context.Context) ([]*file_entity.File, error)
	GetByOwner(ctx context.Context, ownerID, ownerType string) ([]file_entity.File, error)
	GetByOwners(ctx context.Context, ownerIDs []string, ownerType string) ([]file_entity.File, error)

	Create(ctx context.Context, file *file_entity.File) error
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, file *file_entity.File) error
}

type IFileStorage interface {
	Upload(name *string, data []byte) error
	Delete(name string) error
	Get(name string) ([]byte, error)
}

type IFileUsecase interface {
	UploadTemporary(ctx context.Context, file *file_entity.File, ttl time.Duration) (string, error)
	UploadPermanent(ctx context.Context, file *file_entity.File, ownerID, ownerType string) (string, error)
	MakeFilesPermanent(ctx context.Context, fileIDs []string, ownerID, ownerType string) error
	Get(ctx context.Context, name string) (*file_entity.File, error)
	GetByOwner(ctx context.Context, ownerID, ownerType string) ([]file_entity.File, error)
	GetByOwners(ctx context.Context, ownerIDs []string, ownerType string) (map[string][]file_entity.File, error)
}

type FileUsecase struct {
	repository  IFileRepository
	uow         uow.Uow
	fileStorage IFileStorage
	logger      *logger.Logger
	config      *config.Config
}

func NewFileUsecase(
	repository IFileRepository,
	uow uow.Uow,
	fileStorage IFileStorage,
	logger *logger.Logger,
	config *config.Config,
) IFileUsecase {
	return &FileUsecase{
		repository:  repository,
		uow:         uow,
		fileStorage: fileStorage,
		logger:      logger,
		config:      config,
	}
}

func (u *FileUsecase) UploadTemporary(ctx context.Context, file *file_entity.File, ttl time.Duration) (string, error) {
	return u.uploadWithStatus(ctx, file, file_entity.FileStatusTemporary, ttl)
}

func (u *FileUsecase) UploadPermanent(ctx context.Context, file *file_entity.File, ownerID, ownerType string) (string, error) {
	url, err := u.uploadWithStatus(ctx, file, file_entity.FileStatusPermanent, 0)
	if err != nil {
		return "", err
	}

	err = u.uow.Do(ctx, func(ctx context.Context) error {
		repo, err := u.uow.GetRepository(ctx, "file")
		if err != nil {
			return err
		}
		fileRepo := repo.(IFileRepository)

		file.MarkAsPermanent(ownerID, ownerType)
		return fileRepo.Update(ctx, file)
	})

	return url, err
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

func (u *FileUsecase) uploadWithStatus(ctx context.Context, file *file_entity.File, status file_entity.FileStatus, ttl time.Duration) (string, error) {
	var fileName string
	err := u.uow.Do(ctx, func(ctx context.Context) error {
		if err := file.GenerateName(); err != nil {
			u.logger.Errorf("failed to generate file name: %s", err)
			return err
		}

		file.Status = status
		if status == file_entity.FileStatusTemporary && ttl > 0 {
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

		fileName = file.Name

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

	link := fmt.Sprintf("%s/%s/%s", u.config.ApiUrl, u.config.FileLink, fileName)

	return link, err
}

func (u *FileUsecase) Get(ctx context.Context, name string) (*file_entity.File, error) {
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

func (u *FileUsecase) GetByOwner(ctx context.Context, ownerID, ownerType string) ([]file_entity.File, error) {
	files, err := u.repository.GetByOwner(ctx, ownerID, ownerType)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func (u *FileUsecase) GetByOwners(ctx context.Context, ownerIDs []string, ownerType string) (map[string][]file_entity.File, error) {
	files, err := u.repository.GetByOwners(ctx, ownerIDs, ownerType)
	if err != nil {
		return nil, err
	}

	filesByOwner := make(map[string][]file_entity.File)

	for _, file := range files {
		if _, exists := filesByOwner[*file.OwnerID]; !exists {
			filesByOwner[*file.OwnerID] = make([]file_entity.File, 0)
		}

		filesByOwner[*file.OwnerID] = append(filesByOwner[*file.OwnerID], file)
	}

	for _, ownerID := range ownerIDs {
		if _, exists := filesByOwner[ownerID]; !exists {
			filesByOwner[ownerID] = make([]file_entity.File, 0)
		}
	}

	u.logger.Infof("files loaded successfully total_files: %v", len(files))

	return filesByOwner, nil
}
