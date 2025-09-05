package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	Admin  UserRole = "admin"
	Client UserRole = "client"
	Seller UserRole = "seller"
)

type User struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Password  string
	Role      UserRole
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Bussiness rules (constantes)
const (
	minPasswordLength = 6
	minNameLength     = 3
)

// NewUser creates a new user applying the bussiness rules
func NewUser(name, email, password string, role UserRole) (*User, error) {
	if len(name) < minNameLength {
		return nil, errors.New("name must have at least 3 characters")
	}

	if len(password) < minPasswordLength {
		return nil, errors.New("password must have at least 6 characters")
	}

	if role != Admin && role != Client && role != Seller {
		return nil, errors.New("invalid user role")
	}

	now := time.Now()
	return &User{
		ID:        uuid.New(),
		Name:      name,
		Email:     email,
		Password:  password,
		Role:      role,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
