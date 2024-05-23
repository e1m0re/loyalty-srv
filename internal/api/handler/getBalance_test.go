package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"e1m0re/loyalty-srv/internal/models"
	"e1m0re/loyalty-srv/internal/service"
	mockservice "e1m0re/loyalty-srv/internal/service/mocks"
)

func TestHandler_GetBalance(t *testing.T) {
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
			name:   "500 — GetAccountInfoByUserId failed",
			method: http.MethodGet,
			args: args{
				mockServices: func() *service.Services {
					mockAccountsService := mockservice.NewAccountsService(t)
					mockAccountsService.
						On("GetAccountInfoByUserId", mock.Anything, models.UserId(1)).
						Return(nil, errors.New("some error"))
					return &service.Services{
						Accounts: mockAccountsService,
					}
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
			args: args{
				mockServices: func() *service.Services {
					mockAccountsService := mockservice.NewAccountsService(t)
					mockAccountsService.
						On("GetAccountInfoByUserId", mock.Anything, models.UserId(1)).
						Return(&models.AccountInfo{
							CurrentBalance: 500.5,
							Withdrawals:    42,
						}, nil)
					return &service.Services{
						Accounts: mockAccountsService,
					}
				},
			},
			want: want{
				expectedStatusCode:   http.StatusOK,
				expectedResponseBody: "{\"current\":500.5,\"withdrawals\":42}",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			services := test.args.mockServices()
			handler := NewHandler(services)
			router := handler.NewRouter()

			req, err := http.NewRequest(test.method, "/api/user/balance", nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			require.Equal(t, test.want.expectedStatusCode, rr.Code)
			require.Equal(t, test.want.expectedResponseBody, rr.Body.String())
		})
	}
}
