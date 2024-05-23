package service

import (
	"context"
	"crypto/sha1"
	"fmt"

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

func (us *AuthService) CreateUser(ctx context.Context, userInfo *models.UserInfo) (user *models.User, err error) {
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

func (us *AuthService) FindUserByUsername(ctx context.Context, username string) (user *models.User, err error) {
	return nil, err
}

func (us *AuthService) SignIn(ctx context.Context, userInfo *models.UserInfo) (ok bool, err error) {
	return true, nil
}

func (us *AuthService) Verify(ctx context.Context, userInfo *models.UserInfo) (ok bool, err error) {
	return true, nil
}

func getPasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(nil)))
}
