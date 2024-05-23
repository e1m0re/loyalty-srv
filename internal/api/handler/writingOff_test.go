package handler

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"e1m0re/loyalty-srv/internal/models"
	"e1m0re/loyalty-srv/internal/service"
	mockservice "e1m0re/loyalty-srv/internal/service/mocks"
)

func TestHandler_WritingOff(t *testing.T) {
	type args struct {
		inputBody     string
		inputUserInfo models.UserInfo
		mockUser      *models.User
		mockServices  func() *service.Services
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
			name:   "400 — empty body",
			method: http.MethodPost,
			args: args{
				inputBody: "",
				mockServices: func() *service.Services {

					return &service.Services{}
				},
			},
			want: want{
				expectedStatusCode:   http.StatusBadRequest,
				expectedResponseBody: "",
			},
		},
		{
			name:   "500 — ValidateNumber failed",
			method: http.MethodPost,
			args: args{
				inputBody: "{\"order\":\"2377225624\",\"sum\":751}",
				mockServices: func() *service.Services {
					mockOrdersService := mockservice.NewOrdersService(t)
					mockOrdersService.
						On("ValidateNumber", mock.Anything, models.OrderNum("2377225624")).
						Return(false, fmt.Errorf("some error"))

					return &service.Services{
						Orders: mockOrdersService,
					}
				},
			},
			want: want{
				expectedStatusCode:   http.StatusInternalServerError,
				expectedResponseBody: "",
			},
		},
		{
			name:   "422 — invalid order number",
			method: http.MethodPost,
			args: args{
				inputBody: "{\"order\":\"2377225624\",\"sum\":751}",
				mockServices: func() *service.Services {
					mockOrdersService := mockservice.NewOrdersService(t)
					mockOrdersService.
						On("ValidateNumber", mock.Anything, models.OrderNum("2377225624")).
						Return(false, nil)

					return &service.Services{
						Orders: mockOrdersService,
					}
				},
			},
			want: want{
				expectedStatusCode:   http.StatusUnprocessableEntity,
				expectedResponseBody: "",
			},
		},
		{
			name:   "500 — GetAccountByUserId failed",
			method: http.MethodPost,
			args: args{
				inputBody: "{\"order\":\"2377225624\",\"sum\":751}",
				mockServices: func() *service.Services {
					mockOrdersService := mockservice.NewOrdersService(t)
					mockOrdersService.
						On("ValidateNumber", mock.Anything, models.OrderNum("2377225624")).
						Return(true, nil)

					mockAccountsService := mockservice.NewAccountsService(t)
					mockAccountsService.
						On("GetAccountByUserId", mock.Anything, models.UserId(1)).
						Return(nil, fmt.Errorf("some error"))

					return &service.Services{
						Orders:   mockOrdersService,
						Accounts: mockAccountsService,
					}
				},
			},
			want: want{
				expectedStatusCode:   http.StatusInternalServerError,
				expectedResponseBody: "some error",
			},
		},
		{
			name:   "500 — Withdraw failed",
			method: http.MethodPost,
			args: args{
				inputBody: "{\"order\":\"2377225624\",\"sum\":751}",
				mockServices: func() *service.Services {
					mockOrdersService := mockservice.NewOrdersService(t)
					mockOrdersService.
						On("ValidateNumber", mock.Anything, models.OrderNum("2377225624")).
						Return(true, nil)

					stubAccount := &models.Account{ID: models.AccountId(1)}
					mockAccountsService := mockservice.NewAccountsService(t)
					mockAccountsService.
						On("GetAccountByUserId", mock.Anything, models.UserId(1)).
						Return(stubAccount, nil).
						On("Withdraw", mock.Anything, stubAccount.ID, 751, models.OrderNum("2377225624")).
						Return(nil, fmt.Errorf("some error"))

					return &service.Services{
						Orders:   mockOrdersService,
						Accounts: mockAccountsService,
					}
				},
			},
			want: want{
				expectedStatusCode:   http.StatusInternalServerError,
				expectedResponseBody: "some error",
			},
		},
		{
			name:   "200",
			method: http.MethodPost,
			args: args{
				inputBody: "{\"order\":\"2377225624\",\"sum\":751}",
				mockServices: func() *service.Services {
					mockOrdersService := mockservice.NewOrdersService(t)
					mockOrdersService.
						On("ValidateNumber", mock.Anything, models.OrderNum("2377225624")).
						Return(true, nil)

					stubAccount := &models.Account{ID: models.AccountId(1)}
					mockAccountsService := mockservice.NewAccountsService(t)
					mockAccountsService.
						On("GetAccountByUserId", mock.Anything, models.UserId(1)).
						Return(stubAccount, nil).
						On("Withdraw", mock.Anything, stubAccount.ID, 751, models.OrderNum("2377225624")).
						Return(nil, nil)

					return &service.Services{
						Orders:   mockOrdersService,
						Accounts: mockAccountsService,
					}
				},
			},
			want: want{
				expectedStatusCode:   http.StatusOK,
				expectedResponseBody: "",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			services := test.args.mockServices()
			handler := NewHandler(services)
			router := handler.NewRouter()

			req, err := http.NewRequest(test.method, "/api/user/balance/withdraw", bytes.NewReader([]byte(test.args.inputBody)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			require.Equal(t, test.want.expectedStatusCode, rr.Code)
			require.Equal(t, test.want.expectedResponseBody, rr.Body.String())
		})
	}
}
