package mercadopago

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-ecommerce/internal/adapters/mercadopago/mp_dtos"
	"go-ecommerce/internal/core/domain"
	"go-ecommerce/internal/core/ports"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type PaymentProvider struct {
	httpClient  *http.Client
	domain      string
	secretToken string
}

func NewPaymentProvider(client *http.Client, domain, secretToken string) ports.PaymentProvider {
	return &PaymentProvider{
		httpClient:  client,
		domain:      domain,
		secretToken: secretToken,
	}
}

// Helpers funcs
func (ps *PaymentProvider) GeneratePreference(ctx context.Context, order *domain.Order, items []mp_dtos.MpItem, user *domain.User) *mp_dtos.MpPreferenceRequest {
	preference := mp_dtos.MpPreferenceRequest{
		AutoReturn:          "approved",
		StatementDescriptor: "Golang Ecommerce",
		ExternalReference:   fmt.Sprint(order.ID),
		NotificationURL:     fmt.Sprintf("%s/payment/mp/webhook", ps.domain),
		BackUrls: mp_dtos.MpBackUrls{
			Success: fmt.Sprintf("%s/order/%s", ps.domain, order.SecureToken),
			Failure: fmt.Sprintf("%s/order/%s", ps.domain, order.SecureToken),
			Pending: fmt.Sprintf("%s/order/%s", ps.domain, order.SecureToken),
		},
		Items: items,
		Payer: mp_dtos.MpPayer{
			Name:  user.Name,
			Email: user.Email,
			Phone: mp_dtos.Phone{
				AreaCode: "+54",
				Number:   "123456",
			},
		},
		PaymentMethods: mp_dtos.PaymentMethods{
			Installments: 6,
			ExcludedPaymentTypes: []mp_dtos.ExcludedType{
				{ID: "atm"},
				{ID: "ticket"},
				{ID: "bank_transfer"},
			},
			ExcludedPaymentMethods: []mp_dtos.ExcludedMethod{
				{ID: "servipag"},
				{ID: "sencillito"},
			},
		},
	}

	return &preference
}

func (ps *PaymentProvider) handlePayment(ctx context.Context, order *domain.Order, paymentId string) (*domain.Order, error) {
	// prepare request to call mp api
	url := fmt.Sprintf("https://api.mercadopago.com/v1/payments/%s", paymentId)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	// set headers and fetching request
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ps.secretToken))

	resp, err := ps.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// transform body in payment object
	payment := &mp_dtos.MpSimplifiedPayment{}
	err = json.NewDecoder(resp.Body).Decode(payment)
	if err != nil {
		return nil, fmt.Errorf("failed decoding payment: %w", err)

	}

	if payment.ExternalReference == "" {
		return nil, fmt.Errorf("external_reference missing in payment")
	}

	// if order has been updated, return
	if order.PaymentID != nil {
		return nil, nil
	}

	// validations and handling errors
	// external_reference must be equal than order id
	if payment.ExternalReference != order.ID.String() {
		return nil, fmt.Errorf("payment: %s external_reference does not match with order: %s", payment.ExternalReference, order.ID)
	}

	// avoids updating a order with an approved payment but that was never was credited due to account errors or holds
	if payment.TransactionDetails.NetReceivedAmount <= 0 {
		return nil, fmt.Errorf("net received amount is 0 or less for payment: %v", payment.ID)
	}

	// if the order has already been updated with this payment_id, return and do nothing
	strPaymentId := strconv.Itoa(payment.ID)
	if order.PaymentID == &strPaymentId {
		return nil, fmt.Errorf("payment: %v already processed for order: %s", payment.ID, order.ID)
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
		return nil, err
	}

	return order, nil
}

func (ps *PaymentProvider) handleMerchantOrder(ctx context.Context, order *domain.Order, merchantOrderId string) (*domain.Order, error) {
	// prepare request to call mp api
	url := fmt.Sprintf("https://api.mercadopago.com/merchant_orders/%s", merchantOrderId)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	// set headers and fetching request
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ps.secretToken))

	resp, err := ps.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	merchantOrder := &mp_dtos.MpSimplifiedMerchantOrder{}
	err = json.NewDecoder(resp.Body).Decode(&merchantOrder)
	if err != nil {
		return nil, fmt.Errorf("failed decoding merchant order: %w", err)
	}

	// seach and find approved payments inside merchant order, if exist call handlePayment, else return an error
	found := false

	for _, payment := range merchantOrder.Payments {
		if payment.Status == domain.Approved && payment.StatusDetail == domain.Accredited {
			strPaymentId := fmt.Sprint(payment.ID)
			updatedOrder, err := ps.handlePayment(ctx, order, strPaymentId)
			if err != nil {
				slog.Error("error processing mercado pago payment", "error", err)
			}
			found = true
			return updatedOrder, nil
		}
	}

	if !found {
		return nil, fmt.Errorf("merchant order %s received, but no approved/accredited payment found", merchantOrderId)
	}

	return nil, nil
}

// CreatePayment implements ports.PaymentProvider.
func (ps *PaymentProvider) ProcessPayment(ctx context.Context, preference *mp_dtos.MpPreferenceRequest) (*string, error) {
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

	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("failed decoding result: %w", err)
	}

	return &result.InitPoint, nil
}

// VerifyPayment implements ports.PaymentProvider.
func (ps *PaymentProvider) VerifyPayment(ctx context.Context, order *domain.Order, id, topic *string) (*domain.Order, error) {
	if id == nil && topic == nil {
		return nil, errors.New("parameters id or topic not found")
	}

	if id != nil && *topic == "payment" {
		updatedOrder, err := ps.handlePayment(ctx, order, *id)
		if err != nil {
			return nil, err
		}
		return updatedOrder, nil
	}

	if id != nil && *topic == "merchant_order" {
		updatedOrder, err := ps.handleMerchantOrder(ctx, order, *id)
		if err != nil {
			return nil, err
		}
		return updatedOrder, nil
	}

	return nil, nil
}
