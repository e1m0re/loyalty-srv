// Code generated by mockery v2.43.1. DO NOT EDIT.

package mocks

import (
	context "context"
	models "e1m0re/loyalty-srv/internal/models"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// UserRepository is an autogenerated mock type for the UserRepository type
type UserRepository struct {
	mock.Mock
}

// CreateUser provides a mock function with given fields: ctx, userInfo
func (_m *UserRepository) CreateUser(ctx context.Context, userInfo models.UserInfo) (*models.User, error) {
	ret := _m.Called(ctx, userInfo)

	if len(ret) == 0 {
		panic("no return value specified for CreateUser")
	}

	var r0 *models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.UserInfo) (*models.User, error)); ok {
		return rf(ctx, userInfo)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.UserInfo) *models.User); ok {
		r0 = rf(ctx, userInfo)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.UserInfo) error); ok {
		r1 = rf(ctx, userInfo)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserByUsername provides a mock function with given fields: ctx, username
func (_m *UserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	ret := _m.Called(ctx, username)

	if len(ret) == 0 {
		panic("no return value specified for GetUserByUsername")
	}

	var r0 *models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*models.User, error)); ok {
		return rf(ctx, username)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *models.User); ok {
		r0 = rf(ctx, username)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateUsersLastLogin provides a mock function with given fields: ctx, id, t
func (_m *UserRepository) UpdateUsersLastLogin(ctx context.Context, id models.UserId, t time.Time) error {
	ret := _m.Called(ctx, id, t)

	if len(ret) == 0 {
		panic("no return value specified for UpdateUsersLastLogin")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, models.UserId, time.Time) error); ok {
		r0 = rf(ctx, id, t)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewUserRepository creates a new instance of UserRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserRepository {
	mock := &UserRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
