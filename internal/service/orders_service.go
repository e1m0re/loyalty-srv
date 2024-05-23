package service

import (
	"context"

	"e1m0re/loyalty-srv/internal/models"
)

type ordersService struct {
}

func NewOrdersService() OrdersService {
	return &ordersService{}
}

func (os ordersService) ValidateNumber(ctx context.Context, orderNum models.OrderNum) (ok bool, err error) {
	return true, nil
}

func (os ordersService) NewOrder(ctx context.Context, orderNum models.OrderNum) (order *models.Order, isNew bool, err error) {
	return &models.Order{}, false, nil
}

func (os ordersService) GetLoadedOrdersByUserId(ctx context.Context, id models.UserId) (models.OrdersList, error) {
	return models.OrdersList{}, nil
}

func (os ordersService) UpdateOrder(ctx context.Context, id models.OrderId, status models.OrdersStatus, accrual int) (models.Order, error) {
	return models.Order{}, nil
}
