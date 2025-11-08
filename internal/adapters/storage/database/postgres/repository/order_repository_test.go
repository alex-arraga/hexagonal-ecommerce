package repository_test

import (
	"context"
	"go-ecommerce/internal/adapters/storage/database/postgres/repository"
	"go-ecommerce/internal/core/ports/ports_dtos"
	"go-ecommerce/internal/core/services"
	testhelpers "go-ecommerce/internal/test_helpers"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// dependencies
type depToTestingOrderRepo struct {
	userRepo  repository.UserRepo
	prodRepo  repository.ProductRepo
	categRepo repository.CategoryRepo
	orderRepo repository.OrderRepo
}

func newOrderRepoTx(t *testing.T) (*gorm.DB, *depToTestingOrderRepo) {
	db := testhelpers.NewSQLiteTestDB(t)
	tx := db.Begin()
	t.Cleanup(func() { tx.Rollback() })

	userRepo := repository.NewUserRepo(tx).(*repository.UserRepo)
	prodRepo := repository.NewProductRepo(tx).(*repository.ProductRepo)
	categRepo := repository.NewCategoryRepo(tx).(*repository.CategoryRepo)

	orderProdRepo := repository.NewOrderProductRepo(tx).(*repository.OrderProductRepo)
	orderProdSrv := services.NewOrderProductService(orderProdRepo)

	orderRepo := repository.NewOrderRepo(orderProdSrv, tx).(*repository.OrderRepo)

	repos := &depToTestingOrderRepo{
		userRepo:  *userRepo,
		categRepo: *categRepo,
		prodRepo:  *prodRepo,
		orderRepo: *orderRepo,
	}

	return tx, repos
}

func Test_CreateAndGetOrder(t *testing.T) {
	ctx := context.Background()
	_, repos := newOrderRepoTx(t)

	u := testhelpers.NewDomainUser("John", "john@test.com")
	newUser, err := repos.userRepo.SaveUser(ctx, u)
	require.NoError(t, err)
	assert.Equal(t, u.Name, newUser.Name)

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

	// create order
	o := testhelpers.NewDomainOrder(newUser.ID)
	newOrder, err := repos.orderRepo.SaveOrder(ctx, o)
	require.NoError(t, err)
	assert.Equal(t, o.UserID, newOrder.UserID)

	// get order by id
	recoveredOrder, err := repos.orderRepo.GetOrderById(ctx, newOrder.ID)
	require.NoError(t, err)
	require.NotNil(t, recoveredOrder)
	assert.Equal(t, newOrder.ID, recoveredOrder.ID)
}

func Test_ListOrder(t *testing.T) {
	ctx := context.Background()
	_, repos := newOrderRepoTx(t)

	u := testhelpers.NewDomainUser("John", "john@test.com")
	newUser, err := repos.userRepo.SaveUser(ctx, u)
	require.NoError(t, err)
	assert.Equal(t, u.Name, newUser.Name)

	// creates 20 orders
	for range 20 {
		categName := faker.Name()
		prodName := faker.Name()

		c := testhelpers.NewDomainCategory(categName)
		categ, err := repos.categRepo.SaveCategory(ctx, c)
		require.NoError(t, err)

		p := testhelpers.NewDomainProduct(prodName, categ.ID)
		_, err = repos.prodRepo.SaveProduct(ctx, p)
		require.NoError(t, err)

		// create order
		o := testhelpers.NewDomainOrder(newUser.ID)
		newOrder, err := repos.orderRepo.SaveOrder(ctx, o)
		require.NoError(t, err)
		assert.Equal(t, o.UserID, newOrder.UserID)
	}

	orders, err := repos.orderRepo.ListOrders(ctx)
	require.NoError(t, err)
	assert.Len(t, orders, 20)
}

func Test_UpdateOrder(t *testing.T) {
	ctx := context.Background()
	_, repos := newOrderRepoTx(t)

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
