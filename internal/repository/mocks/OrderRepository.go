// Code generated by mockery v2.43.1. DO NOT EDIT.

package mocks

import (
	context "context"
	models "e1m0re/loyalty-srv/internal/models"

	mock "github.com/stretchr/testify/mock"
)

// OrderRepository is an autogenerated mock type for the OrderRepository type
type OrderRepository struct {
	mock.Mock
}

// AddOrder provides a mock function with given fields: ctx, orderInfo
func (_m *OrderRepository) AddOrder(ctx context.Context, orderInfo models.OrderInfo) (*models.Order, error) {
	ret := _m.Called(ctx, orderInfo)

	if len(ret) == 0 {
		panic("no return value specified for AddOrder")
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

// GetLoadedOrdersByUserID provides a mock function with given fields: ctx, userID
func (_m *OrderRepository) GetLoadedOrdersByUserID(ctx context.Context, userID models.UserID) (*[]models.Order, error) {
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

// GetNotCalculatedOrder provides a mock function with given fields: ctx, limit
func (_m *OrderRepository) GetNotCalculatedOrder(ctx context.Context, limit int) (*[]models.Order, error) {
	ret := _m.Called(ctx, limit)

	if len(ret) == 0 {
		panic("no return value specified for GetNotCalculatedOrder")
	}

	var r0 *[]models.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) (*[]models.Order, error)); ok {
		return rf(ctx, limit)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) *[]models.Order); ok {
		r0 = rf(ctx, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]models.Order)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOrderByNumber provides a mock function with given fields: ctx, num
func (_m *OrderRepository) GetOrderByNumber(ctx context.Context, num models.OrderNum) (*models.Order, error) {
	ret := _m.Called(ctx, num)

	if len(ret) == 0 {
		panic("no return value specified for GetOrderByNumber")
	}

	var r0 *models.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.OrderNum) (*models.Order, error)); ok {
		return rf(ctx, num)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.OrderNum) *models.Order); ok {
		r0 = rf(ctx, num)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Order)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.OrderNum) error); ok {
		r1 = rf(ctx, num)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateOrdersCalculated provides a mock function with given fields: ctx, order, calculated
func (_m *OrderRepository) UpdateOrdersCalculated(ctx context.Context, order models.Order, calculated bool) (*models.Order, error) {
	ret := _m.Called(ctx, order, calculated)

	if len(ret) == 0 {
		panic("no return value specified for UpdateOrdersCalculated")
	}

	var r0 *models.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.Order, bool) (*models.Order, error)); ok {
		return rf(ctx, order, calculated)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.Order, bool) *models.Order); ok {
		r0 = rf(ctx, order, calculated)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Order)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.Order, bool) error); ok {
		r1 = rf(ctx, order, calculated)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateOrdersStatus provides a mock function with given fields: ctx, order, status, accrual
func (_m *OrderRepository) UpdateOrdersStatus(ctx context.Context, order models.Order, status models.OrdersStatus, accrual int) (*models.Order, error) {
	ret := _m.Called(ctx, order, status, accrual)

	if len(ret) == 0 {
		panic("no return value specified for UpdateOrdersStatus")
	}

	var r0 *models.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.Order, models.OrdersStatus, int) (*models.Order, error)); ok {
		return rf(ctx, order, status, accrual)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.Order, models.OrdersStatus, int) *models.Order); ok {
		r0 = rf(ctx, order, status, accrual)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Order)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.Order, models.OrdersStatus, int) error); ok {
		r1 = rf(ctx, order, status, accrual)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewOrderRepository creates a new instance of OrderRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewOrderRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *OrderRepository {
	mock := &OrderRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
