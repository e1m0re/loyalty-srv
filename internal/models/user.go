package models

import (
	"time"
)

type UserId int

type User struct {
	ID           UserId     `json:"id"`
	Username     string     `json:"username"`
	PasswordHash string     `json:"password"`
	LastLogin    *time.Time `json:"last_login"`
}

type UserInfo struct {
	Username string `json:"login"`
	Password string `json:"password"`
}
