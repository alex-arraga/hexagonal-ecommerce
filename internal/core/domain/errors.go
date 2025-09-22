package domain

import "errors"

var (
	ErrMinLenghtName     = errors.New("name of user must have at least 3 characters")
	ErrMinLenghtPassword = errors.New("password must have at least 6 characters")
	ErrRoleIsRequire     = errors.New("role is required")
	ErrRoleIsInvalid     = errors.New("invalid user role")
	ErrHashingPassword   = errors.New("error hashing password")
)
