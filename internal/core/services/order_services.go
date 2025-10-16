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

type OrderService struct {
	orderRepo ports.OrderRepository
	ops       ports.OrderProductService
	cart      ports.CartService
	cache     ports.CacheRepository
}

func NewOrderService(orderRepo ports.OrderRepository, ops ports.OrderProductService, cart ports.CartService, cache ports.CacheRepository) ports.OrderService {
	return &OrderService{
		orderRepo: orderRepo,
		ops:       ops,
		cart:      cart,
		cache:     cache,
	}
}

// SaveOrder implements ports.OrderService.
func (os *OrderService) SaveOrder(ctx context.Context, inputs ports.SaveOrderInputs) (*domain.Order, error) {
	var order *domain.Order

	// get all items from order
	cart, err := os.cart.GetCart(ctx, inputs.UserID)
	if err != nil {
		return nil, err
	}

	// get amount of all items
	amount, err := os.cart.CalcItemsAmount(ctx, inputs.UserID)
	if err != nil {
		return nil, err
	}

	if inputs.ID == uuid.Nil {
		// create a new order if inputs.ID doesn't exist
		newOrderInputs := domain.NewOrderInputs{
			UserID:   inputs.UserID,
			Currency: inputs.Currency,
			SubTotal: amount.SubTotal,
			Discount: amount.Discount,
			Total:    amount.Total,
		}
		newOrder, err := domain.NewOrder(newOrderInputs)
		if err != nil {
			return nil, err
		}
		order = newOrder

	} else {
		existingOrder, err := os.orderRepo.GetOrderById(ctx, inputs.ID)
		if err != nil {
			return nil, err
		}

		existingOrder.UpdateOrder(domain.UpdateOrderInputs{
			ExternalReference: *inputs.ExternalReference,
			PaymentID:         *inputs.PaymentID,
			PayStatus:         domain.PayStatus(*inputs.PayStatus),
			PayStatusDetail:   *inputs.PayStatusDetail,
		})
		order = existingOrder
	}

	result, err := os.orderRepo.SaveOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	// if a order was created, creates order-product for each item of cart
	if inputs.ID == uuid.Nil {
		for _, item := range cart.Items {
			_, err := os.ops.AddProductToOrder(ctx, result.ID, item.ProductID, item.Quantity)
			if err != nil {
				return nil, err
			}
		}

		// if all is ok, clean cart
		err := os.cart.Clear(ctx, inputs.UserID)
		if err != nil {
			slog.Error("error cleaning cart", "UserID", inputs.UserID, "error", err)
		}
	}

	// create new order cache key and serialize order created or udpated
	cacheKey := cachekeys.Order(result.ID.String())
	orderSerialized, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	// set order in cache
	err = os.cache.Set(ctx, cacheKey, orderSerialized, cachettl.Order)
	if err != nil {
		slog.Warn("error caching new order created", "order_id", result.ID, "error", err)
	}

	// invalidate order list
	err = os.cache.Delete(ctx, cachekeys.AllOrders())
	if err != nil {
		slog.Warn("error invalidating list of all orders", "error", err)
	}

	return result, nil
}

// GetOrderById implements ports.OrderService.
func (os *OrderService) GetOrderById(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	cacheKey := cachekeys.Order(id.String())

	// check if the order exist in cache, if exist return it
	val, err := os.cache.Get(ctx, cacheKey)
	if err == nil {
		var order domain.Order
		if decodeErr := json.Unmarshal(val, &order); decodeErr != nil {
			return &order, nil
		}
	}

	// else find order in repository
	p, err := os.orderRepo.GetOrderById(ctx, id)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// ListOrders implements ports.OrderService.
func (os *OrderService) ListOrders(ctx context.Context) ([]*domain.Order, error) {
	// check if the products exists in cache
	values, err := os.cache.Get(ctx, cachekeys.AllOrders())
	if err == nil {
		var orders []*domain.Order
		if decodeErr := json.Unmarshal(values, &orders); decodeErr != nil {
			return orders, nil
		}
	}

	// find orders in repository
	orders, err := os.orderRepo.ListOrders(ctx)
	if err != nil {
		return nil, err
	}

	// serialize products
	ordersSerialized, err := json.Marshal(orders)
	if err != nil {
		return nil, err
	}

	// set products of repository in cache
	os.cache.Set(ctx, cachekeys.AllOrders(), ordersSerialized, cachettl.Order)

	return orders, nil
}
