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

func NewCharValueUsecase(
	logger *logger.Logger,
	repository char_value_usecase_contracts.ICharValueRepository,
	uow uow.Uow,
	characteristicUsecase char_value_usecase_contracts.ICharacteristicUsecase,
) ICharValueUsecase {
	return &CharValueUsecase{
		logger:                logger,
		repository:            repository,
		uow:                   uow,
		characteristicUsecase: characteristicUsecase,
	}
}

func (u *CharValueUsecase) CreateMany(ctx context.Context, charValues []product_entity.ProductCharValue) error {
	u.logger.Info("Creating char values")

	characteristicIDs := make([]string, len(charValues))
	for i, val := range charValues {
		characteristicIDs[i] = val.CharacteristicID
	}

	characteristics, err := u.characteristicUsecase.GetByIDs(ctx, characteristicIDs)
	if err != nil {
		u.logger.Errorf("Failed to get characteristics: %v", err)
		return err
	}

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
					if val.ID != "" {
						if err := charValueRepo.Delete(ctx, val.ID); err != nil {
							u.logger.Errorf("Failed to cleanup characteristic  %s: %v", val.ID, err)
						}
					}
				}
			}
		}()

		characteristicsMap := make(map[string]*product_entity.Characteristic)
		optionsMap := make(map[string]map[string]string)

		for _, characteristic := range characteristics {
			characteristicsMap[characteristic.ID] = &characteristic
			// u.logger.Debugf("Characteristci: %+v", characteristic)
			if characteristic.DataType == product_entity.DataTypeSelect {
				optionsMap[characteristic.ID] = make(map[string]string)
				for _, option := range characteristic.Options {
					optionsMap[characteristic.ID][strings.ToLower(option.Value)] = option.ID
				}
			}
		}

		for i := range charValues {
			value := &charValues[i]
			characteristic, exists := characteristicsMap[value.CharacteristicID]
			if !exists {
				return fmt.Errorf("characteristic with ID %s not found: %w",
					value.CharacteristicID, product_constant.ErrCharacteristicNotFound)
			}

			if err := u.validateValueByDataType(value, characteristic, optionsMap); err != nil {
				return fmt.Errorf("validation failed for characteristic %s (index %d): %w",
					characteristic.Name, i, err)
			}

			if characteristic.IsRequired {
				if err := u.validateRequired(value, characteristic.DataType); err != nil {
					return fmt.Errorf("required characteristic %s is empty: %w",
						characteristic.Name, err)
				}
			}
		}

		if err := charValueRepo.CreateMany(ctx, charValues); err != nil {
			u.logger.Errorf("failed to create characteristic values: %v", err)
			return err
		}

		needCleanup = false
		return nil
	})
}

func (u *CharValueUsecase) validateValueByDataType(
	value *product_entity.ProductCharValue,
	characteristic *product_entity.Characteristic,
	optionsMap map[string]map[string]string,
) error {
	if value.StringValue == nil || *value.StringValue == "" {
		return product_constant.ErrInvalidString
	}

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

func (u *CharValueUsecase) validateStringValue(value *product_entity.ProductCharValue) error {
	strVal := *value.StringValue
	if len(strVal) > 1000 {
		return fmt.Errorf("string value too long (max 1000 chars): %w", product_constant.ErrInvalidString)
	}

	if strings.Contains(strVal, "<script>") || strings.Contains(strVal, "javascript:") {
		return fmt.Errorf("string contains dangerous characters: %w", product_constant.ErrInvalidString)
	}

	u.clearOtherValueFields(value, ClearOptions{
		KeepString: true,
	})

	return nil
}

func (u *CharValueUsecase) validateNumberValue(value *product_entity.ProductCharValue, unit *string) error {
	numVal, err := strconv.ParseFloat(*value.StringValue, 64)
	if err != nil {
		return product_constant.ErrInvalidNumber
	}

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

	value.NumberValue = &numVal

	u.clearOtherValueFields(value, ClearOptions{
		KeepNumber: true,
	})

	return nil
}

func (u *CharValueUsecase) validateBooleanValue(value *product_entity.ProductCharValue) error {
	var boolStringMap = map[string]bool{
		"true":  true,
		"1":     true,
		"yes":   true,
		"y":     true,
		"on":    true,
		"false": false,
		"0":     false,
		"no":    false,
		"n":     false,
		"off":   false,
	}

	strVal := strings.TrimSpace(strings.ToLower(*value.StringValue))

	boolVal, exists := boolStringMap[strVal]
	if !exists {
		return fmt.Errorf("invalid boolean string '%s': %w", *value.StringValue, product_constant.ErrInvalidBoolean)
	}

	value.BooleanValue = &boolVal
	u.clearOtherValueFields(value, ClearOptions{
		KeepBoolean: true,
	})

	return nil
}

func (u *CharValueUsecase) validateSelectValue(
	value *product_entity.ProductCharValue,
	characteristicID string,
	optionsMap map[string]map[string]string,
) error {
	strVal := *value.StringValue

	if optMap, exists := optionsMap[characteristicID]; exists {
		if optionID := optMap[strings.ToLower(strVal)]; optionID != "" {
			value.OptionID = &optionID
		} else {
			return fmt.Errorf("value '%s' is not in available options: %w",
				strVal, product_constant.ErrOptionNotFound)
		}
	} else {
		u.logger.Warnf("Options map not found for characteristic %s", characteristicID)
	}

	u.clearOtherValueFields(value, ClearOptions{
		KeepOptionID: true,
	})

	return nil
}

func (u *CharValueUsecase) validateRequired(value *product_entity.ProductCharValue, dataType product_entity.DataType) error {
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

type ClearOptions struct {
	KeepString   bool
	KeepNumber   bool
	KeepBoolean  bool
	KeepOptionID bool
}

func (u *CharValueUsecase) clearOtherValueFields(value *product_entity.ProductCharValue, opts ClearOptions) {
	if !opts.KeepString {
		value.StringValue = nil
	}
	if !opts.KeepNumber {
		value.NumberValue = nil
	}
	if !opts.KeepBoolean {
		value.BooleanValue = nil
	}
	if !opts.KeepOptionID {
		value.OptionID = nil
	}
}
