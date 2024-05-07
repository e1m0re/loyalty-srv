package repository

import (
	"context"

	"github.com/jmoiron/sqlx"

	"e1m0re/loyalty-srv/internal/apperrors"
	"e1m0re/loyalty-srv/internal/models"
)

type userRepository struct {
	DB *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{
		DB: db,
	}
}

func (repo userRepository) CreateUser(ctx context.Context, info models.UserInfo) (*models.User, error) {
	err := apperrors.NewNotImplementedError("userRepository::CreateUser")
	return nil, err
}
