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

	"e1m0re/loyalty-srv/internal/service"
	mockservice "e1m0re/loyalty-srv/internal/service/mocks"
)

func TestHandler_SignIn(t *testing.T) {
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
			name:   "400 — Invalid body",
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
				inputBody: "{\"login\":\"\",\"password\":\"}",
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
				inputBody: "{\"login\":\"\",\"password\":\"\"}",
			},
			want: want{
				expectedStatusCode:   http.StatusBadRequest,
				expectedResponseBody: "",
			},
		},
		{
			name:   "500 — SignIn failed",
			method: http.MethodPost,
			mockServices: func() *service.Services {
				mockSecurityService := mockservice.NewSecurityService(t)
				mockSecurityService.
					On("GenerateAuthToken").
					Return(jwtAuth)

				mockUsersService := mockservice.NewUsersService(t)
				mockUsersService.
					On("SignIn", mock.Anything, mock.AnythingOfType("models.UserInfo")).
					Return("", fmt.Errorf("some error"))

				return &service.Services{
					SecurityService: mockSecurityService,
					UsersService:    mockUsersService,
				}
			},
			args: args{
				inputBody: "{\"login\":\"login\",\"password\":\"password\"}",
			},
			want: want{
				expectedStatusCode:   http.StatusInternalServerError,
				expectedResponseBody: "",
			},
		},
		{
			name:   "401",
			method: http.MethodPost,
			mockServices: func() *service.Services {
				mockSecurityService := mockservice.NewSecurityService(t)
				mockSecurityService.
					On("GenerateAuthToken").
					Return(jwtAuth)

				mockUsersService := mockservice.NewUsersService(t)
				mockUsersService.
					On("SignIn", mock.Anything, mock.AnythingOfType("models.UserInfo")).
					Return("", nil)

				return &service.Services{
					SecurityService: mockSecurityService,
					UsersService:    mockUsersService,
				}
			},
			args: args{
				inputBody: "{\"login\":\"login\",\"password\":\"password\"}",
			},
			want: want{
				expectedStatusCode:   http.StatusUnauthorized,
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

				mockUsersService := mockservice.NewUsersService(t)
				mockUsersService.
					On("SignIn", mock.Anything, mock.AnythingOfType("models.UserInfo")).
					Return("json-token", nil)

				return &service.Services{
					SecurityService: mockSecurityService,
					UsersService:    mockUsersService,
				}
			},
			args: args{
				inputBody: "{\"login\":\"login\",\"password\":\"password\"}",
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

			req, err := http.NewRequest(test.method, "/api/user/login", bytes.NewReader([]byte(test.args.inputBody)))
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
