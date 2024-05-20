package service

import (
	"context"

	"e1m0re/loyalty-srv/internal/models"
	"e1m0re/loyalty-srv/internal/repository"
)

type AuthService struct {
	userRepository repository.UserRepository
}

func NewAuthService(userRepository repository.UserRepository) *AuthService {
	return &AuthService{
		userRepository: userRepository,
	}
}

func (us *AuthService) CreateUser(ctx context.Context, userInfo models.UserInfo) (user *models.User, err error) {
	return us.userRepository.CreateUser(ctx, userInfo)
}

func (us *AuthService) SignIn(ctx context.Context, userInfo models.UserInfo) (ok bool, err error) {
	return true, nil
}

func (us *AuthService) Verify(ctx context.Context, userInfo models.UserInfo) (ok bool, err error) {
	return true, nil
}
