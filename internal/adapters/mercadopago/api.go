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
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type PaymentService struct {
	httpClient  *http.Client
	userRepo    ports.UserRepository
	orderRepo   ports.OrderRepository
	productRepo ports.ProductRepository
	domain      string
	secretToken string
}

func NewPaymentService(
	orderRepo ports.OrderRepository,
	productRepo ports.ProductRepository,
	userRepo ports.UserRepository,
	client *http.Client, domain,
	secretToken string) ports.PaymentProvider {
	return &PaymentService{
		userRepo:    userRepo,
		orderRepo:   orderRepo,
		productRepo: productRepo,
		httpClient:  client,
		domain:      domain,
		secretToken: secretToken,
	}
}

// Helpers funcs
func (ps *PaymentService) generatePreference(ctx context.Context, order *domain.Order, user *domain.User, domain string) (*MpPreferenceRequest, error) {
	items := make([]MpItem, 0)

	for _, orderItem := range order.Items {
		product, err := ps.productRepo.GetProductById(context.Background(), orderItem.ProductID)
		if err != nil {
			return nil, err
		}

		items = append(items, MpItem{
			ID:          orderItem.ProductID.String(),
			Title:       product.Name,
			Description: product.SKU,
			CategoryID:  fmt.Sprint(product.CategoryID),
			CurrencyID:  fmt.Sprint(order.Currency),
			Quantity:    int(orderItem.Quantity),
			UnitPrice:   product.Price,
		})
	}

	preference := MpPreferenceRequest{
		AutoReturn:          "approved",
		StatementDescriptor: "Golang Ecommerce",
		ExternalReference:   fmt.Sprint(order.ID),
		NotificationURL:     fmt.Sprintf("%s/payment/mp/webhook", domain),
		BackUrls: MpBackUrls{
			Success: fmt.Sprintf("%s/order/%s", domain, order.SecureToken),
			Failure: fmt.Sprintf("%s/order/%s", domain, order.SecureToken),
			Pending: fmt.Sprintf("%s/order/%s", domain, order.SecureToken),
		},
		Items: items,
		Payer: MpPayer{
			Name:  user.Name,
			Email: user.Email,
		},
		PaymentMethods: PaymentMethods{
			Installments: 6,
			ExcludedPaymentTypes: []ExcludedType{
				{ID: "atm"},
				{ID: "ticket"},
				{ID: "bank_transfer"},
			},
			ExcludedPaymentMethods: []ExcludedMethod{
				{ID: "servipag"},
				{ID: "sencillito"},
			},
		},
	}

	return &preference, nil
}

func (ps *PaymentService) handlePayment(ctx context.Context, paymentId string) error {
	// prepare request to call mp api
	url := fmt.Sprintf("https://api.mercadopago.com/v1/payments/%s", paymentId)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
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
	order, err := ps.orderRepo.GetOrderById(ctx, parsedID)
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
	strPaymentId := strconv.Itoa(payment.ID)
	if order.PaymentID == &strPaymentId {
		return fmt.Errorf("payment: %v already processed for order: %s", payment.ID, order.ID)
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
		Installments:      payment.Installments,
		ExternalReference: payment.ExternalReference,
		NetReceivedAmount: payment.TransactionDetails.NetReceivedAmount,
		Fee:               fee,
		PaidAt:            &paidAt,
		Paid:              isPaid,
	}

	err = order.UpdateOrder(dataToUpdate)
	if err != nil {
		return err
	}

	_, err = ps.orderRepo.SaveOrder(ctx, order)
	return err
}

func (ps *PaymentService) handleMerchantOrder(ctx context.Context, orderId string) error {
	// prepare request to call mp api
	url := fmt.Sprintf("https://api.mercadopago.com/merchant_orders/%s", orderId)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	// set headers and fetching request
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ps.secretToken))

	resp, err := ps.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// TODO: eliminar solo para ver el body
	jsonBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Print(jsonBytes)

	merchantOrder := &MpSimplifiedMerchantOrder{}
	json.NewDecoder(resp.Body).Decode(merchantOrder)

	// seach and find approved payments inside merchant order, if exist call handlePayment, else return an error
	found := false

	for _, payment := range merchantOrder.Payments {
		if payment.Status == domain.Approved && payment.StatusDetail == domain.Accredited {
			strPaymentId := fmt.Sprint(payment.ID)
			ps.handlePayment(ctx, strPaymentId)
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("merchant order %s received, but no approved/accredited payment found", orderId)
	}

	return nil
}

// CreatePayment implements ports.PaymentProvider.
func (ps *PaymentService) CreatePayment(ctx context.Context, orderId uuid.UUID) (*string, error) {

	order, err := ps.orderRepo.GetOrderById(ctx, orderId)
	if err != nil {
		return nil, err
	}

	user, err := ps.userRepo.GetUserByID(ctx, order.UserID)
	if err != nil {
		return nil, err
	}

	preference, err := ps.generatePreference(ctx, order, user, ps.domain)
	if err != nil {
		return nil, err
	}

	apiUrl := "https://api.mercadopago.com/checkout/preferences"

	jsonBody, _ := json.Marshal(preference)
	req, err := http.NewRequestWithContext(ctx, "POST", apiUrl, bytes.NewBuffer(jsonBody))
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

	if res.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(res.Body)
		slog.Error("Error in mercado pago response", "error", string(bodyBytes), "code", res.StatusCode)
		return nil, fmt.Errorf("bad request to MercadoPago")
	}

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
