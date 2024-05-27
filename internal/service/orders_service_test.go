package service

import (
	"context"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"e1m0re/loyalty-srv/internal/apperrors"
	"e1m0re/loyalty-srv/internal/models"
	"e1m0re/loyalty-srv/internal/repository"
	"e1m0re/loyalty-srv/internal/repository/mocks"
)

func Test_ordersService_ValidateNumber(t *testing.T) {
	type args struct {
		ctx      context.Context
		orderNum models.OrderNum
	}
	type want struct {
		ok     bool
		errMsg string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "empty order number",
			args: args{
				ctx:      context.Background(),
				orderNum: "",
			},
			want: want{
				ok:     false,
				errMsg: apperrors.EmptyOrderNumberError.Error(),
			},
		},
		{
			name: "order number contains not only numbers",
			args: args{
				ctx:      context.Background(),
				orderNum: "123a-123",
			},
			want: want{
				ok:     false,
				errMsg: "strconv.Atoi: parsing \"a\": invalid syntax",
			},
		},
		{
			name: "valid order number",
			args: args{
				ctx:      context.Background(),
				orderNum: "12345678904",
			},
			want: want{
				ok: false,
			},
		},
		{
			name: "valid order number",
			args: args{
				ctx:      context.Background(),
				orderNum: "12345678903",
			},
			want: want{
				ok: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			os := ordersService{
				orderRepository: mocks.NewOrderRepository(t),
			}
			gotOk, gotErr := os.ValidateNumber(test.args.ctx, test.args.orderNum)
			require.Equal(t, test.want.ok, gotOk)
			if len(test.want.errMsg) > 0 {
				require.EqualError(t, gotErr, test.want.errMsg)
			}
		})
	}
}

func Test_ordersService_NewOrder(t *testing.T) {
	type fields struct {
		orderRepository *repository.OrderRepository
	}
	type args struct {
		ctx      context.Context
		orderNum models.OrderNum
	}
	type want struct {
		order  models.Order
		isNew  bool
		errMsg string
	}
	tests := []struct {
		name             string
		mockRepositories func() *repository.Repositories
		args             args
		want             want
	}{
		{
			name: "order added successfully",
			mockRepositories: func() *repository.Repositories {
				mockOrderRepo := mocks.NewOrderRepository(t)
				mockOrderRepo.
					On("AddOrder", mock.Anything, models.OrderInfo{
						UserId:   4,
						OrderNum: "12345678903",
					}).
					Return(&models.Order{
						ID:         1,
						UserID:     4,
						Number:     "12345678903",
						Status:     models.OrderStatusNew,
						Accrual:    0,
						UploadedAt: time.Date(1703, time.May, 27, 12, 0, 0, 0, time.UTC),
					}, nil)
				return &repository.Repositories{
					OrderRepository: mockOrderRepo,
				}
			},
			args: args{
				ctx:      context.Background(),
				orderNum: "12345678903",
			},
			want: want{
				order: models.Order{
					ID:         1,
					Number:     "12345678903",
					UserID:     1,
					Status:     models.OrderStatusNew,
					Accrual:    0,
					UploadedAt: time.Date(1703, time.May, 27, 12, 0, 0, 0, time.UTC),
				},
				isNew:  true,
				errMsg: "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := test.mockRepositories()
			os := ordersService{
				orderRepository: repo.OrderRepository,
			}
			gotOrder, gotIsNew, gotErr := os.NewOrder(test.args.ctx, test.args.orderNum)
			require.Equal(t, test.want.order.ID, gotOrder.ID)
			require.Equal(t, test.want.order.Number, gotOrder.Number)
			require.Equal(t, test.want.isNew, gotIsNew)
			if len(test.want.errMsg) > 0 {
				require.EqualError(t, gotErr, test.want.errMsg)
			}
		})
	}
}
