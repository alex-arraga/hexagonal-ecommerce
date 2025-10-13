package mercadopago

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-ecommerce/internal/core/domain"
	"go-ecommerce/internal/core/ports"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
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
func (ps *PaymentService) CreatePayment(ctx context.Context, orderId uuid.UUID) (*string, error) {

	order, err := ps.repo.GetOrderById(ctx, orderId)
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
func (ps *PaymentService) VerifyPayment(ctx context.Context, paymentId, topic *string) error {
	if paymentId == nil && topic == nil {
		return errors.New("parameters id or topic not found")
	}

	if paymentId != nil && *topic == "payment" {
		return ps.handlePayment(ctx, *paymentId)
	}

	if paymentId != nil && *topic == "merchant_order" {
		return ps.handleMerchantOrder(ctx, *paymentId)
	}

	return nil
}

func (ps *PaymentService) handlePayment(ctx context.Context, id string) error {
	url := fmt.Sprintf("https://api.mercadopago.com/v1/payments/%s", id)
	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return err
	}

	// set headers and fetching request
	req.Header.Set("Authorization", fmt.Sprintf("Bearer: %s", ps.secretToken))

	resp, err := ps.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// TODO: eliminar esto, solo para ver el objeto
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Print(bodyBytes)

	// transform body in payment object
	payment := &MpSimplifiedPayment{}
	json.NewDecoder(resp.Body).Decode(payment)

	parsedID, err := uuid.Parse(payment.ExternalReference)
	if err != nil {
		return err
	}

	// obtain order to update
	order, err := ps.repo.GetOrderById(ctx, parsedID)
	if err != nil {
		return err
	}

	// validations and handling errors
	// external_reference must be equal than order id
	if payment.ExternalReference != order.ID.String() {
		return fmt.Errorf("payment: %s external_reference does not match with order: %s", payment.ExternalReference, order.ID)
	}

	// avoids updating a order with an approved payment but that was never was credited due to account errors or holds
	if payment.TransactionDetails.NetReceivedAmount <= 0 {
		return fmt.Errorf("net received amount is 0 or less for payment: %v", payment.ID)
	}

	// if the order has already been updated with this payment_id, return and do nothing
	if order.PaymentID == &payment.ID {
		return fmt.Errorf("payment: %s already processed for order: %s", payment.ID, order.ID)
	}

	// update order with payment data
	successfullyPayment := payment.Status == domain.Approved && payment.StatusDetail == domain.Accredited
	fee := payment.TransactionAmount - payment.TransactionDetails.NetReceivedAmount

	var paidAt time.Time
	if successfullyPayment {
		paidAt = time.Now()
	}

	var isPaid bool
	if successfullyPayment {
		isPaid = true
	}

	dataToUpdate := domain.UpdateOrderInputs{
		PayStatus:         domain.PayStatus(payment.Status),
		PayStatusDetail:   domain.PayStatusDetail(payment.StatusDetail),
		PaymentID:         fmt.Sprint(payment.ID),
		PayMethod:         payment.PayMethod.ID,
		PayResource:       payment.PayMethod.Type,
		Installments:      &payment.Installments,
		ExternalReference: payment.ExternalReference,
		NetReceivedAmount: &payment.TransactionDetails.NetReceivedAmount,
		Fee:               &fee,
		PaidAt:            &paidAt,
		Paid:              isPaid,
	}

	err = order.UpdateOrder(dataToUpdate)
	if err != nil {
		return err
	}

	_, err = ps.repo.SaveOrder(ctx, order)
	return err
}

func (ps *PaymentService) handleMerchantOrder(ctx context.Context, id string) error {
	return nil
}
