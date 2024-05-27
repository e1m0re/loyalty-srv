package service

import (
	"context"
	"strconv"
	"unicode/utf8"

	"e1m0re/loyalty-srv/internal/apperrors"
	"e1m0re/loyalty-srv/internal/models"
	"e1m0re/loyalty-srv/internal/repository"
)

type ordersService struct {
	orderRepository repository.OrderRepository
}

func NewOrdersService(orderRepository repository.OrderRepository) OrdersService {
	return &ordersService{
		orderRepository: orderRepository,
	}
}

func (os ordersService) ValidateNumber(ctx context.Context, orderNum models.OrderNum) (ok bool, err error) {
	orderLen := len(orderNum)
	if orderLen == 0 {
		return false, apperrors.EmptyOrderNumberError
	}

	var sum int
	parity := orderLen % 2
	for idx, num := range orderNum {
		buf := make([]byte, 1)
		_ = utf8.EncodeRune(buf, num)
		digit, err := strconv.Atoi(string(buf))
		if err != nil {
			return false, err
		}

		if idx%2 == parity {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}

	ok = sum%10 == 0

	return ok, nil
}

func (os ordersService) NewOrder(ctx context.Context, orderNum models.OrderNum) (order *models.Order, isNew bool, err error) {
	ok, err := os.ValidateNumber(ctx, orderNum)
	if !ok {
		return nil, false, err
	}

	orderInfo := models.OrderInfo{
		UserId:   4, // todo get user
		OrderNum: orderNum,
	}

	// todo bad cases
	order, err = os.orderRepository.AddOrder(ctx, orderInfo)
	if err != nil {
		return nil, false, err
	}

	return order, true, nil
}

func (os ordersService) GetLoadedOrdersByUserId(ctx context.Context, userId models.UserId) (*models.OrdersList, error) {
	return os.orderRepository.GetLoadedOrdersByUserId(ctx, userId)
}

func (os ordersService) UpdateOrder(ctx context.Context, id models.OrderId, status models.OrdersStatus, accrual int) (models.Order, error) {
	return models.Order{}, nil
}
