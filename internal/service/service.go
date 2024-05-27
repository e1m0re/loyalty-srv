package service

import (
	"context"
	"github.com/go-chi/jwtauth/v5"

	"e1m0re/loyalty-srv/internal/models"
	"e1m0re/loyalty-srv/internal/repository"
)

//go:generate go run github.com/vektra/mockery/v2@v2.43.1 --name=SecurityService
type SecurityService interface {
	GenerateAuthToken() *jwtauth.JWTAuth
	GenerateToken(user *models.User) (string, error)
	GetPasswordHash(password string) (string, error)
	CheckPassword(hashPassword string, password string) bool
}

//go:generate go run github.com/vektra/mockery/v2@v2.43.1 --name=UsersService
type UsersService interface {
	CreateUser(ctx context.Context, userInfo *models.UserInfo) (user *models.User, err error)
	FindUserByUsername(ctx context.Context, username string) (user *models.User, err error)
	SignIn(ctx context.Context, userInfo *models.UserInfo) (token string, err error)
	Verify(ctx context.Context, userInfo *models.UserInfo) (ok bool, err error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.43.1 --name=OrdersService
type OrdersService interface {
	ValidateNumber(ctx context.Context, orderNum models.OrderNum) (ok bool, err error)
	NewOrder(ctx context.Context, orderInfo models.OrderInfo) (*models.Order, error)
	GetLoadedOrdersByUserID(ctx context.Context, userID models.UserID) (*models.OrdersList, error)
	UpdateOrder(ctx context.Context, id models.OrderID, status models.OrdersStatus, accrual int) (models.Order, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.43.1 --name=AccountsService
type AccountsService interface {
	GetAccountByUserID(ctx context.Context, id models.UserID) (*models.Account, error)
	GetAccountInfoByUserID(ctx context.Context, id models.UserID) (*models.AccountInfo, error)
	Withdraw(ctx context.Context, id models.AccountID, amount int, orderNum models.OrderNum) (*models.Account, error)
	GetWithdrawals(ctx context.Context, id models.UserID) (*models.WithdrawalsList, error)
	UpdateBalance(ctx context.Context, id models.AccountID, amount int) (*models.Account, error)
}

type Services struct {
	UsersService
	OrdersService
	Accounts AccountsService
	SecurityService
}

func NewServices(repo *repository.Repositories, securityService SecurityService) *Services {
	return &Services{
		UsersService:    NewUsersService(repo.UserRepository, securityService),
		OrdersService:   NewOrdersService(repo.OrderRepository),
		SecurityService: securityService,
	}
}
