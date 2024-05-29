package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"

	"e1m0re/loyalty-srv/internal/models"
)

type orderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) OrderRepository {
	return &orderRepository{
		db: db,
	}
}

func (repo orderRepository) AddOrder(ctx context.Context, orderInfo models.OrderInfo) (*models.Order, error) {
	order := &models.Order{
		UserID:     orderInfo.UserID,
		Number:     orderInfo.OrderNum,
		Status:     models.OrderStatusNew,
		UploadedAt: time.Now(),
	}

	query := "INSERT INTO orders (\"user\", created_at, status, number) VALUES ($1, $2, $3, $4) RETURNING id"
	err := repo.db.QueryRowContext(ctx, query, order.UserID, order.UploadedAt, order.Status, order.Number).Scan(&order.ID)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (repo orderRepository) GetOrderByNumber(ctx context.Context, num models.OrderNum) (*models.Order, error) {
	order := models.Order{}
	err := repo.db.GetContext(ctx, &order, "SELECT * FROM orders WHERE number = $1", num)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &order, nil
}

func (repo orderRepository) GetLoadedOrdersByUserID(ctx context.Context, userID models.UserID) (*models.OrdersList, error) {
	orders := models.OrdersList{}
	err := repo.db.SelectContext(ctx, &orders, "SELECT * FROM orders WHERE \"user\" = $1", userID)
	if err != nil {
		return nil, err
	}

	return &orders, nil
}

func (repo orderRepository) UpdateOrdersCalculated(ctx context.Context, order models.Order, calculated bool) (*models.Order, error) {
	_, err := repo.db.ExecContext(ctx, "UPDATE orders SET calculated = $1 WHERE number = $2", calculated, order.Number)
	if err != nil {
		return nil, err
	}

	order.Calculated = calculated
	return &order, nil
}

func (repo orderRepository) UpdateOrdersStatus(ctx context.Context, order models.Order, status models.OrdersStatus, accrual int) (*models.Order, error) {
	_, err := repo.db.ExecContext(ctx, "UPDATE orders set status = $1, accrual = $2 WHERE number = $3", status, accrual, order.Number)
	if err != nil {
		return nil, err
	}

	order.Status = status
	order.Accrual = &accrual
	return &order, nil
}
