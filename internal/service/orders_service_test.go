package service

import (
	"context"
	"fmt"
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
				errMsg: apperrors.ErrEmptyOrderNumber.Error(),
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
						UserID:   4,
						OrderNum: "12345678903",
					}).
					Return(&models.Order{
						ID:         1,
						UserID:     4,
						Number:     "12345678903",
						Status:     models.OrderStatusNew,
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

func Test_ordersService_GetLoadedOrdersByUserID(t *testing.T) {
	accrual := 500
	type args struct {
		ctx    context.Context
		userID models.UserID
	}
	type want struct {
		ordersList models.OrdersList
		errMsg     string
	}
	tests := []struct {
		name             string
		mockRepositories func() *repository.Repositories
		args             args
		want             want
	}{
		{
			name: "Error in repo",
			mockRepositories: func() *repository.Repositories {
				mockOrderRepo := mocks.NewOrderRepository(t)
				mockOrderRepo.
					On("GetLoadedOrdersByUserID", mock.Anything, models.UserID(1)).
					Return(&models.OrdersList{}, fmt.Errorf("some repo error"))
				return &repository.Repositories{
					OrderRepository: mockOrderRepo,
				}
			},
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			want: want{
				ordersList: models.OrdersList{},
				errMsg:     "some repo error",
			},
		},
		{
			name: "Successfully test with empty orders list",
			mockRepositories: func() *repository.Repositories {
				mockOrderRepo := mocks.NewOrderRepository(t)
				mockOrderRepo.
					On("GetLoadedOrdersByUserID", mock.Anything, models.UserID(1)).
					Return(&models.OrdersList{}, nil)
				return &repository.Repositories{
					OrderRepository: mockOrderRepo,
				}
			},
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			want: want{
				ordersList: models.OrdersList{},
				errMsg:     "",
			},
		},
		{
			name: "Successfully test with no empty orders list",
			mockRepositories: func() *repository.Repositories {
				mockOrderRepo := mocks.NewOrderRepository(t)
				mockOrderRepo.
					On("GetLoadedOrdersByUserID", mock.Anything, models.UserID(1)).
					Return(&models.OrdersList{
						{
							ID:         1,
							UserID:     1,
							Number:     "1",
							Status:     "NEW",
							UploadedAt: time.Date(1984, time.April, 1, 12, 13, 0, 0, time.UTC),
						},
						{
							ID:         2,
							UserID:     1,
							Number:     "2",
							Status:     "PROCESSING",
							UploadedAt: time.Date(1984, time.April, 1, 12, 13, 5, 0, time.UTC),
						},
						{
							ID:         3,
							UserID:     1,
							Number:     "3",
							Status:     "INVALID",
							UploadedAt: time.Date(1984, time.April, 1, 12, 13, 10, 0, time.UTC),
						},
						{
							ID:         4,
							UserID:     1,
							Number:     "4",
							Status:     "PROCESSED",
							Accrual:    &accrual,
							UploadedAt: time.Date(1984, time.April, 1, 12, 13, 15, 0, time.UTC),
						},
					}, nil)
				return &repository.Repositories{
					OrderRepository: mockOrderRepo,
				}
			},
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			want: want{
				ordersList: models.OrdersList{
					{
						ID:         1,
						UserID:     1,
						Number:     "1",
						Status:     "NEW",
						UploadedAt: time.Date(1984, time.April, 1, 12, 13, 0, 0, time.UTC),
					},
					{
						ID:         2,
						UserID:     1,
						Number:     "2",
						Status:     "PROCESSING",
						UploadedAt: time.Date(1984, time.April, 1, 12, 13, 5, 0, time.UTC),
					},
					{
						ID:         3,
						UserID:     1,
						Number:     "3",
						Status:     "INVALID",
						UploadedAt: time.Date(1984, time.April, 1, 12, 13, 10, 0, time.UTC),
					},
					{
						ID:         4,
						UserID:     1,
						Number:     "4",
						Status:     "PROCESSED",
						Accrual:    &accrual,
						UploadedAt: time.Date(1984, time.April, 1, 12, 13, 15, 0, time.UTC),
					},
				},
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
			gotOrdersList, gotErr := os.GetLoadedOrdersByUserID(test.args.ctx, test.args.userID)
			require.ElementsMatch(t, test.want.ordersList, *gotOrdersList)
			if len(test.want.errMsg) > 0 {
				require.EqualError(t, gotErr, test.want.errMsg)
			}
		})
	}
}
