package repository_test

import (
	"context"
	"go-ecommerce/internal/adapters/storage/database/postgres/repository"
	testhelpers "go-ecommerce/internal/test_helpers"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func newUserRepoTx(t *testing.T) (*gorm.DB, *repository.UserRepo) {
	db := testhelpers.NewSQLiteTestDB(t)
	tx := db.Begin()
	t.Cleanup(func() { tx.Rollback() })
	return tx, repository.NewUserRepo(tx).(*repository.UserRepo)
}

func Test_CreateUserAndGetUser(t *testing.T) {
	ctx := context.Background()
	_, repo := newUserRepoTx(t)

	// create user
	u := testhelpers.NewDomainUser("John", "john@mail.test")
	userCreated, err := repo.SaveUser(ctx, u)

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

func Test_ListUsers(t *testing.T) {
	ctx := context.Background()
	_, repo := newUserRepoTx(t)

	// creates 20 users
	for range 20 {
		name := faker.FirstName()
		email := faker.Email(options.WithCustomDomain("test"))
		u := testhelpers.NewDomainUser(name, email)
		_, err := repo.SaveUser(ctx, u)
		require.NoError(t, err)
	}

	users1, err := repo.ListUsers(ctx, 0, 10)
	require.NoError(t, err)
	assert.Len(t, users1, 10)

	users2, err := repo.ListUsers(ctx, 10, 10)
	require.NoError(t, err)
	assert.Len(t, users2, 10)
	assert.True(t, len(users2) >= 1) // almost one user rest
}

func Test_UpdateUser(t *testing.T) {
	ctx := context.Background()
	_, repo := newUserRepoTx(t)

	u := testhelpers.NewDomainUser("john", "john@mail.test")
	created, err := repo.SaveUser(ctx, u)
	require.NoError(t, err)

	newName := faker.FirstName()
	newPass := faker.Password()
	newEmail := faker.Email()

	created.Name = newName
	created.Password = newPass
	created.Email = newEmail

	// update all data
	updated, err := repo.SaveUser(ctx, created)
	require.NoError(t, err)
	assert.Equal(t, newName, updated.Name)
	assert.Equal(t, newPass, updated.Password)
	assert.Equal(t, newEmail, updated.Email)
}

func Test_DeleteUser(t *testing.T) {
	ctx := context.Background()
	_, repo := newUserRepoTx(t)

	u := testhelpers.NewDomainUser("john", "john@mail.test")
	created, err := repo.SaveUser(ctx, u)
	require.NoError(t, err)

	// deletes user
	err = repo.DeleteUser(ctx, created.ID)
	require.NoError(t, err)
	assert.Nil(t, err)

	// get by id
	getById, err := repo.GetUserByID(ctx, created.ID)
	require.Nil(t, getById)
	require.ErrorIs(t, err, gorm.ErrRecordNotFound)
}
