package models

import (
	"time"

	"github.com/google/uuid"
)

type ProductModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name      string    `gorm:"size:255;not null"`
	SKU       string    `gorm:"size:255;not null"`
	Stock     int64     `gorm:"not null"`
	Price     float64   `gorm:"not null"`
	Image     string    `gorm:"size:255;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	CategoryID uint64        `gorm:"not null"`
	Category   CategoryModel `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}
