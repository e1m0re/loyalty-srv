package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"e1m0re/loyalty-srv/internal/models"
	"e1m0re/loyalty-srv/internal/service"
	mockservice "e1m0re/loyalty-srv/internal/service/mocks"
)

func TestHandler_GetOrders(t *testing.T) {
	type args struct {
		mockServices func() *service.Services
	}
	type want struct {
		expectedStatusCode   int
		expectedResponseBody string
	}
	tests := []struct {
		name   string
		method string
		args   args
		want   want
	}{
		{
			name:   "500 — GetLoadedOrdersByUserId failed",
			method: http.MethodGet,
			args: args{
				mockServices: func() *service.Services {
					mockOrdersService := mockservice.NewOrdersService(t)
					mockOrdersService.
						On("GetLoadedOrdersByUserId", mock.Anything, models.UserId(1)).
						Return(make(models.OrdersList, 0), errors.New("some error"))

					return &service.Services{
						OrdersService: mockOrdersService,
					}
				},
			},
			want: want{
				expectedStatusCode:   http.StatusInternalServerError,
				expectedResponseBody: "",
			},
		},
		{
			name:   "204 — no data",
			method: http.MethodGet,
			args: args{
				mockServices: func() *service.Services {
					mockOrdersService := mockservice.NewOrdersService(t)
					mockOrdersService.
						On("GetLoadedOrdersByUserId", mock.Anything, models.UserId(1)).
						Return(make(models.OrdersList, 0), nil)

					return &service.Services{
						OrdersService: mockOrdersService,
					}
				},
			},
			want: want{
				expectedStatusCode:   http.StatusNoContent,
				expectedResponseBody: "",
			},
		},
		{
			name:   "200",
			method: http.MethodGet,
			args: args{
				mockServices: func() *service.Services {
					ordersList := models.OrdersList{
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
							Accrual:    500,
							UploadedAt: time.Date(1984, time.April, 1, 12, 13, 15, 0, time.UTC),
						},
					}
					mockOrdersService := mockservice.NewOrdersService(t)
					mockOrdersService.
						On("GetLoadedOrdersByUserId", mock.Anything, models.UserId(1)).
						Return(ordersList, nil)

					return &service.Services{
						OrdersService: mockOrdersService,
					}
				},
			},
			want: want{
				expectedStatusCode:   http.StatusOK,
				expectedResponseBody: "[{\"uploaded_at\":\"1984-04-01T12:13:00Z\",\"ID\":1,\"UserID\":1,\"number\":\"1\",\"status\":\"NEW\"},{\"uploaded_at\":\"1984-04-01T12:13:05Z\",\"ID\":2,\"UserID\":1,\"number\":\"2\",\"status\":\"PROCESSING\"},{\"uploaded_at\":\"1984-04-01T12:13:10Z\",\"ID\":3,\"UserID\":1,\"number\":\"3\",\"status\":\"INVALID\"},{\"uploaded_at\":\"1984-04-01T12:13:15Z\",\"ID\":4,\"UserID\":1,\"number\":\"4\",\"status\":\"PROCESSED\",\"accrual\":500}]",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			services := test.args.mockServices()
			handler := NewHandler(services)
			router := handler.NewRouter()

			req, err := http.NewRequest(test.method, "/api/user/orders", nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			require.Equal(t, test.want.expectedStatusCode, rr.Code)
			require.Equal(t, test.want.expectedResponseBody, rr.Body.String())
		})
	}
}
