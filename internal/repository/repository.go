package repository

import (
	"context"

	"e1m0re/loyalty-srv/internal/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, info models.UserInfo) (*models.User, error)
}
