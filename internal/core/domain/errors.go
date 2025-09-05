package domain

import "errors"

var (
	ErrInternal = errors.New("internal error") // ErrInternal is an error for when an internal service fails to process the request
	ErrConflictingData = errors.New("data conflicts with existing data in unique column") // ErrConflictingData is an error for when data conflicts with existing data
)
