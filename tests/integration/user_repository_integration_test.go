package integration

import (
	"context"
	"go-ecommerce/internal/adapters/storage/database/postgres/repository"
	testhelpers "go-ecommerce/internal/test_helpers"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_UserRepo_With_Postgres_Container(t *testing.T) {
	ctx := context.Background()

	// init postgres container
	cont, err := testhelpers.NewPostgresContainerDB(t)
	require.NoError(t, err)

	tx := cont.DB.Begin()
	t.Cleanup(func() { tx.Rollback() })
	defer cont.Container.Terminate(ctx)

	// verify if database works
	err = cont.DB.Exec("SELECT 1").Error
	require.NoError(t, err)

	repo := repository.NewUserRepo(tx)
	u := testhelpers.NewDomainUser("john", "john@mail.test")

	// create user
	create, err := repo.CreateUser(ctx, u)
	require.NoError(t, err)
	assert.Equal(t, u.ID, create.ID)
	assert.Equal(t, u.Name, create.Name)
	assert.Equal(t, u.Email, create.Email)

	// get by id
	userById, err := repo.GetUserByID(ctx, create.ID)
	require.NoError(t, err)
	assert.Equal(t, create.ID, userById.ID)
	assert.Equal(t, create.Name, userById.Name)
	assert.Equal(t, create.Email, userById.Email)
}
