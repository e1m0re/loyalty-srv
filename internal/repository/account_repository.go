package repository

import (
	"context"
	"math"

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
	account := &models.Account{
		UserID:  userID,
		Balance: 0,
	}
	query := "INSERT INTO accounts (\"user\", current_value) VALUES ($1, $2) RETURNING id"
	err := repo.db.QueryRowContext(ctx, query, account.UserID, account.Balance).Scan(&account.ID)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (repo *accountRepository) GetAccountByUserID(ctx context.Context, userID models.UserID) (*models.Account, error) {
	account := &models.Account{}
	err := repo.db.GetContext(ctx, account, "SELECT * FROM accounts WHERE \"user\" = $1", userID)
	if err != nil {
		return nil, err
	}

	return account, err
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
