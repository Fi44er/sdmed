package category_repository

import (
	"context"

	"github.com/Fi44er/sdmed/internal/module/product/entity"
	"github.com/Fi44er/sdmed/internal/module/product/infrastucture/repository/model"
	"github.com/Fi44er/sdmed/pkg/logger"
	"gorm.io/gorm"
)

type ICategoryRepository interface {
	Create(ctx context.Context, category *product_entity.Category) error
	GetByID(ctx context.Context, id string) (*product_entity.Category, error)
	GetAll(ctx context.Context, offset, limit int) ([]product_entity.Category, error)
	Update(ctx context.Context, category *product_entity.Category) error
	Delete(ctx context.Context, id string) error
	GetByName(ctx context.Context, name string) (*product_entity.Category, error)
}

type CategoryRepository struct {
	logger    *logger.Logger
	db        *gorm.DB
	converter *Converter
}

func NewCategoryRepository(logger *logger.Logger, db *gorm.DB) ICategoryRepository {
	return &CategoryRepository{
		logger:    logger,
		db:        db,
		converter: &Converter{},
	}
}

func (r *CategoryRepository) Create(ctx context.Context, category *product_entity.Category) error {
	r.logger.Infof("Creating category: %v", category)
	categoryModel := r.converter.ToModel(category)
	if err := r.db.WithContext(ctx).Create(categoryModel).Error; err != nil {
		r.logger.Errorf("Failed to create category: %v", err)
		return err
	}
	category.ID = categoryModel.ID
	r.logger.Info("Category created successfully")
	return nil
}

func (r *CategoryRepository) GetByID(ctx context.Context, id string) (*product_entity.Category, error) {
	r.logger.Infof("Getting category: %s", id)
	var categoryModel product_model.Category
	if err := r.db.WithContext(ctx).First(&categoryModel, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Warnf("Category not found: %s", id)
			return nil, nil
		}
		r.logger.Errorf("Failed to get category: %v", err)
		return nil, err
	}
	category := r.converter.Toproduct_entity(&categoryModel)
	r.logger.Info("Category got successfully")
	return category, nil
}

func (r *CategoryRepository) GetAll(ctx context.Context, offset, limit int) ([]product_entity.Category, error) {
	r.logger.Infof("Getting all categories")
	var categoryModels []product_model.Category
	if limit == 0 {
		limit = -1
	}
	if offset == 0 {
		offset = -1
	}
	if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&categoryModels).Error; err != nil {
		r.logger.Errorf("Failed to get categories: %v", err)
		return nil, err
	}
	categories := make([]product_entity.Category, len(categoryModels))
	for i, categoryModel := range categoryModels {
		categories[i] = *r.converter.Toproduct_entity(&categoryModel)
	}
	r.logger.Info("Categories got successfully")
	return categories, nil
}

func (r *CategoryRepository) Update(ctx context.Context, category *product_entity.Category) error {
	r.logger.Infof("Updating category: %v", category)
	categoryModel := r.converter.ToModel(category)
	if err := r.db.WithContext(ctx).Model(&product_model.Category{}).Where("id = ?", category.ID).Updates(categoryModel).Error; err != nil {
		r.logger.Errorf("Failed to update category: %v", err)
		return err
	}
	r.logger.Info("Category updated successfully")
	return nil
}

func (r *CategoryRepository) Delete(ctx context.Context, id string) error {
	r.logger.Infof("Deleting category: %s", id)
	if err := r.db.WithContext(ctx).Delete(&product_model.Category{}, "id = ?", id).Error; err != nil {
		r.logger.Errorf("Failed to delete category: %v", err)
		return err
	}
	r.logger.Info("Category deleted successfully")
	return nil
}

func (r *CategoryRepository) GetByName(ctx context.Context, name string) (*product_entity.Category, error) {
	r.logger.Infof("Getting category: %s", name)
	var categoryModel product_model.Category
	if err := r.db.WithContext(ctx).First(&categoryModel, "name = ?", name).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Warnf("Category not found: %s", name)
			return nil, nil
		}
		r.logger.Errorf("Failed to get category: %v", err)
		return nil, err
	}
	category := r.converter.Toproduct_entity(&categoryModel)
	r.logger.Info("Category got successfully")
	return category, nil
}
