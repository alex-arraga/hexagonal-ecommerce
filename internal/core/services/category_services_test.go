package services_test

import (
	"context"
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

func newCategoryServices(t *testing.T) ports.CategoryService {
	t.Helper()

	db := testhelpers.NewSQLiteTestDB(t)
	tx := db.Begin()
	t.Cleanup(func() { tx.Rollback() })

	redis := mocks.NewMockRedis()

	repo := repository.NewCategoryRepo(tx)
	srv := services.NewCategoryService(repo, redis)

	return srv
}

func Test_CategoryService_Create(t *testing.T) {
	t.Helper()

	ctx := context.Background()
	srv := newCategoryServices(t)

	// factory, create a category
	c := testhelpers.NewDomainCategory("SmartPhones")
	newCateg, err := srv.SaveCategory(ctx, 0, c.Name)

	require.NoError(t, err)
	assert.Equal(t, c.Name, newCateg.Name)
}

func Test_CategoryService_Update(t *testing.T) {
	t.Helper()

	ctx := context.Background()
	srv := newCategoryServices(t)

	// factory, create a category
	c := testhelpers.NewDomainCategory("SmartPhones")
	newCateg, err := srv.SaveCategory(ctx, 0, c.Name)

	require.NoError(t, err)
	assert.Equal(t, c.Name, newCateg.Name)

	// test custom errors
	err = newCateg.UpdateCategory("sm")
	require.Error(t, err)
	assert.Equal(t, err, domain.ErrMinLenghtCategoryNameIsRequire)

	// update category
	name := "smart-phones"
	newCateg.UpdateCategory(name)
	updatedCateg, err := srv.SaveCategory(ctx, 0, name)

	require.NoError(t, err)
	assert.Equal(t, name, updatedCateg.Name)
}

func Test_CategoryService_FindByID(t *testing.T) {
	t.Helper()

	ctx := context.Background()
	srv := newCategoryServices(t)

	// factory, create a category
	c := testhelpers.NewDomainCategory("SmartPhones")

	// create category
	newCateg, err := srv.SaveCategory(ctx, 0, c.Name)
	require.NoError(t, err)
	assert.Equal(t, c.Name, newCateg.Name)

	// find by id
	categ, err := srv.GetCategoryByID(ctx, newCateg.ID)
	require.NoError(t, err)
	assert.Equal(t, newCateg.ID, categ.ID)
	assert.Equal(t, newCateg.Name, categ.Name)
}

func Test_CategoryService_Delete(t *testing.T) {
	t.Helper()

	ctx := context.Background()
	srv := newCategoryServices(t)

	// factory, create a category
	c := testhelpers.NewDomainCategory("SmartPhones")

	// create category
	newCateg, err := srv.SaveCategory(ctx, 0, c.Name)
	require.NoError(t, err)
	assert.Equal(t, c.Name, newCateg.Name)

	// delete
	err = srv.DeleteCategory(ctx, newCateg.ID)
	require.NoError(t, err)

	// find by id
	_, err = srv.GetCategoryByID(ctx, newCateg.ID)
	require.Error(t, domain.ErrCategoriesNotFound, err)
}
