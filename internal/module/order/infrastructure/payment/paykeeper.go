package order_payment

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Fi44er/sdmed/internal/config"
	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/Fi44er/sdmed/pkg/utils"
)

type PayKeeperItem struct {
	ItemType    string  `json:"item_type"`
	PaymentType string  `json:"payment_type"`
	SKU         string  `json:"sku"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	ItemCode    string  `json:"item_code"`
	TruCode     string  `json:"tru_code"`
	Tax         string  `json:"tax"`
	Sum         float64 `json:"sum"`
}

type PayKeeperProvider struct {
	config *config.Config
	logger *logger.Logger
}

func NewPayKeeperProvider(config *config.Config, logger *logger.Logger) *PayKeeperProvider {
	return &PayKeeperProvider{
		config: config,
		logger: logger,
	}
}

func (p *PayKeeperProvider) CreateInvoice(ctx context.Context, email, phone, fio string, amount float64, items []PayKeeperItem) (string, error) {
	user := p.config.PayKeeperUser
	password := p.config.PayKeeperPass
	server := p.config.PayKeeperServer

	auth := base64.StdEncoding.EncodeToString([]byte(user + ":" + password))

	// 1. Get Token
	tokenBody, err := utils.MakeRequest(utils.RequestOptions{
		Method: "GET",
		URL:    server + "/info/settings/token/",
		Headers: map[string]string{
			"Authorization": "Basic " + auth,
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to get PayKeeper token: %v", err)
	}

	var tokenResponse map[string]string
	if err := json.Unmarshal(tokenBody, &tokenResponse); err != nil {
		return "", err
	}

	token, ok := tokenResponse["token"]
	if !ok {
		return "", fmt.Errorf("token not found in PayKeeper response")
	}

	// 2. Create Invoice
	jsonData, _ := json.Marshal(items)
	serviceName := fmt.Sprintf(";PKC|%s|", jsonData)
	expiry := time.Now().AddDate(0, 0, 3).Format("2006-01-02")

	invoiceBody, err := utils.MakeRequest(utils.RequestOptions{
		Method: "POST",
		URL:    server + "/change/invoice/preview/",
		Headers: map[string]string{
			"Content-Type":  "application/x-www-form-urlencoded",
			"Authorization": "Basic " + auth,
		},
		FormData: map[string]string{
			"cart_json":    string(jsonData),
			"client_email": email,
			"client_phone": phone,
			"clientid":     fio,
			"expiry":       expiry,
			"pay_amount":   fmt.Sprintf("%.2f", amount),
			"service_name": serviceName,
			"token":        token,
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to create PayKeeper invoice: %v", err)
	}

	var invoiceResponse map[string]string
	if err := json.Unmarshal(invoiceBody, &invoiceResponse); err != nil {
		return "", err
	}

	invoiceID, ok := invoiceResponse["invoice_id"]
	if !ok {
		return "", fmt.Errorf("invoice_id not found in PayKeeper response: %v", invoiceResponse["msg"])
	}

	return fmt.Sprintf("%s/bill/%s/", server, invoiceID), nil
}
