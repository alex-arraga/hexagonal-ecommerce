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
	ID                int
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

// Merchant Order object
type MerchantItem struct {
	ID          string      `json:"id"`
	CategoryID  string      `json:"category_id"`
	CurrencyID  string      `json:"currency_id"`
	Description string      `json:"description"`
	PictureURL  interface{} `json:"picture_url"`
	Title       string      `json:"title"`
	Quantity    int         `json:"quantity"`
	UnitPrice   int         `json:"unit_price"`
}

type MerchantPayment struct {
	ID                int                    `json:"id"`
	TransactionAmount int                    `json:"transaction_amount"`
	TotalPaidAmount   int                    `json:"total_paid_amount"`
	ShippingCost      int                    `json:"shipping_cost"`
	CurrencyID        string                 `json:"currency_id"`
	Status            domain.PayStatus       `json:"status"`
	StatusDetail      domain.PayStatusDetail `json:"status_detail"`
	OperationType     string                 `json:"operation_type"`
	DateApproved      string                 `json:"date_approved"`
	DateCreated       string                 `json:"date_created"`
	LastModified      string                 `json:"last_modified"`
	AmountRefunded    int                    `json:"amount_refunded"`
}

type MerchantCollector struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
}

type MpSimplifiedMerchantOrder struct {
	ID                int               `json:"id"`
	Status            string            `json:"status"`
	ExternalReference string            `json:"external_reference"`
	PreferenceID      string            `json:"preference_id"`
	Payments          []MerchantPayment `json:"payments"`
	Shipments         []interface{}     `json:"shipments"`
	Payouts           []interface{}     `json:"payouts"`
	Collector         MerchantCollector `json:"collector"`
	Marketplace       string            `json:"marketplace"`
	NotificationURL   string            `json:"notification_url"`
	ShippingCost      int               `json:"shipping_cost"`
	TotalAmount       int               `json:"total_amount"`
	SiteID            string            `json:"site_id"`
	PaidAmount        int               `json:"paid_amount"`
	RefundedAmount    int               `json:"refunded_amount"`
	Payer             MpPayer           `json:"payer"`
	Items             []MerchantItem    `json:"items"`
	Cancelled         bool              `json:"cancelled"`
	OrderStatus       string            `json:"order_status"`
}
