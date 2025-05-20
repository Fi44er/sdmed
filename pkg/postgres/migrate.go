package postgres

import (
	"github.com/Fi44er/sdmedik/backend/internal/module/user/infrastructure/repository/model"
	"github.com/Fi44er/sdmedik/backend/pkg/logger"
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
