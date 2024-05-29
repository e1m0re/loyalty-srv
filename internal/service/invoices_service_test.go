package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"e1m0re/loyalty-srv/internal/apperrors"
	"e1m0re/loyalty-srv/internal/models"
	"e1m0re/loyalty-srv/internal/repository"
	"e1m0re/loyalty-srv/internal/repository/mocks"
)

func Test_invoicesService_CreateAccount(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID models.UserID
	}
	type want struct {
		account *models.Invoice
		errMsg  string
	}
	tests := []struct {
		name             string
		mockRepositories func() *repository.Repositories
		args             args
		want             want
	}{
		{
			name: "AddInvoice failed",
			mockRepositories: func() *repository.Repositories {
				mockAccountRepository := mocks.NewInvoiceRepository(t)
				mockAccountRepository.
					On("AddInvoice", mock.Anything, mock.AnythingOfType("models.UserID")).
					Return(nil, fmt.Errorf("some repos error"))

				return &repository.Repositories{
					InvoiceRepository: mockAccountRepository,
				}
			},
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			want: want{
				account: nil,
				errMsg:  "some repos error",
			},
		},
		{
			name: "Successfully case",
			mockRepositories: func() *repository.Repositories {
				mockAccountRepository := mocks.NewInvoiceRepository(t)
				mockAccountRepository.
					On("AddInvoice", mock.Anything, mock.AnythingOfType("models.UserID")).
					Return(&models.Invoice{
						ID:      1,
						UserID:  1,
						Balance: 0,
					}, nil)

				return &repository.Repositories{
					InvoiceRepository: mockAccountRepository,
				}
			},
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			want: want{
				account: &models.Invoice{
					ID:      1,
					UserID:  1,
					Balance: 0,
				},
				errMsg: "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := test.mockRepositories()
			as := invoicesService{
				invoiceRepository: repo.InvoiceRepository,
			}
			gotAccount, gotErr := as.CreateInvoice(test.args.ctx, test.args.userID)
			require.Equal(t, &test.want.account, &gotAccount)
			if len(test.want.errMsg) > 0 {
				require.EqualError(t, gotErr, test.want.errMsg)
			}
		})
	}
}

func Test_invoicesService_GetAccountByUserID(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID models.UserID
	}
	type want struct {
		account *models.Invoice
		errMsg  string
	}
	tests := []struct {
		name             string
		mockRepositories func() *repository.Repositories
		args             args
		want             want
	}{
		{
			name: "GetInvoiceByUserID failed",
			mockRepositories: func() *repository.Repositories {
				mockAccountRepository := mocks.NewInvoiceRepository(t)
				mockAccountRepository.
					On("GetInvoiceByUserID", mock.Anything, mock.AnythingOfType("models.UserID")).
					Return(nil, fmt.Errorf("some repos error"))

				return &repository.Repositories{
					InvoiceRepository: mockAccountRepository,
				}
			},
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			want: want{
				account: nil,
				errMsg:  "some repos error",
			},
		},
		{
			name: "Successfully case",
			mockRepositories: func() *repository.Repositories {
				mockAccountRepository := mocks.NewInvoiceRepository(t)
				mockAccountRepository.
					On("GetInvoiceByUserID", mock.Anything, mock.AnythingOfType("models.UserID")).
					Return(&models.Invoice{
						ID:      1,
						UserID:  1,
						Balance: 0,
					}, nil)

				return &repository.Repositories{
					InvoiceRepository: mockAccountRepository,
				}
			},
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			want: want{
				account: &models.Invoice{
					ID:      1,
					UserID:  1,
					Balance: 0,
				},
				errMsg: "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := test.mockRepositories()
			as := invoicesService{
				invoiceRepository: repo.InvoiceRepository,
			}
			gotAccount, gotErr := as.GetInvoiceByUserID(test.args.ctx, test.args.userID)
			require.Equal(t, &test.want.account, &gotAccount)
			if len(test.want.errMsg) > 0 {
				require.EqualError(t, gotErr, test.want.errMsg)
			}
		})
	}
}

