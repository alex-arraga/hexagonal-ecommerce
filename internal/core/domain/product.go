package domain

import (
	"time"

	"github.com/google/uuid"
)

// Product rules
const (
	minProductNameLength = 3
	minProductSKULength  = 3
)

type Product struct {
	ID         uuid.UUID
	CategoryID uint64
	SKU        string
	Name       string
	Stock      int64
	Price      float64
	Image      string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Category   *Category
}

func NewProduct(name, sku, image string, stock int64, price float64, categoryID uint64) (*Product, error) {
	if len(name) == 0 {
		return nil, ErrProductNameIsRequire
	}

	if stock <= 0 {
		return nil, ErrProductStockIsRequire
	}

	if price <= 0 {
		return nil, ErrProductPriceIsRequire
	}

	if len(sku) == 0 {
		return nil, ErrProductSKUIsRequire
	}
	if len(image) == 0 {
		return nil, ErrProductImageIsRequire
	}

	if categoryID == 0 {
		return nil, ErrProductCategoryIsRequire
	}

	if len(name) < minProductNameLength {
		return nil, ErrProductMinLenghtName
	}

	if len(sku) < minProductSKULength {
		return nil, ErrProductMinLenghtSKU
	}

	now := time.Now()
	return &Product{
		ID:         uuid.Nil, // repository will asign the id
		Name:       name,
		SKU:        sku,
		Stock:      stock,
		Price:      price,
		Image:      image,
		CategoryID: categoryID,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

func (p *Product) Update(name, sku, image string, stock int64, price float64, categoryID uint64) error {
	if len(name) == 0 {
		return ErrProductNameIsRequire
	}

	if stock <= 0 {
		return ErrProductStockIsRequire
	}

	if price <= 0 {
		return ErrProductPriceIsRequire
	}

	if len(sku) == 0 {
		return ErrProductSKUIsRequire
	}

	if len(image) == 0 {
		return ErrProductImageIsRequire
	}

	if categoryID == 0 {
		return ErrProductCategoryIsRequire
	}

	if len(name) < minProductNameLength {
		return ErrProductMinLenghtName
	}

	if len(sku) < minProductSKULength {
		return ErrProductMinLenghtSKU
	}

	// si pasa todas las validaciones, actualizamos los campos
	p.CategoryID = categoryID
	p.SKU = sku
	p.Name = name
	p.Stock = stock
	p.Price = price
	p.Image = image
	p.UpdatedAt = time.Now()

	return nil
}
