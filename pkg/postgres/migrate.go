package postgres

import (
	"fmt"

	file_model "github.com/Fi44er/sdmed/internal/module/file/infrastructure/repository/model"
	product_model "github.com/Fi44er/sdmed/internal/module/product/infrastructure/repository/model"
	user_model "github.com/Fi44er/sdmed/internal/module/user/infrastructure/repository/model"
	"github.com/Fi44er/sdmed/pkg/logger"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB, trigger bool, log *logger.Logger) error {

	if trigger {
		log.Info("📦 Migrating database...")
		models := []any{
			user_model.Permission{},
			user_model.Role{},
			user_model.User{},
			file_model.File{},
			product_model.Product{},
			product_model.Category{},
		}
		schemas := []string{"user_module", "file_module", "product_module"}

		log.Info("📦 Creating types...")

		db.Exec("CREATE TYPE file_status AS ENUM ('temporary', 'permanent')")
		db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
		for _, schema := range schemas {
			if err := db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %q", schema)).Error; err != nil {
				return fmt.Errorf("failed to create schema %s: %w", schema, err)
			}
		}

		if err := db.AutoMigrate(models...); err != nil {
			log.Errorf("✖ Failed to migrate database: %v", err)
			return err
		}
	}

	log.Info("✅ Database connection successfully")
	return nil
}