func Test_invoicesService_GetAccountInfo(t *testing.T) {
	type args struct {
		ctx     context.Context
		account *models.Invoice
	}
	type want struct {
		accountInfo *models.InvoiceInfo
		errMsg      string
	}
	tests := []struct {
		name             string
		mockRepositories func() *repository.Repositories
		args             args
		want             want
	}{
		{
			name: "Successfully case",
			mockRepositories: func() *repository.Repositories {
				mockAccountRepository := mocks.NewInvoiceRepository(t)
				mockAccountRepository.
					On("GetWithdrawnTotalSum", mock.Anything, mock.AnythingOfType("models.InvoiceID")).
					Return(100, nil)

				return &repository.Repositories{
					InvoiceRepository: mockAccountRepository,
				}
			},
			args: args{
				ctx:     context.Background(),
				account: &models.Invoice{ID: 1, Balance: 777},
			},
			want: want{
				accountInfo: &models.InvoiceInfo{
					CurrentBalance: 777,
					Withdrawals:    100,
				},
				errMsg: "",
			},
		},
		{
			name: "GetWithdrawnTotalSum failed",
			mockRepositories: func() *repository.Repositories {
				mockAccountRepository := mocks.NewInvoiceRepository(t)
				mockAccountRepository.
					On("GetWithdrawnTotalSum", mock.Anything, mock.AnythingOfType("models.InvoiceID")).
					Return(0, fmt.Errorf("some repos error"))

				return &repository.Repositories{
					InvoiceRepository: mockAccountRepository,
				}
			},
			args: args{
				ctx:     context.Background(),
				account: &models.Invoice{ID: 1},
			},
			want: want{
				accountInfo: nil,
				errMsg:      "some repos error",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := test.mockRepositories()
			as := invoicesService{
				invoiceRepository: repo.InvoiceRepository,
			}
			gotAccountInfo, gotErr := as.GetInvoiceInfo(test.args.ctx, test.args.account)
			require.Equal(t, &test.want.accountInfo, &gotAccountInfo)
			if len(test.want.errMsg) > 0 {
				require.EqualError(t, gotErr, test.want.errMsg)
			}
		})
	}
}

func Test_invoicesService_GetWithdrawals(t *testing.T) {
	type args struct {
		ctx     context.Context
		account *models.Invoice
	}
	type want struct {
		withdrawalsList *models.WithdrawalsList
		errMsg          string
	}
	tests := []struct {
		name             string
		mockRepositories func() *repository.Repositories
		args             args
		want             want
	}{
		{
			name: "GetInvoiceByUserID failed",
			mockRepositories: func() *repository.Repositories {
				mockAccountRepository := mocks.NewInvoiceRepository(t)
				mockAccountRepository.
					On("GetWithdrawalsList", mock.Anything, mock.AnythingOfType("models.InvoiceID")).
					Return(nil, fmt.Errorf("some repos error"))

				return &repository.Repositories{
					InvoiceRepository: mockAccountRepository,
				}
			},
			args: args{
				ctx:     context.Background(),
				account: &models.Invoice{},
			},
			want: want{
				withdrawalsList: nil,
				errMsg:          "some repos error",
			},
		},
		{
			name: "Successfully case",
			mockRepositories: func() *repository.Repositories {
				mockAccountRepository := mocks.NewInvoiceRepository(t)
				mockAccountRepository.
					On("GetWithdrawalsList", mock.Anything, mock.AnythingOfType("models.InvoiceID")).
					Return(&models.WithdrawalsList{
						{
							OrderNum:    "2377225624",
							Sum:         500,
							ProcessedAt: time.Date(1703, time.May, 27, 12, 0, 0, 0, time.UTC),
						},
					}, nil)

				return &repository.Repositories{
					InvoiceRepository: mockAccountRepository,
				}
			},
			args: args{
				ctx:     context.Background(),
				account: &models.Invoice{},
			},
			want: want{
				withdrawalsList: &models.WithdrawalsList{
					{
						OrderNum:    "2377225624",
						Sum:         500,
						ProcessedAt: time.Date(1703, time.May, 27, 12, 0, 0, 0, time.UTC),
					},
				},
				errMsg: "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := test.mockRepositories()
			as := invoicesService{
				invoiceRepository: repo.InvoiceRepository,
			}
			gotWithdrawalsList, gotErr := as.GetWithdrawals(test.args.ctx, test.args.account)
			require.Equal(t, &test.want.withdrawalsList, &gotWithdrawalsList)
			if len(test.want.errMsg) > 0 {
				require.EqualError(t, gotErr, test.want.errMsg)
			}
		})
	}
}

