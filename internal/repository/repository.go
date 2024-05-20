package repository

import (
	"context"

	"github.com/jmoiron/sqlx"

	"e1m0re/loyalty-srv/internal/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, info models.UserInfo) (*models.User, error)
}

type Repositories struct {
	UserRepository
}

func NewRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		UserRepository: NewUserRepository(db),
	}
}
