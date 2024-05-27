package repository

import (
	"context"
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
		ID:         0,
		UserID:     orderInfo.UserId,
		Number:     orderInfo.OrderNum,
		Status:     models.OrderStatusNew,
		Accrual:    0,
		UploadedAt: time.Now(),
	}

	query := "INSERT INTO orders (\"user\", created_at, status, number) VALUES ($1, $2, $3, $4) RETURNING id"
	err := repo.db.QueryRowContext(ctx, query, order.UserID, order.UploadedAt, order.Status, order.Number).Scan(&order.ID)
	if err != nil {
		return nil, err
	}

	return order, nil
}
