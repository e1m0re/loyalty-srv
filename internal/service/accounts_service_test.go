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

func Test_accountsService_CreateAccount(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID models.UserID
	}
	type want struct {
		account *models.Account
		errMsg  string
	}
	tests := []struct {
		name             string
		mockRepositories func() *repository.Repositories
		args             args
		want             want
	}{
		{
			name: "AddAccount failed",
			mockRepositories: func() *repository.Repositories {
				mockAccountRepository := mocks.NewAccountRepository(t)
				mockAccountRepository.
					On("AddAccount", mock.Anything, mock.AnythingOfType("models.UserID")).
					Return(nil, fmt.Errorf("some repos error"))

				return &repository.Repositories{
					AccountRepository: mockAccountRepository,
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
				mockAccountRepository := mocks.NewAccountRepository(t)
				mockAccountRepository.
					On("AddAccount", mock.Anything, mock.AnythingOfType("models.UserID")).
					Return(&models.Account{
						ID:      1,
						UserID:  1,
						Balance: 0,
					}, nil)

				return &repository.Repositories{
					AccountRepository: mockAccountRepository,
				}
			},
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			want: want{
				account: &models.Account{
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
			as := accountsService{
				accountRepository: repo.AccountRepository,
			}
			gotAccount, gotErr := as.CreateAccount(test.args.ctx, test.args.userID)
			require.Equal(t, &test.want.account, &gotAccount)
			if len(test.want.errMsg) > 0 {
				require.EqualError(t, gotErr, test.want.errMsg)
			}
		})
	}
}

func Test_accountsService_GetAccountByUserID(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID models.UserID
	}
	type want struct {
		account *models.Account
		errMsg  string
	}
	tests := []struct {
		name             string
		mockRepositories func() *repository.Repositories
		args             args
		want             want
	}{
		{
			name: "GetAccountByUserID failed",
			mockRepositories: func() *repository.Repositories {
				mockAccountRepository := mocks.NewAccountRepository(t)
				mockAccountRepository.
					On("GetAccountByUserID", mock.Anything, mock.AnythingOfType("models.UserID")).
					Return(nil, fmt.Errorf("some repos error"))

				return &repository.Repositories{
					AccountRepository: mockAccountRepository,
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
				mockAccountRepository := mocks.NewAccountRepository(t)
				mockAccountRepository.
					On("GetAccountByUserID", mock.Anything, mock.AnythingOfType("models.UserID")).
					Return(&models.Account{
						ID:      1,
						UserID:  1,
						Balance: 0,
					}, nil)

				return &repository.Repositories{
					AccountRepository: mockAccountRepository,
				}
			},
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			want: want{
				account: &models.Account{
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
			as := accountsService{
				accountRepository: repo.AccountRepository,
			}
			gotAccount, gotErr := as.GetAccountByUserID(test.args.ctx, test.args.userID)
			require.Equal(t, &test.want.account, &gotAccount)
			if len(test.want.errMsg) > 0 {
				require.EqualError(t, gotErr, test.want.errMsg)
			}
		})
	}
}

func Test_accountsService_GetAccountInfo(t *testing.T) {
	type args struct {
		ctx     context.Context
		account *models.Account
	}
	type want struct {
		accountInfo *models.AccountInfo
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
				mockAccountRepository := mocks.NewAccountRepository(t)
				mockAccountRepository.
					On("GetWithdrawnTotalSum", mock.Anything, mock.AnythingOfType("models.AccountID")).
					Return(100, nil)

				return &repository.Repositories{
					AccountRepository: mockAccountRepository,
				}
			},
			args: args{
				ctx:     context.Background(),
				account: &models.Account{ID: 1, Balance: 777},
			},
			want: want{
				accountInfo: &models.AccountInfo{
					CurrentBalance: 777,
					Withdrawals:    100,
				},
				errMsg: "",
			},
		},
		{
			name: "GetWithdrawnTotalSum failed",
			mockRepositories: func() *repository.Repositories {
				mockAccountRepository := mocks.NewAccountRepository(t)
				mockAccountRepository.
					On("GetWithdrawnTotalSum", mock.Anything, mock.AnythingOfType("models.AccountID")).
					Return(0, fmt.Errorf("some repos error"))

				return &repository.Repositories{
					AccountRepository: mockAccountRepository,
				}
			},
			args: args{
				ctx:     context.Background(),
				account: &models.Account{ID: 1},
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
			as := accountsService{
				accountRepository: repo.AccountRepository,
			}
			gotAccountInfo, gotErr := as.GetAccountInfo(test.args.ctx, test.args.account)
			require.Equal(t, &test.want.accountInfo, &gotAccountInfo)
			if len(test.want.errMsg) > 0 {
				require.EqualError(t, gotErr, test.want.errMsg)
			}
		})
	}
}

