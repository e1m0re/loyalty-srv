package handler

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/jwtauth/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"e1m0re/loyalty-srv/internal/apperrors"
	"e1m0re/loyalty-srv/internal/models"
	"e1m0re/loyalty-srv/internal/service"
	mockservice "e1m0re/loyalty-srv/internal/service/mocks"
)

func TestHandler_WritingOff(t *testing.T) {
	jwtAuth := jwtauth.New("HS256", []byte("very secret key"), nil)
	type args struct {
		ctx       context.Context
		headers   map[string]string
		inputBody string
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
			name:   "401",
			method: http.MethodPost,
			mockServices: func() *service.Services {
				mockSecurityService := mockservice.NewSecurityService(t)
				mockSecurityService.
					On("GenerateAuthToken").
					Return(jwtAuth)

				return &service.Services{
					SecurityService: mockSecurityService,
				}
			},
			args: args{
				ctx: context.WithValue(context.Background(), models.CKUserID, 1),
			},
			want: want{
				expectedStatusCode:   http.StatusUnauthorized,
				expectedResponseBody: "",
			},
		},
		{
			name:   "405",
			method: http.MethodGet,
			mockServices: func() *service.Services {
				mockSecurityService := mockservice.NewSecurityService(t)
				mockSecurityService.
					On("GenerateAuthToken").
					Return(jwtAuth)

				return &service.Services{
					SecurityService: mockSecurityService,
				}
			},
			args: args{
				ctx: context.WithValue(context.Background(), models.CKUserID, 1),
				headers: map[string]string{
					"Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6OSwidXNlcm5hbWUiOiJ1c2VyMiJ9.vY8OSC5qvDO-rLLnTUBGevkjIUm2oAjBuSsV75LO1Yw",
				},
			},
			want: want{
				expectedStatusCode:   http.StatusMethodNotAllowed,
				expectedResponseBody: "",
			},
		},
		{
			name:   "400 — empty body",
			method: http.MethodPost,
			mockServices: func() *service.Services {
				mockSecurityService := mockservice.NewSecurityService(t)
				mockSecurityService.
					On("GenerateAuthToken").
					Return(jwtAuth)

				return &service.Services{
					SecurityService: mockSecurityService,
				}
			},
			args: args{
				ctx: context.WithValue(context.Background(), models.CKUserID, 1),
				headers: map[string]string{
					"Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6OSwidXNlcm5hbWUiOiJ1c2VyMiJ9.vY8OSC5qvDO-rLLnTUBGevkjIUm2oAjBuSsV75LO1Yw",
				},
				inputBody: "",
			},
			want: want{
				expectedStatusCode:   http.StatusBadRequest,
				expectedResponseBody: "",
			},
		},
		{
			name:   "500 — ValidateNumber failed",
			method: http.MethodPost,
			mockServices: func() *service.Services {
				mockSecurityService := mockservice.NewSecurityService(t)
				mockSecurityService.
					On("GenerateAuthToken").
					Return(jwtAuth)

				mockOrdersService := mockservice.NewOrdersService(t)
				mockOrdersService.
					On("ValidateNumber", mock.Anything, mock.AnythingOfType("models.OrderNum")).
					Return(false, fmt.Errorf("some error"))

				return &service.Services{
					OrdersService:   mockOrdersService,
					SecurityService: mockSecurityService,
				}
			},
			args: args{
				ctx: context.WithValue(context.Background(), models.CKUserID, 1),
				headers: map[string]string{
					"Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6OSwidXNlcm5hbWUiOiJ1c2VyMiJ9.vY8OSC5qvDO-rLLnTUBGevkjIUm2oAjBuSsV75LO1Yw",
				},
				inputBody: "{\"order\":\"2377225624\",\"sum\":751}",
			},
			want: want{
				expectedStatusCode:   http.StatusInternalServerError,
				expectedResponseBody: "",
			},
		},
		{
			name:   "422 — invalid order number",
			method: http.MethodPost,
			mockServices: func() *service.Services {
				mockSecurityService := mockservice.NewSecurityService(t)
				mockSecurityService.
					On("GenerateAuthToken").
					Return(jwtAuth)

				mockOrdersService := mockservice.NewOrdersService(t)
				mockOrdersService.
					On("ValidateNumber", mock.Anything, mock.AnythingOfType("models.OrderNum")).
					Return(false, nil)

				return &service.Services{
					OrdersService:   mockOrdersService,
					SecurityService: mockSecurityService,
				}
			},
			args: args{
				ctx: context.WithValue(context.Background(), models.CKUserID, 1),
				headers: map[string]string{
					"Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6OSwidXNlcm5hbWUiOiJ1c2VyMiJ9.vY8OSC5qvDO-rLLnTUBGevkjIUm2oAjBuSsV75LO1Yw",
				},
				inputBody: "{\"order\":\"2377225624\",\"sum\":751}",
			},
			want: want{
				expectedStatusCode:   http.StatusUnprocessableEntity,
				expectedResponseBody: "",
			},
		},
		{
			name:   "500 — GetAccountByUserID failed",
			method: http.MethodPost,
			mockServices: func() *service.Services {
				mockSecurityService := mockservice.NewSecurityService(t)
				mockSecurityService.
					On("GenerateAuthToken").
					Return(jwtAuth)

				mockOrdersService := mockservice.NewOrdersService(t)
				mockOrdersService.
					On("ValidateNumber", mock.Anything, mock.AnythingOfType("models.OrderNum")).
					Return(true, nil)

				mockAccountsService := mockservice.NewAccountsService(t)
				mockAccountsService.
					On("GetAccountByUserID", mock.Anything, mock.AnythingOfType("models.UserID")).
					Return(nil, fmt.Errorf("some error"))

				return &service.Services{
					AccountsService: mockAccountsService,
					OrdersService:   mockOrdersService,
					SecurityService: mockSecurityService,
				}
			},
			args: args{
				ctx: context.WithValue(context.Background(), models.CKUserID, 1),
				headers: map[string]string{
					"Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6OSwidXNlcm5hbWUiOiJ1c2VyMiJ9.vY8OSC5qvDO-rLLnTUBGevkjIUm2oAjBuSsV75LO1Yw",
				},
				inputBody: "{\"order\":\"2377225624\",\"sum\":751}",
			},
			want: want{
				expectedStatusCode:   http.StatusInternalServerError,
				expectedResponseBody: "",
			},
		},
		{
			name:   "500 — Withdraw failed",
			method: http.MethodPost,
			mockServices: func() *service.Services {
				mockSecurityService := mockservice.NewSecurityService(t)
				mockSecurityService.
					On("GenerateAuthToken").
					Return(jwtAuth)

				mockOrdersService := mockservice.NewOrdersService(t)
				mockOrdersService.
					On("ValidateNumber", mock.Anything, mock.AnythingOfType("models.OrderNum")).
					Return(true, nil)

				mockAccountsService := mockservice.NewAccountsService(t)
				mockAccountsService.
					On("GetAccountByUserID", mock.Anything, mock.AnythingOfType("models.UserID")).
					Return(&models.Account{}, nil).
					On("UpdateBalance", mock.Anything, models.Account{}, mock.AnythingOfType("float64"), mock.AnythingOfType("models.OrderNum")).
					Return(nil, fmt.Errorf("some error"))

				return &service.Services{
					AccountsService: mockAccountsService,
					OrdersService:   mockOrdersService,
					SecurityService: mockSecurityService,
				}
			},
			args: args{
				ctx: context.WithValue(context.Background(), models.CKUserID, 1),
				headers: map[string]string{
					"Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6OSwidXNlcm5hbWUiOiJ1c2VyMiJ9.vY8OSC5qvDO-rLLnTUBGevkjIUm2oAjBuSsV75LO1Yw",
				},
				inputBody: "{\"order\":\"2377225624\",\"sum\":751}",
			},
			want: want{
				expectedStatusCode:   http.StatusInternalServerError,
				expectedResponseBody: "",
			},
		},
		{
			name:   "402",
			method: http.MethodPost,
			mockServices: func() *service.Services {
				mockSecurityService := mockservice.NewSecurityService(t)
				mockSecurityService.
					On("GenerateAuthToken").
					Return(jwtAuth)

				mockOrdersService := mockservice.NewOrdersService(t)
				mockOrdersService.
					On("ValidateNumber", mock.Anything, mock.AnythingOfType("models.OrderNum")).
					Return(true, nil)

				mockAccountsService := mockservice.NewAccountsService(t)
				mockAccountsService.
					On("GetAccountByUserID", mock.Anything, mock.AnythingOfType("models.UserID")).
					Return(&models.Account{}, nil).
					On("UpdateBalance", mock.Anything, models.Account{}, mock.AnythingOfType("float64"), mock.AnythingOfType("models.OrderNum")).
					Return(nil, apperrors.ErrAccountHasNotEnoughFunds)

				return &service.Services{
					AccountsService: mockAccountsService,
					OrdersService:   mockOrdersService,
					SecurityService: mockSecurityService,
				}
			},
			args: args{
				ctx: context.WithValue(context.Background(), models.CKUserID, 1),
				headers: map[string]string{
					"Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6OSwidXNlcm5hbWUiOiJ1c2VyMiJ9.vY8OSC5qvDO-rLLnTUBGevkjIUm2oAjBuSsV75LO1Yw",
				},
				inputBody: "{\"order\":\"2377225624\",\"sum\":751}",
			},
			want: want{
				expectedStatusCode:   http.StatusPaymentRequired,
				expectedResponseBody: "",
			},
		},
		{
			name:   "200",
			method: http.MethodPost,
			mockServices: func() *service.Services {
				mockSecurityService := mockservice.NewSecurityService(t)
				mockSecurityService.
					On("GenerateAuthToken").
					Return(jwtAuth)

				mockOrdersService := mockservice.NewOrdersService(t)
				mockOrdersService.
					On("ValidateNumber", mock.Anything, mock.AnythingOfType("models.OrderNum")).
					Return(true, nil)

				mockAccountsService := mockservice.NewAccountsService(t)
				mockAccountsService.
					On("GetAccountByUserID", mock.Anything, mock.AnythingOfType("models.UserID")).
					Return(&models.Account{}, nil).
					On("UpdateBalance", mock.Anything, models.Account{}, mock.AnythingOfType("float64"), mock.AnythingOfType("models.OrderNum")).
					Return(nil, nil)

				return &service.Services{
					AccountsService: mockAccountsService,
					OrdersService:   mockOrdersService,
					SecurityService: mockSecurityService,
				}
			},
			args: args{
				ctx: context.WithValue(context.Background(), models.CKUserID, 1),
				headers: map[string]string{
					"Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6OSwidXNlcm5hbWUiOiJ1c2VyMiJ9.vY8OSC5qvDO-rLLnTUBGevkjIUm2oAjBuSsV75LO1Yw",
				},
				inputBody: "{\"order\":\"2377225624\",\"sum\":751}",
			},
			want: want{
				expectedStatusCode:   http.StatusOK,
				expectedResponseBody: "",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			services := test.mockServices()
			handler := NewHandler(services)
			router := handler.NewRouter()

			req, err := http.NewRequestWithContext(test.args.ctx, test.method, "/api/user/balance/withdraw", bytes.NewReader([]byte(test.args.inputBody)))
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
