package service

import (
	"github.com/go-chi/jwtauth/v5"
	"golang.org/x/crypto/bcrypt"

	"e1m0re/loyalty-srv/internal/models"
)

type securityService struct {
	authToken *jwtauth.JWTAuth
}

func NewSecurityService(jwtSecretKey string) SecurityService {
	return &securityService{
		authToken: jwtauth.New("HS256", []byte(jwtSecretKey), nil),
	}
}

func (ss *securityService) GenerateAuthToken() *jwtauth.JWTAuth {
	return ss.authToken
}

func (ss *securityService) GetPasswordHash(password string) (string, error) {
	bytePassword := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (ss *securityService) CheckPassword(hashPassword string, password string) bool {
	return nil == bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
}

func (ss *securityService) GenerateToken(user *models.User) (string, error) {
	claims := map[string]interface{}{"id": user.ID, "username": user.Username}
	_, tokenString, err := ss.authToken.Encode(claims)

	return tokenString, err
}
