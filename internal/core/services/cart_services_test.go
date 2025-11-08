package services_test

import (
	"context"
	"go-ecommerce/internal/adapters/security"
	"go-ecommerce/internal/adapters/storage/database/postgres/repository"
	"go-ecommerce/internal/core/domain"
	"go-ecommerce/internal/core/ports/ports_dtos"
	"go-ecommerce/internal/core/services"
	testhelpers "go-ecommerce/internal/test_helpers"
	"go-ecommerce/internal/test_helpers/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Cart_AddItem(t *testing.T) {
	t.Helper()

	ctx := context.Background()
	db := testhelpers.NewSQLiteTestDB(t)
	tx := db.Begin()
	t.Cleanup(func() { tx.Rollback() })

	redis := mocks.NewMockRedis()
	hasher := &security.Hasher{}

	// repos
	userRepo := repository.NewUserRepo(tx)
	productRepo := repository.NewProductRepo(tx)
	categRepo := repository.NewCategoryRepo(tx)

	// services
	userSrv := services.NewUserService(userRepo, redis, hasher)
	productSrv := services.NewProductService(productRepo, redis)
	categSrv := services.NewCategoryService(categRepo, redis)
	cartSrv := services.NewCartService(redis, productSrv)

	// factory user
	u := testhelpers.NewDomainUser("John", "john@mail.test")
	userInputs := domain.SaveUserInputs{
		Name:     &u.Name,
		Email:    &u.Email,
		Password: &u.Password,
		Role:     &u.Role,
	}

	// save user in db
	newUser, err := userSrv.SaveUser(ctx, userInputs)
	require.NoError(t, err)
	assert.Equal(t, u.Name, newUser.Name)
	assert.Equal(t, u.Email, newUser.Email)

	// factory, create a new category to will use as foreign key of product
	c := testhelpers.NewDomainCategory("Tablets")
	savedCateg, err := categSrv.SaveCategory(ctx, 0, c.Name)
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
	newProd, err := productSrv.SaveProduct(ctx, inputs)
	require.NoError(t, err)
	assert.Equal(t, p.Name, newProd.Name)
	assert.Equal(t, p.CategoryID, newProd.CategoryID)

	// add products to cart before create the order
	err = cartSrv.AddItemToCart(ctx, newUser.ID, newProd.ID, 5)
	require.NoError(t, err)
}

func Test_Cart_GetCart(t *testing.T) {
	t.Helper()

	ctx := context.Background()
	db := testhelpers.NewSQLiteTestDB(t)
	tx := db.Begin()
	t.Cleanup(func() { tx.Rollback() })

	redis := mocks.NewMockRedis()
	hasher := &security.Hasher{}

	// repos
	userRepo := repository.NewUserRepo(tx)
	productRepo := repository.NewProductRepo(tx)
	categRepo := repository.NewCategoryRepo(tx)

	// services
	userSrv := services.NewUserService(userRepo, redis, hasher)
	productSrv := services.NewProductService(productRepo, redis)
	categSrv := services.NewCategoryService(categRepo, redis)
	cartSrv := services.NewCartService(redis, productSrv)

	// factory user
	u := testhelpers.NewDomainUser("John", "john@mail.test")
	userInputs := domain.SaveUserInputs{
		Name:     &u.Name,
		Email:    &u.Email,
		Password: &u.Password,
		Role:     &u.Role,
	}

	// save user in db
	newUser, err := userSrv.SaveUser(ctx, userInputs)
	require.NoError(t, err)
	assert.Equal(t, u.Name, newUser.Name)
	assert.Equal(t, u.Email, newUser.Email)

	// factory, create a new category to will use as foreign key of product
	c := testhelpers.NewDomainCategory("Tablets")
	savedCateg, err := categSrv.SaveCategory(ctx, 0, c.Name)
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
	newProd, err := productSrv.SaveProduct(ctx, inputs)
	require.NoError(t, err)
	assert.Equal(t, p.Name, newProd.Name)
	assert.Equal(t, p.CategoryID, newProd.CategoryID)

	// add products to cart before create the order
	err = cartSrv.AddItemToCart(ctx, newUser.ID, newProd.ID, 5)
	require.NoError(t, err)

	// verify items in cart
	cart, err := cartSrv.GetCart(ctx, newUser.ID)
	require.NoError(t, err)
	assert.Equal(t, newProd.ID, cart.Items[0].ProductID)
	assert.Equal(t, int16(5), cart.Items[0].Quantity)
}

