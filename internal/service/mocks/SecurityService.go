// Code generated by mockery v2.43.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// SecurityService is an autogenerated mock type for the SecurityService type
type SecurityService struct {
	mock.Mock
}

// CheckPassword provides a mock function with given fields: hashPassword, password
func (_m *SecurityService) CheckPassword(hashPassword string, password string) bool {
	ret := _m.Called(hashPassword, password)

	if len(ret) == 0 {
		panic("no return value specified for CheckPassword")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(string, string) bool); ok {
		r0 = rf(hashPassword, password)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// GetPasswordHash provides a mock function with given fields: password
func (_m *SecurityService) GetPasswordHash(password string) (string, error) {
	ret := _m.Called(password)

	if len(ret) == 0 {
		panic("no return value specified for GetPasswordHash")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(password)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(password)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewSecurityService creates a new instance of SecurityService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSecurityService(t interface {
	mock.TestingT
	Cleanup(func())
}) *SecurityService {
	mock := &SecurityService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}