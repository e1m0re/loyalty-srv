package service

import (
	"context"

	"e1m0re/loyalty-srv/internal/apperrors"
	"e1m0re/loyalty-srv/internal/models"
	"e1m0re/loyalty-srv/internal/repository"
)

type invoicesService struct {
	accountRepository repository.InvoiceRepository
}

func NewInvoicesService(accountRepository repository.InvoiceRepository) InvoicesService {
	return &invoicesService{
		accountRepository: accountRepository,
	}
}

func (as invoicesService) GetInvoiceByUserID(ctx context.Context, userID models.UserID) (*models.Invoice, error) {
	return as.accountRepository.GetInvoiceByUserID(ctx, userID)
}

func (as invoicesService) GetInvoiceInfo(ctx context.Context, account *models.Invoice) (*models.InvoiceInfo, error) {
	withdrawnTotalSum, err := as.accountRepository.GetWithdrawnTotalSum(ctx, account.ID)
	if err != nil {
		return nil, err
	}

	return &models.InvoiceInfo{
		CurrentBalance: account.Balance,
		Withdrawals:    withdrawnTotalSum,
	}, nil
}

func (as invoicesService) CreateInvoice(ctx context.Context, id models.UserID) (*models.Invoice, error) {
	return as.accountRepository.AddInvoice(ctx, id)
}

func (as invoicesService) UpdateBalance(ctx context.Context, account models.Invoice, amount float64, orderNum models.OrderNum) (*models.Invoice, error) {
	if account.Balance < amount {
		return nil, apperrors.ErrAccountHasNotEnoughFunds
	}

	_, err := as.accountRepository.AddInvoiceChange(ctx, account.ID, amount, orderNum)
	if err != nil {
		return nil, err
	}

	return as.accountRepository.UpdateBalance(ctx, account, amount)
}

func (as invoicesService) GetWithdrawals(ctx context.Context, account *models.Invoice) (*models.WithdrawalsList, error) {
	return as.accountRepository.GetWithdrawalsList(ctx, account.ID)
}
