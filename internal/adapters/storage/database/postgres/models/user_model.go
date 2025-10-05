package models

import (
	"go-ecommerce/internal/core/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserModel struct {
	ID        uuid.UUID       `gorm:"type:uuid;primaryKey"`
	Name      string          `gorm:"size:255;not null"`
	Email     string          `gorm:"size:255;unique;not null"`
	Password  string          `gorm:"size:255;not null"`
	Role      domain.UserRole `gorm:"size:10;not null;default:client"`
	CreatedAt time.Time       `gorm:"autoCreateTime"`
	UpdatedAt time.Time       `gorm:"autoUpdateTime"`
}

// This function will be executed before to create a new product model
func (u *UserModel) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}
