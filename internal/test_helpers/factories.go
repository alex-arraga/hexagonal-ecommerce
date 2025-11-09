package testhelpers

import (
	"go-ecommerce/internal/core/domain"
	"time"

	"github.com/google/uuid"
)

func NewDomainUser(name, email string) *domain.User {
	return &domain.User{
		Name:      name,
		Email:     email,
		Password:  "password",
		Role:      domain.Client,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func NewDomainCategory(name string) *domain.Category {
	return &domain.Category{
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func NewDomainProduct(name string, categoryId uint64) *domain.Product {
	return &domain.Product{
		Name:       name,
		CategoryID: categoryId,
		SKU:        "product-test",
		Stock:      100,
		Price:      10,
		Image:      "product-image-test",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

func NewDomainOrder(userId uuid.UUID) *domain.Order {
	externalRef := uuid.New().String()

	return &domain.Order{
		Providers:         domain.MercadoPago,
		UserID:            userId,
		PaymentID:         nil,
		ExternalReference: &externalRef,
		Currency:          domain.ARS,
		Fee:               nil,
		Installments:      nil,
		PayStatus:         domain.Approved,
		NetReceivedAmount: nil,
		PayMethod:         nil,
		PayResource:       nil,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
}

func NewDomainOrderProduct(orderId, productId uuid.UUID, quantity int16) *domain.OrderProduct {
	return &domain.OrderProduct{
		OrderID:   orderId,
		ProductID: productId,
		Quantity:  quantity,
	}
}
