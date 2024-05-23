package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"

	"e1m0re/loyalty-srv/internal/models"
)

type UserRepository interface {
	GetUserByUsername(ctx context.Context, username string) (user *models.User, err error)
	CreateUser(ctx context.Context, user models.User) (models.UserId, error)
	UpdateUsersLastLogin(ctx context.Context, id models.UserId, t time.Time) error
}

type Repositories struct {
	UserRepository
}

func NewRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		UserRepository: NewUserRepository(db),
	}
}
