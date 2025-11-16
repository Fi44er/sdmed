package category_repository

import (
	"context"

	product_entity "github.com/Fi44er/sdmed/internal/module/product/entity"
	product_model "github.com/Fi44er/sdmed/internal/module/product/infrastructure/repository/model"
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
	r.logger.Infof("Creating category: %s", category.Name)

	categoryModel := r.converter.ToModel(category)
	if err := r.db.WithContext(ctx).Create(categoryModel).Error; err != nil {
		r.logger.Errorf("Failed to create category '%s': %v", category.Name, err)
		return err
	}
	category.ID = categoryModel.ID

	r.logger.Infof("Category created successfully: %s (ID: %s)", category.Name, category.ID)
	return nil
}

func (r *CategoryRepository) GetByID(ctx context.Context, id string) (*product_entity.Category, error) {
	r.logger.Debugf("Getting category by ID: %s", id)

	var categoryModel product_model.Category
	if err := r.db.WithContext(ctx).First(&categoryModel, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Debugf("Category not found: %s", id)
			return nil, nil
		}
		r.logger.Errorf("Failed to get category by ID %s: %v", id, err)
		return nil, err
	}
	category := r.converter.Toproduct_entity(&categoryModel)

	r.logger.Debugf("Category retrieved successfully: %s", id)
	return category, nil
}

func (r *CategoryRepository) GetAll(ctx context.Context, offset, limit int) ([]product_entity.Category, error) {
	r.logger.Debugf("Getting all categories (offset: %d, limit: %d)", offset, limit)

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

	r.logger.Debugf("Retrieved %d categories", len(categories))
	return categories, nil
}

func (r *CategoryRepository) Update(ctx context.Context, category *product_entity.Category) error {
	r.logger.Infof("Updating category: %s (ID: %s)", category.Name, category.ID)

	categoryModel := r.converter.ToModel(category)
	if err := r.db.WithContext(ctx).Model(&product_model.Category{}).Where("id = ?", category.ID).Updates(categoryModel).Error; err != nil {
		r.logger.Errorf("Failed to update category %s: %v", category.ID, err)
		return err
	}

	r.logger.Infof("Category updated successfully: %s", category.ID)
	return nil
}

func (r *CategoryRepository) Delete(ctx context.Context, id string) error {
	r.logger.Infof("Deleting category: %s", id)

	if err := r.db.WithContext(ctx).Delete(&product_model.Category{}, "id = ?", id).Error; err != nil {
		r.logger.Errorf("Failed to delete category %s: %v", id, err)
		return err
	}

	r.logger.Infof("Category deleted successfully: %s", id)
	return nil
}

func (r *CategoryRepository) GetByName(ctx context.Context, name string) (*product_entity.Category, error) {
	r.logger.Debugf("Getting category by name: %s", name)

	var categoryModel product_model.Category
	if err := r.db.WithContext(ctx).First(&categoryModel, "name = ?", name).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Debugf("Category not found by name: %s", name)
			return nil, nil
		}
		r.logger.Errorf("Failed to get category by name %s: %v", name, err)
		return nil, err
	}
	category := r.converter.Toproduct_entity(&categoryModel)

	r.logger.Debugf("Category found by name: %s (ID: %s)", name, category.ID)
	return category, nil
}
