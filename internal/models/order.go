package models

import (
	"encoding/json"
	"time"
)

type OrderID int

type OrdersStatus string

const (
	OrderStatusNew        = OrdersStatus("NEW")        // Заказ загружен в систему, но не попал в обработку.
	OrderStatusProcessing = OrdersStatus("PROCESSING") // Вознаграждение за заказ рассчитывается.
	OrderStatusInvalid    = OrdersStatus("INVALID")    // Система расчёта вознаграждений отказала в расчёте.
	OrderStatusProcessed  = OrdersStatus("PROCESSED")  // Данные по заказу проверены и информация о расчёте успешно получена.
)

type OrderNum string

type Order struct {
	ID         OrderID      `db:"id" json:"-"`
	UserID     UserID       `db:"user" json:"-"`
	Number     OrderNum     `db:"number" json:"number"`
	Status     OrdersStatus `db:"status" json:"status"`
	Accrual    *float64     `db:"accrual" json:"accrual,omitempty"`
	UploadedAt time.Time    `db:"created_at" json:"uploaded_at"`
	Calculated bool         `db:"calculated" json:"-"`
}

func (o *Order) MarshalJSON() ([]byte, error) {
	type Alias Order
	return json.Marshal(&struct {
		UploadedAt string `json:"uploaded_at"`
		*Alias
	}{
		UploadedAt: o.UploadedAt.Format(time.RFC3339),
		Alias:      (*Alias)(o),
	})
}

type OrdersList = []Order

type OrderInfo struct {
	UserID
	OrderNum
}

type OrdersStatusInfo struct {
	OrderNumber OrderNum     `json:"order"`
	Status      OrdersStatus `json:"status"`
	Accrual     float64      `json:"accrual"`
}
