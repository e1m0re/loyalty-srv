package service

import (
	"context"

	"github.com/go-chi/jwtauth/v5"

	"e1m0re/loyalty-srv/internal/models"
	"e1m0re/loyalty-srv/internal/repository"
)

//go:generate go run github.com/vektra/mockery/v2@v2.43.1 --name=InvoicesService
type InvoicesService interface {
	GetInvoiceByID(ctx context.Context, invoiceID models.InvoiceID) (*models.Invoice, error)
	GetInvoiceByUserID(ctx context.Context, userID models.UserID) (*models.Invoice, error)
	GetInvoiceInfo(ctx context.Context, invoice *models.Invoice) (*models.InvoiceInfo, error)
	GetWithdrawals(ctx context.Context, invoice *models.Invoice) (*models.WithdrawalsList, error)
	CreateInvoice(ctx context.Context, id models.UserID) (*models.Invoice, error)
	UpdateBalance(ctx context.Context, invoice models.Invoice, amount float64, orderNum models.OrderNum) (*models.Invoice, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.43.1 --name=OrdersService
type OrdersService interface {
	GetLoadedOrdersByUserID(ctx context.Context, userID models.UserID) (*models.OrdersList, error)
	GetNotCalculatedOrder(ctx context.Context, limit int) (*models.OrdersList, error)
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

type OrdersProcessor interface {
	RecalculateProcessedOrders(ctx context.Context) error
	CheckProcessingOrders(ctx context.Context) error
}

type Services struct {
	InvoicesService
	OrdersService
	SecurityService
	UsersService
}

func NewServices(repo *repository.Repositories, securityService SecurityService) *Services {
	return &Services{
		InvoicesService: NewInvoicesService(repo.InvoiceRepository),
		OrdersService:   NewOrdersService(repo.OrderRepository),
		SecurityService: securityService,
		UsersService:    NewUsersService(repo.UserRepository, securityService),
	}
}