func Test_invoicesService_UpdateBalance(t *testing.T) {
	type args struct {
		ctx      context.Context
		account  models.Invoice
		amount   float64
		orderNum models.OrderNum
	}
	type want struct {
		account *models.Invoice
		errMsg  string
	}
	tests := []struct {
		name             string
		mockRepositories func() *repository.Repositories
		args             args
		want             want
	}{
		{
			name: "ErrInvoiceHasNotEnoughFunds",
			mockRepositories: func() *repository.Repositories {

				return &repository.Repositories{}
			},
			args: args{
				ctx: context.Background(),
				account: models.Invoice{
					Balance: 10,
				},
				amount:   100,
				orderNum: "123",
			},
			want: want{
				account: nil,
				errMsg:  apperrors.ErrInvoiceHasNotEnoughFunds.Error(),
			},
		},
		{
			name: "AddInvoiceChange failed",
			mockRepositories: func() *repository.Repositories {
				mockAccountRepository := mocks.NewInvoiceRepository(t)
				mockAccountRepository.
					On("AddInvoiceChange", mock.Anything, mock.AnythingOfType("models.InvoiceID"), mock.AnythingOfType("float64"), mock.Anything).
					Return(nil, fmt.Errorf("some repos error"))

				return &repository.Repositories{
					InvoiceRepository: mockAccountRepository,
				}
			},
			args: args{
				ctx: context.Background(),
				account: models.Invoice{
					Balance: 1000,
				},
				amount:   100,
				orderNum: "123",
			},
			want: want{
				account: nil,
				errMsg:  "some repos error",
			},
		},
		{
			name: "UpdateAccount failed",
			mockRepositories: func() *repository.Repositories {
				mockAccountRepository := mocks.NewInvoiceRepository(t)
				mockAccountRepository.
					On("AddInvoiceChange", mock.Anything, mock.AnythingOfType("models.InvoiceID"), mock.AnythingOfType("float64"), mock.Anything).
					Return(nil, nil).
					On("UpdateBalance", mock.Anything, mock.AnythingOfType("models.Invoice"), mock.AnythingOfType("float64")).
					Return(nil, fmt.Errorf("some repos error"))

				return &repository.Repositories{
					InvoiceRepository: mockAccountRepository,
				}
			},
			args: args{
				ctx: context.Background(),
				account: models.Invoice{
					Balance: 1000,
				},
				amount:   100,
				orderNum: "123",
			},
			want: want{
				account: nil,
				errMsg:  "some repos error",
			},
		},
		{
			name: "Successfully case",
			mockRepositories: func() *repository.Repositories {
				mockAccountRepository := mocks.NewInvoiceRepository(t)
				mockAccountRepository.
					On("AddInvoiceChange", mock.Anything, mock.AnythingOfType("models.InvoiceID"), mock.AnythingOfType("float64"), mock.Anything).
					Return(nil, nil).
					On("UpdateBalance", mock.Anything, mock.AnythingOfType("models.Invoice"), mock.AnythingOfType("float64")).
					Return(&models.Invoice{
						Balance: 900,
					}, nil)

				return &repository.Repositories{
					InvoiceRepository: mockAccountRepository,
				}
			},
			args: args{
				ctx: context.Background(),
				account: models.Invoice{
					Balance: 1000,
				},
				amount:   100,
				orderNum: "123",
			},
			want: want{
				account: &models.Invoice{
					Balance: 900,
				},
				errMsg: "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := test.mockRepositories()
			as := invoicesService{
				invoiceRepository: repo.InvoiceRepository,
			}
			gotAccount, gotErr := as.UpdateBalance(test.args.ctx, test.args.account, test.args.amount, test.args.orderNum)
			require.Equal(t, &test.want.account, &gotAccount)
			if len(test.want.errMsg) > 0 {
				require.EqualError(t, gotErr, test.want.errMsg)
			}
		})
	}
}
