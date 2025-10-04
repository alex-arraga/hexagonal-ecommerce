package services_test

import (
	"context"
	"go-ecommerce/internal/adapters/security"
	"go-ecommerce/internal/adapters/storage/database/postgres/repository"
	"go-ecommerce/internal/core/domain"
	"go-ecommerce/internal/core/ports"
	"go-ecommerce/internal/core/services"
	testhelpers "go-ecommerce/internal/test_helpers"
	"go-ecommerce/internal/test_helpers/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newRepoServices(t *testing.T) ports.UserService {
	t.Helper()

	db := testhelpers.NewSQLiteTestDB(t)
	tx := db.Begin()
	t.Cleanup(func() { tx.Rollback() })

	redis := mocks.NewMockRedis()
	hasher := &security.Hasher{}

	repo := repository.NewUserRepo(tx)
	srv := services.NewUserService(repo, redis, hasher)

	return srv
}

func Test_UserService_Register(t *testing.T) {
	t.Helper()

	ctx := context.Background()
	srv := newRepoServices(t)

	u := testhelpers.NewDomainUser("John", "john@mail.test")

	inputs := domain.SaveUserInputs{
		Name:     &u.Name,
		Email:    &u.Email,
		Password: &u.Password,
		Role:     &u.Role,
	}

	registered, err := srv.SaveUser(ctx, inputs)
	require.NoError(t, err)

	assert.Equal(t, u.Name, registered.Name)
	assert.Equal(t, u.Email, registered.Email)
}
