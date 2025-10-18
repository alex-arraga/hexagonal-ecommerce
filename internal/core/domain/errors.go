package domain

import "errors"

var (
	// User errors
	ErrEmailExist        = errors.New("an account with this email already exist, try loggin")
	ErrEmailIsRequire    = errors.New("email is required")
	ErrCreatingUser      = errors.New("couldn't create the user")
	ErrNameIsRequire     = errors.New("name is required")
	ErrPasswordIsRequire = errors.New("password is required")
	ErrMinLenghtName     = errors.New("name of user must have at least 3 characters")
	ErrMinLenghtPassword = errors.New("password must have at least 6 characters")
	ErrRoleIsRequire     = errors.New("role is required")
	ErrRoleIsInvalid     = errors.New("invalid user role")
	ErrHashingPassword   = errors.New("error hashing password")
	ErrUserNotFound      = errors.New("user not found")
	ErrUsersNotFound     = errors.New("list of users not found")
)

var (
	// Category errors
	ErrCategoryNameIsRequire = errors.New("name of category is required")
	ErrCategoryNotFound      = errors.New("category not found")
	ErrCategoriesNotFound    = errors.New("list of categories not found")
)

var (
	// Product errors
	ErrProductStockIsRequire    = errors.New("stock of product is required")
	ErrProductPriceIsRequire    = errors.New("price of product is required")
	ErrProductImageIsRequire    = errors.New("image of product is required")
	ErrProductCategoryIsRequire = errors.New("category of product is required")
	ErrProductSKUIsRequire      = errors.New("sku of product is required")
	ErrProductNameIsRequire     = errors.New("name of category is required")
	ErrProductMinLenghtName     = errors.New("name of product must have at least 3 characters")
	ErrProductMinLenghtSKU      = errors.New("sku of product must have at least 3 characters")
	ErrProductNotFound          = errors.New("product not found")
	ErrProductsNotFound         = errors.New("list of products not found")
)

// Order-Product errors
var (
	ErrOrderProductNotFound  = errors.New("order-product not found")
	ErrOrdersProductNotFound = errors.New("list of orders-product not found")
)

// Cart errors
var (
	ErrAlreadyEmptyCart                    = errors.New("the cart is already empty")
	ErrProductNotFoundCart                 = errors.New("product not found in cart")
	ErrNegativeQuantityNonExistProductCart = errors.New("product not exist in cart, quantity must be a positive number")
)

var (
	ErrOrderNotFound = errors.New("order not found")
	ErrOrdersNotFound = errors.New("list of orders not found")
)
