package models

import (
	"time"
)

type UserId int

type User struct {
	ID        UserId     `db:"id"`
	Username  string     `db:"username"`
	Password  string     `db:"password"`
	LastLogin *time.Time `db:"last_login"`
}

type UserInfo struct {
	Username string `json:"login"`
	Password string `json:"password"`
}
