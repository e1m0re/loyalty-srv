package apperrors

import (
	"errors"
)

var (
	ErrEntityNotFound = errors.New("entity not found")

	ErrBusyLogin = errors.New("busy login")

	ErrInvalidOrderNumber = errors.New("invalid order number")
	ErrEmptyOrderNumber   = errors.New("empty order number")
	ErrOrderIsLoaded      = errors.New("order with this number loaded yet")
	ErrOtherUsersOrder    = errors.New("the order number has already been uploaded by another user")
)
