package repository

import (
	"context"

	"github.com/Fi44er/sdmed/internal/module/file/entity"
	"github.com/Fi44er/sdmed/internal/module/file/infrastucture/repository/model"
	"github.com/Fi44er/sdmed/pkg/logger"
	"gorm.io/gorm"
)

type FileRepository struct {
	logger    *logger.Logger
	db        *gorm.DB
	converter *Converter
}

func NewFileRepository(logger *logger.Logger, db *gorm.DB) *FileRepository {
	return &FileRepository{
		logger:    logger,
		db:        db,
		converter: &Converter{},
	}
}

func (r *FileRepository) Create(ctx context.Context, file *entity.File) error {
	r.logger.Info("creating file...")
	fileModel := r.converter.ToModel(file)
	if err := r.db.WithContext(ctx).Create(fileModel).Error; err != nil {
		return err
	}
	file.ID = fileModel.ID
	return nil
}

func (r *FileRepository) GetByID(ctx context.Context, id string) (*entity.File, error) {
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

func (r *FileRepository) Delete(ctx context.Context, id string) error {
	r.logger.Info("deleting file by id...")
	if err := r.db.WithContext(ctx).Delete(&model.File{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (r *FileRepository) GetByName(ctx context.Context, name string) (*entity.File, error) {
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

func (r *FileRepository) GetByOwner(ctx context.Context, ownerID, ownerType string) (*entity.File, error) {
	r.logger.Info("getting file by owner...")
	var fileModel model.File
	if err := r.db.WithContext(ctx).Where("owner_id = ? AND owner_type = ?", ownerID, ownerType).First(&fileModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Warnf("File not found: %s", ownerID)
			return nil, nil
		}
		return nil, err
	}

	file := r.converter.ToEntity(&fileModel)

	return file, nil
}
