package repository_test

import (
	"context"
	"go-ecommerce/internal/adapters/storage/database/postgres/repository"
	testhelpers "go-ecommerce/internal/test_helpers"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func newRepoTx(t *testing.T) (*gorm.DB, *repository.UserRepo) {
	db := testhelpers.NewSQLiteTestDB(t)
	tx := db.Begin()
	t.Cleanup(func() { tx.Rollback() })
	return tx, repository.NewUserRepo(tx).(*repository.UserRepo)
}

func Test_CreateUserAndGetUser(t *testing.T) {
	ctx := context.Background()
	_, repo := newRepoTx(t)

	// create user
	u := testhelpers.NewDomainUser("John", "john@gmail.com")
	userCreated, err := repo.CreateUser(ctx, u)

	require.NoError(t, err)
	require.NotNil(t, userCreated)
	assert.Equal(t, u.Name, userCreated.Name)
	assert.Equal(t, u.Email, userCreated.Email)

	// get by id
	getById, err := repo.GetUserByID(ctx, userCreated.ID)
	require.NoError(t, err)
	require.NotNil(t, getById)
	assert.Equal(t, userCreated.ID, getById.ID)

	// get by email
	getByEmail, err := repo.GetUserByEmail(ctx, userCreated.Email)
	require.NoError(t, err)
	require.NotNil(t, getByEmail)
	assert.Equal(t, userCreated.Email, getByEmail.Email)
}
