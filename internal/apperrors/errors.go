package apperrors

import (
	"errors"
	"fmt"
)

func NewNotImplementedError(method string) error {
	return fmt.Errorf("method %s not implemented", method)
}

func NewInvalidRequestFormat() error {
	return fmt.Errorf("invalid reguest")
}

var (
	BusyLoginError          = errors.New("busy login")
	InvalidOrderNumberError = errors.New("invalid order number")
	EmptyOrderNumberError   = errors.New("empty order number")
	OtherUsersOrderError    = errors.New("the order number has already been uploaded by another user")
	EntityNotFoundError     = errors.New("entity not found")
)