func Test_Cart_DeleteProductByID(t *testing.T) {
	t.Helper()

	ctx := context.Background()
	db := testhelpers.NewSQLiteTestDB(t)
	tx := db.Begin()
	t.Cleanup(func() { tx.Rollback() })

	redis := mocks.NewMockRedis()
	hasher := &security.Hasher{}

	// repos
	userRepo := repository.NewUserRepo(tx)
	productRepo := repository.NewProductRepo(tx)
	categRepo := repository.NewCategoryRepo(tx)

	// services
	userSrv := services.NewUserService(userRepo, redis, hasher)
	productSrv := services.NewProductService(productRepo, redis)
	categSrv := services.NewCategoryService(categRepo, redis)
	cartSrv := services.NewCartService(redis, productSrv)

	// factory user
	u := testhelpers.NewDomainUser("John", "john@mail.test")
	userInputs := domain.SaveUserInputs{
		Name:     &u.Name,
		Email:    &u.Email,
		Password: &u.Password,
		Role:     &u.Role,
	}

	// save user in db
	newUser, err := userSrv.SaveUser(ctx, userInputs)
	require.NoError(t, err)
	assert.Equal(t, u.Name, newUser.Name)
	assert.Equal(t, u.Email, newUser.Email)

	// factory, create a new category to will use as foreign key of product
	c := testhelpers.NewDomainCategory("Tablets")
	savedCateg, err := categSrv.SaveCategory(ctx, 0, c.Name)
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
	newProd, err := productSrv.SaveProduct(ctx, inputs)
	require.NoError(t, err)
	assert.Equal(t, p.Name, newProd.Name)
	assert.Equal(t, p.CategoryID, newProd.CategoryID)

	// add products to cart before create the order
	err = cartSrv.AddItemToCart(ctx, newUser.ID, newProd.ID, 5)
	require.NoError(t, err)

	// verify items in cart
	cart, err := cartSrv.GetCart(ctx, newUser.ID)
	require.NoError(t, err)
	assert.Equal(t, newProd.ID, cart.Items[0].ProductID)

	// delete product in cart and verify
	err = cartSrv.RemoveItem(ctx, newUser.ID, newProd.ID)
	require.NoError(t, err)

	// reload the cart
	cart, err = cartSrv.GetCart(ctx, newUser.ID)
	require.NoError(t, err)

	// ensure the cart is empty
	assert.Empty(t, cart.Items)
}

func Test_Cart_Clear(t *testing.T) {
	t.Helper()

	ctx := context.Background()
	db := testhelpers.NewSQLiteTestDB(t)
	tx := db.Begin()
	t.Cleanup(func() { tx.Rollback() })

	redis := mocks.NewMockRedis()
	hasher := &security.Hasher{}

	// repos
	userRepo := repository.NewUserRepo(tx)
	productRepo := repository.NewProductRepo(tx)
	categRepo := repository.NewCategoryRepo(tx)

	// services
	userSrv := services.NewUserService(userRepo, redis, hasher)
	productSrv := services.NewProductService(productRepo, redis)
	categSrv := services.NewCategoryService(categRepo, redis)
	cartSrv := services.NewCartService(redis, productSrv)

	// factory user
	u := testhelpers.NewDomainUser("John", "john@mail.test")
	userInputs := domain.SaveUserInputs{
		Name:     &u.Name,
		Email:    &u.Email,
		Password: &u.Password,
		Role:     &u.Role,
	}

	// save user in db
	newUser, err := userSrv.SaveUser(ctx, userInputs)
	require.NoError(t, err)
	assert.Equal(t, u.Name, newUser.Name)
	assert.Equal(t, u.Email, newUser.Email)

	// factory, create a new category to will use as foreign key of product
	c := testhelpers.NewDomainCategory("Tablets")
	savedCateg, err := categSrv.SaveCategory(ctx, 0, c.Name)
	require.NoError(t, err)
	assert.Equal(t, c.Name, savedCateg.Name)

	// factory, create a new product
	p1 := testhelpers.NewDomainProduct("Ipad 14 pro", savedCateg.ID)
	inputsP1 := ports_dtos.SaveProductInputs{
		Name:       &p1.Name,
		Image:      &p1.Image,
		SKU:        &p1.SKU,
		Price:      &p1.Price,
		Stock:      &p1.Stock,
		CategoryID: &savedCateg.ID,
	}

	p2 := testhelpers.NewDomainProduct("Ipad 16 pro", savedCateg.ID)
	inputsP2 := ports_dtos.SaveProductInputs{
		Name:       &p2.Name,
		Image:      &p2.Image,
		SKU:        &p2.SKU,
		Price:      &p2.Price,
		Stock:      &p2.Stock,
		CategoryID: &savedCateg.ID,
	}

	// save 2 products in db
	newProd1, err := productSrv.SaveProduct(ctx, inputsP1)
	require.NoError(t, err)
	assert.Equal(t, p1.Name, newProd1.Name)
	assert.Equal(t, p1.CategoryID, newProd1.CategoryID)

	// save 2 products in db
	newProd2, err := productSrv.SaveProduct(ctx, inputsP2)
	require.NoError(t, err)
	assert.Equal(t, p2.Name, newProd2.Name)
	assert.Equal(t, p2.CategoryID, newProd2.CategoryID)

	// add product-1 to cart
	err = cartSrv.AddItemToCart(ctx, newUser.ID, newProd1.ID, 5)
	require.NoError(t, err)

	// add product-2 to cart
	err = cartSrv.AddItemToCart(ctx, newUser.ID, newProd2.ID, 3)
	require.NoError(t, err)

	// verify items in cart
	cart, err := cartSrv.GetCart(ctx, newUser.ID)
	require.NoError(t, err)
	assert.Equal(t, newProd1.ID, cart.Items[0].ProductID)
	assert.Equal(t, newProd2.ID, cart.Items[1].ProductID)

	// delete product in cart and verify
	err = cartSrv.Clear(ctx, newUser.ID)
	require.NoError(t, err)

	// reload the cart
	cart, err = cartSrv.GetCart(ctx, newUser.ID)
	require.NoError(t, err)

	// ensure the cart is empty
	assert.Empty(t, cart.Items)
}
