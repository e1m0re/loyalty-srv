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

func TestHandler_SignIn(t *testing.T) {
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
			name:   "405",
			method: http.MethodGet,
			args: args{
				mockServices: func() *service.Services {

					return &service.Services{}
				},
			},
			want: want{
				expectedStatusCode:   http.StatusMethodNotAllowed,
				expectedResponseBody: "",
			},
		},
		{
			name:   "400 — Invalid body",
			method: http.MethodPost,
			args: args{
				inputBody: "{\"login\":\"\",\"password\":\"}",
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
			name:   "400 — Empty login and password",
			method: http.MethodPost,
			args: args{
				inputBody: "{\"login\":\"\",\"password\":\"\"}",
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
			name:   "500 — SignIn failed",
			method: http.MethodPost,
			args: args{
				inputBody: "{\"login\":\"login\",\"password\":\"password\"}",
				mockServices: func() *service.Services {
					userInfo := &models.UserInfo{
						Username: "login",
						Password: "password",
					}
					mockUsersService := mockservice.NewUsersService(t)
					mockUsersService.
						On("SignIn", mock.Anything, userInfo).
						Return(false, fmt.Errorf("some error"))

					return &service.Services{
						UsersService: mockUsersService,
					}
				},
			},
			want: want{
				expectedStatusCode:   http.StatusInternalServerError,
				expectedResponseBody: "some error",
			},
		},
		{
			name:   "401",
			method: http.MethodPost,
			args: args{
				inputBody: "{\"login\":\"login\",\"password\":\"password\"}",
				mockServices: func() *service.Services {
					userInfo := &models.UserInfo{
						Username: "login",
						Password: "password",
					}
					mockUsersService := mockservice.NewUsersService(t)
					mockUsersService.
						On("SignIn", mock.Anything, userInfo).
						Return(false, nil)

					return &service.Services{
						UsersService: mockUsersService,
					}
				},
			},
			want: want{
				expectedStatusCode:   http.StatusUnauthorized,
				expectedResponseBody: "",
			},
		},
		{
			name:   "200",
			method: http.MethodPost,
			args: args{
				inputBody: "{\"login\":\"login\",\"password\":\"password\"}",
				mockServices: func() *service.Services {
					userInfo := &models.UserInfo{
						Username: "login",
						Password: "password",
					}
					mockUsersService := mockservice.NewUsersService(t)
					mockUsersService.
						On("SignIn", mock.Anything, userInfo).
						Return(true, nil)

					return &service.Services{
						UsersService: mockUsersService,
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

			req, err := http.NewRequest(test.method, "/api/user/login", bytes.NewReader([]byte(test.args.inputBody)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			require.Equal(t, test.want.expectedStatusCode, rr.Code)
			require.Equal(t, test.want.expectedResponseBody, rr.Body.String())
		})
	}
}
