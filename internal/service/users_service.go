package service

import (
	"context"

	"e1m0re/loyalty-srv/internal/models"
	"e1m0re/loyalty-srv/internal/repository"
)

type usersService struct {
	UserRepository repository.UserRepository
}

func NewUsersService(userRepository repository.UserRepository) UsersService {
	return &usersService{
		UserRepository: userRepository,
	}
}

func (us usersService) SignUp(ctx context.Context, userInfo models.UserInfo) (user *models.User, err error) {
	return us.UserRepository.CreateUser(ctx, userInfo)
}

func (us usersService) SignIn(ctx context.Context, userInfo models.UserInfo) (ok bool, err error) {
	return true, nil
}

func (us usersService) Verify(ctx context.Context, userInfo models.UserInfo) (ok bool, err error) {
	return true, nil
}
