package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/jwtauth/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"e1m0re/loyalty-srv/internal/models"
	"e1m0re/loyalty-srv/internal/service"
	mockservice "e1m0re/loyalty-srv/internal/service/mocks"
)

func TestHandler_GetBalance(t *testing.T) {
	jwtAuth := jwtauth.New("HS256", []byte("very secret key"), nil)
	type args struct {
		headers map[string]string
		ctx     context.Context
	}
	type want struct {
		expectedStatusCode   int
		expectedResponseBody string
	}
	tests := []struct {
		name         string
		mockServices func() *service.Services
		method       string
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
				ctx: context.WithValue(context.Background(), models.CKUserID, 1),
				headers: map[string]string{
					"Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6OSwidXNlcm5hbWUiOiJ1c2VyMiJ9.vY8OSC5qvDO-rLLnTUBGevkjIUm2oAjBuSsV75LO1Yw",
				},
			},
			want: want{
				expectedStatusCode:   http.StatusMethodNotAllowed,
				expectedResponseBody: "",
			},
		},
		{
			name:   "500 — GetInvoiceByUserID failed",
			method: http.MethodGet,
			mockServices: func() *service.Services {
				mockSecurityService := mockservice.NewSecurityService(t)
				mockSecurityService.
					On("GenerateAuthToken").
					Return(jwtAuth)

				mockInvoicesService := mockservice.NewInvoicesService(t)
				mockInvoicesService.
					On("GetInvoiceByUserID", mock.Anything, mock.AnythingOfType("models.UserID")).
					Return(nil, errors.New("some error"))

				return &service.Services{
					InvoicesService: mockInvoicesService,
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
			name:   "500 — GetInvoiceInfo failed",
			method: http.MethodGet,
			mockServices: func() *service.Services {
				mockSecurityService := mockservice.NewSecurityService(t)
				mockSecurityService.
					On("GenerateAuthToken").
					Return(jwtAuth)

				mockInvoicesService := mockservice.NewInvoicesService(t)
				mockInvoicesService.
					On("GetInvoiceByUserID", mock.Anything, mock.AnythingOfType("models.UserID")).
					Return(&models.Invoice{}, errors.New("some error"))

				return &service.Services{
					InvoicesService: mockInvoicesService,
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

				mockInvoicesService := mockservice.NewInvoicesService(t)
				mockInvoicesService.
					On("GetInvoiceByUserID", mock.Anything, mock.AnythingOfType("models.UserID")).
					Return(&models.Invoice{}, nil).
					On("GetInvoiceInfo", mock.Anything, &models.Invoice{}).
					Return(&models.InvoiceInfo{
						CurrentBalance: 500.5,
						Withdrawals:    42,
					}, nil)

				return &service.Services{
					InvoicesService: mockInvoicesService,
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
				expectedResponseBody: "{\"current\":500.5,\"withdrawals\":42}",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			services := test.mockServices()
			handler := NewHandler(services)
			router := handler.NewRouter()

			req, err := http.NewRequestWithContext(test.args.ctx, test.method, "/api/user/balance", nil)
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