func Test_accountsService_GetWithdrawals(t *testing.T) {
	type args struct {
		ctx     context.Context
		account *models.Account
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
			name: "GetAccountByUserID failed",
			mockRepositories: func() *repository.Repositories {
				mockAccountRepository := mocks.NewAccountRepository(t)
				mockAccountRepository.
					On("GetWithdrawalsList", mock.Anything, mock.AnythingOfType("models.AccountID")).
					Return(nil, fmt.Errorf("some repos error"))

				return &repository.Repositories{
					AccountRepository: mockAccountRepository,
				}
			},
			args: args{
				ctx:     context.Background(),
				account: &models.Account{},
			},
			want: want{
				withdrawalsList: nil,
				errMsg:          "some repos error",
			},
		},
		{
			name: "Successfully case",
			mockRepositories: func() *repository.Repositories {
				mockAccountRepository := mocks.NewAccountRepository(t)
				mockAccountRepository.
					On("GetWithdrawalsList", mock.Anything, mock.AnythingOfType("models.AccountID")).
					Return(&models.WithdrawalsList{
						{
							OrderNum:    "2377225624",
							Sum:         500,
							ProcessedAt: time.Date(1703, time.May, 27, 12, 0, 0, 0, time.UTC),
						},
					}, nil)

				return &repository.Repositories{
					AccountRepository: mockAccountRepository,
				}
			},
			args: args{
				ctx:     context.Background(),
				account: &models.Account{},
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
			as := accountsService{
				accountRepository: repo.AccountRepository,
			}
			gotWithdrawalsList, gotErr := as.GetWithdrawals(test.args.ctx, test.args.account)
			require.Equal(t, &test.want.withdrawalsList, &gotWithdrawalsList)
			if len(test.want.errMsg) > 0 {
				require.EqualError(t, gotErr, test.want.errMsg)
			}
		})
	}
}

func Test_accountsService_UpdateBalance(t *testing.T) {
	type args struct {
		ctx      context.Context
		account  models.Account
		amount   float64
		orderNum models.OrderNum
	}
	type want struct {
		account *models.Account
		errMsg  string
	}
	tests := []struct {
		name             string
		mockRepositories func() *repository.Repositories
		args             args
		want             want
	}{
		{
			name: "ErrAccountHasNotEnoughFunds",
			mockRepositories: func() *repository.Repositories {

				return &repository.Repositories{}
			},
			args: args{
				ctx: context.Background(),
				account: models.Account{
					Balance: 10,
				},
				amount:   100,
				orderNum: "123",
			},
			want: want{
				account: nil,
				errMsg:  apperrors.ErrAccountHasNotEnoughFunds.Error(),
			},
		},
		{
			name: "AddAccountChange failed",
			mockRepositories: func() *repository.Repositories {
				mockAccountRepository := mocks.NewAccountRepository(t)
				mockAccountRepository.
					On("AddAccountChange", mock.Anything, mock.AnythingOfType("models.AccountID"), mock.AnythingOfType("float64"), mock.Anything).
					Return(nil, fmt.Errorf("some repos error"))

				return &repository.Repositories{
					AccountRepository: mockAccountRepository,
				}
			},
			args: args{
				ctx: context.Background(),
				account: models.Account{
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
				mockAccountRepository := mocks.NewAccountRepository(t)
				mockAccountRepository.
					On("AddAccountChange", mock.Anything, mock.AnythingOfType("models.AccountID"), mock.AnythingOfType("float64"), mock.Anything).
					Return(nil, nil).
					On("UpdateBalance", mock.Anything, mock.AnythingOfType("models.Account"), mock.AnythingOfType("float64")).
					Return(nil, fmt.Errorf("some repos error"))

				return &repository.Repositories{
					AccountRepository: mockAccountRepository,
				}
			},
			args: args{
				ctx: context.Background(),
				account: models.Account{
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
				mockAccountRepository := mocks.NewAccountRepository(t)
				mockAccountRepository.
					On("AddAccountChange", mock.Anything, mock.AnythingOfType("models.AccountID"), mock.AnythingOfType("float64"), mock.Anything).
					Return(nil, nil).
					On("UpdateBalance", mock.Anything, mock.AnythingOfType("models.Account"), mock.AnythingOfType("float64")).
					Return(&models.Account{
						Balance: 900,
					}, nil)

				return &repository.Repositories{
					AccountRepository: mockAccountRepository,
				}
			},
			args: args{
				ctx: context.Background(),
				account: models.Account{
					Balance: 1000,
				},
				amount:   100,
				orderNum: "123",
			},
			want: want{
				account: &models.Account{
					Balance: 900,
				},
				errMsg: "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := test.mockRepositories()
			as := accountsService{
				accountRepository: repo.AccountRepository,
			}
			gotAccount, gotErr := as.UpdateBalance(test.args.ctx, test.args.account, test.args.amount, test.args.orderNum)
			require.Equal(t, &test.want.account, &gotAccount)
			if len(test.want.errMsg) > 0 {
				require.EqualError(t, gotErr, test.want.errMsg)
			}
		})
	}
}
