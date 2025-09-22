package domain

import "errors"

var (
	// User errors
	ErrEmailIsRequire    = errors.New("email is required")
	ErrNameIsRequire     = errors.New("name is required")
	ErrPasswordIsRequire = errors.New("password is required")
	ErrMinLenghtName     = errors.New("name of user must have at least 3 characters")
	ErrMinLenghtPassword = errors.New("password must have at least 6 characters")
	ErrRoleIsRequire     = errors.New("role is required")
	ErrRoleIsInvalid     = errors.New("invalid user role")
	ErrHashingPassword   = errors.New("error hashing password")
)

var (
	// Category errors
	ErrCategoryNameIsRequire = errors.New("name of category is required")
	ErrCategoryNotFound      = errors.New("category not found")
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
	ErrProductNotFound          = errors.New("category not found")
)
