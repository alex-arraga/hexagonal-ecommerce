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

func Test_UserService_Update(t *testing.T) {
	t.Helper()

	ctx := context.Background()
	srv := newRepoServices(t)

	// createa a new user
	u := testhelpers.NewDomainUser("John", "john@mail.test")

	inputs := domain.SaveUserInputs{
		Name:     &u.Name,
		Email:    &u.Email,
		Password: &u.Password,
		Role:     &u.Role,
	}

	newUser, err := srv.SaveUser(ctx, inputs)
	require.NoError(t, err)

	assert.Equal(t, u.Name, newUser.Name)
	assert.Equal(t, u.Email, newUser.Email)

	// generate update data
	name := "John Alan"
	email := "john_alan@mail.test"

	updateInputs := domain.SaveUserInputs{
		ID:    newUser.ID,
		Name:  &name,
		Email: &email,
	}

	// apply domain rules
	err = newUser.UpdateUser(updateInputs, &security.Hasher{})
	require.NoError(t, err)

	resultUpdate, err := srv.SaveUser(ctx, newUser.ToInputs())

	require.NoError(t, err)
	assert.Equal(t, name, resultUpdate.Name)
	assert.Equal(t, email, resultUpdate.Email)
}

func Test_UserService_GetByIDAndEmail(t *testing.T) {
	t.Helper()

	ctx := context.Background()
	srv := newRepoServices(t)

	// createa a new user
	u := testhelpers.NewDomainUser("John", "john@mail.test")

	inputs := domain.SaveUserInputs{
		Name:     &u.Name,
		Email:    &u.Email,
		Password: &u.Password,
		Role:     &u.Role,
	}

	newUser, err := srv.SaveUser(ctx, inputs)
	require.NoError(t, err)

	assert.Equal(t, u.Name, newUser.Name)
	assert.Equal(t, u.Email, newUser.Email)

	// obtain by id
	userById, err := srv.GetUserByID(ctx, newUser.ID)
	require.NoError(t, err)
	assert.Equal(t, newUser.Email, userById.Email)

	// obtain by email
	userByEmail, err := srv.GetUserByEmail(ctx, newUser.Email)
	require.NoError(t, err)
	assert.Equal(t, newUser.Email, userByEmail.Email)
}
