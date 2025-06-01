package postgres

import (
	"github.com/Fi44er/sdmed/internal/module/user/infrastructure/repository/model"
	"github.com/Fi44er/sdmed/pkg/logger"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB, trigger bool, log *logger.Logger) error {

	if trigger {
		log.Info("📦 Migrating database...")
		models := []interface{}{
			model.Permission{},
			model.Role{},
			model.User{},
		}

		log.Info("📦 Creating types...")

		db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
		if err := db.Exec("CREATE SCHEMA IF NOT EXISTS \"user_module\"").Error; err != nil {
			return err
		}

		if err := db.AutoMigrate(models...); err != nil {
			log.Errorf("✖ Failed to migrate database: %v", err)
			return err
		}
	}

	log.Info("✅ Database connection successfully")
	return nil
}
