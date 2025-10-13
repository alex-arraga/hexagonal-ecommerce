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
	Approved    PayStatus = "approved"     //
	Pending     PayStatus = "pending"      //
	InProcess   PayStatus = "in_process"   //
	Authorized  PayStatus = "authorized"   //
	Cancelled   PayStatus = "cancelled"    //
	Refunded    PayStatus = "refunded"     //
	ChargedBack PayStatus = "charged_back" //
	Rejected    PayStatus = "rejected"     //
	Expired     PayStatus = "expired"      //
	SoftDelete  PayStatus = "soft_delete"  //
	NonExistent PayStatus = "non-existent" //
)

const (
	Accredited        PayStatusDetail = "accredited"            //
	PendingCapture    PayStatusDetail = "pending_capture"       //
	PartiallyRefunded PayStatusDetail = "partially_refunded"    //
	InProcessDetail   PayStatusDetail = "in_process"            //
	ExpiredDetail     PayStatusDetail = "expired"               //
	BankError         PayStatusDetail = "bank_error"            //
	Blacklist         PayStatusDetail = "cc_rejected_blacklist" //
	NonExistentDetail PayStatusDetail = "non-existent"          //
)

type Order struct {
	ID                uuid.UUID
	Providers         Providers
	UserID            uuid.UUID
	PaymentID         *string
	SecureToken       *uuid.UUID // token that allows the user to view the status of the order
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
		SecureToken:       nil,
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
