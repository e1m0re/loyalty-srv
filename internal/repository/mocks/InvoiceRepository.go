// Code generated by mockery v2.43.1. DO NOT EDIT.

package mocks

import (
	context "context"
	models "e1m0re/loyalty-srv/internal/models"

	mock "github.com/stretchr/testify/mock"
)

// InvoiceRepository is an autogenerated mock type for the InvoiceRepository type
type InvoiceRepository struct {
	mock.Mock
}

// AddInvoice provides a mock function with given fields: ctx, userID
func (_m *InvoiceRepository) AddInvoice(ctx context.Context, userID models.UserID) (*models.Invoice, error) {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for AddInvoice")
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

// AddInvoiceChange provides a mock function with given fields: ctx, invoiceID, amount, orderNum
func (_m *InvoiceRepository) AddInvoiceChange(ctx context.Context, invoiceID models.InvoiceID, amount float64, orderNum models.OrderNum) (*models.InvoiceChanges, error) {
	ret := _m.Called(ctx, invoiceID, amount, orderNum)

	if len(ret) == 0 {
		panic("no return value specified for AddInvoiceChange")
	}

	var r0 *models.InvoiceChanges
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.InvoiceID, float64, models.OrderNum) (*models.InvoiceChanges, error)); ok {
		return rf(ctx, invoiceID, amount, orderNum)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.InvoiceID, float64, models.OrderNum) *models.InvoiceChanges); ok {
		r0 = rf(ctx, invoiceID, amount, orderNum)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.InvoiceChanges)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.InvoiceID, float64, models.OrderNum) error); ok {
		r1 = rf(ctx, invoiceID, amount, orderNum)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetInvoiceByID provides a mock function with given fields: ctx, invoiceID
func (_m *InvoiceRepository) GetInvoiceByID(ctx context.Context, invoiceID models.InvoiceID) (*models.Invoice, error) {
	ret := _m.Called(ctx, invoiceID)

	if len(ret) == 0 {
		panic("no return value specified for GetInvoiceByID")
	}

	var r0 *models.Invoice
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.InvoiceID) (*models.Invoice, error)); ok {
		return rf(ctx, invoiceID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.InvoiceID) *models.Invoice); ok {
		r0 = rf(ctx, invoiceID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Invoice)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.InvoiceID) error); ok {
		r1 = rf(ctx, invoiceID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetInvoiceByUserID provides a mock function with given fields: ctx, userID
func (_m *InvoiceRepository) GetInvoiceByUserID(ctx context.Context, userID models.UserID) (*models.Invoice, error) {
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

// GetWithdrawalsList provides a mock function with given fields: ctx, invoiceID
func (_m *InvoiceRepository) GetWithdrawalsList(ctx context.Context, invoiceID models.InvoiceID) (*[]models.Withdrawal, error) {
	ret := _m.Called(ctx, invoiceID)

	if len(ret) == 0 {
		panic("no return value specified for GetWithdrawalsList")
	}

	var r0 *[]models.Withdrawal
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.InvoiceID) (*[]models.Withdrawal, error)); ok {
		return rf(ctx, invoiceID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.InvoiceID) *[]models.Withdrawal); ok {
		r0 = rf(ctx, invoiceID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]models.Withdrawal)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.InvoiceID) error); ok {
		r1 = rf(ctx, invoiceID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetWithdrawnTotalSum provides a mock function with given fields: ctx, invoiceID
func (_m *InvoiceRepository) GetWithdrawnTotalSum(ctx context.Context, invoiceID models.InvoiceID) (int, error) {
	ret := _m.Called(ctx, invoiceID)

	if len(ret) == 0 {
		panic("no return value specified for GetWithdrawnTotalSum")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.InvoiceID) (int, error)); ok {
		return rf(ctx, invoiceID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.InvoiceID) int); ok {
		r0 = rf(ctx, invoiceID)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.InvoiceID) error); ok {
		r1 = rf(ctx, invoiceID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateBalance provides a mock function with given fields: ctx, invoice, amount
func (_m *InvoiceRepository) UpdateBalance(ctx context.Context, invoice models.Invoice, amount float64) (*models.Invoice, error) {
	ret := _m.Called(ctx, invoice, amount)

	if len(ret) == 0 {
		panic("no return value specified for UpdateBalance")
	}

	var r0 *models.Invoice
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.Invoice, float64) (*models.Invoice, error)); ok {
		return rf(ctx, invoice, amount)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.Invoice, float64) *models.Invoice); ok {
		r0 = rf(ctx, invoice, amount)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Invoice)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.Invoice, float64) error); ok {
		r1 = rf(ctx, invoice, amount)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewInvoiceRepository creates a new instance of InvoiceRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewInvoiceRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *InvoiceRepository {
	mock := &InvoiceRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}