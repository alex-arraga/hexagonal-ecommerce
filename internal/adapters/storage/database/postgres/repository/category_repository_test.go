package repository_test

import (
	"context"
	"go-ecommerce/internal/adapters/storage/database/postgres/repository"
	testhelpers "go-ecommerce/internal/test_helpers"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// dependencies
type depToTestingCategoryRepo struct {
	categRepo repository.CategoryRepo
}

func newCategoryRepoTx(t *testing.T) (*gorm.DB, *depToTestingCategoryRepo) {
	db := testhelpers.NewSQLiteTestDB(t)
	tx := db.Begin()
	t.Cleanup(func() { tx.Rollback() })

	categRepo := repository.NewCategoryRepo(tx).(*repository.CategoryRepo)

	repos := &depToTestingCategoryRepo{
		categRepo: *categRepo,
	}

	return tx, repos
}

func Test_CreateAndGetCategory(t *testing.T) {
	ctx := context.Background()
	_, repos := newCategoryRepoTx(t)

	c := testhelpers.NewDomainCategory("SmartPhones")
	newCateg, err := repos.categRepo.SaveCategory(ctx, c)
	require.NoError(t, err)
	assert.Equal(t, c.Name, newCateg.Name)
}

func Test_ListCategories(t *testing.T) {
	ctx := context.Background()
	_, repos := newCategoryRepoTx(t)

	// creates 20 orders
	for range 5 {
		categName := faker.Name()

		c := testhelpers.NewDomainCategory(categName)
		categ, err := repos.categRepo.SaveCategory(ctx, c)
		require.NoError(t, err)
		assert.Equal(t, c.Name, categ.Name)
	}

	categories, err := repos.categRepo.ListCategories(ctx)
	require.NoError(t, err)
	assert.Len(t, categories, 5)
}

func Test_UpdateCategory(t *testing.T) {
	ctx := context.Background()
	_, repos := newCategoryRepoTx(t)

	// create new category
	c := testhelpers.NewDomainCategory("SmartPhones")
	newCateg, err := repos.categRepo.SaveCategory(ctx, c)
	require.NoError(t, err)
	assert.Equal(t, c.Name, newCateg.Name)

	// update category
	newCateg.UpdateCategory("smart-phones")
	updatedCateg, err := repos.categRepo.SaveCategory(ctx, newCateg)
	require.NoError(t, err)
	assert.Equal(t, newCateg.ID, updatedCateg.ID)
	assert.Equal(t, newCateg.Name, updatedCateg.Name)
}
