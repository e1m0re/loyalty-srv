package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
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

func (repo userRepository) CreateUser(ctx context.Context, userInfo models.UserInfo) (user *models.User, err error) {
	user = &models.User{
		Username: userInfo.Username,
		Password: userInfo.Password,
	}

	query := "INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id"
	err = repo.db.QueryRowxContext(ctx, query, userInfo.Username, userInfo.Password).Scan(&user.ID)
	if err != nil {
		if err.(*pgconn.PgError).Code == "23505" {
			return nil, apperrors.ErrBusyLogin
		}
		return nil, err
	}

	return user, nil
}

func (repo userRepository) GetUserByUsername(ctx context.Context, username string) (user *models.User, err error) {
	user = &models.User{}
	query := "SELECT * FROM users WHERE username = $1 LIMIT 1"
	err = repo.db.GetContext(ctx, user, query, username)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, apperrors.ErrEntityNotFound
	case err != nil:
		return nil, err
	default:
		return
	}
}

func (repo userRepository) UpdateUsersLastLogin(ctx context.Context, id models.UserID, t time.Time) error {
	query := "UPDATE users SET last_login = $1 WHERE id = $2"
	_, err := repo.db.ExecContext(ctx, query, t, id)

	return err
}
