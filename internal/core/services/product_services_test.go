package services_test

import (
	"context"
	"go-ecommerce/internal/adapters/storage/database/postgres/repository"
	"go-ecommerce/internal/core/ports"
	"go-ecommerce/internal/core/services"
	testhelpers "go-ecommerce/internal/test_helpers"
	"go-ecommerce/internal/test_helpers/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ProductService_Create(t *testing.T) {
	t.Helper()

	ctx := context.Background()
	db := testhelpers.NewSQLiteTestDB(t)
	tx := db.Begin()
	t.Cleanup(func() { tx.Rollback() })

	redis := mocks.NewMockRedis()

	// repos
	prodRepo := repository.NewProductRepo(tx)
	prodSrv := services.NewProductService(prodRepo, redis)

	// services
	catRepo := repository.NewCategoryRepo(tx)
	categSrv := services.NewCategoryService(catRepo, redis)

	// factory, create a category as foreign key
	c := testhelpers.NewDomainCategory("Tablets")
	savedCateg, err := categSrv.SaveCategory(ctx, 0, c.Name)
	require.NoError(t, err)
	assert.Equal(t, c.Name, savedCateg.Name)

	// factory, create new product
	p := testhelpers.NewDomainProduct("Ipad 14 pro", savedCateg.ID)
	inputs := ports.SaveProductInputs{
		Name:       p.Name,
		Image:      p.Image,
		SKU:        p.SKU,
		Price:      p.Price,
		Stock:      p.Stock,
		CategoryID: savedCateg.ID,
	}

	savedProd, err := prodSrv.SaveProduct(ctx, inputs)
	require.NoError(t, err)

	assert.Equal(t, p.Name, savedProd.Name)
	assert.Equal(t, p.CategoryID, savedProd.CategoryID)
}
