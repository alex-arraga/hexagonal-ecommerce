package domain

import (
	"time"

	"github.com/google/uuid"
)

// OrderProduct is an entity that represents pivot table between order and product
type OrderProduct struct {
	ID        uuid.UUID
	OrderID   uuid.UUID
	ProductID uuid.UUID
	Quantity  int16
	CreatedAt time.Time
	UpdatedAt time.Time

	// Relations
	Order *Order
	Items *Product
}

func NewOrderProduct(orderID, productID uuid.UUID, quantity int16) *OrderProduct {
	return &OrderProduct{
		ID:        uuid.Nil, // repository will asign the id
		OrderID:   orderID,
		ProductID: productID,
		Quantity:  quantity,
	}
}
