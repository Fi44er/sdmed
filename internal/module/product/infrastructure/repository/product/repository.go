package product_repository

import (
	"context"
	"fmt"

	product_entity "github.com/Fi44er/sdmed/internal/module/product/entity"
	product_model "github.com/Fi44er/sdmed/internal/module/product/infrastructure/repository/model"
	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/Fi44er/sdmed/pkg/utils"
	"gorm.io/gorm"
)

type IProductRepository interface {
	Create(ctx context.Context, product *product_entity.Product) error
	GetByID(ctx context.Context, id string) (*product_entity.Product, error)
	GetAll(ctx context.Context, params product_entity.ProductFilterParams) ([]product_entity.Product, int64, error)
	Update(ctx context.Context, product *product_entity.Product) error
	Delete(ctx context.Context, id string) error
	GetByArticle(ctx context.Context, article string) (*product_entity.Product, error)
	GetBySlug(ctx context.Context, slug string) (*product_entity.Product, error)
	Count(ctx context.Context) (int64, error)
	GetFiltersByCategory(ctx context.Context, categoryID string) ([]product_entity.Filter, error)
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

func (r *ProductRepository) GetAll(ctx context.Context, params product_entity.ProductFilterParams) ([]product_entity.Product, int64, error) {
	r.logger.Infof("Getting products with filters")

	var productModels []product_model.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&product_model.Product{})

	if params.CategoryID != "" {
		query = query.Where("category_id = ?", params.CategoryID)
	}

	if params.MinPrice != nil {
		query = query.Where("manual_price >= ?", *params.MinPrice)
	}
	if params.MaxPrice != nil {
		query = query.Where("manual_price <= ?", *params.MaxPrice)
	}

	for charID, value := range params.Characteristics {
		subQuery := r.db.Table("characteristic_values").
			Select("1").
			Joins("LEFT JOIN char_options ON char_options.id = characteristic_values.option_id").
			Where("characteristic_values.product_id = products.id").
			Where("characteristic_values.characteristic_id = ?", charID).
			Where(
				"(characteristic_values.string_value = ? OR CAST(characteristic_values.number_value AS TEXT) = ? OR char_options.value = ?)",
				value, value, value,
			)

		query = query.Where("EXISTS (?)", subQuery)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	switch params.Sort {
	case "price_asc":
		query = query.Order("manual_price ASC")
	case "price_desc":
		query = query.Order("manual_price DESC")
	case "newest":
		query = query.Order("created_at DESC")
	default:
		query = query.Order("name ASC")
	}

	offset, limit := utils.SafeCalculateForPostgres(params.Page, params.PageSize)

	if limit <= 0 {
		limit = -1
	}
	if offset <= 0 {
		offset = -1
	}

	err := query.
		Limit(limit).
		Offset(offset).
		Preload("Characteristics", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("JOIN characteristics ON characteristics.id = characteristic_values.characteristic_id").
				Select("characteristic_values.*, characteristics.name as characteristic_name")
		}).
		Preload("Characteristics.Option").
		Find(&productModels).Error

	if err != nil {
		r.logger.Errorf("Error getting products: %v", err)
		return nil, 0, err
	}

	products := make([]product_entity.Product, len(productModels))
	for i, productModel := range productModels {
		products[i] = *r.converter.ToEntity(&productModel)
	}

	r.logger.Infof("Successfully retrieved %d products, total count: %d", len(products), total)
	return products, total, nil
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

func (r *ProductRepository) GetBySlug(ctx context.Context, slug string) (*product_entity.Product, error) {
	r.logger.Debugf("Getting product by Slug: %s", slug)

	var productModel product_model.Product
	err := r.db.WithContext(ctx).
		Preload("Characteristics", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("JOIN characteristics ON characteristics.id = characteristic_values.characteristic_id").
				Select("characteristic_values.*, characteristics.name as characteristic_name")
		}).
		Preload("Characteristics.Option").
		First(&productModel, "slug = ?", slug).
		Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Debugf("Product not found: %s", slug)
			return nil, nil
		}
		r.logger.Errorf("Failed to get product by Slug %s: %v", slug, err)
		return nil, err
	}

	product := r.converter.ToEntity(&productModel)
	r.logger.Debugf("Product retrieved successfully: %s", slug)
	return product, nil
}

func (r *ProductRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&product_model.Product{}).Count(&count).Error
	if err != nil {
		r.logger.Errorf("Failed to count products: %v", err)
		return 0, err
	}
	return count, nil
}

func (r *ProductRepository) GetFiltersByCategory(ctx context.Context, categoryID string) ([]product_entity.Filter, error) {
	var results []struct {
		CharacteristicID   string
		CharacteristicName string
		DataType           string
		Unit               string
		StringValue        string
		NumberValue        float64
		BooleanValue       bool
		OptionValue        string
	}

	err := r.db.WithContext(ctx).
		Table("characteristic_values").
		Select(`
			characteristics.id as characteristic_id,
			characteristics.name as characteristic_name,
			characteristics.data_type,
			characteristics.unit,
			characteristic_values.string_value,
			characteristic_values.number_value,
			characteristic_values.boolean_value,
			char_options.value as option_value
		`).
		Joins("JOIN products ON products.id = characteristic_values.product_id").
		Joins("JOIN characteristics ON characteristics.id = characteristic_values.characteristic_id").
		Joins("LEFT JOIN char_options ON char_options.id = characteristic_values.option_id").
		Where("products.category_id = ? AND products.is_active = ?", categoryID, true).
		Group(`
			characteristics.id, characteristics.name, characteristics.data_type, characteristics.unit,
			characteristic_values.string_value, characteristic_values.number_value,
			characteristic_values.boolean_value, char_options.value
		`).
		Order("characteristics.name ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	filterMap := make(map[string]*product_entity.Filter)
	finalFilters := make([]product_entity.Filter, 0, len(filterMap))

	for _, res := range results {
		if _, ok := filterMap[res.CharacteristicID]; !ok {
			f := &product_entity.Filter{
				CharacteristicID:   res.CharacteristicID,
				CharacteristicName: res.CharacteristicName,
				DataType:           res.DataType,
				Unit:               res.Unit,
				Options:            []string{},
			}
			filterMap[res.CharacteristicID] = f
		}

		var valStr string
		switch res.DataType {
		case string(product_entity.DataTypeSelect):
			valStr = res.OptionValue
		case string(product_entity.DataTypeNumber):
			valStr = fmt.Sprintf("%v", res.NumberValue)
		case string(product_entity.DataTypeBoolean):
			valStr = fmt.Sprintf("%v", res.BooleanValue)
		default:
			valStr = res.StringValue
		}

		if valStr != "" {
			filterMap[res.CharacteristicID].Options = append(filterMap[res.CharacteristicID].Options, valStr)
		}
	}

	for _, f := range filterMap {
		finalFilters = append(finalFilters, *f)
	}

	return finalFilters, nil
}
