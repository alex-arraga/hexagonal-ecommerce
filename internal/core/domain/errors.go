package domain

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid email or password")                          // ErrInvalidCredentials is an error for when the credentials are invalid
	ErrDataNotFound       = errors.New("data not found")                                     // ErrDataNotFound is an error for when requested data is not found
	ErrInternal           = errors.New("internal error")                                     // ErrInternal is an error for when an internal service fails to process the request
	ErrConflictingData    = errors.New("data conflicts with existing data in unique column") // ErrConflictingData is an error for when data conflicts with existing data
)
