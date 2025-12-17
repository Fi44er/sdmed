package char_value_usecase

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"

	product_entity "github.com/Fi44er/sdmed/internal/module/product/entity"
	product_constant "github.com/Fi44er/sdmed/internal/module/product/pkg"
	char_value_usecase_contracts "github.com/Fi44er/sdmed/internal/module/product/usecase/char_value/contracts"
	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/Fi44er/sdmed/pkg/postgres/uow"
)

var ownerType = "char_value"

type ICharValueUsecase interface {
	CreateMany(ctx context.Context, charValues []product_entity.ProductCharValue) error
}

type CharValueUsecase struct {
	logger                *logger.Logger
	repository            char_value_usecase_contracts.ICharValueRepository
	uow                   uow.Uow
	characteristicUsecase char_value_usecase_contracts.ICharacteristicUsecase
}

func NewCharValueUsecase(logger *logger.Logger, repository char_value_usecase_contracts.ICharValueRepository, uow uow.Uow) ICharValueUsecase {
	return &CharValueUsecase{
		logger:     logger,
		repository: repository,
		uow:        uow,
	}
}

func (u *CharValueUsecase) CreateMany(ctx context.Context, charValues []product_entity.ProductCharValue) error {
	u.logger.Info("Creating char values")

	return u.uow.Do(ctx, func(ctx context.Context) error {
		repo, err := u.uow.GetRepository(ctx, ownerType)
		if err != nil {
			u.logger.Errorf("Failed to get repository: %v", err)
			return err
		}

		charValueRepo := repo.(char_value_usecase_contracts.ICharValueRepository)
		needCleanup := true
		defer func() {
			if needCleanup {
				for _, val := range charValues {
					u.logger.Warnf("Cleaning up characteristic value due to failed creation: %s", val.ID)
					if err := charValueRepo.Delete(ctx, val.ID); err != nil {
						u.logger.Errorf("Failed to cleanup characteristic  %s: %v", val.ID, err)
					}
				}
			}
		}()

		characteristicIDs := make([]string, len(charValues))
		for i, val := range charValues {
			characteristicIDs[i] = val.CharacteristicID
		}

		characteristics, err := u.characteristicUsecase.GetByIDs(ctx, characteristicIDs)
		if err != nil {
			u.logger.Errorf("Failed to get characteristics: %v", err)
			return err
		}

		characteristicsMap := make(map[string]*product_entity.Characteristic)
		optionsMap := make(map[string]map[string]bool) // characteristicID -> optionValue -> exists

		for _, characteristic := range characteristics {
			characteristicsMap[characteristic.ID] = &characteristic

			if characteristic.DataType == product_entity.DataTypeSelect {
				optionsMap[characteristic.ID] = make(map[string]bool)
				for _, option := range characteristic.Options {
					optionsMap[characteristic.ID][strings.ToLower(option.Value)] = true
				}
			}
		}

		for i, value := range charValues {
			// Проверяем существует ли характеристика
			characteristic, exists := characteristicsMap[value.CharacteristicID]
			if !exists {
				return fmt.Errorf("characteristic with ID %s not found: %w",
					value.CharacteristicID, product_constant.ErrCharacteristicNotFound)
			}

			// Валидируем в зависимости от типа данных
			if err := u.validateValueByDataType(value, characteristic, optionsMap); err != nil {
				return fmt.Errorf("validation failed for characteristic %s (index %d): %w",
					characteristic.Name, i, err)
			}

			// Проверка обязательных полей
			if characteristic.IsRequired {
				if err := u.validateRequired(value, characteristic.DataType); err != nil {
					return fmt.Errorf("required characteristic %s is empty: %w",
						characteristic.Name, err)
				}
			}
		}

		if err := u.repository.CreateMany(ctx, charValues); err != nil {
			u.logger.Errorf("failed to create characteristic values: %v", err)
			return err
		}

		needCleanup = false
		// u.logger.Infof("Characteristic value created successfully: %s", charValue.ID)
		return nil
	})
}

func (u *CharValueUsecase) validateValueByDataType(
	value product_entity.ProductCharValue,
	characteristic *product_entity.Characteristic,
	optionsMap map[string]map[string]bool,
) error {
	switch characteristic.DataType {
	case product_entity.DataTypeString:
		return u.validateStringValue(value)

	case product_entity.DataTypeNumber:
		return u.validateNumberValue(value, characteristic.Unit)

	case product_entity.DataTypeBoolean:
		return u.validateBooleanValue(value)

	case product_entity.DataTypeSelect:
		return u.validateSelectValue(value, characteristic.ID, optionsMap)

	default:
		return fmt.Errorf("unknown data type %s: %w",
			string(characteristic.DataType), product_constant.ErrInvalidDataType)
	}
}

