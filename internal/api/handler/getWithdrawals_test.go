package handler

import (
	"context"
	"errors"
	"fmt"
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

func TestHandler_GetWithdrawals(t *testing.T) {
	jwtAuth := jwtauth.New("HS256", []byte("very secret key"), nil)
	type args struct {
		ctx     context.Context
		headers map[string]string
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
				ctx:     context.WithValue(context.Background(), models.CKUserID, 1),
				headers: make(map[string]string),
			},
			want: want{
				expectedStatusCode:   http.StatusUnauthorized,
				expectedResponseBody: "",
			},
		},
		{
			name:   "405",
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
				ctx:     context.WithValue(context.Background(), models.CKUserID, 1),
				headers: make(map[string]string),
			},
			want: want{
				expectedStatusCode:   http.StatusMethodNotAllowed,
				expectedResponseBody: "",
			},
		},
		{
			name:   "500 — GetAccountByUserID failed",
			method: http.MethodGet,
			mockServices: func() *service.Services {
				mockSecurityService := mockservice.NewSecurityService(t)
				mockSecurityService.
					On("GenerateAuthToken").
					Return(jwtAuth)

				mockAccountsService := mockservice.NewAccountsService(t)
				mockAccountsService.
					On("GetAccountByUserID", mock.Anything, mock.AnythingOfType("models.UserID")).
					Return(nil, errors.New("some error"))

				return &service.Services{
					AccountsService: mockAccountsService,
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
				expectedStatusCode:   http.StatusInternalServerError,
				expectedResponseBody: "",
			},
		},
		{
			name:   "500 — GetWithdrawals failed",
			method: http.MethodGet,
			mockServices: func() *service.Services {
				mockSecurityService := mockservice.NewSecurityService(t)
				mockSecurityService.
					On("GenerateAuthToken").
					Return(jwtAuth)

				mockAccountsService := mockservice.NewAccountsService(t)
				mockAccountsService.
					On("GetAccountByUserID", mock.Anything, mock.AnythingOfType("models.UserID")).
					Return(&models.Account{}, nil).
					On("GetWithdrawals", mock.Anything, &models.Account{}).
					Return(nil, fmt.Errorf("some error"))

				return &service.Services{
					AccountsService: mockAccountsService,
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
				expectedStatusCode:   http.StatusInternalServerError,
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

				withdrawalsList := &models.WithdrawalsList{
					{
						OrderNum:    "2377225624",
						Sum:         500,
						ProcessedAt: time.Date(1984, time.April, 5, 6, 7, 8, 0, time.UTC),
					},
				}
				mockAccountsService := mockservice.NewAccountsService(t)
				mockAccountsService.
					On("GetAccountByUserID", mock.Anything, mock.AnythingOfType("models.UserID")).
					Return(&models.Account{}, nil).
					On("GetWithdrawals", mock.Anything, &models.Account{}).
					Return(withdrawalsList, nil)

				return &service.Services{
					AccountsService: mockAccountsService,
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
				expectedStatusCode:   http.StatusOK,
				expectedResponseBody: "[{\"order\":\"2377225624\",\"sum\":500,\"processed_at\":\"1984-04-05T06:07:08Z\"}]",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			services := test.mockServices()
			handler := NewHandler(services)
			router := handler.NewRouter()

			req, err := http.NewRequestWithContext(test.args.ctx, test.method, "/api/user/withdrawals", nil)
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
