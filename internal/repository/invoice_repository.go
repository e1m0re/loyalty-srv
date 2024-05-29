package repository

import (
	"context"
	"log/slog"
	"math"
	"time"

	"github.com/jmoiron/sqlx"

	"e1m0re/loyalty-srv/internal/models"
)

type invoiceRepository struct {
	db *sqlx.DB
}

func NewInvoiceRepository(db *sqlx.DB) InvoiceRepository {
	return &invoiceRepository{
		db: db,
	}
}

func (repo *invoiceRepository) AddInvoice(ctx context.Context, userID models.UserID) (*models.Invoice, error) {
	a := &models.Invoice{
		UserID:  userID,
		Balance: 0,
	}
	query := "INSERT INTO invoices (\"user\", balance) VALUES ($1, $2) RETURNING id"
	err := repo.db.QueryRowContext(ctx, query, a.UserID, a.Balance).Scan(&a.ID)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (repo *invoiceRepository) AddInvoiceChange(ctx context.Context, accountID models.InvoiceID, amount float64, orderNum models.OrderNum) (*models.InvoiceChanges, error) {
	ac := &models.InvoiceChanges{
		InvoiceID: accountID,
		OrderNum:  orderNum,
		Amount:    amount,
		TS:        time.Now(),
	}
	query := "INSERT INTO invoices_changes (account, amount, ts, \"order\") VALUES ($1, $2, $3, $4) RETURNING id"
	err := repo.db.QueryRowContext(ctx, query, ac.InvoiceID, ac.Amount, ac.TS, ac.OrderNum).Scan(&ac.ID)
	if err != nil {
		return nil, err
	}

	return ac, err
}

func (repo *invoiceRepository) GetInvoiceByUserID(ctx context.Context, userID models.UserID) (*models.Invoice, error) {
	account := &models.Invoice{}
	err := repo.db.GetContext(ctx, account, "SELECT id, \"user\", balance::numeric FROM invoices WHERE \"user\" = $1", userID)
	if err != nil {
		return nil, err
	}

	return account, err
}

func (repo *invoiceRepository) GetWithdrawalsList(ctx context.Context, accountID models.InvoiceID) (*models.WithdrawalsList, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT \"order\", ABS(delta::numeric), ts FROM invoices_changes WHERE account = $1 and amount::numeric < 0", accountID)
	if err != nil {
		return nil, err
	}

	withdrawals := make(models.WithdrawalsList, 0)
	defer rows.Close()
	for rows.Next() {
		w := models.Withdrawal{}
		err := rows.Scan(&w.OrderNum, &w.Sum, &w.ProcessedAt)
		if err != nil {
			slog.Error("GetWithdrawalsList", slog.Any("accountID", accountID), slog.String("err", err.Error()))
			continue
		}

		withdrawals = append(withdrawals, w)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &withdrawals, nil
}

func (repo *invoiceRepository) GetWithdrawnTotalSum(ctx context.Context, accountID models.InvoiceID) (int, error) {
	var sum float64
	query := "SELECT sum(amount) FROM invoices_changes WHERE account = $1 AND amount::numeric < 0"
	err := repo.db.QueryRowContext(ctx, query, accountID).Scan(&sum)
	if err != nil {
		return 0, err
	}

	return int(math.Abs(sum)), nil
}

func (repo *invoiceRepository) UpdateBalance(ctx context.Context, account models.Invoice, amount float64) (*models.Invoice, error) {
	account.Balance = account.Balance - amount
	query := "UPDATE invoices SET balance = $1 WHERE id = $2"
	_, err := repo.db.ExecContext(ctx, query, account.Balance, account.ID)
	if err != nil {
		return nil, err
	}

	return &account, nil
}
