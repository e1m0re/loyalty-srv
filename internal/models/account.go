package models

import "time"

type AccountID int

type Account struct {
	ID      AccountID `json:"id" db:"id"`
	UserID  UserID    `json:"user" db:"user"`
	Balance float64   `json:"balance,omitempty" db:"balance"`
}

type Withdrawal struct {
	OrderNum    OrderNum  `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

type AccountInfo struct {
	CurrentBalance float64 `json:"current"`
	Withdrawals    int     `json:"withdrawals"`
}

type WithdrawalsList = []Withdrawal
