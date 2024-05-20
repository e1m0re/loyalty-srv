package models

import "time"

type AccountId int

type Account struct {
	ID      AccountId `json:"id"`
	User    UserId    `json:"user"`
	Balance float64   `json:"balance,omitempty"`
}

type Withdrawal struct {
	OrderNum    OrderNum   `json:"order"`
	Sum         int        `json:"sum"`
	ProcessedAt *time.Time `json:"processed_at"`
}

type AccountInfo struct {
	CurrentBalance float64 `json:"current"`
	Withdrawals    int     `json:"withdrawals"`
}

type WithdrawalsList = []Withdrawal
