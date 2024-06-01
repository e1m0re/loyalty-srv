package service

import (
	"context"
	"e1m0re/loyalty-srv/internal/models"
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"e1m0re/loyalty-srv/internal/service/mocks"
)

func Test_orderProcessor_RecalculateProcessedOrders(t *testing.T) {
	zeroAccrual := float64(0)
	noZeroAccrual := float64(100)
	type args struct {
		ctx context.Context
	}
	type want struct {
		errMsg string
	}
	tests := []struct {
		name        string
		mockProcess func() *orderProcessor
		args        args
		want        want
	}{
		{
			name: "GetNotCalculatedOrder failed",
			args: args{ctx: context.Background()},
			mockProcess: func() *orderProcessor {
				mockOrderService := mocks.NewOrdersService(t)
				mockOrderService.
					On("GetNotCalculatedOrder", mock.Anything).
					Return(nil, fmt.Errorf("some repos error"))

				return &orderProcessor{
					ordersService: mockOrderService,
				}
			},
			want: want{errMsg: "some repos error"},
		},
		{
			name: "All order calculated",
			args: args{ctx: context.Background()},
			mockProcess: func() *orderProcessor {
				mockOrderService := mocks.NewOrdersService(t)
				mockOrderService.
					On("GetNotCalculatedOrder", mock.Anything).
					Return(nil, nil)

				return &orderProcessor{
					ordersService: mockOrderService,
				}
			},
			want: want{errMsg: ""},
		},
		{
			name: "Order is not calculated and has not accrual",
			args: args{ctx: context.Background()},
			mockProcess: func() *orderProcessor {
				mockOrderService := mocks.NewOrdersService(t)
				mockOrderService.
					On("GetNotCalculatedOrder", mock.Anything).
					Return(&models.Order{Accrual: nil}, nil)

				return &orderProcessor{
					ordersService: mockOrderService,
				}
			},
			want: want{errMsg: ""},
		},
		{
			name: "Order is not calculated and has 0 accrual",
			args: args{ctx: context.Background()},
			mockProcess: func() *orderProcessor {
				mockOrderService := mocks.NewOrdersService(t)
				mockOrderService.
					On("GetNotCalculatedOrder", mock.Anything).
					Return(&models.Order{Accrual: &zeroAccrual}, nil)

				return &orderProcessor{
					ordersService: mockOrderService,
				}
			},
			want: want{errMsg: ""},
		},
		{
			name: "GetInvoiceByUserID failed",
			args: args{ctx: context.Background()},
			mockProcess: func() *orderProcessor {
				mockOrderService := mocks.NewOrdersService(t)
				mockOrderService.
					On("GetNotCalculatedOrder", mock.Anything).
					Return(&models.Order{Accrual: &noZeroAccrual}, nil)

				mockInvoiceService := mocks.NewInvoicesService(t)
				mockInvoiceService.
					On("GetInvoiceByUserID", mock.Anything, mock.AnythingOfType("models.UserID")).
					Return(nil, fmt.Errorf("some repos error"))

				return &orderProcessor{
					invoicesService: mockInvoiceService,
					ordersService:   mockOrderService,
				}
			},
			want: want{errMsg: "some repos error"},
		},
		{
			name: "UpdateBalance failed",
			args: args{ctx: context.Background()},
			mockProcess: func() *orderProcessor {
				mockOrderService := mocks.NewOrdersService(t)
				mockOrderService.
					On("GetNotCalculatedOrder", mock.Anything).
					Return(&models.Order{Accrual: &noZeroAccrual}, nil)

				mockInvoiceService := mocks.NewInvoicesService(t)
				mockInvoiceService.
					On("GetInvoiceByUserID", mock.Anything, mock.AnythingOfType("models.UserID")).
					Return(&models.Invoice{}, nil).
					On("UpdateBalance", mock.Anything, mock.AnythingOfType("models.Invoice"), mock.AnythingOfType("float64"), mock.Anything).
					Return(nil, fmt.Errorf("some error"))

				return &orderProcessor{
					invoicesService: mockInvoiceService,
					ordersService:   mockOrderService,
				}
			},
			want: want{errMsg: "some error"},
		},
		{
			name: "UpdateOrdersCalculated failed",
			args: args{ctx: context.Background()},
			mockProcess: func() *orderProcessor {
				mockOrderService := mocks.NewOrdersService(t)
				mockOrderService.
					On("GetNotCalculatedOrder", mock.Anything).
					Return(&models.Order{Accrual: &noZeroAccrual}, nil).
					On("UpdateOrdersCalculated", mock.Anything, mock.AnythingOfType("models.Order"), mock.AnythingOfType("bool")).
					Return(nil, fmt.Errorf("some error"))

				mockInvoiceService := mocks.NewInvoicesService(t)
				mockInvoiceService.
					On("GetInvoiceByUserID", mock.Anything, mock.AnythingOfType("models.UserID")).
					Return(&models.Invoice{}, nil).
					On("UpdateBalance", mock.Anything, mock.AnythingOfType("models.Invoice"), mock.AnythingOfType("float64"), mock.Anything).
					Return(&models.Invoice{}, nil)

				return &orderProcessor{
					invoicesService: mockInvoiceService,
					ordersService:   mockOrderService,
				}
			},
			want: want{errMsg: "some error"},
		},
		{
			name: "Successfully case",
			args: args{ctx: context.Background()},
			mockProcess: func() *orderProcessor {
				mockOrderService := mocks.NewOrdersService(t)
				mockOrderService.
					On("GetNotCalculatedOrder", mock.Anything).
					Return(&models.Order{Accrual: &noZeroAccrual}, nil).
					On("UpdateOrdersCalculated", mock.Anything, mock.AnythingOfType("models.Order"), mock.AnythingOfType("bool")).
					Return(nil, nil)

				mockInvoiceService := mocks.NewInvoicesService(t)
				mockInvoiceService.
					On("GetInvoiceByUserID", mock.Anything, mock.AnythingOfType("models.UserID")).
					Return(&models.Invoice{}, nil).
					On("UpdateBalance", mock.Anything, mock.AnythingOfType("models.Invoice"), mock.AnythingOfType("float64"), mock.Anything).
					Return(&models.Invoice{}, nil)

				return &orderProcessor{
					invoicesService: mockInvoiceService,
					ordersService:   mockOrderService,
				}
			},
			want: want{errMsg: ""},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := test.mockProcess()
			gotErr := p.RecalculateProcessedOrders(test.args.ctx)

			if len(test.want.errMsg) > 0 {
				require.EqualError(t, gotErr, test.want.errMsg)
			} else {
				require.Equal(t, nil, gotErr)
			}
		})
	}
}
