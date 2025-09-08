package testhelpers

import (
	"go-ecommerce/internal/core/domain"
	"time"

	"github.com/google/uuid"
)

func NewDomainUser(name, email string) *domain.User {
	return &domain.User{
		ID:        uuid.New(),
		Name:      name,
		Email:     email,
		Password:  "password",
		Role:      domain.Client,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
