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

func TestHandler_GetWithdrawals(t *testing.T) {
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
			name:   "500 — GetWithdrawals failed",
			method: http.MethodGet,
			args: args{
				mockServices: func() *service.Services {
					mockAccountsService := mockservice.NewAccountsService(t)
					mockAccountsService.
						On("GetWithdrawals", mock.Anything, models.UserId(1)).
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
					withdrawalsList := &models.WithdrawalsList{
						{
							OrderNum:    "2377225624",
							Sum:         500,
							ProcessedAt: time.Date(1984, time.April, 5, 6, 7, 8, 0, time.UTC),
						},
					}
					mockAccountsService := mockservice.NewAccountsService(t)
					mockAccountsService.
						On("GetWithdrawals", mock.Anything, models.UserId(1)).
						Return(withdrawalsList, nil)

					return &service.Services{
						Accounts: mockAccountsService,
					}
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
			services := test.args.mockServices()
			handler := NewHandler(services)
			router := handler.NewRouter()

			req, err := http.NewRequest(test.method, "/api/user/withdrawals", nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			require.Equal(t, test.want.expectedStatusCode, rr.Code)
			require.Equal(t, test.want.expectedResponseBody, rr.Body.String())
		})
	}
}
