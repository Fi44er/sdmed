package postgres

import (
	"fmt"

	file_model "github.com/Fi44er/sdmed/internal/module/file/infrastucture/repository/model"
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
			file_model.File{},
		}
		schemas := []string{"user_module", "file_module"}

		log.Info("📦 Creating types...")

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
