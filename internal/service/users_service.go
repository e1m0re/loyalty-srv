package service

import (
	"context"
	"time"

	"e1m0re/loyalty-srv/internal/models"
	"e1m0re/loyalty-srv/internal/repository"
)

type usersService struct {
	userRepository  repository.UserRepository
	securityService SecurityService
}

func NewUsersService(userRepository repository.UserRepository, securityService SecurityService) UsersService {
	return &usersService{
		userRepository:  userRepository,
		securityService: securityService,
	}
}

func (us *usersService) CreateUser(ctx context.Context, userInfo *models.UserInfo) (user *models.User, err error) {
	passwordHash, err := us.securityService.GetPasswordHash(userInfo.Password)
	if err != nil {
		return nil, err
	}

	user, err = us.userRepository.CreateUser(ctx, models.UserInfo{
		Username: userInfo.Username,
		Password: passwordHash,
	})

	return user, err
}

func (us *usersService) FindUserByUsername(ctx context.Context, username string) (user *models.User, err error) {
	return us.userRepository.GetUserByUsername(ctx, username)
}

func (us *usersService) SignIn(ctx context.Context, userInfo *models.UserInfo) (ok bool, err error) {
	user, err := us.userRepository.GetUserByUsername(ctx, userInfo.Username)
	if err != nil {
		return false, err
	}

	if !us.securityService.CheckPassword(user.Password, userInfo.Password) {
		return false, nil
	}

	err = us.userRepository.UpdateUsersLastLogin(ctx, user.ID, time.Now())
	if err != nil {
		return false, err
	}

	return true, nil
}

func (us *usersService) Verify(ctx context.Context, userInfo *models.UserInfo) (ok bool, err error) {
	return true, nil
}
