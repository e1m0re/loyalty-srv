package repository

import (
	"context"

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
