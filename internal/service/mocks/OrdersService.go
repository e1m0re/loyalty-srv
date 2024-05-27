// Code generated by mockery v2.43.1. DO NOT EDIT.

package mocks

import (
	context "context"
	models "e1m0re/loyalty-srv/internal/models"

	mock "github.com/stretchr/testify/mock"
)

// OrdersService is an autogenerated mock type for the OrdersService type
type OrdersService struct {
	mock.Mock
}

// GetLoadedOrdersByUserID provides a mock function with given fields: ctx, userID
func (_m *OrdersService) GetLoadedOrdersByUserID(ctx context.Context, userID models.UserID) (*[]models.Order, error) {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for GetLoadedOrdersByUserID")
	}

	var r0 *[]models.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.UserID) (*[]models.Order, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.UserID) *[]models.Order); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]models.Order)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.UserID) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewOrder provides a mock function with given fields: ctx, orderInfo
func (_m *OrdersService) NewOrder(ctx context.Context, orderInfo models.OrderInfo) (*models.Order, error) {
	ret := _m.Called(ctx, orderInfo)

	if len(ret) == 0 {
		panic("no return value specified for NewOrder")
	}

	var r0 *models.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.OrderInfo) (*models.Order, error)); ok {
		return rf(ctx, orderInfo)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.OrderInfo) *models.Order); ok {
		r0 = rf(ctx, orderInfo)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Order)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.OrderInfo) error); ok {
		r1 = rf(ctx, orderInfo)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateOrder provides a mock function with given fields: ctx, id, status, accrual
func (_m *OrdersService) UpdateOrder(ctx context.Context, id models.OrderID, status models.OrdersStatus, accrual int) (models.Order, error) {
	ret := _m.Called(ctx, id, status, accrual)

	if len(ret) == 0 {
		panic("no return value specified for UpdateOrder")
	}

	var r0 models.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.OrderID, models.OrdersStatus, int) (models.Order, error)); ok {
		return rf(ctx, id, status, accrual)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.OrderID, models.OrdersStatus, int) models.Order); ok {
		r0 = rf(ctx, id, status, accrual)
	} else {
		r0 = ret.Get(0).(models.Order)
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.OrderID, models.OrdersStatus, int) error); ok {
		r1 = rf(ctx, id, status, accrual)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ValidateNumber provides a mock function with given fields: ctx, orderNum
func (_m *OrdersService) ValidateNumber(ctx context.Context, orderNum models.OrderNum) (bool, error) {
	ret := _m.Called(ctx, orderNum)

	if len(ret) == 0 {
		panic("no return value specified for ValidateNumber")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.OrderNum) (bool, error)); ok {
		return rf(ctx, orderNum)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.OrderNum) bool); ok {
		r0 = rf(ctx, orderNum)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.OrderNum) error); ok {
		r1 = rf(ctx, orderNum)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewOrdersService creates a new instance of OrdersService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewOrdersService(t interface {
	mock.TestingT
	Cleanup(func())
}) *OrdersService {
	mock := &OrdersService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
