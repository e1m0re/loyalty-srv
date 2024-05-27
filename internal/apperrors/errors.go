package apperrors

import (
	"errors"
)

var (
	ErrBusyLogin          = errors.New("busy login")
	ErrInvalidOrderNumber = errors.New("invalid order number")
	ErrEmptyOrderNumber   = errors.New("empty order number")
	ErrOtherUsersOrder    = errors.New("the order number has already been uploaded by another user")
	ErrEntityNotFound     = errors.New("entity not found")
)
