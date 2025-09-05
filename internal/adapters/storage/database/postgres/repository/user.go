package repository

import (
	"context"
	"go-ecommerce/internal/adapters/storage/database/postgres/database_dtos"
	"go-ecommerce/internal/adapters/storage/database/postgres/models"
	"go-ecommerce/internal/core/domain"
	"go-ecommerce/internal/core/ports"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

/**
 * UserRepository implements port.UserRepository interface
 * and provides an access to the postgres database
 */
type UserRepo struct {
	db *gorm.DB
}

// instance the UserRepo struct and return ports.UserRepository interface
func NewUserRepo(db *gorm.DB) ports.UserRepository {
	return &UserRepo{db: db}
}

// CreateUser inserts a new user into the database
func (repo *UserRepo) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	userdb := database_dtos.CovertToDBUser(user)

	if result := repo.db.Create(userdb); result.Error != nil {
		return nil, result.Error
	}

	domainUser := database_dtos.CovertToDomainUser(userdb)
	return domainUser, nil
}

// GetUserByID selects a user by id
func (repo *UserRepo) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var dbUser *models.UserModel

	if result := repo.db.First(dbUser, "id = ?", id); result.Error != nil {
		return nil, result.Error
	}

	domainUser := database_dtos.CovertToDomainUser(dbUser)
	return domainUser, nil
}

// GetUserByEmail selects a user by email
func (repo *UserRepo) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var dbUser *models.UserModel

	if result := repo.db.First(&dbUser, "email = ?", email); result.Error != nil {
		return nil, result.Error
	}

	domainUser := database_dtos.CovertToDomainUser(dbUser)
	return domainUser, nil
}

// ListUsers selects a list of users with pagination
func (repo *UserRepo) ListUsers(ctx context.Context, skip, limit uint64) ([]*domain.User, error) {
	var dbUsers []*models.UserModel

	if result := repo.db.Find(&dbUsers); result.Error != nil {
		return []*domain.User{}, result.Error
	}

	domainUsers := database_dtos.CovertToDomainUsers(dbUsers)
	return domainUsers, nil
}

// UpdateUser updates a user
func (repo *UserRepo) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	var dbUser *models.UserModel

	result := repo.db.Model(dbUser).Where("id = ?", user.ID).Updates(user)
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
	var dbUser *models.UserModel

	if result := repo.db.Delete(&dbUser, "id = ?", id); result.Error != nil {
		return result.Error
	}
	return nil
}
