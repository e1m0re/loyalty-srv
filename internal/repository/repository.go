package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"

	"e1m0re/loyalty-srv/internal/models"
)

//go:generate go run github.com/vektra/mockery/v2@v2.43.1 --name=InvoiceRepository
type InvoiceRepository interface {
	AddInvoice(ctx context.Context, userID models.UserID) (*models.Invoice, error)
	AddInvoiceChange(ctx context.Context, invoiceID models.InvoiceID, amount float64, orderNum models.OrderNum) (*models.InvoiceChanges, error)
	GetInvoiceByUserID(ctx context.Context, userID models.UserID) (*models.Invoice, error)
	GetWithdrawalsList(ctx context.Context, invoiceID models.InvoiceID) (*models.WithdrawalsList, error)
	GetWithdrawnTotalSum(ctx context.Context, invoiceID models.InvoiceID) (int, error)
	UpdateBalance(ctx context.Context, invoice models.Invoice, amount float64) (*models.Invoice, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.43.1 --name=OrderRepository
type OrderRepository interface {
	AddOrder(ctx context.Context, orderInfo models.OrderInfo) (*models.Order, error)
	GetLoadedOrdersByUserID(ctx context.Context, userID models.UserID) (*models.OrdersList, error)
	GetOrderByNumber(ctx context.Context, num models.OrderNum) (*models.Order, error)
	UpdateOrdersCalculated(ctx context.Context, order models.Order, calculated bool) (*models.Order, error)
	UpdateOrdersStatus(ctx context.Context, order models.Order, status models.OrdersStatus, accrual int) (*models.Order, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.43.1 --name=UserRepository
type UserRepository interface {
	CreateUser(ctx context.Context, userInfo models.UserInfo) (user *models.User, err error)
	GetUserByUsername(ctx context.Context, username string) (user *models.User, err error)
	UpdateUsersLastLogin(ctx context.Context, id models.UserID, t time.Time) error
}

type Repositories struct {
	InvoiceRepository
	OrderRepository
	UserRepository
}

func NewRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		InvoiceRepository: NewInvoiceRepository(db),
		OrderRepository:   NewOrderRepository(db),
		UserRepository:    NewUserRepository(db),
	}
}
