package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"e1m0re/loyalty-srv/internal/apperrors"
	"e1m0re/loyalty-srv/internal/models"
	"e1m0re/loyalty-srv/internal/repository"
	"e1m0re/loyalty-srv/internal/repository/mocks"
)

func TestOrdersService_ValidateNumber(t *testing.T) {
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

func TestOrdersService_NewOrder(t *testing.T) {
	type args struct {
		ctx       context.Context
		orderInfo models.OrderInfo
	}
	type want struct {
		order  *models.Order
		errMsg string
	}
	tests := []struct {
		name             string
		mockRepositories func() *repository.Repositories
		args             args
		want             want
	}{
		{
			name: "Invalid order number",
			mockRepositories: func() *repository.Repositories {

				return &repository.Repositories{}
			},
			args: args{
				ctx: context.Background(),
				orderInfo: models.OrderInfo{
					UserID:   1,
					OrderNum: "12345678904",
				},
			},
			want: want{
				order:  nil,
				errMsg: "",
			},
		},
		{
			name: "GetOrderByNumber failed",
			mockRepositories: func() *repository.Repositories {
				mockOrderRepo := mocks.NewOrderRepository(t)
				mockOrderRepo.
					On("GetOrderByNumber", mock.Anything, mock.AnythingOfType("models.OrderNum")).
					Return(nil, fmt.Errorf("some repos error"))

				return &repository.Repositories{
					OrderRepository: mockOrderRepo,
				}
			},
			args: args{
				ctx: context.Background(),
				orderInfo: models.OrderInfo{
					UserID:   1,
					OrderNum: "12345678903",
				},
			},
			want: want{
				order:  nil,
				errMsg: "some repos error",
			},
		},
		{
			name: "order was loaded",
			mockRepositories: func() *repository.Repositories {
				mockOrderRepo := mocks.NewOrderRepository(t)
				mockOrderRepo.
					On("GetOrderByNumber", mock.Anything, mock.AnythingOfType("models.OrderNum")).
					Return(&models.Order{UserID: 1}, nil)

				return &repository.Repositories{
					OrderRepository: mockOrderRepo,
				}
			},
			args: args{
				ctx: context.Background(),
				orderInfo: models.OrderInfo{
					UserID:   1,
					OrderNum: "12345678903",
				},
			},
			want: want{
				order:  nil,
				errMsg: apperrors.ErrOrderWasLoaded.Error(),
			},
		},
		{
			name: "order was loaded by other user",
			mockRepositories: func() *repository.Repositories {
				mockOrderRepo := mocks.NewOrderRepository(t)
				mockOrderRepo.
					On("GetOrderByNumber", mock.Anything, mock.AnythingOfType("models.OrderNum")).
					Return(&models.Order{UserID: 2}, nil)

				return &repository.Repositories{
					OrderRepository: mockOrderRepo,
				}
			},
			args: args{
				ctx: context.Background(),
				orderInfo: models.OrderInfo{
					UserID:   1,
					OrderNum: "12345678903",
				},
			},
			want: want{
				order:  nil,
				errMsg: apperrors.ErrOrderWasLoadedByAnotherUser.Error(),
			},
		},
		{
			name: "AddOrder failed",
			mockRepositories: func() *repository.Repositories {
				mockOrderRepo := mocks.NewOrderRepository(t)
				mockOrderRepo.
					On("GetOrderByNumber", mock.Anything, mock.AnythingOfType("models.OrderNum")).
					Return(nil, nil).
					On("AddOrder", mock.Anything, mock.AnythingOfType("models.OrderInfo")).
					Return(nil, fmt.Errorf("some repos error"))

				return &repository.Repositories{
					OrderRepository: mockOrderRepo,
				}
			},
			args: args{
				ctx: context.Background(),
				orderInfo: models.OrderInfo{
					UserID:   1,
					OrderNum: "12345678903",
				},
			},
			want: want{
				order:  nil,
				errMsg: "some repos error",
			},
		},
		{
			name: "order added successfully",
			mockRepositories: func() *repository.Repositories {
				orderInfo := models.OrderInfo{
					UserID:   1,
					OrderNum: "12345678903",
				}
				mockOrderRepo := mocks.NewOrderRepository(t)
				mockOrderRepo.
					On("GetOrderByNumber", mock.Anything, orderInfo.OrderNum).
					Return(nil, nil).
					On("AddOrder", mock.Anything, orderInfo).
					Return(&models.Order{
						ID:         1,
						UserID:     1,
						Number:     "12345678903",
						Status:     models.OrderStatusNew,
						UploadedAt: time.Date(1703, time.May, 27, 12, 0, 0, 0, time.UTC),
					}, nil)

				return &repository.Repositories{
					OrderRepository: mockOrderRepo,
				}
			},
			args: args{
				ctx: context.Background(),
				orderInfo: models.OrderInfo{
					UserID:   1,
					OrderNum: "12345678903",
				},
			},
			want: want{
				order: &models.Order{
					ID:         1,
					Number:     "12345678903",
					UserID:     1,
					Status:     models.OrderStatusNew,
					UploadedAt: time.Date(1703, time.May, 27, 12, 0, 0, 0, time.UTC),
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
			gotOrder, gotErr := os.NewOrder(test.args.ctx, test.args.orderInfo)
			require.Equal(t, test.want.order, gotOrder)
			if len(test.want.errMsg) > 0 {
				require.EqualError(t, gotErr, test.want.errMsg)
			}
		})
	}
}

func TestOrdersService_GetLoadedOrdersByUserID(t *testing.T) {
	accrual := float64(500)
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
					Return(&models.OrdersList{}, fmt.Errorf("some repos error"))
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
				errMsg:     "some repos error",
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

func Test_ordersService_UpdateOrdersCalculated(t *testing.T) {
	type args struct {
		ctx        context.Context
		order      models.Order
		calculated bool
	}
	type want struct {
		order  *models.Order
		errMsg string
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
					On("UpdateOrdersCalculated", mock.Anything, mock.AnythingOfType("models.Order"), mock.AnythingOfType("bool")).
					Return(nil, fmt.Errorf("some repos error"))

				return &repository.Repositories{
					OrderRepository: mockOrderRepo,
				}
			},
			args: args{
				ctx:        context.Background(),
				order:      models.Order{},
				calculated: true,
			},
			want: want{
				order:  nil,
				errMsg: "some repos error",
			},
		},
		{
			name: "Successfully case",
			mockRepositories: func() *repository.Repositories {
				mockOrderRepo := mocks.NewOrderRepository(t)
				mockOrderRepo.
					On("UpdateOrdersCalculated", mock.Anything, mock.AnythingOfType("models.Order"), mock.AnythingOfType("bool")).
					Return(&models.Order{Calculated: true}, nil)

				return &repository.Repositories{
					OrderRepository: mockOrderRepo,
				}
			},
			args: args{
				ctx:        context.Background(),
				order:      models.Order{},
				calculated: true,
			},
			want: want{
				order:  &models.Order{Calculated: true},
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
			gotOrder, gotErr := os.UpdateOrdersCalculated(test.args.ctx, test.args.order, test.args.calculated)
			require.Equal(t, test.want.order, gotOrder)
			if len(test.want.errMsg) > 0 {
				require.EqualError(t, gotErr, test.want.errMsg)
			}
		})
	}
}

func Test_ordersService_UpdateOrdersStatus(t *testing.T) {
	accrual := float64(500)
	type args struct {
		ctx     context.Context
		order   models.Order
		status  models.OrdersStatus
		accrual float64
	}
	type want struct {
		order  *models.Order
		errMsg string
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
					On("UpdateOrdersStatus", mock.Anything, mock.AnythingOfType("models.Order"), mock.Anything, mock.AnythingOfType("float64")).
					Return(nil, fmt.Errorf("some repos error"))

				return &repository.Repositories{
					OrderRepository: mockOrderRepo,
				}
			},
			args: args{
				ctx:     context.Background(),
				order:   models.Order{},
				status:  models.OrderStatusProcessed,
				accrual: accrual,
			},
			want: want{
				order:  nil,
				errMsg: "some repos error",
			},
		},
		{
			name: "Successfully case",
			mockRepositories: func() *repository.Repositories {
				mockOrderRepo := mocks.NewOrderRepository(t)
				mockOrderRepo.
					On("UpdateOrdersStatus", mock.Anything, mock.AnythingOfType("models.Order"), mock.Anything, mock.AnythingOfType("float64")).
					Return(&models.Order{
						Status:  models.OrderStatusProcessed,
						Accrual: &accrual,
					}, nil)

				return &repository.Repositories{
					OrderRepository: mockOrderRepo,
				}
			},
			args: args{
				ctx:     context.Background(),
				order:   models.Order{},
				status:  models.OrderStatusProcessed,
				accrual: accrual,
			},
			want: want{
				order: &models.Order{
					Status:  models.OrderStatusProcessed,
					Accrual: &accrual,
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
			gotOrder, gotErr := os.UpdateOrdersStatus(test.args.ctx, test.args.order, test.args.status, test.args.accrual)
			require.Equal(t, test.want.order, gotOrder)
			if len(test.want.errMsg) > 0 {
				require.EqualError(t, gotErr, test.want.errMsg)
			}
		})
	}
}

