package domain

import (
	"time"

	"github.com/google/uuid"
)

const ( // Product errors
	minProductNameLength = 3
)

type Product struct {
	ID         uint64
	CategoryID uint64
	SKU        uuid.UUID
	Name       string
	Stock      int64
	Price      float64
	Image      string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Category   *Category
}

func NewProduct(categoryID uint64, name string, stock int64, price float64, image string) (*Product, error) {
	if len(name) == 0 {
		return nil, ErrProductNameIsRequire
	}

	if len(name) < minProductNameLength {
		return nil, ErrProductMinLenghtName
	}

	if stock <= 0 {
		return nil, ErrProductStockIsRequire
	}

	if price <= 0 {
		return nil, ErrProductPriceIsRequire
	}

	if len(image) == 0 {
		return nil, ErrProductImageIsRequire
	}

	if categoryID == 0 {
		return nil, ErrProductCategoryIsRequire
	}

	now := time.Now()
	return &Product{
		CategoryID: categoryID,
		SKU:        uuid.New(),
		Name:       name,
		Stock:      stock,
		Price:      price,
		Image:      image,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}
