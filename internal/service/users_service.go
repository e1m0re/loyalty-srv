package service

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"e1m0re/loyalty-srv/internal/apperrors"
	"e1m0re/loyalty-srv/internal/models"
	"e1m0re/loyalty-srv/internal/repository"
)

type usersService struct {
	userRepository repository.UserRepository
}

func NewUsersService(userRepository repository.UserRepository) UsersService {
	return &usersService{
		userRepository: userRepository,
	}
}

func (us *usersService) CreateUser(ctx context.Context, userInfo *models.UserInfo) (user *models.User, err error) {
	user, err = us.userRepository.GetUserByUsername(ctx, userInfo.Username)
	if err != nil && errors.Is(err, apperrors.EntityNotFoundError) != true {
		return nil, err
	}

	if user != nil {
		return nil, apperrors.BusyLoginError
	}

	hash, err := us.getPasswordHash(userInfo.Password)
	if err != nil {
		return nil, err
	}

	user = &models.User{
		Username: userInfo.Username,
		Password: hash,
	}

	id, err := us.userRepository.CreateUser(ctx, *user)
	if err != nil {
		return nil, err
	}

	user.ID = id
	return user, nil
}

func (us *usersService) FindUserByUsername(ctx context.Context, username string) (user *models.User, err error) {
	return us.userRepository.GetUserByUsername(ctx, username)
}

func (us *usersService) SignIn(ctx context.Context, userInfo *models.UserInfo) (ok bool, err error) {
	user, err := us.userRepository.GetUserByUsername(ctx, userInfo.Username)
	if err != nil {
		return false, err
	}

	if !us.checkPassword(user.Password, userInfo.Password) {
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

func (us *usersService) getPasswordHash(password string) (string, error) {
	bytePassword := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (us *usersService) checkPassword(hashPassword string, password string) bool {
	return nil == bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
}
