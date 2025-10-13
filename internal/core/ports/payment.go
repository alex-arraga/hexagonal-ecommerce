package ports

import (
	"context"

	"github.com/google/uuid"
)

type CreatePaymentParams struct {
	OrderID uuid.UUID
}

type PaymentProvider interface {
	CreatePayment(ctx context.Context, params CreatePaymentParams) (*string, error)
	VerifyPayment(ctx context.Context, paymentId string) (string, error)
}
