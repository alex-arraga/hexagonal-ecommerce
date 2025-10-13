package mercadopago

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go-ecommerce/internal/core/domain"
	"go-ecommerce/internal/core/ports"
	"net/http"
)

type PaymentService struct {
	httpClient  *http.Client
	repo        ports.OrderRepository
	domain      string
	secretToken string
}

func NewPaymentService(repo ports.OrderRepository, client *http.Client, domain, secretToken string) ports.PaymentProvider {
	return &PaymentService{
		repo:        repo,
		httpClient:  client,
		domain:      domain,
		secretToken: secretToken,
	}
}

// Helper func to create a preference
func NewPreference(order *domain.Order, domain string) MpPreferenceRequest {
	items := make([]MpItem, len(order.Items))

	for _, orderItem := range order.Items {
		items = append(items, MpItem{
			ID:          orderItem.ProductID.String(),
			Title:       orderItem.Product.Name,
			Description: orderItem.Product.SKU,
			CategoryID:  orderItem.Product.Category.Name,
			CurrencyID:  string(orderItem.Order.Currency),
			Quantity:    int(orderItem.Quantity),
			UnitPrice:   orderItem.Product.Price,
		})
	}

	preference := MpPreferenceRequest{
		AutoReturn:          "approved",
		StatementDescriptor: "Golang Ecommerce",
		ExternalReference:   order.ID.String(),
		NotificationURL:     fmt.Sprintf("%s/order/api/payment/mp/webhook", domain),
		BackUrls: MpBackUrls{
			Success: fmt.Sprintf("%s/order/%s", domain, order.SecureToken),
			Failure: fmt.Sprintf("%s/order/%s", domain, order.SecureToken),
			Pending: fmt.Sprintf("%s/order/%s", domain, order.SecureToken),
		},
		Items: items,
		Payer: MpPayer{
			Name:  order.User.Name,
			Email: order.User.Email,
		},
		PaymentMethods: PaymentMethods{
			Installments: 6,
			ExcludedPaymentTypes: []ExcludedType{
				ExcludedType{ID: "atm"},
				ExcludedType{ID: "ticket"},
				ExcludedType{ID: "bank_transfer"},
			},
			ExcludedPaymentMethods: []ExcludedMethod{
				ExcludedMethod{ID: "servipag"},
				ExcludedMethod{ID: "sencillito"},
			},
		},
	}

	return preference
}

// CreatePayment implements ports.PaymentProvider.
func (ps *PaymentService) CreatePayment(ctx context.Context, params ports.CreatePaymentParams) (*string, error) {

	order, err := ps.repo.GetOrderById(ctx, params.OrderID)
	if err != nil {
		return nil, err
	}

	preference := NewPreference(order, ps.domain)

	jsonBody, _ := json.Marshal(preference)
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.com/checkout/preferences", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ps.secretToken))

	res, err := ps.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result struct {
		InitPoint string `json:"init_point"`
	}

	json.NewDecoder(res.Body).Decode(&result)
	return &result.InitPoint, nil
}

// VerifyPayment implements ports.PaymentProvider.
func (p *PaymentService) VerifyPayment(ctx context.Context, paymentId string) (string, error) {
	panic("unimplemented")
}
