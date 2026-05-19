package cart_usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	cart_dto "github.com/Fi44er/sdmed/internal/module/cart/dto"
	cart_entity "github.com/Fi44er/sdmed/internal/module/cart/entity"
	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/go-playground/validator/v10"
)

type CartUsecase struct {
	repo           ICartRepository
	productAdapter IProductUsecaseAdapter
	truAdapter     ITRUCodeUsecaseAdapter
	promoAdapter   IPromotionUsecaseAdapter
	logger         *logger.Logger
	validator      *validator.Validate
}

func NewCartUsecase(
	repo ICartRepository,
	productAdapter IProductUsecaseAdapter,
	truAdapter ITRUCodeUsecaseAdapter,
	promoAdapter IPromotionUsecaseAdapter,
	logger *logger.Logger,
	validator *validator.Validate,
) *CartUsecase {
	return &CartUsecase{
		repo:           repo,
		productAdapter: productAdapter,
		truAdapter:     truAdapter,
		promoAdapter:   promoAdapter,
		logger:         logger,
		validator:      validator,
	}
}

func (u *CartUsecase) GetCart(ctx context.Context, userID string, cartID string) (*cart_entity.Cart, error) {
	var cart *cart_entity.Cart
	var err error

	if userID != "" {
		cart, err = u.repo.GetByUserID(ctx, userID)
	} else if cartID != "" {
		cart, err = u.repo.GetByID(ctx, cartID)
	}

	if err != nil {
		return nil, err
	}

	if cart == nil {
		cart = &cart_entity.Cart{
			UserID: userID,
			Items:  []cart_entity.CartItem{},
		}
	}

	// Calculate TotalPrice
	cart.TotalPrice = 0
	for _, item := range cart.Items {
		cart.TotalPrice += item.UnitPrice * float64(item.Quantity)
	}

	return cart, nil
}

func (u *CartUsecase) AddItem(ctx context.Context, data *cart_dto.AddCartItemRequest, userID string, cartID string) (*cart_entity.Cart, error) {
	if err := u.validator.Struct(data); err != nil {
		return nil, err
	}

	cart, err := u.GetCart(ctx, userID, cartID)
	if err != nil {
		return nil, err
	}

	if cart.ID == "" {
		cart.UserID = userID
		if err := u.repo.Create(ctx, cart); err != nil {
			return nil, err
		}
	}

	product, err := u.productAdapter.GetByID(ctx, data.ProductID)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, fmt.Errorf("product not found")
	}

	// Validate options and calculate price
	// Note: In current product module, prices are not yet tied to options.
	// Using manual price as base.
	selectedPrice := 0.0
	priceResolved := false

	// Если у продукта привязан TRUCodeID, пытаемся найти региональную цену из ТРУ-кода
	if product.TRUCodeID != nil && *product.TRUCodeID != "" && data.Iso != "" {
		truCode, err := u.truAdapter.GetByID(ctx, *product.TRUCodeID)
		if err == nil && truCode != nil {
			for _, truPrice := range truCode.Prices {
				if strings.EqualFold(truPrice.RegionIso, data.Iso) {
					selectedPrice = truPrice.Price
					priceResolved = true
					u.logger.Infof("Resolved regional price for product %s in region %s: %f", product.Article, data.Iso, selectedPrice)
					break
				}
			}
		} else if err != nil {
			u.logger.Warnf("Failed to fetch TRU code %s: %v", *product.TRUCodeID, err)
		}
	}

	// Фолбек на ручную цену, если региональная цена не найдена или не привязана
	if !priceResolved {
		if product.ManualPrice != nil {
			selectedPrice = *product.ManualPrice
		}
	}

	// Simplify SelectedOptions for storage
	selectedOptions := make([]cart_entity.SelectedOptionResponse, 0)
	for _, opt := range data.DynamicOptions {
		// Find option name from product characteristics
		for _, char := range product.CharValues {
			if char.CharacteristicID == opt.ID {
				selectedOptions = append(selectedOptions, cart_entity.SelectedOptionResponse{
					Name:  char.CharacteristicName,
					Value: opt.Value,
				})
				break
			}
		}
	}

	optionsBytes, _ := json.Marshal(selectedOptions)
	existingItem, err := u.repo.GetItemByCriteria(ctx, data.ProductID, cart.ID, data.Iso, string(optionsBytes))
	if err != nil {
		return nil, err
	}

	if existingItem != nil {
		existingItem.Quantity += data.Quantity
		if existingItem.Quantity <= 0 {
			if err := u.repo.DeleteItem(ctx, existingItem.ID, cart.ID); err != nil {
				return nil, err
			}
		} else {
			existingItem.UnitPrice = selectedPrice
			if err := u.repo.UpdateItem(ctx, existingItem); err != nil {
				return nil, err
			}
		}
	} else {
		if data.Quantity > 0 {
			newItem := &cart_entity.CartItem{
				CartID:          cart.ID,
				ProductID:       data.ProductID,
				Article:         product.Article,
				Quantity:        data.Quantity,
				UnitPrice:       selectedPrice,
				SelectedOptions: selectedOptions,
				Iso:             data.Iso,
			}
			if err := u.repo.CreateItem(ctx, newItem); err != nil {
				return nil, err
			}
		}
	}

	// Update Cart TotalPrice in DB
	updatedCart, err := u.GetCart(ctx, userID, cart.ID)
	if err == nil && updatedCart != nil {
		_ = u.repo.Update(ctx, updatedCart)
	}

	return updatedCart, nil
}

