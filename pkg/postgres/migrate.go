package postgres

import (
	auth_models "github.com/Fi44er/sdmed/internal/module/auth/infrastucture/repository/models"
	cart_model "github.com/Fi44er/sdmed/internal/module/cart/infrastructure/repository/model"
	file_model "github.com/Fi44er/sdmed/internal/module/file/infrastructure/repository/model"
	order_model "github.com/Fi44er/sdmed/internal/module/order/infrastructure/repository/model"
	product_model "github.com/Fi44er/sdmed/internal/module/product/infrastructure/repository/model"
	scraper_models "github.com/Fi44er/sdmed/internal/module/scraper/infrastructure/repository/models"
	tru_model "github.com/Fi44er/sdmed/internal/module/tru/infrastructure/repository/model"
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
			scraper_models.ScraperSettings{},

			tru_model.TRUCode{},
			tru_model.TRUCodePrice{},
			cart_model.Cart{},
			cart_model.CartItem{},

			order_model.Order{},
			order_model.OrderItem{},
		}

		log.Info("📦 Creating types...")

		db.Exec("CREATE TYPE file_status AS ENUM ('temporary', 'permanent')")
		db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")

		if err := db.AutoMigrate(models...); err != nil {
			log.Errorf("✖ Failed to migrate database: %v", err)
			return err
		}
	}

	// Always seed default roles and permissions if they do not exist
	if err := SeedRolesAndPermissions(db, log); err != nil {
		log.Errorf("✖ Failed to seed roles and permissions: %v", err)
		return err
	}

	log.Info("✅ Database connection successfully")
	return nil
}

func SeedRolesAndPermissions(db *gorm.DB, log *logger.Logger) error {
	var count int64
	if err := db.Model(&user_model.Role{}).Count(&count).Error; err != nil {
		log.Errorf("Failed to count roles: %v", err)
		return err
	}
	if count > 0 {
		log.Info("Roles already seeded")
		return nil
	}

	log.Info("🌱 Seeding roles and permissions...")

	// 1. Create all unique permissions
	permissionsList := []string{
		"*:*",
		"orders:all",
		"orders:my",
		"products:write",
		"products:read",
		"categories:write",
		"categories:read",
		"scraper:read",
		"scraper:write",
		"users:all",
	}

	permMap := make(map[string]user_model.Permission)
	for _, name := range permissionsList {
		perm := user_model.Permission{
			Name: name,
		}
		if err := db.Create(&perm).Error; err != nil {
			log.Errorf("Failed to create permission %s: %v", name, err)
			return err
		}
		permMap[name] = perm
	}

	// 2. Define roles and their associated permissions
	rolePermissions := map[string][]string{
		"admin": {
			"*:*",
		},
		"manager": {
			"orders:all",
			"products:write",
			"products:read",
			"categories:write",
			"categories:read",
			"scraper:read",
			"scraper:write",
		},
		"user": {
			"orders:my",
			"products:read",
			"categories:read",
		},
		"guest": {
			"products:read",
			"categories:read",
		},
	}

	for roleName, permNames := range rolePermissions {
		var rolePerms []user_model.Permission
		for _, name := range permNames {
			if perm, ok := permMap[name]; ok {
				rolePerms = append(rolePerms, perm)
			}
		}

		role := user_model.Role{
			Name:        roleName,
			Permissions: rolePerms,
		}
		if err := db.Create(&role).Error; err != nil {
			log.Errorf("Failed to create role %s: %v", roleName, err)
			return err
		}
	}

	log.Info("✅ Seeding roles and permissions completed successfully")
	return nil
}

