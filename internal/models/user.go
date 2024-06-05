package models

import (
	"time"
)

type UserID int

type User struct {
	ID        UserID     `db:"id"`
	Username  string     `db:"username"`
	Password  string     `db:"password"`
	LastLogin *time.Time `db:"last_login"`
}

type UserInfo struct {
	Username string `json:"login"`
	Password string `json:"password"`
}