func (u *CartUsecase) DeleteItem(ctx context.Context, itemID string, userID string, cartID string) (*cart_entity.Cart, error) {
	cart, err := u.GetCart(ctx, userID, cartID)
	if err != nil {
		return nil, err
	}

	if err := u.repo.DeleteItem(ctx, itemID, cart.ID); err != nil {
		return nil, err
	}

	return u.GetCart(ctx, userID, cart.ID)
}

func (u *CartUsecase) MoveByUserID(ctx context.Context, fromUserID string, toUserID string) error {
	if fromUserID == toUserID || fromUserID == "" || toUserID == "" {
		return nil
	}

	fromCart, err := u.repo.GetByUserID(ctx, fromUserID)
	if err != nil || fromCart == nil {
		return err
	}

	toCart, err := u.repo.GetByUserID(ctx, toUserID)
	if err != nil {
		return err
	}

	if toCart == nil {
		// Если у целевого пользователя нет корзины, просто перепривязываем гостевую
		fromCart.UserID = toUserID
		return u.repo.Update(ctx, fromCart)
	}

	// Если у целевого пользователя уже есть корзина, переносим элементы
	for _, item := range fromCart.Items {
		optionsBytes, _ := json.Marshal(item.SelectedOptions)
		existingItem, err := u.repo.GetItemByCriteria(ctx, item.ProductID, toCart.ID, item.Iso, string(optionsBytes))
		if err != nil {
			continue
		}

		if existingItem != nil {
			// Обновляем количество если такой товар уже есть
			existingItem.Quantity += item.Quantity
			_ = u.repo.UpdateItem(ctx, existingItem)
		} else {
			// Перепривязываем элемент к новой корзине
			item.CartID = toCart.ID
			_ = u.repo.UpdateItem(ctx, &item)
		}
	}

	if err := u.repo.Delete(ctx, fromCart.ID); err != nil {
		return err
	}

	// Update target cart TotalPrice
	updatedToCart, err := u.GetCart(ctx, toUserID, toCart.ID)
	if err == nil && updatedToCart != nil {
		_ = u.repo.Update(ctx, updatedToCart)
	}

	return nil
}

func (u *CartUsecase) Move(ctx context.Context, userID string, sessionCartID string) error {
	if sessionCartID == "" || userID == "" {
		return nil
	}

	sessCart, err := u.repo.GetByID(ctx, sessionCartID)
	if err != nil || sessCart == nil {
		return err
	}

	for _, item := range sessCart.Items {
		options, _ := json.Marshal(item.SelectedOptions)
		var dynOpts []cart_entity.DynamicOption
		// This is a bit simplified, but follows the logic of merging
		_ = json.Unmarshal(options, &dynOpts)

		_, err := u.AddItem(ctx, &cart_dto.AddCartItemRequest{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Iso:       item.Iso,
			// DynamicOptions would need to be reconstructed if we want full fidelity,
			// but for Move we can often just reparent the item or merge by criteria.
		}, userID, "")
		if err != nil {
			u.logger.Errorf("failed to move item %s: %v", item.ID, err)
		}
	}

	return u.repo.Delete(ctx, sessionCartID)
}

func (u *CartUsecase) Clear(ctx context.Context, cartID string) error {
	if cartID == "" {
		return nil
	}
	return u.repo.Delete(ctx, cartID)
}
