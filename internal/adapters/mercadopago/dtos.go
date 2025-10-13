package mercadopago

import "go-ecommerce/internal/core/domain"

// Preference request params
type MpPreferenceRequest struct {
	AutoReturn          string         `json:"auto_return"`
	StatementDescriptor string         `json:"statement_descriptor"`
	ExternalReference   string         `json:"external_reference"`
	NotificationURL     string         `json:"notification_url"`
	BackUrls            MpBackUrls     `json:"back_urls"`
	Items               []MpItem       `json:"items"`
	Payer               MpPayer        `json:"payer"`
	PaymentMethods      PaymentMethods `json:"payment_methods"`
}

type MpBackUrls struct {
	Success string `json:"success"`
	Failure string `json:"failure"`
	Pending string `json:"pending"`
}

type MpItem struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	CategoryID  string  `json:"category_id"`
	CurrencyID  string  `json:"currency_id"`
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
}

type MpPayer struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Email   string `json:"email"`
	Phone   Phone  `json:"phone"`
}

type Phone struct {
	AreaCode string `json:"area_code"`
	Number   string `json:"number"`
}

type PaymentMethods struct {
	Installments           int              `json:"installments"`
	ExcludedPaymentTypes   []ExcludedType   `json:"excluded_payment_types"`
	ExcludedPaymentMethods []ExcludedMethod `json:"excluded_payment_methods"`
}

type ExcludedType struct {
	ID string `json:"id"`
}

type ExcludedMethod struct {
	ID string `json:"id"`
}

// Payments object
type TransactionDetails struct {
	TotalPaidAmount   float64
	NetReceivedAmount float64
	InstallmentAmount float64
}

type Order struct {
	ID   string
	Type string
}

type PayMethod struct {
	ID   *string
	Type *string
}

type MpSimplifiedPayment struct {
	ID                string
	Status            domain.PayStatus
	StatusDetail      domain.PayStatusDetail
	DateApproved      *string
	TransactionAmount float64
	CurrencyID        string
	Installments      uint8
	ExternalReference string
	PayMethod         PayMethod

	Payer              MpPayer
	TransactionDetails TransactionDetails
	Order              *Order
}

/*
   const updateOrder: UpdateOrderFromWebhookType = {
     pay_status: payment.status as PayStatusType,
     pay_status_detail: payment.status_detail,
     payment_id: payment.id.toString(),
     merchant_order_id: merchantOrderId,
     fee: fee,
     installments: payment.installments,
     net_received_amount: payment.transaction_details.net_received_amount,
     pay_method: payment.payment_method.id,
     pay_resource: payment.payment_method.type,
     paid_at: paidAt,
     updated_at: new Date(Date.now()).toISOString(),
     expires_at: expiresAt,
   }
*/
