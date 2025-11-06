package services_test

import (
	"context"
	"go-ecommerce/internal/adapters/storage/database/postgres/repository"
	"go-ecommerce/internal/core/ports/ports_dtos"
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
	inputs := ports_dtos.SaveProductInputs{
		Name:       &p.Name,
		Image:      &p.Image,
		SKU:        &p.SKU,
		Price:      &p.Price,
		Stock:      &p.Stock,
		CategoryID: &savedCateg.ID,
	}

	savedProd, err := prodSrv.SaveProduct(ctx, inputs)
	require.NoError(t, err)

	assert.Equal(t, p.Name, savedProd.Name)
	assert.Equal(t, p.CategoryID, savedProd.CategoryID)
}

func Test_ProductService_Update(t *testing.T) {
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
	inputs := ports_dtos.SaveProductInputs{
		Name:       &p.Name,
		Image:      &p.Image,
		SKU:        &p.SKU,
		Price:      &p.Price,
		Stock:      &p.Stock,
		CategoryID: &savedCateg.ID,
	}

	// save new product
	newProd, err := prodSrv.SaveProduct(ctx, inputs)
	require.NoError(t, err)

	assert.Equal(t, p.Name, newProd.Name)
	assert.Equal(t, p.CategoryID, newProd.CategoryID)

	// update product with bussiness rules
	name := "Ipad 16 pro max"
	image := "image-ipad16-test"
	price := 199.99
	var stock int64 = 48

	updateInputs := ports_dtos.SaveProductInputs{
		ID:    newProd.ID,
		Name:  &name,
		Image: &image,
		Price: &price,
		Stock: &stock,
	}
	err = newProd.Update(updateInputs)
	require.NoError(t, err)

	lal := newProd.ToInputs()

	// save the updated product
	updatedProd, err := prodSrv.SaveProduct(ctx, lal)
	require.NoError(t, err)

	assert.Equal(t, name, updatedProd.Name)
	assert.Equal(t, image, updatedProd.Image)
	assert.Equal(t, price, updatedProd.Price)
	assert.Equal(t, stock, updatedProd.Stock)
}
