package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"

	"e1m0re/loyalty-srv/internal/apperrors"
	"e1m0re/loyalty-srv/internal/models"
)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (repo *userRepository) CreateUser(ctx context.Context, user models.User) (models.UserId, error) {
	var id int
	query := "INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id"
	row := repo.db.QueryRowxContext(ctx, query, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return models.UserId(id), nil
}

func (repo *userRepository) GetUserByUsername(ctx context.Context, username string) (user *models.User, err error) {
	user = &models.User{}
	query := "SELECT * FROM users WHERE username = $1 LIMIT 1"
	err = repo.db.GetContext(ctx, user, query, username)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, apperrors.EntityNotFoundError
	case err != nil:
		return nil, err
	default:
		return
	}
}

func (repo *userRepository) UpdateUsersLastLogin(ctx context.Context, id models.UserId, t time.Time) error {
	query := "UPDATE users SET last_login = $1 WHERE id = $2"
	_, err := repo.db.ExecContext(ctx, query, t, id)

	return err
}
