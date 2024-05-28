package apperrors

import (
	"errors"
)

var (
	ErrEntityNotFound = errors.New("entity not found")

	ErrBusyLogin = errors.New("busy login")

	ErrInvalidOrderNumber          = errors.New("invalid order number")
	ErrEmptyOrderNumber            = errors.New("empty order number")
	ErrOrderWasLoaded              = errors.New("order with this number was loaded")
	ErrOrderWasLoadedByAnotherUser = errors.New("order with this number was loaded by another user")

	ErrAccountHasNotEnoughFunds = errors.New("there are insufficient funds in the account")
)