func Test_ordersService_GetNotCalculatedOrder(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type want struct {
		order  *models.Order
		errMsg string
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
					On("GetNotCalculatedOrder", mock.Anything).
					Return(nil, fmt.Errorf("some repos error"))

				return &repository.Repositories{
					OrderRepository: mockOrderRepo,
				}
			},
			args: args{
				ctx: context.Background(),
			},
			want: want{
				order:  nil,
				errMsg: "some repos error",
			},
		},
		{
			name: "All orders are calculated",
			mockRepositories: func() *repository.Repositories {
				mockOrderRepo := mocks.NewOrderRepository(t)
				mockOrderRepo.
					On("GetNotCalculatedOrder", mock.Anything).
					Return(nil, nil)

				return &repository.Repositories{
					OrderRepository: mockOrderRepo,
				}
			},
			args: args{
				ctx: context.Background(),
			},
			want: want{
				order:  nil,
				errMsg: "",
			},
		},
		{
			name: "Successfully case",
			mockRepositories: func() *repository.Repositories {
				mockOrderRepo := mocks.NewOrderRepository(t)
				mockOrderRepo.
					On("GetNotCalculatedOrder", mock.Anything).
					Return(&models.Order{Calculated: false}, nil)

				return &repository.Repositories{
					OrderRepository: mockOrderRepo,
				}
			},
			args: args{
				ctx: context.Background(),
			},
			want: want{
				order:  &models.Order{Calculated: false},
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
			gotOrder, gotErr := os.GetNotCalculatedOrder(test.args.ctx)
			require.Equal(t, test.want.order, gotOrder)
			if len(test.want.errMsg) > 0 {
				require.EqualError(t, gotErr, test.want.errMsg)
			}
		})
	}
}

