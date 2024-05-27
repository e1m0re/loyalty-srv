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

func (as accountsService) GetAccountByUserID(ctx context.Context, id models.UserID) (*models.Account, error) {
	return &models.Account{}, nil
}

func (as accountsService) GetAccountInfoByUserID(ctx context.Context, id models.UserID) (*models.AccountInfo, error) {
	return &models.AccountInfo{}, nil
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
