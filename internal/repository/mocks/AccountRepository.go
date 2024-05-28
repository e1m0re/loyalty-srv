// Code generated by mockery v2.43.1. DO NOT EDIT.

package mocks

import (
	context "context"
	models "e1m0re/loyalty-srv/internal/models"

	mock "github.com/stretchr/testify/mock"
)

// AccountRepository is an autogenerated mock type for the AccountRepository type
type AccountRepository struct {
	mock.Mock
}

// AddAccount provides a mock function with given fields: ctx, userID
func (_m *AccountRepository) AddAccount(ctx context.Context, userID models.UserID) (*models.Account, error) {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for AddAccount")
	}

	var r0 *models.Account
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.UserID) (*models.Account, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.UserID) *models.Account); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Account)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.UserID) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewAccountRepository creates a new instance of AccountRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAccountRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *AccountRepository {
	mock := &AccountRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}