package models

import "time"

type CategoryModel struct {
	ID        uint64    `gorm:"primaryKey"`
	Name      string    `gorm:"size:255;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