func (u *CharValueUsecase) validateStringValue(value product_entity.ProductCharValue) error {
	if value.StringValue != nil && *value.StringValue != "" {
		strVal := *value.StringValue
		if len(strVal) > 1000 {
			return fmt.Errorf("string value too long (max 1000 chars): %w", product_constant.ErrInvalidString)
		}

		if strings.Contains(strVal, "<script>") || strings.Contains(strVal, "javascript:") {
			return fmt.Errorf("string contains dangerous characters: %w", product_constant.ErrInvalidString)
		}
	}

	if value.NumberValue != nil || value.BooleanValue != nil || value.OptionID != nil {
		return fmt.Errorf("string characteristic should not have number/bool/option values: %w", product_constant.ErrInvalidString)
	}

	return nil
}

func (u *CharValueUsecase) validateNumberValue(value product_entity.ProductCharValue, unit *string) error {
	if value.NumberValue != nil {
		numVal := *value.NumberValue

		if math.IsNaN(numVal) || math.IsInf(numVal, 0) {
			return fmt.Errorf("invalid number value (NaN or Infinity): %w", product_constant.ErrInvalidNumber)
		}

		if numVal < -1000000 || numVal > 1000000 {
			return fmt.Errorf("number value out of range (-1,000,000 to 1,000,000): %w", product_constant.ErrInvalidNumber)
		}

		// if unit != nil && *unit != "" {
		// 	switch *unit {
		// 	case "кг", "г", "л", "м", "см":
		// 		if numVal < 0 {
		// 			return fmt.Errorf("negative value not allowed for unit %s: %w", *unit, product_constant.ErrInvalidNumber)
		// 		}
		// 	}
		// }
	}

	if value.StringValue != nil && *value.StringValue != "" {
		if _, err := strconv.ParseFloat(*value.StringValue, 64); err == nil {
			return fmt.Errorf("number value should be in NumberValue field, not StringValue: %w", product_constant.ErrInvalidNumber)
		}
	}

	if value.BooleanValue != nil || value.OptionID != nil {
		return fmt.Errorf("number characteristic should not have bool/option values: %w", product_constant.ErrInvalidNumber)
	}

	return nil
}

func (u *CharValueUsecase) validateBooleanValue(value product_entity.ProductCharValue) error {
	if value.StringValue != nil && *value.StringValue != "" {
		strVal := strings.ToLower(*value.StringValue)
		if strVal == "true" || strVal == "false" || strVal == "1" || strVal == "0" {
			return fmt.Errorf("boolean value should be in BooleanValue field, not StringValue: %w", product_constant.ErrInvalidBoolean)
		}
	}

	if value.NumberValue != nil || value.OptionID != nil {
		return fmt.Errorf("boolean characteristic should not have number/option values: %w", product_constant.ErrInvalidBoolean)
	}

	return nil
}

func (u *CharValueUsecase) validateSelectValue(
	value product_entity.ProductCharValue,
	characteristicID string,
	optionsMap map[string]map[string]bool,
) error {
	if value.OptionID != nil && *value.OptionID != "" {
		u.logger.Warnf("Option with ID %s not validated in-depth", *value.OptionID)
	}

	if value.StringValue != nil && *value.StringValue != "" {
		strVal := *value.StringValue

		if optMap, exists := optionsMap[characteristicID]; exists {
			if !optMap[strings.ToLower(strVal)] {
				return fmt.Errorf("value '%s' is not in available options: %w",
					strVal, product_constant.ErrOptionNotFound)
			}
		} else {
			u.logger.Warnf("Options map not found for characteristic %s", characteristicID)
		}
	}

	if value.NumberValue != nil || value.BooleanValue != nil {
		return fmt.Errorf("select characteristic should not have number/bool values: %w", product_constant.ErrOptionNotFound)
	}

	return nil
}

func (u *CharValueUsecase) validateRequired(value product_entity.ProductCharValue, dataType product_entity.DataType) error {
	switch dataType {
	case product_entity.DataTypeString:
		if value.StringValue == nil || *value.StringValue == "" {
			return product_constant.ErrValueRequired
		}

	case product_entity.DataTypeNumber:
		if value.NumberValue == nil {
			return product_constant.ErrValueRequired
		}

	case product_entity.DataTypeBoolean:
		if value.BooleanValue == nil {
			return product_constant.ErrValueRequired
		}

	case product_entity.DataTypeSelect:
		if (value.OptionID == nil || *value.OptionID == "") &&
			(value.StringValue == nil || *value.StringValue == "") {
			return product_constant.ErrValueRequired
		}
	}

	return nil
}
