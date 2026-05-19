package order_usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/Fi44er/sdmed/internal/config"
	order_dto "github.com/Fi44er/sdmed/internal/module/order/dto"
	order_entity "github.com/Fi44er/sdmed/internal/module/order/entity"
	order_payment "github.com/Fi44er/sdmed/internal/module/order/infrastructure/payment"
	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/go-playground/validator/v10"
)

type OrderUsecase struct {
	repo                IOrderRepository
	paymentProvider     IPaymentProvider
	cartAdapter         ICartAdapter
	productAdapter      IProductAdapter
	truAdapter          ITRUAdapter
	notificationAdapter INotificationAdapter
	chatAdapter         IChatAdapter
	logger              *logger.Logger
	validator           *validator.Validate
	config              *config.Config
}

func NewOrderUsecase(
	repo IOrderRepository,
	paymentProvider IPaymentProvider,
	cartAdapter ICartAdapter,
	productAdapter IProductAdapter,
	truAdapter ITRUAdapter,
	notificationAdapter INotificationAdapter,
	chatAdapter IChatAdapter,
	logger *logger.Logger,
	validator *validator.Validate,
	config *config.Config,
) *OrderUsecase {
	return &OrderUsecase{
		repo:                repo,
		paymentProvider:     paymentProvider,
		cartAdapter:         cartAdapter,
		productAdapter:      productAdapter,
		truAdapter:          truAdapter,
		notificationAdapter: notificationAdapter,
		chatAdapter:         chatAdapter,
		logger:              logger,
		validator:           validator,
		config:              config,
	}
}

func (u *OrderUsecase) CreateOrder(ctx context.Context, data *order_dto.CreateOrderRequest, userID string) (string, error) {
	if err := u.validator.Struct(data); err != nil {
		return "", err
	}

	cart, err := u.cartAdapter.GetCart(ctx, userID)
	if err != nil {
		return "", err
	}
	if cart == nil || len(cart.Items) == 0 {
		return "", fmt.Errorf("cart is empty")
	}

	// Prepare PayKeeper items
	pkItems := make([]order_payment.PayKeeperItem, 0, len(cart.Items))
	for _, item := range cart.Items {
		categoryArticle := strings.Split(item.Article, ".")[0]
		truCode, _ := u.truAdapter.GetByCode(ctx, categoryArticle)
		if truCode != "" {
			truCode += "00000000643"
		}

		itemTotalPrice := item.UnitPrice * float64(item.Quantity)
		pkItems = append(pkItems, order_payment.PayKeeperItem{
			ItemType:    "goods",
			PaymentType: "full",
			Name:        item.ProductID,
			Price:       item.UnitPrice,
			Quantity:    item.Quantity,
			TruCode:     truCode,
			Tax:         "none",
			Sum:         itemTotalPrice,
		})
	}

	// Create PayKeeper invoice
	link, err := u.paymentProvider.CreateInvoice(ctx, data.Email, data.PhoneNumber, data.FIO, cart.TotalPrice, pkItems)
	if err != nil {
		return "", err
	}

	// Chat integration
	chatID := userID
	fragmentLink := ""
	if u.chatAdapter != nil {
		fragmentID, err := u.chatAdapter.AddEndMsgID(ctx, chatID)
		if err == nil {
			fragmentLink = fmt.Sprintf("%s/admin/admin_chat?fragment=%s&chat_id=%s", u.config.ClienUrl, fragmentID, chatID)
		}
	}

	// Create order in DB
	order := &order_entity.Order{
		UserID:       &userID,
		TotalPrice:   cart.TotalPrice,
		Email:        data.Email,
		Phone:        data.PhoneNumber,
		FIO:          data.FIO,
		Address:      data.Address,
		Status:       order_entity.OrderStatusPending,
		FragmentLink: fragmentLink,
	}
	if userID == "" {
		order.UserID = nil
	}

	if err := u.repo.Create(ctx, order); err != nil {
		return "", err
	}

	// Add order items
	orderItems := make([]order_entity.OrderItem, len(cart.Items))
	for i, item := range cart.Items {
		itemTotalPrice := item.UnitPrice * float64(item.Quantity)
		orderItems[i] = order_entity.OrderItem{
			OrderID:         order.ID,
			ProductID:       item.ProductID,
			Name:            item.ProductID, // Should be product name
			Article:         item.Article,
			Price:           item.UnitPrice,
			Quantity:        item.Quantity,
			TotalPrice:      itemTotalPrice,
			SelectedOptions: item.SelectedOptions,
		}
	}

	if err := u.repo.AddItems(ctx, orderItems); err != nil {
		return "", err
	}

	// Clear cart
	_ = u.cartAdapter.ClearCart(ctx, cart.ID)

	// Send notification
	u.notificationAdapter.SendOrderNotification(ctx, order)

	return link, nil
}

func (u *OrderUsecase) CreateDirectOrder(ctx context.Context, data *order_dto.CreateDirectOrderRequest) (string, error) {
	if err := u.validator.Struct(data); err != nil {
		return "", err
	}

	product, err := u.productAdapter.GetByID(ctx, data.ProductID)
	if err != nil || product == nil {
		return "", fmt.Errorf("product not found")
	}

	price := *product.ManualPrice

	categoryArticle := strings.Split(product.Article, ".")[0]
	truCode, _ := u.truAdapter.GetByCode(ctx, categoryArticle)
	if truCode != "" {
		truCode += "00000000643"
	}

	pkItems := []order_payment.PayKeeperItem{
		{
			ItemType:    "goods",
			PaymentType: "full",
			Name:        product.Name,
			Price:       price,
			Quantity:    1,
			TruCode:     truCode,
			Tax:         "none",
			Sum:         price,
		},
	}

	link, err := u.paymentProvider.CreateInvoice(ctx, data.Email, data.PhoneNumber, data.FIO, price, pkItems)
	if err != nil {
		return "", err
	}

	order := &order_entity.Order{
		TotalPrice: price,
		Email:      data.Email,
		Phone:      data.PhoneNumber,
		FIO:        data.FIO,
		Address:    data.Address,
		Status:     order_entity.OrderStatusPending,
	}

	if err := u.repo.Create(ctx, order); err != nil {
		return "", err
	}

	orderItem := order_entity.OrderItem{
		OrderID:    order.ID,
		ProductID:  product.ID,
		Name:       product.Name,
		Article:    product.Article,
		Price:      price,
		Quantity:   1,
		TotalPrice: price,
	}

	if err := u.repo.AddItems(ctx, []order_entity.OrderItem{orderItem}); err != nil {
		return "", err
	}

	u.notificationAdapter.SendOrderNotification(ctx, order)

	return link, nil
}

func (u *OrderUsecase) GetMyOrders(ctx context.Context, userID string) ([]order_entity.Order, error) {
	return u.repo.GetByUserID(ctx, userID)
}

func (u *OrderUsecase) GetAllOrders(ctx context.Context, offset, limit int) ([]order_entity.Order, error) {
	return u.repo.GetAll(ctx, offset, limit)
}

func (u *OrderUsecase) ChangeStatus(ctx context.Context, data *order_dto.ChangeOrderStatusRequest) error {
	if err := u.validator.Struct(data); err != nil {
		return err
	}

	order := &order_entity.Order{
		ID:     data.OrderID,
		Status: data.Status,
	}

	return u.repo.Update(ctx, order)
}
