package ports

import (
	"context"

	"github.com/google/uuid"
)

type PaymentProvider interface {
	CreatePayment(ctx context.Context, orderId uuid.UUID) (*string, error)
	VerifyPayment(ctx context.Context, paymentId, topic *string) error
}
