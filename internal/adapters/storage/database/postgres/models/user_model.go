package models

import (
	"go-ecommerce/internal/core/domain"
	"time"

	"github.com/google/uuid"
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
