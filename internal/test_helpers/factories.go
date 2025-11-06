package testhelpers

import (
	"go-ecommerce/internal/core/domain"
	"time"
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
		Price:      15,
		Image:      "product-image-test",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}
