package domain

import (
	"time"

	"github.com/google/uuid"
)

type Providers string
type Currencies string
type DisscountTypes string
type PayStatus string
type PayStatusDetail string

const (
	MercadoPago Providers = "mercado-pago"
)

const (
	ARS Currencies = "ARS"
	USD Currencies = "USD"
)

const (
	Percentage DisscountTypes = "percentage" // Percentage = subtotal * (15 / 100)
	Bundle     DisscountTypes = "bundle"     // Bundle = 2x1, 3x2, Pack of 3 services $X
	Fixed      DisscountTypes = "fixed"      // Fixed = $500
)

const (
	Approved    PayStatus = "approved"     // The payment has been approved
	Pending     PayStatus = "pending"      // The payment has been initiated but has not yet been processed by the payment method
	InProcess   PayStatus = "in_process"   // The payment is being validated (e.g., manual review)
	Authorized  PayStatus = "authorized"   // The payment was authorized by the issuer (bank or card) but has not yet been captured; the funds are reserved, not debited.
	Cancelled   PayStatus = "cancelled"    // The payment was canceled before completion
	Refunded    PayStatus = "refunded"     // The payment was returned to the user
	ChargedBack PayStatus = "charged_back" // The buyer disputed the payment and the issuer reversed the transaction (chargeback). The money is returned to the buyer.
	Rejected    PayStatus = "rejected"     // The payment was declined. This may be due to insufficient funds, incorrect details, etc
	Expired     PayStatus = "expired"      // The pay order has expired
	SoftDelete  PayStatus = "soft_delete"  // Order that remained unpaid for a long period of time and will be deleted
	NonExistent PayStatus = "non-existent" // Order does not exist, the user created their reservation but did not attempt to pay for it
)

// PayStatusDetail represents detailed payment status information.
const (
	Accredited        PayStatusDetail = "accredited"            // Payment approved and successfully credited; funds are now available.
	PendingCapture    PayStatusDetail = "pending_capture"       // Payment authorized but pending manual capture; funds are reserved, not yet charged.
	PartiallyRefunded PayStatusDetail = "partially_refunded"    // Payment was partially refunded; only part of the amount was returned to the buyer.
	InProcessDetail   PayStatusDetail = "in_process"            // Payment is under review or being processed; not yet completed.
	ExpiredDetail     PayStatusDetail = "expired"               // Payment request expired before completion (e.g., buyer didn't finish in time).
	BankError         PayStatusDetail = "bank_error"            // Payment failed due to a bank or issuer error.
	Blacklist         PayStatusDetail = "cc_rejected_blacklist" // Payment rejected because the card or user is blacklisted.
	NonExistentDetail PayStatusDetail = "non-existent"          // No payment information found (invalid or deleted payment ID).
)

type Order struct {
	ID                uuid.UUID
	Providers         Providers
	UserID            uuid.UUID
	PaymentID         *string
	SecureToken       uuid.UUID // Allows the user to view the order status. Will be automatically generate by gorm
	ExternalReference *string
	Currency          Currencies
	Fee               *float64
	Installments      *uint8
	NetReceivedAmount *float64
	PayMethod         *string
	PayResource       *string
	SubTotal          float64
	Disscount         *float64
	DisscountType     *DisscountTypes
	Total             float64
	Paid              bool
	PayStatus         PayStatus
	PayStatusDetail   *PayStatusDetail
	CreatedAt         time.Time
	UpdatedAt         time.Time
	PaidAt            *time.Time
	ExpiresAt         *time.Time

	// Relations
	User  *User
	Items []OrderProduct
}

type NewOrderInputs struct {
	UserID         uuid.UUID
	Currency       Currencies
	SubTotal       float64
	Disscount      *float64
	DisscountTypes *DisscountTypes
}

func NewOrder(inputs NewOrderInputs) (*Order, error) {
	now := time.Now()
	expireInTreeDays := now.AddDate(0, 0, 3)

	return &Order{
		ID:                uuid.Nil, // repository will asign the id
		UserID:            inputs.UserID,
		Currency:          inputs.Currency,
		PayStatus:         Pending,
		Providers:         MercadoPago,
		ExternalReference: nil,
		SecureToken:       uuid.Nil,
		PaymentID:         nil,
		PayStatusDetail:   nil,
		Paid:              false,
		SubTotal:          inputs.SubTotal,
		Disscount:         inputs.Disscount,
		DisscountType:     inputs.DisscountTypes,
		Total:             inputs.SubTotal - *inputs.Disscount,
		CreatedAt:         now,
		UpdatedAt:         now,
		ExpiresAt:         &expireInTreeDays, // the order is created with 3 days to pay it
	}, nil
}

type UpdateOrderInputs struct {
	PaymentID         string
	PayStatus         PayStatus
	PayStatusDetail   PayStatusDetail
	PayMethod         *string
	PayResource       *string
	Installments      *uint8
	Paid              bool
	Fee               *float64
	NetReceivedAmount *float64
	ExternalReference string
	PaidAt            *time.Time
}

func (o *Order) UpdateOrder(inputs UpdateOrderInputs) error {
	o.ExternalReference = &inputs.ExternalReference
	o.PaymentID = &inputs.PaymentID

	o.PayStatus = inputs.PayStatus
	o.PayStatusDetail = &inputs.PayStatusDetail
	o.UpdatedAt = time.Now()

	if o.PayStatus == Approved {
		o.ExpiresAt = nil
	}

	return nil
}
