package integration

import (
	"context"
	"go-ecommerce/internal/adapters/storage/database/postgres/repository"
	testhelpers "go-ecommerce/internal/test_helpers"
	"go-ecommerce/internal/test_helpers/test_containers"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_UserRepo_With_Postgres_Container(t *testing.T) {
	ctx := context.Background()

	// init postgres container
	cont, err := test_containers.NewPostgresContainerDB(t)
	require.NoError(t, err)
	defer cont.Container.Terminate(ctx)

	// init transaction and when the test ends, execute a rollback
	tx := cont.DB.Begin()
	t.Cleanup(func() { tx.Rollback() })

	// verify if database works
	err = cont.DB.Exec("SELECT 1").Error
	require.NoError(t, err)

	repo := repository.NewUserRepo(tx)
	u := testhelpers.NewDomainUser("john", "john@mail.test")

	// repo testing - create user
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
