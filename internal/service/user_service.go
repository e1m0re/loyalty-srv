package service

import (
	"context"

	"e1m0re/loyalty-srv/internal/models"
	"e1m0re/loyalty-srv/internal/repository"
)

type userService struct {
	UserRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{
		UserRepository: userRepository,
	}
}

func (us userService) SignUp(ctx context.Context, userInfo models.UserInfo) (user *models.User, err error) {
	return us.UserRepository.CreateUser(ctx, userInfo)
}

func (us userService) SignIn(ctx context.Context, userInfo models.UserInfo) (ok bool, err error) {
	return true, nil
}

func (us userService) Verify(ctx context.Context, userInfo models.UserInfo) (ok bool, err error) {
	return true, nil
}
