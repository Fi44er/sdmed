package postgres

import (
	auth_models "github.com/Fi44er/sdmed/internal/module/auth/infrastucture/repository/models"
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
			user_model.Role{},
			user_model.User{},
			user_model.Permission{},

			auth_models.UserSession{},
			file_model.File{},
			product_model.Product{},
			product_model.Category{},
			product_model.Characteristic{},
			product_model.CharacteristicValue{},
			product_model.CharOption{},
		}

		log.Info("📦 Creating types...")

		db.Exec("CREATE TYPE file_status AS ENUM ('temporary', 'permanent')")
		db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")

		if err := db.AutoMigrate(models...); err != nil {
			log.Errorf("✖ Failed to migrate database: %v", err)
			return err
		}
	}

	log.Info("✅ Database connection successfully")
	return nil
}
