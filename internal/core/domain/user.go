package domain

import (
	"fmt"
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

type SaveUserInputs struct {
	ID       uuid.UUID
	Name     *string
	Email    *string
	Password *string
	Role     *UserRole
}

// NewUser creates a new user applying the bussiness rules
func NewUser(i SaveUserInputs, hasher PasswordHasher) (*User, error) {
	if i.ID != uuid.Nil {
		return nil, fmt.Errorf("%w: ID must not be provided when creating a user", ErrCreatingUser)
	}

	if i.Name != nil && len(*i.Name) == 0 {
		return nil, ErrNameIsRequire
	}

	if i.Email != nil && len(*i.Email) == 0 {
		return nil, ErrEmailIsRequire
	}

	if i.Password != nil && len(*i.Password) == 0 {
		return nil, ErrPasswordIsRequire
	}

	if i.Name != nil && len(*i.Name) < minNameLength {
		return nil, ErrMinLenghtName
	}

	if i.Password != nil && len(*i.Password) < minPasswordLength {
		return nil, ErrMinLenghtPassword
	}

	if i.Role != nil && len(*i.Role) == 0 {
		return nil, ErrRoleIsRequire
	}

	if *i.Role != Admin && *i.Role != Client && *i.Role != Seller {
		return nil, ErrRoleIsInvalid
	}

	// hash password with the provided hasher
	hashedPassword, err := hasher.Hash(*i.Password)
	if err != nil {
		return nil, ErrHashingPassword
	}

	now := time.Now()
	return &User{
		Name:      *i.Name,
		Email:     *i.Email,
		Password:  hashedPassword,
		Role:      *i.Role,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (u *User) UpdateUser(i SaveUserInputs, hasher PasswordHasher) error {
	if i.Name != nil {
		if len(*i.Name) == 0 {
			return ErrNameIsRequire
		}
		if len(*i.Name) < minNameLength {
			return ErrMinLenghtName
		}
		u.Name = *i.Name
	}

	if i.Password != nil {
		if len(*i.Password) == 0 {
			return ErrPasswordIsRequire
		}
		if len(*i.Password) < minPasswordLength {
			return ErrMinLenghtPassword
		}
		// hash password with the provided hasher if password exist
		hashed, err := hasher.Hash(*i.Password)
		if err != nil {
			return ErrHashingPassword
		}
		u.Password = hashed
	}

	if i.Role != nil {
		if len(*i.Role) == 0 {
			return ErrRoleIsRequire
		}
		if *i.Role != Admin && *i.Role != Client && *i.Role != Seller {
			return ErrRoleIsInvalid
		}
		u.Role = *i.Role
	}

	if i.Email != nil {
		if len(*i.Email) == 0 {
			return ErrEmailIsRequire
		}
		u.Email = *i.Email
	}

	u.UpdatedAt = time.Now()
	return nil
}
