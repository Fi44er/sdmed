package repository

import (
	"context"

	"github.com/Fi44er/sdmed/internal/module/file/entity"
	"github.com/Fi44er/sdmed/pkg/logger"
	"gorm.io/gorm"
)

type FileRepository struct {
	logger *logger.Logger
	db     *gorm.DB
}

func NewFileRepository(logger *logger.Logger, db *gorm.DB) *FileRepository {
	return &FileRepository{
		logger: logger,
		db:     db,
	}
}

func (r *FileRepository) Create(ctx context.Context, file *entity.File) error {
	r.logger.Info("creating file...")
	if err := r.db.WithContext(ctx).Create(file).Error; err != nil {
		return err
	}
	return nil
}

func (r *FileRepository) GetByID(ctx context.Context, id string) (*entity.File, error) {
	r.logger.Info("getting file by id...")
	var file entity.File
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&file).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Warnf("File not found: %s", id)
			return nil, nil
		}
		return nil, err
	}
	return &file, nil
}

func (r *FileRepository) GetByOwner(ctx context.Context, ownerID, ownerType string) (*entity.File, error) {
	r.logger.Info("getting file by owner...")
	var file entity.File
	if err := r.db.WithContext(ctx).Where("owner_id = ? AND owner_type = ?", ownerID, ownerType).First(&file).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Warnf("File not found: %s", ownerID)
			return nil, nil
		}
		return nil, err
	}
	return &file, nil
}
