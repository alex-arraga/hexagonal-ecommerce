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

// Bussiness rules (constantes)
const (
	minPasswordLength = 6
	minNameLength     = 3
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type User struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Password  string
	Role      UserRole
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser creates a new user applying the bussiness rules
func NewUser(name, email, password string, role UserRole, hasher PasswordHasher) (*User, error) {
	if role != Admin && role != Client && role != Seller {
		return nil, errors.New("role is required")
	}

	if len(name) < minNameLength {
		return nil, errors.New("name must have at least 3 characters")
	}

	if len(password) < minPasswordLength {
		return nil, errors.New("password must have at least 6 characters")
	}

	if role != Admin && role != Client && role != Seller {
		return nil, errors.New("invalid user role")
	}

	// hash password with the provided hasher
	hashedPassword, err := hasher.Hash(password)
	if err != nil {
		return nil, errors.New("could not hash password")
	}

	now := time.Now()
	return &User{
		ID:        uuid.New(),
		Name:      name,
		Email:     email,
		Password:  hashedPassword,
		Role:      role,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
