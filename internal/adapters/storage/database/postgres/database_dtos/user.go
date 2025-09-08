package database_dtos

import (
	"go-ecommerce/internal/adapters/storage/database/postgres/models"
	"go-ecommerce/internal/core/domain"
)

// domain.User -> DB model
func CovertToDBUser(u *domain.User) *models.UserModel {
	return &models.UserModel{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Password:  u.Password,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// domain.Users -> DB models
func CovertToDBUsers(users []*domain.User) []*models.UserModel {
	dbUsers := make([]*models.UserModel, len(users))

	for _, u := range users {
		dbUsers = append(dbUsers, &models.UserModel{
			ID:        u.ID,
			Name:      u.Name,
			Email:     u.Email,
			Password:  u.Password,
			Role:      u.Role,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		})
	}

	return dbUsers
}

// DB model -> domain.User
func CovertToDomainUser(u *models.UserModel) *domain.User {
	return &domain.User{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Password:  u.Password,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// DB models -> domain.Users
func CovertToDomainUsers(users []*models.UserModel) []*domain.User {
	var domainUsers []*domain.User

	for _, u := range users {
		domainUsers = append(domainUsers, &domain.User{
			ID:        u.ID,
			Name:      u.Name,
			Email:     u.Email,
			Password:  u.Password,
			Role:      u.Role,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		})
	}

	return domainUsers
}
