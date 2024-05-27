package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"e1m0re/loyalty-srv/internal/models"
	"e1m0re/loyalty-srv/internal/service"
	mockservice "e1m0re/loyalty-srv/internal/service/mocks"
)

func TestHandler_GetOrders(t *testing.T) {
	jwtAuth := jwtauth.New("HS256", []byte("very secret key"), nil)
	type args struct {
		headers map[string]string
		ctx     context.Context
	}
	type want struct {
		expectedStatusCode   int
		expectedResponseBody string
	}
	tests := []struct {
		name         string
		method       string
		mockServices func() *service.Services
		args         args
		want         want
	}{
		{
			name:   "401 — unauthorized",
			method: http.MethodPost,
			mockServices: func() *service.Services {
				mockSecurityService := mockservice.NewSecurityService(t)
				mockSecurityService.
					On("GenerateAuthToken").
					Return(jwtAuth)

				mockOrdersService := mockservice.NewOrdersService(t)

				return &service.Services{
					OrdersService:   mockOrdersService,
					SecurityService: mockSecurityService,
				}
			},
			args: args{
				ctx:     context.WithValue(context.Background(), "userID", 1),
				headers: make(map[string]string),
			},
			want: want{
				expectedStatusCode:   http.StatusUnauthorized,
				expectedResponseBody: "",
			},
		},
		{
			name:   "500 — GetLoadedOrdersByUserID failed",
			method: http.MethodGet,
			mockServices: func() *service.Services {
				mockSecurityService := mockservice.NewSecurityService(t)
				mockSecurityService.
					On("GenerateAuthToken").
					Return(jwtAuth)

				mockOrdersService := mockservice.NewOrdersService(t)
				mockOrdersService.
					On("GetLoadedOrdersByUserID", mock.Anything, mock.AnythingOfType("models.UserID")).
					Return(nil, errors.New("some error"))

				return &service.Services{
					OrdersService:   mockOrdersService,
					SecurityService: mockSecurityService,
				}
			},
			args: args{
				ctx: context.WithValue(context.Background(), "userID", 1),
				headers: map[string]string{
					"Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6OSwidXNlcm5hbWUiOiJ1c2VyMiJ9.vY8OSC5qvDO-rLLnTUBGevkjIUm2oAjBuSsV75LO1Yw",
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
			mockServices: func() *service.Services {
				mockSecurityService := mockservice.NewSecurityService(t)
				mockSecurityService.
					On("GenerateAuthToken").
					Return(jwtAuth)

				mockOrdersService := mockservice.NewOrdersService(t)
				mockOrdersService.
					On("GetLoadedOrdersByUserID", mock.Anything, mock.AnythingOfType("models.UserID")).
					Return(&models.OrdersList{}, nil)

				return &service.Services{
					OrdersService:   mockOrdersService,
					SecurityService: mockSecurityService,
				}
			},
			args: args{
				ctx: context.WithValue(context.Background(), "userID", 1),
				headers: map[string]string{
					"Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6OSwidXNlcm5hbWUiOiJ1c2VyMiJ9.vY8OSC5qvDO-rLLnTUBGevkjIUm2oAjBuSsV75LO1Yw",
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
			mockServices: func() *service.Services {
				mockSecurityService := mockservice.NewSecurityService(t)
				mockSecurityService.
					On("GenerateAuthToken").
					Return(jwtAuth)

				accrual := 500
				ordersList := &models.OrdersList{
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
				}
				mockOrdersService := mockservice.NewOrdersService(t)
				mockOrdersService.
					On("GetLoadedOrdersByUserID", mock.Anything, mock.AnythingOfType("models.UserID")).
					Return(ordersList, nil)

				return &service.Services{
					OrdersService:   mockOrdersService,
					SecurityService: mockSecurityService,
				}
			},
			args: args{
				ctx: context.WithValue(context.Background(), "userID", 1),
				headers: map[string]string{
					"Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6OSwidXNlcm5hbWUiOiJ1c2VyMiJ9.vY8OSC5qvDO-rLLnTUBGevkjIUm2oAjBuSsV75LO1Yw",
				},
			},
			want: want{
				expectedStatusCode:   http.StatusOK,
				expectedResponseBody: "[{\"uploaded_at\":\"1984-04-01T12:13:00Z\",\"number\":\"1\",\"status\":\"NEW\"},{\"uploaded_at\":\"1984-04-01T12:13:05Z\",\"number\":\"2\",\"status\":\"PROCESSING\"},{\"uploaded_at\":\"1984-04-01T12:13:10Z\",\"number\":\"3\",\"status\":\"INVALID\"},{\"uploaded_at\":\"1984-04-01T12:13:15Z\",\"number\":\"4\",\"status\":\"PROCESSED\",\"accrual\":500}]",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			services := test.mockServices()
			handler := NewHandler(services)
			router := handler.NewRouter()

			req, err := http.NewRequestWithContext(test.args.ctx, test.method, "/api/user/orders", nil)
			for k, v := range test.args.headers {
				req.Header.Add(k, v)
			}
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			require.Equal(t, test.want.expectedStatusCode, rr.Code)
			require.Equal(t, test.want.expectedResponseBody, rr.Body.String())
		})
	}
}
