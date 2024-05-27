package service

import (
	"context"

	"e1m0re/loyalty-srv/internal/models"
	"e1m0re/loyalty-srv/internal/repository"
)

//go:generate go run github.com/vektra/mockery/v2@v2.43.1 --name=UsersService
type UsersService interface {
	CreateUser(ctx context.Context, userInfo *models.UserInfo) (user *models.User, err error)
	FindUserByUsername(ctx context.Context, username string) (user *models.User, err error)
	SignIn(ctx context.Context, userInfo *models.UserInfo) (ok bool, err error)
	Verify(ctx context.Context, userInfo *models.UserInfo) (ok bool, err error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.43.1 --name=OrdersService
type OrdersService interface {
	ValidateNumber(ctx context.Context, orderNum models.OrderNum) (ok bool, err error)
	NewOrder(ctx context.Context, orderNum models.OrderNum) (order *models.Order, isNew bool, err error)
	GetLoadedOrdersByUserId(ctx context.Context, id models.UserId) (models.OrdersList, error)
	UpdateOrder(ctx context.Context, id models.OrderId, status models.OrdersStatus, accrual int) (models.Order, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.43.1 --name=AccountsService
type AccountsService interface {
	GetAccountByUserId(ctx context.Context, id models.UserId) (*models.Account, error)
	GetAccountInfoByUserId(ctx context.Context, id models.UserId) (*models.AccountInfo, error)
	Withdraw(ctx context.Context, id models.AccountId, amount int, orderNum models.OrderNum) (*models.Account, error)
	GetWithdrawals(ctx context.Context, id models.UserId) (*models.WithdrawalsList, error)
	UpdateBalance(ctx context.Context, id models.AccountId, amount int) (*models.Account, error)
}

type Services struct {
	UsersService
	OrdersService
	Accounts AccountsService
}

func NewServices(repo *repository.Repositories) *Services {
	return &Services{
		UsersService:  NewUsersService(repo.UserRepository),
		OrdersService: NewOrdersService(repo.OrderRepository),
	}
}
