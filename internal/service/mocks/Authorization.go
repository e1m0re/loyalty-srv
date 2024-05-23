// Code generated by mockery v2.43.1. DO NOT EDIT.

package mocks

import (
	context "context"
	models "e1m0re/loyalty-srv/internal/models"

	mock "github.com/stretchr/testify/mock"
)

// Authorization is an autogenerated mock type for the Authorization type
type Authorization struct {
	mock.Mock
}

// CreateUser provides a mock function with given fields: ctx, userInfo
func (_m *Authorization) CreateUser(ctx context.Context, userInfo *models.UserInfo) (*models.User, error) {
	ret := _m.Called(ctx, userInfo)

	if len(ret) == 0 {
		panic("no return value specified for CreateUser")
	}

	var r0 *models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.UserInfo) (*models.User, error)); ok {
		return rf(ctx, userInfo)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *models.UserInfo) *models.User); ok {
		r0 = rf(ctx, userInfo)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *models.UserInfo) error); ok {
		r1 = rf(ctx, userInfo)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindUserByUsername provides a mock function with given fields: ctx, username
func (_m *Authorization) FindUserByUsername(ctx context.Context, username string) (*models.User, error) {
	ret := _m.Called(ctx, username)

	if len(ret) == 0 {
		panic("no return value specified for FindUserByUsername")
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

// SignIn provides a mock function with given fields: ctx, userInfo
func (_m *Authorization) SignIn(ctx context.Context, userInfo *models.UserInfo) (bool, error) {
	ret := _m.Called(ctx, userInfo)

	if len(ret) == 0 {
		panic("no return value specified for SignIn")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.UserInfo) (bool, error)); ok {
		return rf(ctx, userInfo)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *models.UserInfo) bool); ok {
		r0 = rf(ctx, userInfo)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *models.UserInfo) error); ok {
		r1 = rf(ctx, userInfo)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Verify provides a mock function with given fields: ctx, userInfo
func (_m *Authorization) Verify(ctx context.Context, userInfo *models.UserInfo) (bool, error) {
	ret := _m.Called(ctx, userInfo)

	if len(ret) == 0 {
		panic("no return value specified for Verify")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.UserInfo) (bool, error)); ok {
		return rf(ctx, userInfo)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *models.UserInfo) bool); ok {
		r0 = rf(ctx, userInfo)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *models.UserInfo) error); ok {
		r1 = rf(ctx, userInfo)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewAuthorization creates a new instance of Authorization. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAuthorization(t interface {
	mock.TestingT
	Cleanup(func())
}) *Authorization {
	mock := &Authorization{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
