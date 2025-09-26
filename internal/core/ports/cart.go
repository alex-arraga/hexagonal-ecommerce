package ports

import (
	"context"
	"go-ecommerce/internal/core/domain"

	"github.com/google/uuid"
)

type CartService interface {
	GetCart(ctx context.Context, userId uuid.UUID) (*domain.Cart, error)
	AddItemToCart(ctx context.Context, userId, productId uuid.UUID, quantity uint8) error
	RemoveItem(ctx context.Context, userId, productId uuid.UUID) error
	Clear(ctx context.Context, userId uuid.UUID) error
}
