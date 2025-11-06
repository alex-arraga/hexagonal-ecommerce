package ports_dtos

import "github.com/google/uuid"

// ProductService is an interface for interacting with product-related business logic
type SaveProductInputs struct {
	ID         uuid.UUID
	Name       *string
	Image      *string
	SKU        *string
	Price      *float64
	Stock      *int64
	CategoryID *uint64
}
