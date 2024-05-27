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

//go:generate go run github.com/vektra/mockery/v2@v2.43.1 --name=OrderRepository
type OrderRepository interface {
	AddOrder(ctx context.Context, orderInfo models.OrderInfo) (*models.Order, error)
}

type Repositories struct {
	UserRepository
	OrderRepository
}

func NewRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		UserRepository:  NewUserRepository(db),
		OrderRepository: NewOrderRepository(db),
	}
}
