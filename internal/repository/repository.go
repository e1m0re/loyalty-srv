package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"

	"e1m0re/loyalty-srv/internal/models"
)

//go:generate go run github.com/vektra/mockery/v2@v2.43.1 --name=AccountRepository
type AccountRepository interface {
	AddAccount(ctx context.Context, userID models.UserID) (*models.Account, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.43.1 --name=OrderRepository
type OrderRepository interface {
	AddOrder(ctx context.Context, orderInfo models.OrderInfo) (*models.Order, error)
	GetLoadedOrdersByUserID(ctx context.Context, userID models.UserID) (*models.OrdersList, error)
	GetOrderByNumber(ctx context.Context, num models.OrderNum) (*models.Order, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.43.1 --name=UserRepository
type UserRepository interface {
	CreateUser(ctx context.Context, userInfo models.UserInfo) (user *models.User, err error)
	GetUserByUsername(ctx context.Context, username string) (user *models.User, err error)
	UpdateUsersLastLogin(ctx context.Context, id models.UserID, t time.Time) error
}

type Repositories struct {
	AccountRepository
	OrderRepository
	UserRepository
}

func NewRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		AccountRepository: NewAccountRepository(db),
		OrderRepository:   NewOrderRepository(db),
		UserRepository:    NewUserRepository(db),
	}
}
