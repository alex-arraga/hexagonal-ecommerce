package domain

import (
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
	if len(name) < minNameLength {
		return nil, ErrMinLenghtName
	}

	if len(password) < minPasswordLength {
		return nil, ErrMinLenghtPassword
	}

	if len(role) == 0 {
		return nil, ErrRoleIsRequire
	}

	if role != Admin && role != Client && role != Seller {
		return nil, ErrRoleIsInvalid
	}

	// hash password with the provided hasher
	hashedPassword, err := hasher.Hash(password)
	if err != nil {
		return nil, ErrHashingPassword
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