func Test_ordersService_GetNotProcessedOrder(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type want struct {
		order  *models.Order
		errMsg string
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
					On("GetNotProcessedOrder", mock.Anything).
					Return(nil, fmt.Errorf("some repos error"))

				return &repository.Repositories{
					OrderRepository: mockOrderRepo,
				}
			},
			args: args{
				ctx: context.Background(),
			},
			want: want{
				order:  nil,
				errMsg: "some repos error",
			},
		},
		{
			name: "All orders are processed",
			mockRepositories: func() *repository.Repositories {
				mockOrderRepo := mocks.NewOrderRepository(t)
				mockOrderRepo.
					On("GetNotProcessedOrder", mock.Anything).
					Return(nil, nil)

				return &repository.Repositories{
					OrderRepository: mockOrderRepo,
				}
			},
			args: args{
				ctx: context.Background(),
			},
			want: want{
				order:  nil,
				errMsg: "",
			},
		},
		{
			name: "Successfully case",
			mockRepositories: func() *repository.Repositories {
				mockOrderRepo := mocks.NewOrderRepository(t)
				mockOrderRepo.
					On("GetNotProcessedOrder", mock.Anything).
					Return(&models.Order{Status: models.OrderStatusNew}, nil)

				return &repository.Repositories{
					OrderRepository: mockOrderRepo,
				}
			},
			args: args{
				ctx: context.Background(),
			},
			want: want{
				order:  &models.Order{Status: models.OrderStatusNew},
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
			gotOrder, gotErr := os.GetNotProcessedOrder(test.args.ctx)
			require.Equal(t, test.want.order, gotOrder)
			if len(test.want.errMsg) > 0 {
				require.EqualError(t, gotErr, test.want.errMsg)
			}
		})
	}
}
