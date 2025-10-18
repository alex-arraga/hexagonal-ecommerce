package ports

import (
	"context"
	"go-ecommerce/internal/adapters/mercadopago/mp_dtos"
	"go-ecommerce/internal/core/domain"

	"github.com/google/uuid"
)

type PaymentService interface {
	StartPayment(ctx context.Context, orderId uuid.UUID) (*string, error)
	VerifyPayment(ctx context.Context, paymentId, topic *string) error
}

type PaymentProvider interface {
	GeneratePreference(ctx context.Context, order *domain.Order, items []mp_dtos.MpItem, user *domain.User) *mp_dtos.MpPreferenceRequest
	GenerateNewPayment(ctx context.Context, preference *mp_dtos.MpPreferenceRequest) (*string, error)
	VerifyPayment(ctx context.Context, id, topic *string) (*mp_dtos.MpSimplifiedPayment, error)
}
