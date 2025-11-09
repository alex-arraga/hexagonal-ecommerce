package services_test

import (
	"context"
	"go-ecommerce/internal/adapters/security"
	"go-ecommerce/internal/adapters/storage/database/postgres/repository"
	"go-ecommerce/internal/core/domain"
	"go-ecommerce/internal/core/ports"
	"go-ecommerce/internal/core/ports/ports_dtos"
	"go-ecommerce/internal/core/services"
	testhelpers "go-ecommerce/internal/test_helpers"
	"go-ecommerce/internal/test_helpers/mocks"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type depToTestingOrderSrv struct {
	userSrv    ports.UserService
	opSrv      ports.OrderProductService
	productSrv ports.ProductService
	categSrv   ports.CategoryService
	cartSrv    ports.CartService
	orderSrv   ports.OrderService
}

func newOrderSrvTest(t *testing.T) *depToTestingOrderSrv {
	db := testhelpers.NewSQLiteTestDB(t)
	tx := db.Begin()
	t.Cleanup(func() { tx.Rollback() })

	redis := mocks.NewMockRedis()
	hasher := &security.Hasher{}

	userRepo := repository.NewUserRepo(tx).(*repository.UserRepo)
	prodRepo := repository.NewProductRepo(tx).(*repository.ProductRepo)
	categRepo := repository.NewCategoryRepo(tx).(*repository.CategoryRepo)
	orderProdRepo := repository.NewOrderProductRepo(tx).(*repository.OrderProductRepo)

	orderProdSrv := services.NewOrderProductService(orderProdRepo)
	orderRepo := repository.NewOrderRepo(orderProdSrv, tx).(*repository.OrderRepo)

	// services
	userSrv := services.NewUserService(userRepo, redis, hasher)
	opSrv := services.NewOrderProductService(orderProdRepo)
	productSrv := services.NewProductService(prodRepo, redis)
	categSrv := services.NewCategoryService(categRepo, redis)
	cartSrv := services.NewCartService(redis, productSrv)
	orderSrv := services.NewOrderService(orderRepo, orderProdSrv, cartSrv, redis)

	srvs := &depToTestingOrderSrv{
		userSrv:    userSrv,
		opSrv:      opSrv,
		productSrv: productSrv,
		categSrv:   categSrv,
		cartSrv:    cartSrv,
		orderSrv:   orderSrv,
	}

	return srvs
}

func Test_OrderServices_Create(t *testing.T) {
	t.Helper()

	srv := newOrderSrvTest(t)
	ctx := context.Background()

	// factory user
	u := testhelpers.NewDomainUser("John", "john@mail.test")
	userInputs := domain.SaveUserInputs{
		Name:     &u.Name,
		Email:    &u.Email,
		Password: &u.Password,
		Role:     &u.Role,
	}

	// save user in db
	newUser, err := srv.userSrv.SaveUser(ctx, userInputs)
	require.NoError(t, err)
	assert.Equal(t, u.Name, newUser.Name)
	assert.Equal(t, u.Email, newUser.Email)

	// factory, create a new category to will use as foreign key of product
	c := testhelpers.NewDomainCategory("Tablets")
	savedCateg, err := srv.categSrv.SaveCategory(ctx, 0, c.Name)
	require.NoError(t, err)
	assert.Equal(t, c.Name, savedCateg.Name)

	// factory, create a new product
	p := testhelpers.NewDomainProduct("Ipad 14 pro", savedCateg.ID)
	inputs := ports_dtos.SaveProductInputs{
		Name:       &p.Name,
		Image:      &p.Image,
		SKU:        &p.SKU,
		Price:      &p.Price,
		Stock:      &p.Stock,
		CategoryID: &savedCateg.ID,
	}

	// save product in db
	newProd, err := srv.productSrv.SaveProduct(ctx, inputs)
	require.NoError(t, err)
	assert.Equal(t, p.Name, newProd.Name)
	assert.Equal(t, p.CategoryID, newProd.CategoryID)

	// add products to cart before create the order
	err = srv.cartSrv.AddItemToCart(ctx, newUser.ID, newProd.ID, 5)
	require.NoError(t, err)

	// factory, create a new order
	o := testhelpers.NewDomainOrder(newUser.ID)
	orderInputs := ports.SaveOrderInputs{
		UserID:            o.UserID,
		Currency:          o.Currency,
		ExternalReference: o.ExternalReference,
		PaymentID:         o.PaymentID,
		PayStatus:         &o.PayStatus,
		PayStatusDetail:   o.PayStatusDetail,
	}

	// save order in db
	newOrder, err := srv.orderSrv.SaveOrder(ctx, orderInputs)
	require.NoError(t, err)
	assert.Equal(t, o.UserID, newOrder.UserID)

	// verify if the order contains the correct items
	orderItems, err := srv.opSrv.GetByOrderID(ctx, newOrder.ID)
	require.NoError(t, err)
	require.Len(t, orderItems, 1)

	item := orderItems[0]
	assert.Equal(t, newProd.ID, item.ProductID)
	assert.Equal(t, int16(5), item.Quantity)
}

func Test_OrderServices_Update(t *testing.T) {
	t.Helper()

	srv := newOrderSrvTest(t)
	ctx := context.Background()

	// factory user
	u := testhelpers.NewDomainUser("John", "john@mail.test")
	userInputs := domain.SaveUserInputs{
		Name:     &u.Name,
		Email:    &u.Email,
		Password: &u.Password,
		Role:     &u.Role,
	}

	// save user in db
	newUser, err := srv.userSrv.SaveUser(ctx, userInputs)
	require.NoError(t, err)
	assert.Equal(t, u.Name, newUser.Name)
	assert.Equal(t, u.Email, newUser.Email)

	// factory, create a new category to will use as foreign key of product
	c := testhelpers.NewDomainCategory("Tablets")
	savedCateg, err := srv.categSrv.SaveCategory(ctx, 0, c.Name)
	require.NoError(t, err)
	assert.Equal(t, c.Name, savedCateg.Name)

	// factory, create a new product
	p := testhelpers.NewDomainProduct("Ipad 14 pro", savedCateg.ID)
	inputs := ports_dtos.SaveProductInputs{
		Name:       &p.Name,
		Image:      &p.Image,
		SKU:        &p.SKU,
		Price:      &p.Price,
		Stock:      &p.Stock,
		CategoryID: &savedCateg.ID,
	}

	// save product in db
	newProd, err := srv.productSrv.SaveProduct(ctx, inputs)
	require.NoError(t, err)
	assert.Equal(t, p.Name, newProd.Name)
	assert.Equal(t, p.CategoryID, newProd.CategoryID)

	// add products to cart before create the order
	err = srv.cartSrv.AddItemToCart(ctx, newUser.ID, newProd.ID, 5)
	require.NoError(t, err)

	// factory, create a new order
	o := testhelpers.NewDomainOrder(newUser.ID)
	orderInputs := ports.SaveOrderInputs{
		UserID:            o.UserID,
		Currency:          o.Currency,
		ExternalReference: o.ExternalReference,
		PaymentID:         o.PaymentID,
		PayStatus:         &o.PayStatus,
		PayStatusDetail:   o.PayStatusDetail,
	}

	// save new order in db
	newOrder, err := srv.orderSrv.SaveOrder(ctx, orderInputs)
	require.NoError(t, err)
	assert.Equal(t, o.UserID, newOrder.UserID)

	// verify if the order contains the correct items
	orderItems, err := srv.opSrv.GetByOrderID(ctx, newOrder.ID)
	require.NoError(t, err)
	require.Len(t, orderItems, 1)

	item := orderItems[0]
	assert.Equal(t, newProd.ID, item.ProductID)
	assert.Equal(t, int16(5), item.Quantity)

	// update order
	payStatus := domain.Approved
	payStatusDetail := domain.Accredited
	extRef := uuid.NewString()
	paymentId := uuid.NewString()

	updateData := ports.SaveOrderInputs{
		ID:                newOrder.ID,
		UserID:            newUser.ID,
		Currency:          domain.ARS,
		PayStatus:         &payStatus,
		PayStatusDetail:   &payStatusDetail,
		ExternalReference: &extRef,
		PaymentID:         &paymentId,
	}

	// add products again because after save, the cart is cleaned
	err = srv.cartSrv.AddItemToCart(ctx, newUser.ID, newProd.ID, 5)
	require.NoError(t, err)

	// save updated order in db
	updatedOrder, err := srv.orderSrv.SaveOrder(ctx, updateData)
	require.NoError(t, err)
	assert.Equal(t, updateData.ID, updatedOrder.ID)
	assert.Equal(t, *updateData.PayStatus, updatedOrder.PayStatus)
	assert.Equal(t, updateData.ExternalReference, updatedOrder.ExternalReference)
	assert.Equal(t, updateData.PaymentID, updatedOrder.PaymentID)
}

func Test_OrderServices_GetByID(t *testing.T) {
	t.Helper()

	srv := newOrderSrvTest(t)
	ctx := context.Background()

	// factory user
	u := testhelpers.NewDomainUser("John", "john@mail.test")
	userInputs := domain.SaveUserInputs{
		Name:     &u.Name,
		Email:    &u.Email,
		Password: &u.Password,
		Role:     &u.Role,
	}

	// save user in db
	newUser, err := srv.userSrv.SaveUser(ctx, userInputs)
	require.NoError(t, err)
	assert.Equal(t, u.Name, newUser.Name)
	assert.Equal(t, u.Email, newUser.Email)

	// factory, create a new category to will use as foreign key of product
	c := testhelpers.NewDomainCategory("Tablets")
	savedCateg, err := srv.categSrv.SaveCategory(ctx, 0, c.Name)
	require.NoError(t, err)
	assert.Equal(t, c.Name, savedCateg.Name)

	// factory, create a new product
	p := testhelpers.NewDomainProduct("Ipad 14 pro", savedCateg.ID)
	inputs := ports_dtos.SaveProductInputs{
		Name:       &p.Name,
		Image:      &p.Image,
		SKU:        &p.SKU,
		Price:      &p.Price,
		Stock:      &p.Stock,
		CategoryID: &savedCateg.ID,
	}

	// save product in db
	newProd, err := srv.productSrv.SaveProduct(ctx, inputs)
	require.NoError(t, err)
	assert.Equal(t, p.Name, newProd.Name)
	assert.Equal(t, p.CategoryID, newProd.CategoryID)

	// add products to cart before create the order
	err = srv.cartSrv.AddItemToCart(ctx, newUser.ID, newProd.ID, 5)
	require.NoError(t, err)

	// factory, create a new order
	o := testhelpers.NewDomainOrder(newUser.ID)
	orderInputs := ports.SaveOrderInputs{
		UserID:            o.UserID,
		Currency:          o.Currency,
		ExternalReference: o.ExternalReference,
		PaymentID:         o.PaymentID,
		PayStatus:         &o.PayStatus,
		PayStatusDetail:   o.PayStatusDetail,
	}

	// save order in db
	newOrder, err := srv.orderSrv.SaveOrder(ctx, orderInputs)
	require.NoError(t, err)
	assert.Equal(t, o.UserID, newOrder.UserID)

	// verify if the order contains the correct items
	orderItems, err := srv.opSrv.GetByOrderID(ctx, newOrder.ID)
	require.NoError(t, err)
	require.Len(t, orderItems, 1)

	item := orderItems[0]
	assert.Equal(t, newProd.ID, item.ProductID)
	assert.Equal(t, int16(5), item.Quantity)

	// find order by id
	order, err := srv.orderSrv.GetOrderById(ctx, newOrder.ID)
	require.NoError(t, err)
	assert.Equal(t, order.ID, newOrder.ID)
	assert.Equal(t, order.Total, newOrder.Total)
	assert.Equal(t, order.UserID, newOrder.UserID)
}
