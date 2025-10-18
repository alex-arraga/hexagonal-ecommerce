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

// Returns the mercado pago payment object
func (ps *PaymentProvider) handlePayment(ctx context.Context, paymentId string) (*mp_dtos.MpSimplifiedPayment, error) {
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

	return payment, nil
}

func (ps *PaymentProvider) handleMerchantOrder(ctx context.Context, merchantOrderId string) (*mp_dtos.MpSimplifiedPayment, error) {
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
			mpPayment, err := ps.handlePayment(ctx, strPaymentId)
			if err != nil {
				slog.Error("error processing mercado pago payment", "error", err)
			}
			found = true
			return mpPayment, nil
		}
	}

	if !found {
		return nil, fmt.Errorf("merchant order %s received, but no approved/accredited payment found", merchantOrderId)
	}

	return nil, nil
}

// GenerateNewPayment implements ports.PaymentProvider.
func (ps *PaymentProvider) GenerateNewPayment(ctx context.Context, preference *mp_dtos.MpPreferenceRequest) (*string, error) {
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
func (ps *PaymentProvider) VerifyPayment(ctx context.Context, id, topic *string) (*mp_dtos.MpSimplifiedPayment, error) {
	if id == nil && topic == nil {
		return nil, errors.New("parameters id or topic not found")
	}

	if id != nil && *topic == "payment" {
		payment, err := ps.handlePayment(ctx, *id)
		if err != nil {
			return nil, err
		}
		return payment, nil
	}

	if id != nil && *topic == "merchant_order" {
		payment, err := ps.handleMerchantOrder(ctx, *id)
		if err != nil {
			return nil, err
		}
		return payment, nil
	}

	return nil, nil
}
