package handler

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"e1m0re/loyalty-srv/internal/apperrors"
	"e1m0re/loyalty-srv/internal/models"
	"e1m0re/loyalty-srv/internal/service"
	mockservice "e1m0re/loyalty-srv/internal/service/mocks"
)

func TestHandler_AddOrder(t *testing.T) {
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
			name:   "400",
			method: http.MethodPost,
			args: args{
				inputBody: "",
				mockServices: func() *service.Services {
					mockOrdersService := mockservice.NewOrdersService(t)
					return &service.Services{
						Orders: mockOrdersService,
					}
				},
			},
			want: want{
				expectedStatusCode:   http.StatusBadRequest,
				expectedResponseBody: "",
			},
		},
		{
			name:   "422",
			method: http.MethodPost,
			args: args{
				inputBody: "1984",
				mockServices: func() *service.Services {
					mockOrdersService := mockservice.NewOrdersService(t)
					mockOrdersService.
						On("NewOrder", mock.Anything, models.OrderNum("1984")).
						Return(nil, false, apperrors.InvalidOrderNumberError)
					return &service.Services{
						Orders: mockOrdersService,
					}
				},
			},
			want: want{
				expectedStatusCode:   http.StatusUnprocessableEntity,
				expectedResponseBody: "invalid order number",
			},
		},
		{
			name:   "409",
			method: http.MethodPost,
			args: args{
				inputBody: "1984",
				mockServices: func() *service.Services {
					mockOrdersService := mockservice.NewOrdersService(t)
					mockOrdersService.
						On("NewOrder", mock.Anything, models.OrderNum("1984")).
						Return(nil, false, apperrors.OtherUsersOrderError)
					return &service.Services{
						Orders: mockOrdersService,
					}
				},
			},
			want: want{
				expectedStatusCode:   http.StatusConflict,
				expectedResponseBody: "the order number has already been uploaded by another user",
			},
		},
		{
			name:   "500",
			method: http.MethodPost,
			args: args{
				inputBody: "1984",
				mockServices: func() *service.Services {
					mockOrdersService := mockservice.NewOrdersService(t)
					mockOrdersService.
						On("NewOrder", mock.Anything, models.OrderNum("1984")).
						Return(nil, false, fmt.Errorf("some error"))
					return &service.Services{
						Orders: mockOrdersService,
					}
				},
			},
			want: want{
				expectedStatusCode:   http.StatusInternalServerError,
				expectedResponseBody: "some error",
			},
		},
		{
			name:   "202",
			method: http.MethodPost,
			args: args{
				inputBody: "1984",
				mockServices: func() *service.Services {
					mockOrdersService := mockservice.NewOrdersService(t)
					mockOrdersService.
						On("NewOrder", mock.Anything, models.OrderNum("1984")).
						Return(nil, true, nil)
					return &service.Services{
						Orders: mockOrdersService,
					}
				},
			},
			want: want{
				expectedStatusCode:   http.StatusAccepted,
				expectedResponseBody: "",
			},
		},
		{
			name:   "200",
			method: http.MethodPost,
			args: args{
				inputBody: "1984",
				mockServices: func() *service.Services {
					mockOrdersService := mockservice.NewOrdersService(t)
					mockOrdersService.
						On("NewOrder", mock.Anything, models.OrderNum("1984")).
						Return(nil, false, nil)
					return &service.Services{
						Orders: mockOrdersService,
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

			req, err := http.NewRequest(test.method, "/api/user/orders", bytes.NewReader([]byte(test.args.inputBody)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			require.Equal(t, test.want.expectedStatusCode, rr.Code)
			require.Equal(t, test.want.expectedResponseBody, rr.Body.String())
		})
	}
}
