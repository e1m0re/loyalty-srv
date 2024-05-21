package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"e1m0re/loyalty-srv/internal/models"
	"e1m0re/loyalty-srv/internal/service"
	mock_service "e1m0re/loyalty-srv/internal/service/mocks"
)

func TestHandler_SignUp(t *testing.T) {
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
		name   string
		method string
		args   args
		want   want
	}{
		//{
		//	name:   "405",
		//	method: http.MethodGet,
		//	args: args{
		//		inputBody:     "",
		//		inputUserInfo: models.UserInfo{},
		//		mockUser:      &models.User{},
		//	},
		//	want: want{
		//		expectedStatusCode:   405,
		//		expectedResponseBody: "",
		//	},
		//},
		{
			name:   "400 — Invalid JSON body",
			method: http.MethodPost,
			args: args{
				inputBody: `{login:login,password:password}`,
				mockServices: func() service.Services {
					mockAuthorization := mock_service.NewAuthorization(t)

					return service.Services{
						Authorization: mockAuthorization,
					}
				},
			},
			want: want{
				expectedStatusCode:   http.StatusBadRequest,
				expectedResponseBody: "invalid character 'l' looking for beginning of object key string",
			},
		},
		{
			name:   "400 — Empty login and password",
			method: http.MethodPost,
			args: args{
				inputBody: `{"login":"","password":""}`,
				mockServices: func() service.Services {
					mockAuthorization := mock_service.NewAuthorization(t)

					return service.Services{
						Authorization: mockAuthorization,
					}
				},
			},
			want: want{
				expectedStatusCode:   http.StatusBadRequest,
				expectedResponseBody: "",
			},
		},
		{
			name:   "500 — FindUserByUsername failed",
			method: http.MethodPost,
			args: args{
				inputBody: `{"login":"test","password":"password"}`,
				mockServices: func() service.Services {
					mockAuthorization := mock_service.NewAuthorization(t)
					mockAuthorization.
						On("FindUserByUsername", mock.Anything, "test").
						Return(nil, errors.New("user not found"))

					return service.Services{
						Authorization: mockAuthorization,
					}
				},
			},
			want: want{
				expectedStatusCode:   http.StatusInternalServerError,
				expectedResponseBody: "",
			},
		},
		{
			name:   "409 — username busy",
			method: http.MethodPost,
			args: args{
				inputBody: `{"login":"test","password":"password"}`,
				mockServices: func() service.Services {
					mockAuthorization := mock_service.NewAuthorization(t)
					mockAuthorization.
						On("FindUserByUsername", mock.Anything, "test").
						Return(&models.User{ID: 1, Username: "test"}, nil)

					return service.Services{
						Authorization: mockAuthorization,
					}
				},
			},
			want: want{
				expectedStatusCode:   http.StatusConflict,
				expectedResponseBody: "",
			},
		},
		{
			name:   "500 — CreateUser failed",
			method: http.MethodPost,
			args: args{
				inputBody: `{"login":"test","password":"password"}`,
				mockServices: func() service.Services {
					mockAuthorization := mock_service.NewAuthorization(t)
					mockAuthorization.
						On("FindUserByUsername", mock.Anything, "test").
						Return(nil, nil).
						On("CreateUser", mock.Anything, models.UserInfo{Username: "test", Password: "password"}).
						Return(nil, errors.New("create failed"))

					return service.Services{
						Authorization: mockAuthorization,
					}
				},
			},
			want: want{
				expectedStatusCode:   http.StatusInternalServerError,
				expectedResponseBody: "",
			},
		},
		{
			name:   "500 — SignIn failed",
			method: http.MethodPost,
			args: args{
				inputBody: `{"login":"test","password":"password"}`,
				mockServices: func() service.Services {
					userInfo := models.UserInfo{Username: "test", Password: "password"}
					user := &models.User{ID: 1, Username: "test"}
					mockAuthorization := mock_service.NewAuthorization(t)
					mockAuthorization.
						On("FindUserByUsername", mock.Anything, "test").
						Return(nil, nil).
						On("CreateUser", mock.Anything, userInfo).
						Return(user, nil).
						On("SignIn", mock.Anything, userInfo).
						Return(false, errors.New("signin failed"))

					return service.Services{
						Authorization: mockAuthorization,
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
			method: http.MethodPost,
			args: args{
				inputBody: `{"login":"test","password":"password"}`,
				mockServices: func() service.Services {
					userInfo := models.UserInfo{Username: "test", Password: "password"}
					user := &models.User{ID: 1, Username: "test"}
					mockAuthorization := mock_service.NewAuthorization(t)
					mockAuthorization.
						On("FindUserByUsername", mock.Anything, "test").
						Return(nil, nil).
						On("CreateUser", mock.Anything, userInfo).
						Return(user, nil).
						On("SignIn", mock.Anything, userInfo).
						Return(true, nil)

					return service.Services{
						Authorization: mockAuthorization,
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
			handler := NewHandler(&services)
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
