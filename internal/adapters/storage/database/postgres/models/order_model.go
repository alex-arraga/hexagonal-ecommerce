package models

import (
	"go-ecommerce/internal/core/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderModel struct {
	ID                uuid.UUID               `gorm:"type:uuid;primaryKey"`
	Providers         domain.Providers        `gorm:"type:varchar(50)"`
	UserID            uuid.UUID               `gorm:"type:uuid"`
	PaymentID         *string                 `gorm:"type:varchar(255)"`
	SecureToken       *uuid.UUID              `gorm:"type:uuid"`
	ExternalReference *string                 `gorm:"type:varchar(255)"`
	Currency          domain.Currencies       `gorm:"type:varchar(10)"`
	SubTotal          float64                 `gorm:"type:numeric"`
	Disscount         *float64                `gorm:"type:numeric"`
	DisscountType     *domain.DisscountTypes  `gorm:"type:varchar(50)"`
	Total             float64                 `gorm:"type:numeric"`
	Paid              bool                    `gorm:"type:boolean"`
	PayStatus         domain.PayStatus        `gorm:"type:varchar(50)"`
	PayStatusDetail   *domain.PayStatusDetail `gorm:"type:varchar(100)"`
	CreatedAt         time.Time               `gorm:"autoCreateTime"`
	UpdatedAt         time.Time               `gorm:"autoUpdateTime"`
	ExpiresAt         time.Time               `gorm:"type:timestamp"`

	// Relations
	User *UserModel `gorm:"foreignKey:UserID;references:ID"`
}

func (o *OrderModel) BeforeCreate(tx *gorm.DB) (err error) {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return
}
