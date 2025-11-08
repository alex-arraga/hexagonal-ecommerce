package repository_test

import (
	"context"
	"go-ecommerce/internal/adapters/storage/database/postgres/repository"
	"go-ecommerce/internal/core/domain"
	"go-ecommerce/internal/core/ports/ports_dtos"
	testhelpers "go-ecommerce/internal/test_helpers"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type reposToProductTest struct {
	prodRepo  repository.ProductRepo
	categRepo repository.CategoryRepo
}

func newProductRepoTx(t *testing.T) (*gorm.DB, *reposToProductTest) {
	db := testhelpers.NewSQLiteTestDB(t)
	tx := db.Begin()
	t.Cleanup(func() { tx.Rollback() })
	prodRepo := repository.NewProductRepo(tx).(*repository.ProductRepo)
	categRepo := repository.NewCategoryRepo(tx).(*repository.CategoryRepo)

	repos := &reposToProductTest{
		prodRepo:  *prodRepo,
		categRepo: *categRepo,
	}

	return tx, repos
}

func Test_CreateAndGetProduct(t *testing.T) {
	ctx := context.Background()
	_, repos := newProductRepoTx(t)

	c := testhelpers.NewDomainCategory("SmartPhones")
	newCateg, err := repos.categRepo.SaveCategory(ctx, c)
	require.NoError(t, err)
	assert.Equal(t, c.Name, newCateg.Name)

	// create product
	p := testhelpers.NewDomainProduct("Iphone 15 Pro Max", newCateg.ID)
	newProduct, err := repos.prodRepo.SaveProduct(ctx, p)

	require.NoError(t, err)
	require.NotNil(t, newProduct)
	assert.Equal(t, p.Name, newProduct.Name)
	assert.Equal(t, p.Price, newProduct.Price)
	assert.Equal(t, p.SKU, newProduct.SKU)

	// get by id
	getById, err := repos.prodRepo.GetProductById(ctx, newProduct.ID)
	require.NoError(t, err)
	require.NotNil(t, getById)
	assert.Equal(t, newProduct.ID, getById.ID)
}

func Test_ListProducts(t *testing.T) {
	ctx := context.Background()
	_, repos := newProductRepoTx(t)

	// creates 20 products and categories
	for range 20 {
		categName := faker.Name()
		prodName := faker.Name()

		c := testhelpers.NewDomainCategory(categName)
		categ, err := repos.categRepo.SaveCategory(ctx, c)
		require.NoError(t, err)

		p := testhelpers.NewDomainProduct(prodName, categ.ID)
		_, err = repos.prodRepo.SaveProduct(ctx, p)
		require.NoError(t, err)
	}

	products, err := repos.prodRepo.ListProducts(ctx)
	require.NoError(t, err)
	assert.Len(t, products, 20)
}

func Test_UpdateProduct(t *testing.T) {
	ctx := context.Background()
	_, repos := newProductRepoTx(t)

	c := testhelpers.NewDomainCategory("SmartPhones")
	newCateg, err := repos.categRepo.SaveCategory(ctx, c)
	require.NoError(t, err)
	assert.Equal(t, c.Name, newCateg.Name)

	// create product
	p := testhelpers.NewDomainProduct("Iphone 15 Pro Max", newCateg.ID)
	newProduct, err := repos.prodRepo.SaveProduct(ctx, p)

	require.NoError(t, err)
	require.NotNil(t, newProduct)
	assert.Equal(t, p.Name, newProduct.Name)
	assert.Equal(t, p.Price, newProduct.Price)
	assert.Equal(t, p.SKU, newProduct.SKU)

	// get by id
	recoveredProduct, err := repos.prodRepo.GetProductById(ctx, newProduct.ID)
	require.NoError(t, err)
	require.NotNil(t, recoveredProduct)
	assert.Equal(t, newProduct.ID, recoveredProduct.ID)

	// update data
	newName := "Iphone 15 Pro Max - Rosa"
	newPrice := 120.99
	var newStock int64 = 12

	updateData := ports_dtos.SaveProductInputs{
		Name:  &newName,
		Price: &newPrice,
		Stock: &newStock,
	}
	newProduct.Update(updateData)

	// save the new data of product
	updatedProduct, err := repos.prodRepo.SaveProduct(ctx, newProduct)
	require.NoError(t, err)
	assert.Equal(t, newProduct.Name, updatedProduct.Name)
	assert.Equal(t, newProduct.Price, updatedProduct.Price)
	assert.Equal(t, newProduct.Stock, updatedProduct.Stock)
}

func Test_DeleteProduct(t *testing.T) {
	ctx := context.Background()
	_, repos := newProductRepoTx(t)

	c := testhelpers.NewDomainCategory("SmartPhones")
	newCateg, err := repos.categRepo.SaveCategory(ctx, c)
	require.NoError(t, err)
	assert.Equal(t, c.Name, newCateg.Name)

	// create product
	p := testhelpers.NewDomainProduct("Iphone 15 Pro Max", newCateg.ID)
	newProduct, err := repos.prodRepo.SaveProduct(ctx, p)

	require.NoError(t, err)
	require.NotNil(t, newProduct)
	assert.Equal(t, p.Name, newProduct.Name)
	assert.Equal(t, p.Price, newProduct.Price)
	assert.Equal(t, p.SKU, newProduct.SKU)

	// get by id
	recoveredProduct, err := repos.prodRepo.GetProductById(ctx, newProduct.ID)
	require.NoError(t, err)
	require.NotNil(t, recoveredProduct)
	assert.Equal(t, newProduct.ID, recoveredProduct.ID)

	err = repos.prodRepo.DeleteProduct(ctx, recoveredProduct.ID)
	require.NoError(t, err)

	// get by id
	_, err = repos.prodRepo.GetProductById(ctx, newProduct.ID)
	require.Error(t, domain.ErrProductNotFound, err)
}
