package service

import (
	"context"

	"e1m0re/loyalty-srv/internal/models"
)

type accountsService struct {
}

func NewAccountsService() AccountsService {
	return &accountsService{}
}

func (as accountsService) GetAccountByUserId(ctx context.Context, id models.UserId) (models.Account, error) {
	return models.Account{}, nil
}

func (as accountsService) GetAccountInfoByUserId(ctx context.Context, id models.UserId) (models.AccountInfo, error) {
	return models.AccountInfo{}, nil
}

func (as accountsService) Withdraw(ctx context.Context, id models.AccountId, amount int, orderNum models.OrderNum) (models.Account, error) {
	return models.Account{}, nil
}

func (as accountsService) GetWithdrawals(ctx context.Context, id models.UserId) (models.WithdrawalsList, error) {
	return models.WithdrawalsList{}, nil
}

func (as accountsService) UpdateBalance(ctx context.Context, id models.AccountId, amount int) (models.Account, error) {
	return models.Account{}, nil
}
