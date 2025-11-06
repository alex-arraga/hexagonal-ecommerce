package domain

import (
	"go-ecommerce/internal/core/ports/ports_dtos"
	"time"

	"github.com/google/uuid"
)

// Product rules
const (
	minProductNameLength = 3
	minProductSKULength  = 3
)

type DisscountTypes string

const (
	Percentage DisscountTypes = "percentage" // Percentage = subtotal * (15 / 100)
	Bundle     DisscountTypes = "bundle"     // Bundle = 2x1, 3x2, Pack of 3 services $X
	Fixed      DisscountTypes = "fixed"      // Fixed = $500
)

type Product struct {
	ID            uuid.UUID
	CategoryID    uint64
	SKU           string
	Name          string
	Stock         int64
	Price         float64
	Disscount     float64
	DisscountType DisscountTypes
	Image         string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Category      *Category
}

// TODO -> Cambiar esto por un port_dto
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

	return &Product{
		ID:         uuid.Nil, // repository will asign the id
		Name:       name,
		SKU:        sku,
		Stock:      stock,
		Price:      price,
		Image:      image,
		CategoryID: categoryID,
	}, nil
}

func (p *Product) Update(inputs ports_dtos.SaveProductInputs) error {
	// validations
	if inputs.Name != nil && len(*inputs.Name) == 0 {
		return ErrProductNameIsRequire
	}

	if inputs.Stock != nil && *inputs.Stock <= 0 {
		return ErrProductStockIsRequire
	}

	if inputs.Price != nil && *inputs.Price <= 0 {
		return ErrProductPriceIsRequire
	}

	if inputs.SKU != nil && len(*inputs.SKU) == 0 {
		return ErrProductSKUIsRequire
	}

	if inputs.Image != nil && len(*inputs.Image) == 0 {
		return ErrProductImageIsRequire
	}

	if inputs.CategoryID != nil && *inputs.CategoryID == 0 {
		return ErrProductCategoryIsRequire
	}

	if inputs.Name != nil && len(*inputs.Name) < minProductNameLength {
		return ErrProductMinLenghtName
	}

	if inputs.SKU != nil && len(*inputs.SKU) < minProductSKULength {
		return ErrProductMinLenghtSKU
	}

	// update the existing fields
	if inputs.CategoryID != nil {
		p.CategoryID = *inputs.CategoryID
	}
	if inputs.SKU != nil {
		p.SKU = *inputs.SKU
	}
	if inputs.Name != nil {
		p.Name = *inputs.Name
	}
	if inputs.Stock != nil {
		p.Stock = *inputs.Stock
	}
	if inputs.Price != nil {
		p.Price = *inputs.Price
	}
	if inputs.Image != nil {
		p.Image = *inputs.Image
	}
	p.UpdatedAt = time.Now()

	return nil
}

func (p *Product) ToInputs() ports_dtos.SaveProductInputs {
	return ports_dtos.SaveProductInputs{
		ID:         p.ID,
		Name:       &p.Name,
		Image:      &p.Image,
		SKU:        &p.SKU,
		Price:      &p.Price,
		Stock:      &p.Stock,
		CategoryID: &p.CategoryID,
	}
}
