package repository_test

import (
	"context"
	"go-ecommerce/internal/adapters/storage/database/postgres/repository"
	"go-ecommerce/internal/core/services"
	testhelpers "go-ecommerce/internal/test_helpers"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// dependencies
type depToTestingOrderProductRepo struct {
	userRepo         repository.UserRepo
	prodRepo         repository.ProductRepo
	categRepo        repository.CategoryRepo
	orderRepo        repository.OrderRepo
	orderProductRepo repository.OrderProductRepo
}

func newOrderProductRepoTx(t *testing.T) (*gorm.DB, *depToTestingOrderProductRepo) {
	db := testhelpers.NewSQLiteTestDB(t)
	tx := db.Begin()
	t.Cleanup(func() { tx.Rollback() })

	userRepo := repository.NewUserRepo(tx).(*repository.UserRepo)
	prodRepo := repository.NewProductRepo(tx).(*repository.ProductRepo)
	categRepo := repository.NewCategoryRepo(tx).(*repository.CategoryRepo)

	orderProdRepo := repository.NewOrderProductRepo(tx).(*repository.OrderProductRepo)
	orderProdSrv := services.NewOrderProductService(orderProdRepo)

	orderRepo := repository.NewOrderRepo(orderProdSrv, tx).(*repository.OrderRepo)

	repos := &depToTestingOrderProductRepo{
		userRepo:         *userRepo,
		categRepo:        *categRepo,
		prodRepo:         *prodRepo,
		orderRepo:        *orderRepo,
		orderProductRepo: *orderProdRepo,
	}

	return tx, repos
}

func Test_CreateAndGetOrderProduct(t *testing.T) {
	ctx := context.Background()
	_, repos := newOrderProductRepoTx(t)

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

	op := testhelpers.NewDomainOrderProduct(newOrder.ID, newProduct.ID, 4)
	newOrderProduct, err := repos.orderProductRepo.SaveOrderProduct(ctx, op)
	require.NoError(t, err)
	assert.Equal(t, op.OrderID, newOrderProduct.OrderID)
	assert.Equal(t, op.ProductID, newOrderProduct.ProductID)
	assert.Equal(t, op.Quantity, newOrderProduct.Quantity)
}

func Test_ListOrderProduct(t *testing.T) {
	ctx := context.Background()
	_, repos := newOrderProductRepoTx(t)

	u := testhelpers.NewDomainUser("John", "john@test.com")
	newUser, err := repos.userRepo.SaveUser(ctx, u)
	require.NoError(t, err)
	assert.Equal(t, u.Name, newUser.Name)

	categName := faker.Name()
	prodName := faker.Name()

	c := testhelpers.NewDomainCategory(categName)
	categ, err := repos.categRepo.SaveCategory(ctx, c)
	require.NoError(t, err)

	p := testhelpers.NewDomainProduct(prodName, categ.ID)
	newProduct, err := repos.prodRepo.SaveProduct(ctx, p)
	require.NoError(t, err)

	// create order
	o := testhelpers.NewDomainOrder(newUser.ID)
	newOrder1, err := repos.orderRepo.SaveOrder(ctx, o)
	require.NoError(t, err)
	assert.Equal(t, o.UserID, newOrder1.UserID)

	newOrder2, err := repos.orderRepo.SaveOrder(ctx, o)
	require.NoError(t, err)
	assert.Equal(t, o.UserID, newOrder2.UserID)

	// add 4 order-products to order-1
	for range 20 {
		op := testhelpers.NewDomainOrderProduct(newOrder1.ID, newProduct.ID, 6)
		newOrderProduct, err := repos.orderProductRepo.SaveOrderProduct(ctx, op)
		require.NoError(t, err)
		assert.Equal(t, op.ProductID, newOrderProduct.ProductID)
		assert.Equal(t, op.OrderID, newOrderProduct.OrderID)
		assert.Equal(t, op.Quantity, newOrderProduct.Quantity)
	}

	// add 4 order-products to order-2
	for range 4 {
		op := testhelpers.NewDomainOrderProduct(newOrder2.ID, newProduct.ID, 6)
		newOrderProduct, err := repos.orderProductRepo.SaveOrderProduct(ctx, op)
		require.NoError(t, err)
		assert.Equal(t, op.ProductID, newOrderProduct.ProductID)
		assert.Equal(t, op.OrderID, newOrderProduct.OrderID)
		assert.Equal(t, op.Quantity, newOrderProduct.Quantity)
	}

	// find all order-products
	orderProducts, err := repos.orderProductRepo.ListOrderProducts(ctx, uuid.Nil)
	require.NoError(t, err)
	assert.Len(t, orderProducts, 24)

	// find all order-products by orderID
	opsByOderId1, err := repos.orderProductRepo.ListOrderProducts(ctx, newOrder1.ID)
	require.NoError(t, err)
	assert.Len(t, opsByOderId1, 20)

	// find all order-products by orderID
	opsByOderId2, err := repos.orderProductRepo.ListOrderProducts(ctx, newOrder2.ID)
	require.NoError(t, err)
	assert.Len(t, opsByOderId2, 4)
}

func Test_UpdateOrderProduct(t *testing.T) {
	ctx := context.Background()
	_, repos := newOrderProductRepoTx(t)

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

	// factory new order-product and save in db
	op := testhelpers.NewDomainOrderProduct(newOrder.ID, newProduct.ID, 4)
	newOrderProduct, err := repos.orderProductRepo.SaveOrderProduct(ctx, op)
	require.NoError(t, err)
	assert.Equal(t, op.OrderID, newOrderProduct.OrderID)
	assert.Equal(t, op.ProductID, newOrderProduct.ProductID)
	assert.Equal(t, op.Quantity, newOrderProduct.Quantity)

	err = newOrderProduct.UpdateOrderProduct(2)
	require.NoError(t, err)

	// update order-product and verify
	updateOrderProduct, err := repos.orderProductRepo.SaveOrderProduct(ctx, newOrderProduct)
	require.NoError(t, err)
	assert.Equal(t, newOrderProduct.OrderID, updateOrderProduct.OrderID)
	assert.Equal(t, newOrderProduct.Quantity, updateOrderProduct.Quantity)
}
