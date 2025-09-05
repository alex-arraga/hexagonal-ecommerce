package repository

import (
	"context"
	"go-ecommerce/internal/core/domain"
	"go-ecommerce/internal/core/ports"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

// instance the UserRepo struct and return ports.UserRepository interface
func NewUserRepo(db *gorm.DB) ports.UserRepository {
	return &UserRepo{db: db}
}

// CreateUser inserts a new user into the database
func (repo *UserRepo) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	if result := repo.db.Create(user); result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

// GetUserByID selects a user by id
func (repo *UserRepo) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User

	if result := repo.db.First(&user, "id = ?", id); result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// GetUserByEmail selects a user by email
func (repo *UserRepo) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User

	if result := repo.db.First(&user, "email = ?", email); result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// ListUsers selects a list of users with pagination
func (repo *UserRepo) ListUsers(ctx context.Context, skip, limit uint64) ([]domain.User, error) {
	var users []domain.User

	if result := repo.db.Find(&users); result.Error != nil {
		return []domain.User{}, result.Error
	}

	return users, nil
}

// UpdateUser updates a user
func (repo *UserRepo) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	result := repo.db.Model(&domain.User{}).Where("id = ?", user.ID).Updates(user)
	if result.Error != nil {
		return nil, result.Error
	}

	updatedUser, err := repo.GetUserByID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

// DeleteUser deletes a user
func (repo *UserRepo) DeleteUser(ctx context.Context, id uuid.UUID) error {
	var user domain.User

	if result := repo.db.Delete(&user, "id = ?", id); result.Error != nil {
		return result.Error
	}
	return nil
}
