package service

import (
	"context"
	"e1m0re/loyalty-srv/internal/models"
	"fmt"
	"io"
	"net/http"
	"strings"
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
					Return(&models.Order{Accrual: nil}, nil).
					On("UpdateOrdersCalculated", mock.Anything, mock.AnythingOfType("models.Order"), true).
					Return(nil, nil)

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
					Return(&models.Order{Accrual: &zeroAccrual}, nil).
					On("UpdateOrdersCalculated", mock.Anything, mock.AnythingOfType("models.Order"), true).
					Return(nil, nil)

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

func Test_orderProcessor_RequestOrdersStatus(t *testing.T) {

	type args struct {
		ctx      context.Context
		orderNum models.OrderNum
	}
	type want struct {
		osi     *models.OrdersStatusInfo
		timeout int64
		errMsg  string
	}
	tests := []struct {
		name        string
		mockService func() OrdersProcessor
		args        args
		want        want
	}{
		{
			name: "httpClient.Get failed",
			mockService: func() OrdersProcessor {
				mockHTTPClient := mocks.NewHTTPClient(t)
				mockHTTPClient.
					On("Get", mock.Anything).
					Return(nil, fmt.Errorf("some error of request contractor"))

				return &orderProcessor{
					httpClient: mockHTTPClient,
				}
			},
			args: args{
				ctx:      context.Background(),
				orderNum: "1",
			},
			want: want{
				osi:     nil,
				timeout: 0,
				errMsg:  "some error of request contractor",
			},
		},
		{
			name: "Response with 204 code (StatusNoContent)",
			mockService: func() OrdersProcessor {
				mockHTTPClient := mocks.NewHTTPClient(t)
				response := &http.Response{
					Body:       io.NopCloser(strings.NewReader("")),
					StatusCode: http.StatusNoContent,
				}
				mockHTTPClient.
					On("Get", mock.Anything).
					Return(response, nil)

				return &orderProcessor{
					httpClient: mockHTTPClient,
				}
			},
			args: args{
				ctx:      context.Background(),
				orderNum: "1",
			},
			want: want{
				osi:     nil,
				timeout: 0,
				errMsg:  "",
			},
		},
		{
			name: "Response with 202 code (StatusOK) and bad body",
			mockService: func() OrdersProcessor {
				mockHTTPClient := mocks.NewHTTPClient(t)
				response := &http.Response{
					Body:       io.NopCloser(strings.NewReader("{\"order\":\"<number>\",\"status\":\"PROCESSED\",")),
					StatusCode: http.StatusOK,
				}
				mockHTTPClient.
					On("Get", mock.Anything).
					Return(response, nil)

				return &orderProcessor{
					httpClient: mockHTTPClient,
				}
			},
			args: args{
				ctx:      context.Background(),
				orderNum: "1",
			},
			want: want{
				osi:     nil,
				timeout: 0,
				errMsg:  "unexpected EOF",
			},
		},
		{
			name: "Response with 202 code (StatusOK)",
			mockService: func() OrdersProcessor {
				mockHTTPClient := mocks.NewHTTPClient(t)
				response := &http.Response{
					Body:       io.NopCloser(strings.NewReader("{\"order\":\"1\",\"status\":\"PROCESSED\",\"accrual\":500}")),
					StatusCode: http.StatusOK,
				}
				mockHTTPClient.
					On("Get", mock.Anything).
					Return(response, nil)

				return &orderProcessor{
					httpClient: mockHTTPClient,
				}
			},
			args: args{
				ctx:      context.Background(),
				orderNum: "1",
			},
			want: want{
				osi: &models.OrdersStatusInfo{
					OrderNumber: "1",
					Status:      "PROCESSED",
					Accrual:     500,
				},
				timeout: 0,
				errMsg:  "",
			},
		},
		{
			name: "Response with 429 code (StatusTooManyRequests) and wrung Retry-After header",
			mockService: func() OrdersProcessor {
				mockHTTPClient := mocks.NewHTTPClient(t)
				response := &http.Response{
					Body:       io.NopCloser(strings.NewReader("")),
					StatusCode: http.StatusTooManyRequests,
					Header: http.Header{
						"Retry-After": []string{""},
					},
				}
				mockHTTPClient.
					On("Get", mock.Anything).
					Return(response, nil)

				return &orderProcessor{
					httpClient: mockHTTPClient,
				}
			},
			args: args{
				ctx:      context.Background(),
				orderNum: "1",
			},
			want: want{
				osi:     nil,
				timeout: 0,
				errMsg:  "",
			},
		},
		{
			name: "Response with 429 code (StatusTooManyRequests) and without Retry-After header",
			mockService: func() OrdersProcessor {
				mockHTTPClient := mocks.NewHTTPClient(t)
				response := &http.Response{
					Body:       io.NopCloser(strings.NewReader("")),
					StatusCode: http.StatusTooManyRequests,
				}
				mockHTTPClient.
					On("Get", mock.Anything).
					Return(response, nil)

				return &orderProcessor{
					httpClient: mockHTTPClient,
				}
			},
			args: args{
				ctx:      context.Background(),
				orderNum: "1",
			},
			want: want{
				osi:     nil,
				timeout: 0,
				errMsg:  "",
			},
		},
		{
			name: "Response with 429 code (StatusTooManyRequests) and good Retry-After header",
			mockService: func() OrdersProcessor {
				mockHTTPClient := mocks.NewHTTPClient(t)
				response := &http.Response{
					Body:       io.NopCloser(strings.NewReader("")),
					StatusCode: http.StatusTooManyRequests,
					Header: http.Header{
						"Retry-After": []string{"30"},
					},
				}
				mockHTTPClient.
					On("Get", mock.Anything).
					Return(response, nil)

				return &orderProcessor{
					httpClient: mockHTTPClient,
				}
			},
			args: args{
				ctx:      context.Background(),
				orderNum: "1",
			},
			want: want{
				osi:     nil,
				timeout: 30,
				errMsg:  "",
			},
		},
		{
			name: "Response with 500 code (StatusInternalServerError)",
			mockService: func() OrdersProcessor {
				mockHTTPClient := mocks.NewHTTPClient(t)
				response := &http.Response{
					Body:       io.NopCloser(strings.NewReader("something wrong")),
					StatusCode: http.StatusInternalServerError,
				}
				mockHTTPClient.
					On("Get", mock.Anything).
					Return(response, nil)

				return &orderProcessor{
					httpClient: mockHTTPClient,
				}
			},
			args: args{
				ctx:      context.Background(),
				orderNum: "1",
			},
			want: want{
				osi:     nil,
				timeout: 0,
				errMsg:  "internal server error: something wrong",
			},
		},
		{
			name: "Response with unknown status code",
			mockService: func() OrdersProcessor {
				mockHTTPClient := mocks.NewHTTPClient(t)
				response := &http.Response{
					Body:       io.NopCloser(strings.NewReader("")),
					StatusCode: http.StatusMethodNotAllowed,
				}
				mockHTTPClient.
					On("Get", mock.Anything).
					Return(response, nil)

				return &orderProcessor{
					httpClient: mockHTTPClient,
				}
			},
			args: args{
				ctx:      context.Background(),
				orderNum: "1",
			},
			want: want{
				osi:     nil,
				timeout: 0,
				errMsg:  "unknown response type",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			op := test.mockService()
			gotOsi, gotTimeout, gotErr := op.RequestOrdersStatus(test.args.ctx, test.args.orderNum)
			require.Equal(t, test.want.osi, gotOsi)
			require.Equal(t, test.want.timeout, gotTimeout)
			if len(test.want.errMsg) > 0 {
				require.EqualError(t, gotErr, test.want.errMsg)
			} else {
				require.Equal(t, nil, gotErr)
			}
		})
	}
}
