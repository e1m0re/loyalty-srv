package repository

import (
	"context"
	"log/slog"
	"math"
	"time"

	"github.com/jmoiron/sqlx"

	"e1m0re/loyalty-srv/internal/models"
)

type accountRepository struct {
	db *sqlx.DB
}

func NewAccountRepository(db *sqlx.DB) AccountRepository {
	return &accountRepository{
		db: db,
	}
}

func (repo *accountRepository) AddAccount(ctx context.Context, userID models.UserID) (*models.Account, error) {
	a := &models.Account{
		UserID:  userID,
		Balance: 0,
	}
	query := "INSERT INTO accounts (\"user\", current_value) VALUES ($1, $2) RETURNING id"
	err := repo.db.QueryRowContext(ctx, query, a.UserID, a.Balance).Scan(&a.ID)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (repo *accountRepository) AddAccountChange(ctx context.Context, accountID models.AccountID, amount float64, orderNum models.OrderNum) (*models.AccountChanges, error) {
	ac := &models.AccountChanges{
		AccountID: accountID,
		OrderNum:  orderNum,
		Amount:    amount,
		TS:        time.Now(),
	}
	query := "INSERT INTO accounts_changes (account, amount, ts, \"order\") VALUES ($1, $2, $3, $4) RETURNING id"
	err := repo.db.QueryRowContext(ctx, query, ac.AccountID, ac.Amount, ac.TS, ac.OrderNum).Scan(&ac.ID)
	if err != nil {
		return nil, err
	}

	return ac, err
}

func (repo *accountRepository) GetAccountByUserID(ctx context.Context, userID models.UserID) (*models.Account, error) {
	account := &models.Account{}
	err := repo.db.GetContext(ctx, account, "SELECT id, \"user\", balance::numeric FROM accounts WHERE \"user\" = $1", userID)
	if err != nil {
		return nil, err
	}

	return account, err
}

func (repo *accountRepository) GetWithdrawalsList(ctx context.Context, accountID models.AccountID) (*models.WithdrawalsList, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT \"order\", ABS(delta::numeric), ts FROM accounts_changes WHERE account = $1 and delta::numeric < 0", accountID)
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

func (repo *accountRepository) GetWithdrawnTotalSum(ctx context.Context, accountID models.AccountID) (int, error) {
	var sum float64
	query := "SELECT sum(delta) FROM accounts_changes WHERE account = $1 AND delta::numeric < 0"
	err := repo.db.QueryRowContext(ctx, query, accountID).Scan(&sum)
	if err != nil {
		return 0, err
	}

	return int(math.Abs(sum)), nil
}

func (repo *accountRepository) UpdateAccount(ctx context.Context, account *models.Account) error {
	query := "UPDATE accounts SET balance=$1 WHERE id = $1"
	_, err := repo.db.ExecContext(ctx, query, account.Balance, account.ID)

	return err
}
