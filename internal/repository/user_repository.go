package repository

import (
	"context"

	"github.com/jmoiron/sqlx"

	"e1m0re/loyalty-srv/internal/models"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (repo *UserRepo) CreateUser(ctx context.Context, user models.User) (models.UserId, error) {

	var id int
	query := "INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id"
	row := repo.db.QueryRowxContext(ctx, query, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return models.UserId(id), nil
}
