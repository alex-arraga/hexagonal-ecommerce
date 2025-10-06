package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderProductModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	OrderID   uuid.UUID `gorm:"type:uuid;not null"`
	ProductID uuid.UUID `gorm:"type:uuid;not null"`
	Quantity  int16     `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	// Relations
	Order *OrderModel   `gorm:"foreignKey:OrderID;references:ID"`
	Items *ProductModel `gorm:"foreignKey:ProductID;references:ID"`
}

// This function will be executed before to create a new order-product model
func (op *OrderProductModel) BeforeCreate(tx *gorm.DB) (err error) {
	if op.ID == uuid.Nil {
		op.ID = uuid.New()
	}
	return
}
