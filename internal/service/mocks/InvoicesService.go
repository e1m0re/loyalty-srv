// Code generated by mockery v2.43.1. DO NOT EDIT.

package mocks

import (
	context "context"
	models "e1m0re/loyalty-srv/internal/models"

	mock "github.com/stretchr/testify/mock"
)

// InvoicesService is an autogenerated mock type for the InvoicesService type
type InvoicesService struct {
	mock.Mock
}

// CreateInvoice provides a mock function with given fields: ctx, id
func (_m *InvoicesService) CreateInvoice(ctx context.Context, id models.UserID) (*models.Invoice, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for CreateInvoice")
	}

	var r0 *models.Invoice
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.UserID) (*models.Invoice, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.UserID) *models.Invoice); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Invoice)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.UserID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetInvoiceByUserID provides a mock function with given fields: ctx, userID
func (_m *InvoicesService) GetInvoiceByUserID(ctx context.Context, userID models.UserID) (*models.Invoice, error) {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for GetInvoiceByUserID")
	}

	var r0 *models.Invoice
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.UserID) (*models.Invoice, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.UserID) *models.Invoice); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Invoice)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.UserID) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetInvoiceInfo provides a mock function with given fields: ctx, invoice
func (_m *InvoicesService) GetInvoiceInfo(ctx context.Context, invoice *models.Invoice) (*models.InvoiceInfo, error) {
	ret := _m.Called(ctx, invoice)

	if len(ret) == 0 {
		panic("no return value specified for GetInvoiceInfo")
	}

	var r0 *models.InvoiceInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Invoice) (*models.InvoiceInfo, error)); ok {
		return rf(ctx, invoice)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *models.Invoice) *models.InvoiceInfo); ok {
		r0 = rf(ctx, invoice)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.InvoiceInfo)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *models.Invoice) error); ok {
		r1 = rf(ctx, invoice)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetWithdrawals provides a mock function with given fields: ctx, invoice
func (_m *InvoicesService) GetWithdrawals(ctx context.Context, invoice *models.Invoice) (*[]models.Withdrawal, error) {
	ret := _m.Called(ctx, invoice)

	if len(ret) == 0 {
		panic("no return value specified for GetWithdrawals")
	}

	var r0 *[]models.Withdrawal
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Invoice) (*[]models.Withdrawal, error)); ok {
		return rf(ctx, invoice)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *models.Invoice) *[]models.Withdrawal); ok {
		r0 = rf(ctx, invoice)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]models.Withdrawal)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *models.Invoice) error); ok {
		r1 = rf(ctx, invoice)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateBalance provides a mock function with given fields: ctx, invoice, amount, orderNum
func (_m *InvoicesService) UpdateBalance(ctx context.Context, invoice models.Invoice, amount float64, orderNum models.OrderNum) (*models.Invoice, error) {
	ret := _m.Called(ctx, invoice, amount, orderNum)

	if len(ret) == 0 {
		panic("no return value specified for UpdateBalance")
	}

	var r0 *models.Invoice
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.Invoice, float64, models.OrderNum) (*models.Invoice, error)); ok {
		return rf(ctx, invoice, amount, orderNum)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.Invoice, float64, models.OrderNum) *models.Invoice); ok {
		r0 = rf(ctx, invoice, amount, orderNum)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Invoice)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.Invoice, float64, models.OrderNum) error); ok {
		r1 = rf(ctx, invoice, amount, orderNum)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewInvoicesService creates a new instance of InvoicesService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewInvoicesService(t interface {
	mock.TestingT
	Cleanup(func())
}) *InvoicesService {
	mock := &InvoicesService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
