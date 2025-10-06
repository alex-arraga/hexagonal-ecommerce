package services

import (
	"context"
	"encoding/json"
	cachekeys "go-ecommerce/internal/adapters/storage/cache/cache_keys"
	cachettl "go-ecommerce/internal/adapters/storage/cache/cache_ttl"
	"go-ecommerce/internal/core/domain"
	"go-ecommerce/internal/core/ports"
	"log/slog"

	"github.com/google/uuid"
)

type CartService struct {
	cache ports.CacheRepository
}

func NewCartService(cache ports.CacheRepository) ports.CartService {
	return &CartService{cache: cache}
}

// helper func
func (c *CartService) loadCart(ctx context.Context, userId uuid.UUID) *domain.Cart {
	cacheKey := cachekeys.Cart(userId.String())

	data, err := c.cache.Get(ctx, cacheKey)
	if err != nil {
		slog.Warn("error obtaining cart in cache", "user_id", userId, "error", err)
	}

	// if cart doesn't exist, create one and return
	if len(data) == 0 {
		return domain.NewCart(userId)
	}

	var cart *domain.Cart
	err = json.Unmarshal(data, &cart)
	if err != nil {
		slog.Warn("error deserializing items of cart", "error", err)
		domain.NewCart(userId)
	}

	// return cart with values
	return cart
}

// helper func
func (c *CartService) saveCart(ctx context.Context, cart *domain.Cart) error {
	cacheKey := cachekeys.Cart(cart.UserID.String())

	// serialized data before caching
	data, err := json.Marshal(&cart)
	if err != nil {
		return err
	}

	// set cache with new values of cart
	return c.cache.Set(ctx, cacheKey, data, cachettl.Cart)
}

// AddItemToCart implements ports.CartService.
func (c *CartService) AddItemToCart(ctx context.Context, userId, productId uuid.UUID, quantity int16) error {
	cart := c.loadCart(ctx, userId)
	err := cart.AddItem(productId, quantity)
	if err != nil {
		return err
	}
	return c.saveCart(ctx, cart)
}

// GetCart implements ports.CartService.
func (c *CartService) GetCart(ctx context.Context, userId uuid.UUID) (*domain.Cart, error) {
	cart := c.loadCart(ctx, userId)
	return cart, nil
}

// RemoveItem implements ports.CartService.
func (c *CartService) RemoveItem(ctx context.Context, userId, productId uuid.UUID) error {
	cart := c.loadCart(ctx, userId)
	cart.RemoveItem(productId)
	return c.saveCart(ctx, cart)
}

// Clear implements ports.CartService.
func (c *CartService) Clear(ctx context.Context, userId uuid.UUID) error {
	cacheKey := cachekeys.Cart(userId.String())
	return c.cache.Delete(ctx, cacheKey)
}
