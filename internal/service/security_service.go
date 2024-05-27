package service

import "golang.org/x/crypto/bcrypt"

type securityService struct{}

func NewSecurityService() SecurityService {
	return &securityService{}
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
