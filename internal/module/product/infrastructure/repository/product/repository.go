package product_repository

import (
	"context"

	product_entity "github.com/Fi44er/sdmed/internal/module/product/entity"
	product_model "github.com/Fi44er/sdmed/internal/module/product/infrastructure/repository/model"
	"github.com/Fi44er/sdmed/pkg/logger"
	"gorm.io/gorm"
)

type IProductRepository interface {
	Create(ctx context.Context, product *product_entity.Product) error
	GetByID(ctx context.Context, id string) (*product_entity.Product, error)
	GetAll(ctx context.Context, limit, offset int) ([]product_entity.Product, error)
	Update(ctx context.Context, product *product_entity.Product) error
	Delete(ctx context.Context, id string) error
	GetByArticle(ctx context.Context, article string) (*product_entity.Product, error)
}

type ProductRepository struct {
	logger    *logger.Logger
	db        *gorm.DB
	converter *Converter
}

func NewProductRepository(logger *logger.Logger, db *gorm.DB) IProductRepository {
	return &ProductRepository{
		logger:    logger,
		db:        db,
		converter: &Converter{},
	}
}

func (r *ProductRepository) Create(ctx context.Context, product *product_entity.Product) error {
	r.logger.Infof("Creating product: %+v", product)
	productModel := r.converter.ToModel(product)
	if err := r.db.WithContext(ctx).Create(productModel).Error; err != nil {
		r.logger.Errorf("Error creating product: %v", err)
		return err
	}
	product.ID = productModel.ID
	r.logger.Info("Product created successfully")
	return nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id string) (*product_entity.Product, error) {
	r.logger.Infof("Getting product: %s", id)
	var productModel product_model.Product
	if err := r.db.WithContext(ctx).First(&productModel, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Warnf("Product not found: %s", id)
			return nil, nil
		}
		r.logger.Errorf("Error getting product: %v", err)
		return nil, err
	}
	product := r.converter.ToEntity(&productModel)
	r.logger.Info("Product got successfully")
	return product, nil
}

func (r *ProductRepository) GetAll(ctx context.Context, limit, offset int) ([]product_entity.Product, error) {
	r.logger.Infof("Getting all products")
	var productModels []product_model.Product
	if limit == 0 {
		limit = -1
	}
	if offset == 0 {
		offset = -1
	}
	if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&productModels).Error; err != nil {
		r.logger.Errorf("Error getting products: %v", err)
		return nil, err
	}
	products := make([]product_entity.Product, len(productModels))
	for i, productModel := range productModels {
		products[i] = *r.converter.ToEntity(&productModel)
	}
	r.logger.Info("Products got successfully")
	return products, nil
}

func (r *ProductRepository) Update(ctx context.Context, product *product_entity.Product) error {
	r.logger.Infof("Updating product: %+v", product)
	productModel := r.converter.ToModel(product)
	if err := r.db.WithContext(ctx).Model(&product_model.Product{}).Where("id = ?", product.ID).Updates(productModel).Error; err != nil {
		r.logger.Errorf("Error updating product: %v", err)
		return err
	}
	r.logger.Info("Product updated successfully")
	return nil
}

func (r *ProductRepository) Delete(ctx context.Context, id string) error {
	r.logger.Infof("Deleting product: %s", id)
	if err := r.db.WithContext(ctx).Delete(&product_model.Product{}, "id = ?", id).Error; err != nil {
		r.logger.Errorf("Error deleting product: %v", err)
		return err
	}
	r.logger.Info("Product deleted successfully")
	return nil
}

func (r *ProductRepository) GetByArticle(ctx context.Context, article string) (*product_entity.Product, error) {
	r.logger.Infof("Getting product by article: %s", article)
	var productModel product_model.Product
	if err := r.db.WithContext(ctx).Where("article = ?", article).First(&productModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Warnf("Product not found: %s", article)
			return nil, nil
		}
		r.logger.Errorf("Error getting product by article: %v", err)
		return nil, err
	}
	product := r.converter.ToEntity(&productModel)
	r.logger.Info("Product got successfully")
	return product, nil
}
