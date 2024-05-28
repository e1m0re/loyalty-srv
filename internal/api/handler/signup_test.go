package handler

import (
	"bytes"
	"errors"
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

func TestHandler_SignUp(t *testing.T) {
	jwtAuth := jwtauth.New("HS256", []byte("very secret key"), nil)
	type args struct {
		inputBody     string
		inputUserInfo models.UserInfo
		mockUser      *models.User
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
			args: args{},
			want: want{
				expectedStatusCode:   http.StatusMethodNotAllowed,
				expectedResponseBody: "",
			},
		},
		{
			name:   "400 — Invalid JSON body",
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
				inputBody: `{login:login,password:password}`,
			},
			want: want{
				expectedStatusCode:   http.StatusBadRequest,
				expectedResponseBody: "",
			},
		},
		{
			name:   "400 — Empty login and password",
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
				inputBody: `{"login":"","password":""}`,
			},
			want: want{
				expectedStatusCode:   http.StatusBadRequest,
				expectedResponseBody: "",
			},
		},
		{
			name:   "409 — username busy",
			method: http.MethodPost,
			mockServices: func() *service.Services {
				mockSecurityService := mockservice.NewSecurityService(t)
				mockSecurityService.
					On("GenerateAuthToken").
					Return(jwtAuth)

				mockUsersService := mockservice.NewUsersService(t)
				mockUsersService.
					On("CreateUser", mock.Anything, mock.AnythingOfType("models.UserInfo")).
					Return(nil, apperrors.ErrBusyLogin)

				return &service.Services{
					SecurityService: mockSecurityService,
					UsersService:    mockUsersService,
				}
			},
			args: args{
				inputBody: `{"login":"test","password":"password"}`,
			},
			want: want{
				expectedStatusCode:   http.StatusConflict,
				expectedResponseBody: "",
			},
		},
		{
			name:   "500 — CreateUser failed",
			method: http.MethodPost,
			mockServices: func() *service.Services {
				mockSecurityService := mockservice.NewSecurityService(t)
				mockSecurityService.
					On("GenerateAuthToken").
					Return(jwtAuth)

				mockUsersService := mockservice.NewUsersService(t)
				mockUsersService.
					On("CreateUser", mock.Anything, mock.AnythingOfType("models.UserInfo")).
					Return(nil, errors.New("create failed"))

				return &service.Services{
					SecurityService: mockSecurityService,
					UsersService:    mockUsersService,
				}
			},
			args: args{
				inputBody: `{"login":"test","password":"password"}`,
			},
			want: want{
				expectedStatusCode:   http.StatusInternalServerError,
				expectedResponseBody: "",
			},
		},
		{
			name:   "500 — CreateAccount failed",
			method: http.MethodPost,
			mockServices: func() *service.Services {
				mockAccountService := mockservice.NewAccountsService(t)
				mockAccountService.
					On("CreateAccount", mock.Anything, mock.AnythingOfType("models.UserID")).
					Return(nil, fmt.Errorf("some repo error"))

				mockSecurityService := mockservice.NewSecurityService(t)
				mockSecurityService.
					On("GenerateAuthToken").
					Return(jwtAuth)

				mockUsersService := mockservice.NewUsersService(t)
				mockUsersService.
					On("CreateUser", mock.Anything, mock.AnythingOfType("models.UserInfo")).
					Return(&models.User{}, nil)

				return &service.Services{
					AccountsService: mockAccountService,
					SecurityService: mockSecurityService,
					UsersService:    mockUsersService,
				}
			},
			args: args{
				inputBody: `{"login":"test","password":"password"}`,
			},
			want: want{
				expectedStatusCode:   http.StatusInternalServerError,
				expectedResponseBody: "",
			},
		},
		{
			name:   "500 — SignIn failed",
			method: http.MethodPost,
			mockServices: func() *service.Services {
				mockAccountService := mockservice.NewAccountsService(t)
				mockAccountService.
					On("CreateAccount", mock.Anything, mock.AnythingOfType("models.UserID")).
					Return(nil, nil)

				mockSecurityService := mockservice.NewSecurityService(t)
				mockSecurityService.
					On("GenerateAuthToken").
					Return(jwtAuth)

				mockUsersService := mockservice.NewUsersService(t)
				mockUsersService.
					On("CreateUser", mock.Anything, mock.AnythingOfType("models.UserInfo")).
					Return(&models.User{}, nil).
					On("SignIn", mock.Anything, mock.AnythingOfType("models.UserInfo")).
					Return("", errors.New("signin failed"))

				return &service.Services{
					AccountsService: mockAccountService,
					SecurityService: mockSecurityService,
					UsersService:    mockUsersService,
				}
			},
			args: args{
				inputBody: `{"login":"test","password":"password"}`,
			},
			want: want{
				expectedStatusCode:   http.StatusInternalServerError,
				expectedResponseBody: "",
			},
		},
		{
			name:   "200",
			method: http.MethodPost,
			mockServices: func() *service.Services {
				mockAccountService := mockservice.NewAccountsService(t)
				mockAccountService.
					On("CreateAccount", mock.Anything, mock.AnythingOfType("models.UserID")).
					Return(nil, nil)

				mockSecurityService := mockservice.NewSecurityService(t)
				mockSecurityService.
					On("GenerateAuthToken").
					Return(jwtAuth)

				mockUsersService := mockservice.NewUsersService(t)
				mockUsersService.
					On("CreateUser", mock.Anything, mock.AnythingOfType("models.UserInfo")).
					Return(&models.User{}, nil).
					On("SignIn", mock.Anything, mock.AnythingOfType("models.UserInfo")).
					Return("json-token", nil)

				return &service.Services{
					AccountsService: mockAccountService,
					SecurityService: mockSecurityService,
					UsersService:    mockUsersService,
				}
			},
			args: args{
				inputBody: `{"login":"test","password":"password"}`,
			},
			want: want{
				expectedStatusCode:   http.StatusOK,
				expectedResponseBody: "json-token",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			services := test.mockServices()
			handler := NewHandler(services)
			router := handler.NewRouter()

			req, err := http.NewRequest(test.method, "/api/user/register", bytes.NewReader([]byte(test.args.inputBody)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			require.Equal(t, test.want.expectedStatusCode, rr.Code)
			require.Equal(t, test.want.expectedResponseBody, rr.Body.String())
		})
	}
}
