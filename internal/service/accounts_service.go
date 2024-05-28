package service

import (
	"context"

	"e1m0re/loyalty-srv/internal/apperrors"
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

func (as accountsService) Withdraw(ctx context.Context, account models.Account, amount float64, orderNum models.OrderNum) (*models.Account, error) {
	if account.Balance < amount {
		return nil, apperrors.ErrAccountHasNotEnoughFunds
	}

	_, err := as.accountRepository.AddAccountChange(ctx, account.ID, amount, orderNum)
	if err != nil {
		return nil, err
	}

	account.Balance = account.Balance - amount
	err = as.accountRepository.UpdateAccount(ctx, &account)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (as accountsService) GetWithdrawals(ctx context.Context, account *models.Account) (*models.WithdrawalsList, error) {
	return as.accountRepository.GetWithdrawalsList(ctx, account.ID)
}

func (as accountsService) UpdateBalance(ctx context.Context, id models.AccountID, amount int) (*models.Account, error) {
	return &models.Account{}, nil
}
