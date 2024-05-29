package service

import (
	"context"

	"e1m0re/loyalty-srv/internal/apperrors"
	"e1m0re/loyalty-srv/internal/models"
	"e1m0re/loyalty-srv/internal/repository"
)

type invoicesService struct {
	invoiceRepository repository.InvoiceRepository
}

func NewInvoicesService(accountRepository repository.InvoiceRepository) InvoicesService {
	return &invoicesService{
		invoiceRepository: accountRepository,
	}
}

func (is invoicesService) GetInvoiceByID(ctx context.Context, invoiceID models.InvoiceID) (*models.Invoice, error) {
	return is.invoiceRepository.GetInvoiceByID(ctx, invoiceID)
}

func (is invoicesService) GetInvoiceByUserID(ctx context.Context, userID models.UserID) (*models.Invoice, error) {
	return is.invoiceRepository.GetInvoiceByUserID(ctx, userID)
}

func (is invoicesService) GetInvoiceInfo(ctx context.Context, account *models.Invoice) (*models.InvoiceInfo, error) {
	withdrawnTotalSum, err := is.invoiceRepository.GetWithdrawnTotalSum(ctx, account.ID)
	if err != nil {
		return nil, err
	}

	return &models.InvoiceInfo{
		CurrentBalance: account.Balance,
		Withdrawals:    withdrawnTotalSum,
	}, nil
}

func (is invoicesService) CreateInvoice(ctx context.Context, id models.UserID) (*models.Invoice, error) {
	return is.invoiceRepository.AddInvoice(ctx, id)
}

func (is invoicesService) UpdateBalance(ctx context.Context, account models.Invoice, amount float64, orderNum models.OrderNum) (*models.Invoice, error) {
	if account.Balance < amount {
		return nil, apperrors.ErrInvoiceHasNotEnoughFunds
	}

	_, err := is.invoiceRepository.AddInvoiceChange(ctx, account.ID, amount, orderNum)
	if err != nil {
		return nil, err
	}

	return is.invoiceRepository.UpdateBalance(ctx, account, amount)
}

func (is invoicesService) GetWithdrawals(ctx context.Context, account *models.Invoice) (*models.WithdrawalsList, error) {
	return is.invoiceRepository.GetWithdrawalsList(ctx, account.ID)
}
