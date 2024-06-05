package handler

import (
	"bytes"
	"context"
	"e1m0re/loyalty-srv/internal/apperrors"
	"fmt"
	"github.com/go-chi/jwtauth/v5"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"e1m0re/loyalty-srv/internal/models"
	"e1m0re/loyalty-srv/internal/service"
	mockservice "e1m0re/loyalty-srv/internal/service/mocks"
)

func TestHandler_AddOrder(t *testing.T) {
	jwtAuth := jwtauth.New("HS256", []byte("very secret key"), nil)
	type args struct {
		inputBody     string
		headers       map[string]string
		inputUserInfo models.UserInfo
		ctx           context.Context
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
			name:   "400 — empty body",
			method: http.MethodPost,
			mockServices: func() *service.Services {
				mockSecurityService := mockservice.NewSecurityService(t)
				mockSecurityService.
					On("GenerateAuthToken").
					Return(jwtAuth)

				mockOrdersService := mockservice.NewOrdersService(t)

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
				inputBody: "",
			},
			want: want{
				expectedStatusCode:   http.StatusBadRequest,
				expectedResponseBody: "",
			},
		},
		{
			name:   "422",
			method: http.MethodPost,
			mockServices: func() *service.Services {
				mockSecurityService := mockservice.NewSecurityService(t)
				mockSecurityService.
					On("GenerateAuthToken").
					Return(jwtAuth)

				mockOrdersService := mockservice.NewOrdersService(t)
				mockOrdersService.
					On("NewOrder", mock.Anything, mock.AnythingOfType("models.OrderInfo")).
					Return(nil, apperrors.ErrInvalidOrderNumber)

				return &service.Services{
					OrdersService:   mockOrdersService,
					SecurityService: mockSecurityService,
				}
			},
			args: args{
				inputBody: "1984",
				ctx:       context.WithValue(context.Background(), models.CKUserID, 1),
				headers: map[string]string{
					"Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6OSwidXNlcm5hbWUiOiJ1c2VyMiJ9.vY8OSC5qvDO-rLLnTUBGevkjIUm2oAjBuSsV75LO1Yw",
				},
			},
			want: want{
				expectedStatusCode:   http.StatusUnprocessableEntity,
				expectedResponseBody: "",
			},
		},
		{
			name:   "409",
			method: http.MethodPost,
			mockServices: func() *service.Services {
				mockSecurityService := mockservice.NewSecurityService(t)
				mockSecurityService.
					On("GenerateAuthToken").
					Return(jwtAuth)

				mockOrdersService := mockservice.NewOrdersService(t)
				mockOrdersService.
					On("NewOrder", mock.Anything, mock.AnythingOfType("models.OrderInfo")).
					Return(nil, apperrors.ErrOrderWasLoadedByAnotherUser)

				return &service.Services{
					OrdersService:   mockOrdersService,
					SecurityService: mockSecurityService,
				}
			},
			args: args{
				inputBody: "1984",
				ctx:       context.WithValue(context.Background(), models.CKUserID, 1),
				headers: map[string]string{
					"Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6OSwidXNlcm5hbWUiOiJ1c2VyMiJ9.vY8OSC5qvDO-rLLnTUBGevkjIUm2oAjBuSsV75LO1Yw",
				},
			},
			want: want{
				expectedStatusCode:   http.StatusConflict,
				expectedResponseBody: "",
			},
		},
		{
			name:   "500",
			method: http.MethodPost,
			mockServices: func() *service.Services {
				mockSecurityService := mockservice.NewSecurityService(t)
				mockSecurityService.
					On("GenerateAuthToken").
					Return(jwtAuth)

				mockOrdersService := mockservice.NewOrdersService(t)
				mockOrdersService.
					On("NewOrder", mock.Anything, mock.AnythingOfType("models.OrderInfo")).
					Return(nil, fmt.Errorf("some error"))

				return &service.Services{
					OrdersService:   mockOrdersService,
					SecurityService: mockSecurityService,
				}
			},
			args: args{
				inputBody: "1984",
				ctx:       context.WithValue(context.Background(), models.CKUserID, 1),
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
			name:   "202",
			method: http.MethodPost,
			mockServices: func() *service.Services {
				mockSecurityService := mockservice.NewSecurityService(t)
				mockSecurityService.
					On("GenerateAuthToken").
					Return(jwtAuth)

				mockOrdersService := mockservice.NewOrdersService(t)
				mockOrdersService.
					On("NewOrder", mock.Anything, mock.AnythingOfType("models.OrderInfo")).
					Return(nil, nil)

				return &service.Services{
					OrdersService:   mockOrdersService,
					SecurityService: mockSecurityService,
				}
			},
			args: args{
				inputBody: "12345678903",
				ctx:       context.WithValue(context.Background(), models.CKUserID, 1),
				headers: map[string]string{
					"Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6OSwidXNlcm5hbWUiOiJ1c2VyMiJ9.vY8OSC5qvDO-rLLnTUBGevkjIUm2oAjBuSsV75LO1Yw",
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
			mockServices: func() *service.Services {
				mockSecurityService := mockservice.NewSecurityService(t)
				mockSecurityService.
					On("GenerateAuthToken").
					Return(jwtAuth)

				mockOrdersService := mockservice.NewOrdersService(t)
				mockOrdersService.
					On("NewOrder", mock.Anything, mock.AnythingOfType("models.OrderInfo")).
					Return(nil, apperrors.ErrOrderWasLoaded)

				return &service.Services{
					OrdersService:   mockOrdersService,
					SecurityService: mockSecurityService,
				}
			},
			args: args{
				inputBody: "12345678903",
				ctx:       context.WithValue(context.Background(), models.CKUserID, 1),
				headers: map[string]string{
					"Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6OSwidXNlcm5hbWUiOiJ1c2VyMiJ9.vY8OSC5qvDO-rLLnTUBGevkjIUm2oAjBuSsV75LO1Yw",
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

			services := test.mockServices()
			handler := NewHandler(services)
			router := handler.NewRouter()

			req, err := http.NewRequestWithContext(test.args.ctx, test.method, "/api/user/orders", bytes.NewReader([]byte(test.args.inputBody)))
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
