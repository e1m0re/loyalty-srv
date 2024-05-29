package service

import (
	"context"

	"github.com/go-chi/jwtauth/v5"

	"e1m0re/loyalty-srv/internal/models"
	"e1m0re/loyalty-srv/internal/repository"
)

//go:generate go run github.com/vektra/mockery/v2@v2.43.1 --name=AccountsService
type AccountsService interface {
	GetAccountByUserID(ctx context.Context, userID models.UserID) (*models.Account, error)
	GetAccountInfo(ctx context.Context, account *models.Account) (*models.AccountInfo, error)
	GetWithdrawals(ctx context.Context, account *models.Account) (*models.WithdrawalsList, error)
	CreateAccount(ctx context.Context, id models.UserID) (*models.Account, error)
	UpdateBalance(ctx context.Context, id models.AccountID, amount int) (*models.Account, error)
	Withdraw(ctx context.Context, account models.Account, amount float64, orderNum models.OrderNum) (*models.Account, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.43.1 --name=OrdersService
type OrdersService interface {
	GetLoadedOrdersByUserID(ctx context.Context, userID models.UserID) (*models.OrdersList, error)
	NewOrder(ctx context.Context, orderInfo models.OrderInfo) (*models.Order, error)
	UpdateOrdersCalculated(ctx context.Context, order models.Order, calculated bool) (*models.Order, error)
	UpdateOrdersStatus(ctx context.Context, order models.Order, status models.OrdersStatus, accrual int) (*models.Order, error)
	ValidateNumber(ctx context.Context, orderNum models.OrderNum) (ok bool, err error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.43.1 --name=SecurityService
type SecurityService interface {
	CheckPassword(hashPassword string, password string) bool
	GenerateAuthToken() *jwtauth.JWTAuth
	GenerateToken(user *models.User) (string, error)
	GetPasswordHash(password string) (string, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.43.1 --name=UsersService
type UsersService interface {
	CreateUser(ctx context.Context, userInfo models.UserInfo) (user *models.User, err error)
	FindUserByUsername(ctx context.Context, username string) (user *models.User, err error)
	SignIn(ctx context.Context, userInfo models.UserInfo) (token string, err error)
}

type Services struct {
	AccountsService
	OrdersService
	SecurityService
	UsersService
}

func NewServices(repo *repository.Repositories, securityService SecurityService) *Services {
	return &Services{
		AccountsService: NewAccountsService(repo.AccountRepository),
		OrdersService:   NewOrdersService(repo.OrderRepository),
		SecurityService: securityService,
		UsersService:    NewUsersService(repo.UserRepository, securityService),
	}
}
