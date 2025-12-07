package characteristic_usecase

import (
	"context"

	product_entity "github.com/Fi44er/sdmed/internal/module/product/entity"
	product_constant "github.com/Fi44er/sdmed/internal/module/product/pkg"
	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/Fi44er/sdmed/pkg/postgres/uow"
)

var ownerType string = "characteristic"

type ICharacteristicRepository interface {
	Create(ctx context.Context, characteristic *product_entity.Characteristic) error
	CreateMany(ctx context.Context, characteristics []product_entity.Characteristic) error
	Update(ctx context.Context, characteristic *product_entity.Characteristic) error
	Delete(ctx context.Context, id string) error
	DeleteByCategory(ctx context.Context, categoryID string) error
	GetByID(ctx context.Context, id string) (*product_entity.Characteristic, error)
	GetByCategoryID(ctx context.Context, categoryID string) ([]product_entity.Characteristic, error)
	GetByCategoryAndName(ctx context.Context, categoryID, name string) (*product_entity.Characteristic, error)
}

type ICharacteristicUsecase interface {
	Create(ctx context.Context, characteristic *product_entity.Characteristic) error
	CreateMany(ctx context.Context, characteristics []product_entity.Characteristic) error
	Delete(ctx context.Context, id string) error
	DeleteByCategory(ctx context.Context, categoryID string) error
}

type CharacteristicUsecase struct {
	repository ICharacteristicRepository
	uow        uow.Uow
	logger     *logger.Logger
}

func NewCharacteristicUsecase(repository ICharacteristicRepository, uow uow.Uow, logger *logger.Logger) ICharacteristicUsecase {
	return &CharacteristicUsecase{
		repository: repository,
		uow:        uow,
		logger:     logger,
	}
}

func (u *CharacteristicUsecase) Delete(ctx context.Context, id string) error {
	return u.uow.Do(ctx, func(ctx context.Context) error {
		repo, err := u.uow.GetRepository(ctx, ownerType)
		if err != nil {
			u.logger.Errorf("Failed to get repository: %v", err)
			return err
		}

		characteristicRepo := repo.(ICharacteristicRepository)

		needCleanup := true
		defer func() {
			if needCleanup {
				u.logger.Warnf("Cleaning up characteristic due to failed creation: %s", id)
				if err := characteristicRepo.Delete(ctx, id); err != nil {
					u.logger.Errorf("Failed to delete characteristic: %v", err)
				}
			}
		}()

		if err := characteristicRepo.Delete(ctx, id); err != nil {
			u.logger.Errorf("Failed to delete characteristic %s: %v", id, err)
			return err
		}

		return nil
	})
}

func (u *CharacteristicUsecase) DeleteByCategory(ctx context.Context, categoryID string) error {
	u.logger.Infof("Deleting characteristics by category: %s", categoryID)

	return u.uow.Do(ctx, func(ctx context.Context) error {
		repo, err := u.uow.GetRepository(ctx, ownerType)
		if err != nil {
			u.logger.Errorf("Failed to get repository: %v", err)
			return err
		}

		characteristicRepo := repo.(ICharacteristicRepository)

		if err := characteristicRepo.DeleteByCategory(ctx, categoryID); err != nil {
			u.logger.Errorf("Failed to delete characteristics by category %s: %v", categoryID, err)
			return err
		}

		return nil
	})
}

func (u *CharacteristicUsecase) Create(ctx context.Context, characteristic *product_entity.Characteristic) error {
	u.logger.Infof("Creating characteristic: %s", characteristic.Name)

	return u.uow.Do(ctx, func(ctx context.Context) error {
		repo, err := u.uow.GetRepository(ctx, ownerType)
		if err != nil {
			u.logger.Errorf("Failed to get repository: %v", err)
			return err
		}

		characteristicRepo := repo.(ICharacteristicRepository)

		needCleanup := true
		defer func() {
			if needCleanup {
				u.logger.Warnf("Cleaning up characteristic due to failed creation: %s", characteristic.ID)
				if err := characteristicRepo.Delete(ctx, characteristic.ID); err != nil {
					u.logger.Errorf("Failed to delete characteristic %s: %v", characteristic.ID, err)
				}
			}
		}()

		existCharacteristic, err := characteristicRepo.GetByCategoryAndName(ctx, characteristic.CategoryID, characteristic.Name)
		if err != nil {
			u.logger.Errorf("Failed to check characteristic existence by category %s and name %s: %v", characteristic.CategoryID, characteristic.Name, err)
			return err
		}

		if existCharacteristic != nil {
			u.logger.Warnf("Characteristic already exists: %s", existCharacteristic.Name)
			return product_constant.ErrCharacteristicAlreadyExists
		}

		if err := characteristicRepo.Create(ctx, characteristic); err != nil {
			u.logger.Errorf("Failed to create characteristic in repository: %v", err)
			return err
		}

		needCleanup = false
		u.logger.Infof("Characteristic created successfully: %s (ID: %s)", characteristic.Name, characteristic.ID)

		return nil
	})
}

func (u *CharacteristicUsecase) CreateMany(ctx context.Context, characteristics []product_entity.Characteristic) error {
	u.logger.Info("Creating characteristics")

	return u.uow.Do(ctx, func(ctx context.Context) error {
		repo, err := u.uow.GetRepository(ctx, ownerType)
		if err != nil {
			u.logger.Errorf("Failed to get repository: %v", err)
			return err
		}

		characteristicRepo := repo.(ICharacteristicRepository)

		needCleanup := true
		newCharacteristics := make([]product_entity.Characteristic, 0, len(characteristics))

		defer func() {
			if needCleanup {
				for _, characteristic := range newCharacteristics {
					u.logger.Warnf("Cleaning up characteristic due to failed creation: %s", characteristic.ID)
					if err := characteristicRepo.Delete(ctx, characteristic.ID); err != nil {
						u.logger.Errorf("Failed to delete characteristic %s: %v", characteristic.ID, err)
					}
				}

			}
		}()

		for _, characteristic := range characteristics {
			existCharacteristic, err := characteristicRepo.GetByCategoryAndName(ctx, characteristic.CategoryID, characteristic.Name)
			if err != nil {
				u.logger.Errorf("Failed to check characteristic existence by category %s and name %s: %v", characteristic.CategoryID, characteristic.Name, err)
				return err
			}
			if existCharacteristic == nil {
				newCharacteristics = append(newCharacteristics, characteristic)
			}
		}

		if err := characteristicRepo.CreateMany(ctx, newCharacteristics); err != nil {
			u.logger.Errorf("Failed to create characteristics in repository: %v", err)
			return err
		}

		needCleanup = false

		u.logger.Infof("Characteristics created successfully: %d", len(newCharacteristics))

		return nil
	})
}
