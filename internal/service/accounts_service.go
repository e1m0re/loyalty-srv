package service

import (
	"context"

	"e1m0re/loyalty-srv/internal/models"
	"e1m0re/loyalty-srv/internal/repository"
)

type accountsService struct {
	accountRepository repository.AccountRepository
}

func NewAccountsService(accountRepository repository.AccountRepository) AccountsService {
	return &accountsService{
		accountRepository: accountRepository,
	}
}

func (as accountsService) GetAccountByUserID(ctx context.Context, userID models.UserID) (*models.Account, error) {
	return as.accountRepository.GetAccountByUserID(ctx, userID)
}

func (as accountsService) GetAccountInfo(ctx context.Context, account *models.Account) (*models.AccountInfo, error) {
	withdrawnTotalSum, err := as.accountRepository.GetWithdrawnTotalSum(ctx, account.ID)
	if err != nil {
		return nil, err
	}

	return &models.AccountInfo{
		CurrentBalance: account.Balance,
		Withdrawals:    withdrawnTotalSum,
	}, nil
}

func (as accountsService) CreateAccount(ctx context.Context, id models.UserID) (*models.Account, error) {
	return as.accountRepository.AddAccount(ctx, id)
}

func (as accountsService) Withdraw(ctx context.Context, id models.AccountID, amount int, orderNum models.OrderNum) (*models.Account, error) {
	return &models.Account{}, nil
}

func (as accountsService) GetWithdrawals(ctx context.Context, id models.UserID) (*models.WithdrawalsList, error) {
	return &models.WithdrawalsList{}, nil
}

func (as accountsService) UpdateBalance(ctx context.Context, id models.AccountID, amount int) (*models.Account, error) {
	return &models.Account{}, nil
}
