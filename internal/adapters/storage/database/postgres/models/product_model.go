package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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

	CategoryID uint64         `gorm:"not null"`
	Category   *CategoryModel `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// This function will be executed before to create a new product model
func (p *ProductModel) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return
}
