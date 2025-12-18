package repository

import (
	"context"
	"fmt"
	"time"

	file_entity "github.com/Fi44er/sdmed/internal/module/file/entity"
	"github.com/Fi44er/sdmed/internal/module/file/infrastructure/repository/model"
	"github.com/Fi44er/sdmed/pkg/logger"
	"gorm.io/gorm"
)

type IFileRepository interface {
	Delete(ctx context.Context, id string) error
	DeleteByOwner(ctx context.Context, ownerID, ownerType string) error
	Update(ctx context.Context, file *file_entity.File) error
	Create(ctx context.Context, file *file_entity.File) error

	GetByName(ctx context.Context, name string) (*file_entity.File, error)
	GetExpiredTemporaryFiles(ctx context.Context) ([]*file_entity.File, error)
	GetByID(ctx context.Context, id string) (*file_entity.File, error)
	GetByOwner(ctx context.Context, ownerID, ownerType string) ([]file_entity.File, error)
	GetByOwners(ctx context.Context, ownerIDs []string, ownerType string) ([]file_entity.File, error)
}

type FileRepository struct {
	logger    *logger.Logger
	db        *gorm.DB
	converter *Converter
}

func NewFileRepository(logger *logger.Logger, db *gorm.DB) IFileRepository {
	return &FileRepository{
		logger:    logger,
		db:        db,
		converter: &Converter{},
	}
}

func (r *FileRepository) Create(ctx context.Context, file *file_entity.File) error {
	r.logger.Info("creating file...")
	fileModel := r.converter.ToModel(file)
	if err := r.db.WithContext(ctx).Create(fileModel).Error; err != nil {
		return err
	}
	file.ID = fileModel.ID
	return nil
}

func (r *FileRepository) GetByID(ctx context.Context, id string) (*file_entity.File, error) {
	r.logger.Info("getting file by id...")
	var fileModel model.File
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&fileModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Warnf("File not found: %s", id)
			return nil, nil
		}
		return nil, err
	}

	file := r.converter.ToEntity(&fileModel)
	return file, nil
}

func (r *FileRepository) GetByName(ctx context.Context, name string) (*file_entity.File, error) {
	r.logger.Info("getting file by name...")
	var fileModel model.File
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&fileModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Warnf("File not found: %s", name)
			return nil, nil
		}
		return nil, err
	}

	file := r.converter.ToEntity(&fileModel)
	return file, nil
}

func (r *FileRepository) GetByOwner(ctx context.Context, ownerID, ownerType string) ([]file_entity.File, error) {
	r.logger.Info("getting files by owner...")
	filesModels := make([]model.File, 0)
	if err := r.db.WithContext(ctx).Where("owner_id = ? AND owner_type = ?", ownerID, ownerType).Find(&filesModels).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Warnf("Files not found: %s", ownerID)
			return nil, nil
		}
		return nil, err
	}

	filesEntity := make([]file_entity.File, 0)
	for _, model := range filesModels {
		file := r.converter.ToEntity(&model)
		filesEntity = append(filesEntity, *file)
	}

	return filesEntity, nil
}

func (r *FileRepository) GetByOwners(ctx context.Context, ownerIDs []string, ownerType string) ([]file_entity.File, error) {
	r.logger.Infof("getting files by multiple owners... owner_count: %v; owner_type: %s", len(ownerIDs), ownerType)

	if len(ownerIDs) == 0 {
		return make([]file_entity.File, 0), nil
	}

	var filesModels []model.File
	if err := r.db.WithContext(ctx).
		Where("owner_id IN (?) AND owner_type = ?", ownerIDs, ownerType).
		Find(&filesModels).Error; err != nil {
		return nil, fmt.Errorf("database query failed: %w", err)
	}

	filesEntity := make([]file_entity.File, 0)
	for _, model := range filesModels {
		file := r.converter.ToEntity(&model)
		filesEntity = append(filesEntity, *file)
	}

	return filesEntity, nil
}

func (r *FileRepository) GetExpiredTemporaryFiles(ctx context.Context) ([]*file_entity.File, error) {
	r.logger.Info("getting expired temporary files...")
	var fileModels []model.File

	now := time.Now()

	if err := r.db.WithContext(ctx).
		Where("status = ? AND expires_at IS NOT NULL AND expires_at < ?",
			file_entity.FileStatusTemporary, now).
		Find(&fileModels).Error; err != nil {
		return nil, err
	}

	files := make([]*file_entity.File, len(fileModels))
	for i, fileModel := range fileModels {
		files[i] = r.converter.ToEntity(&fileModel)
	}

	return files, nil
}
func (r *FileRepository) Update(ctx context.Context, file *file_entity.File) error {
	r.logger.Info("updating file...")

	fileModel := r.converter.ToModel(file)
	if err := r.db.WithContext(ctx).Model(&fileModel).Updates(fileModel).Error; err != nil {
		return err
	}
	return nil
}

func (r *FileRepository) Delete(ctx context.Context, id string) error {
	r.logger.Info("deleting file by id...")
	if err := r.db.WithContext(ctx).Delete(&model.File{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (r *FileRepository) DeleteByOwner(ctx context.Context, ownerID, ownerType string) error {
	r.logger.Infof("deleting files by owner; ownerID: %s; ownerType: %s", ownerID, ownerType)
	if err := r.db.WithContext(ctx).Delete(&model.File{}, "owner_id = ? AND owner_type = ?", ownerID, ownerType).Error; err != nil {
		return err
	}
	return nil
}
