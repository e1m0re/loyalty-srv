package repository

import (
	"context"

	"github.com/jmoiron/sqlx"

	"e1m0re/loyalty-srv/internal/apperrors"
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

func (repo *UserRepo) CreateUser(ctx context.Context, info models.UserInfo) (*models.User, error) {
	err := apperrors.NewNotImplementedError("userRepository::CreateUser")
	return nil, err
}
