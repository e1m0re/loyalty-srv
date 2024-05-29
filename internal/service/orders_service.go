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
		return false, apperrors.ErrEmptyOrderNumber
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

func (os ordersService) NewOrder(ctx context.Context, orderInfo models.OrderInfo) (*models.Order, error) {
	ok, err := os.ValidateNumber(ctx, orderInfo.OrderNum)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, apperrors.ErrInvalidOrderNumber
	}

	order, err := os.orderRepository.GetOrderByNumber(ctx, orderInfo.OrderNum)
	if err != nil {
		return nil, err
	}

	if order != nil {
		if order.UserID != orderInfo.UserID {
			return nil, apperrors.ErrOrderWasLoadedByAnotherUser
		}

		return nil, apperrors.ErrOrderWasLoaded
	}

	order, err = os.orderRepository.AddOrder(ctx, orderInfo)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (os ordersService) GetLoadedOrdersByUserID(ctx context.Context, userID models.UserID) (*models.OrdersList, error) {
	return os.orderRepository.GetLoadedOrdersByUserID(ctx, userID)
}

func (os ordersService) UpdateOrdersCalculated(ctx context.Context, order models.Order, calculated bool) (*models.Order, error) {
	return os.orderRepository.UpdateOrdersCalculated(ctx, order, calculated)
}

func (os ordersService) UpdateOrdersStatus(ctx context.Context, order models.Order, status models.OrdersStatus, accrual int) (*models.Order, error) {
	return os.orderRepository.UpdateOrdersStatus(ctx, order, status, accrual)
}
