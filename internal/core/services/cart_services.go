package services

import (
	"context"
	"encoding/json"
	"fmt"
	cachekeys "go-ecommerce/internal/adapters/storage/cache/cache_keys"
	cachettl "go-ecommerce/internal/adapters/storage/cache/cache_ttl"
	"go-ecommerce/internal/core/domain"
	"go-ecommerce/internal/core/ports"
	"log/slog"

	"github.com/google/uuid"
)

type CartService struct {
	ps    ports.ProductService
	cache ports.CacheRepository
}

func NewCartService(cache ports.CacheRepository, ps ports.ProductService) ports.CartService {
	return &CartService{cache: cache, ps: ps}
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

	var cart domain.Cart
	err = json.Unmarshal(data, &cart)
	if err != nil {
		slog.Warn("error deserializing items of cart", "error", err)
		domain.NewCart(userId)
	}

	// return cart with values
	return &cart
}

// helper func
func (c *CartService) saveCart(ctx context.Context, cart *domain.Cart) error {
	cacheKey := cachekeys.Cart(cart.UserID.String())

	// serialized data before caching
	data, err := json.Marshal(cart)
	if err != nil {
		return err
	}

	// If cart hasn't items, delete it
	if len(cart.Items) <= 0 {
		err := c.cache.Delete(ctx, cacheKey)
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

// CalcItemsAmount implements ports.CartService.
func (c *CartService) CalcItemsAmount(ctx context.Context, userId uuid.UUID) (*ports.Amount, error) {
	cart, err := c.GetCart(ctx, userId)
	if err != nil {
		return nil, err
	}

	var subTotal float64 = 0
	var discount float64 = 0
	var total float64 = 0

	if len(cart.Items) <= 0 {
		return nil, fmt.Errorf("items not found in cart")
	}

	for _, item := range cart.Items {
		prod, err := c.ps.GetProductById(ctx, item.ProductID)
		if err != nil {
			return nil, err
		}
		subTotal += prod.Price * float64(item.Quantity)
		discount += prod.Disscount * float64(item.Quantity)
		total = subTotal - discount
	}

	amount := &ports.Amount{
		SubTotal: subTotal,
		Discount: discount,
		Total:    total,
	}

	return amount, nil
}

// RemoveItem implements ports.CartService.
func (c *CartService) RemoveItem(ctx context.Context, userId, productId uuid.UUID) error {
	cart := c.loadCart(ctx, userId)
	err := cart.RemoveItem(productId)
	if err != nil {
		return err
	}
	return c.saveCart(ctx, cart)
}

// Clear implements ports.CartService.
func (c *CartService) Clear(ctx context.Context, userId uuid.UUID) error {
	cart := c.loadCart(ctx, userId)
	err := cart.Clear()
	if err != nil {
		return err
	}
	cacheKey := cachekeys.Cart(userId.String())
	return c.cache.Delete(ctx, cacheKey)
}
