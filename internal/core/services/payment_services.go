package services

import (
	"context"
	"fmt"
	"go-ecommerce/internal/adapters/mercadopago/mp_dtos"
	"go-ecommerce/internal/core/ports"

	"github.com/google/uuid"
)

type PaymentService struct {
	userRepo    ports.UserRepository
	orderRepo   ports.OrderRepository
	productRepo ports.ProductRepository
	mp          ports.PaymentProvider
}

func NewPaymentService(userRepo ports.UserRepository, orderRepo ports.OrderRepository, productRepo ports.ProductRepository, mp ports.PaymentProvider) ports.PaymentService {
	return &PaymentService{
		userRepo:    userRepo,
		orderRepo:   orderRepo,
		productRepo: productRepo,
		mp:          mp,
	}
}

// CreatePayment implements ports.PaymentProvider.
func (p *PaymentService) StartPayment(ctx context.Context, orderId uuid.UUID) (*string, error) {
	order, err := p.orderRepo.GetOrderById(ctx, orderId)
	if err != nil {
		return nil, err
	}

	user, err := p.userRepo.GetUserByID(ctx, order.UserID)
	if err != nil {
		return nil, err
	}

	// generate mp items
	items := make([]mp_dtos.MpItem, 0)

	for _, orderItem := range order.Items {
		product, err := p.productRepo.GetProductById(ctx, orderItem.ProductID)
		if err != nil {
			return nil, err
		}

		items = append(items, mp_dtos.MpItem{
			ID:          orderItem.ProductID.String(),
			Title:       product.Name,
			Description: product.SKU,
			CategoryID:  fmt.Sprint(product.CategoryID),
			CurrencyID:  fmt.Sprint(order.Currency),
			Quantity:    int(orderItem.Quantity),
			UnitPrice:   product.Price,
		})
	}

	preferece := p.mp.GeneratePreference(ctx, order, items, user)

	redirectUrl, err := p.mp.ProcessPayment(ctx, preferece)
	if err != nil {
		return nil, err
	}
	return redirectUrl, nil
}

// VerifyPayment implements ports.PaymentProvider.
func (p *PaymentService) VerifyPayment(ctx context.Context, paymentId *string, topic *string) error {
	panic("unimplemented")
}
