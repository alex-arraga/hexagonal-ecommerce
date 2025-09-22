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
