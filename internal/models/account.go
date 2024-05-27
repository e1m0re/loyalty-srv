package models

import "time"

type AccountID int

type Account struct {
	ID      AccountID `json:"id"`
	User    UserID    `json:"user"`
	Balance float64   `json:"balance,omitempty"`
}

type Withdrawal struct {
	OrderNum    OrderNum  `json:"order"`
	Sum         int       `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

type AccountInfo struct {
	CurrentBalance float64 `json:"current"`
	Withdrawals    int     `json:"withdrawals"`
}

type WithdrawalsList = []Withdrawal
