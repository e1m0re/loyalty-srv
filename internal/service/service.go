package service

import (
	"context"

	"e1m0re/loyalty-srv/internal/models"
)

type UsersService interface {
	SignUp(ctx context.Context, userInfo models.UserInfo) (user *models.User, err error)
	SignIn(ctx context.Context, userInfo models.UserInfo) (ok bool, err error)
	Verify(ctx context.Context, userInfo models.UserInfo) (ok bool, err error)
}

type OrdersService interface {
	ValidateNumber(ctx context.Context, orderNum models.OrderNum) (ok bool, err error)
	LoadOrder(ctx context.Context, orderNum models.OrderNum) (order *models.Order, isNew bool, err error)
	GetLoadedOrdersByUserId(ctx context.Context, id models.UserId) (models.OrdersList, error)
	UpdateOrder(ctx context.Context, id models.OrderId, status models.OrdersStatus, accrual int) (models.Order, error)
}

type AccountsService interface {
	GetAccountByUserId(ctx context.Context, id models.UserId) (models.Account, error)
	GetAccountInfoByUserId(ctx context.Context, id models.UserId) (models.AccountInfo, error)
	Withdraw(ctx context.Context, id models.AccountId, amount int, orderNum models.OrderNum) (models.Account, error)
	GetWithdrawals(ctx context.Context, id models.UserId) (models.WithdrawalsList, error)
	UpdateBalance(ctx context.Context, id models.AccountId, amount int) (models.Account, error)
}
