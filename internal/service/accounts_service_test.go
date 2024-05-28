package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

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
					Return(nil, fmt.Errorf("some repo error"))

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
				errMsg:  "some repo error",
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
