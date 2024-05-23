package service

import (
	"context"
	"crypto/sha1"
	"fmt"

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
	user = &models.User{
		Username: userInfo.Username,
		Password: getPasswordHash(userInfo.Password),
	}

	id, err := us.userRepository.CreateUser(ctx, *user)
	if err != nil {
		return nil, err
	}

	user.ID = id
	return user, nil
}

func (us *usersService) FindUserByUsername(ctx context.Context, username string) (user *models.User, err error) {
	return nil, err
}

func (us *usersService) SignIn(ctx context.Context, userInfo *models.UserInfo) (ok bool, err error) {
	return true, nil
}

func (us *usersService) Verify(ctx context.Context, userInfo *models.UserInfo) (ok bool, err error) {
	return true, nil
}

func getPasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(nil)))
}
