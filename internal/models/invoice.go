package models

import "time"

type InvoiceID int

type Invoice struct {
	ID      InvoiceID `json:"id" db:"id"`
	UserID  UserID    `json:"user" db:"user"`
	Balance float64   `json:"balance,omitempty" db:"balance"`
}

type InvoiceChanges struct {
	ID        int
	InvoiceID InvoiceID
	OrderNum  OrderNum
	Amount    float64
	TS        time.Time
}

type Withdrawal struct {
	OrderNum    OrderNum  `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

type InvoiceInfo struct {
	CurrentBalance float64 `json:"current"`
	Withdrawals    int     `json:"withdrawals"`
}

type WithdrawalsList = []Withdrawal
