package models

import (
	"encoding/json"
	"time"
)

type OrderId int

type OrdersStatus string

const (
	OrderStatusNew        = OrdersStatus("NEW")        // Заказ загружен в систему, но не попал в обработку.
	OrderStatusProcessing = OrdersStatus("PROCESSING") // Вознаграждение за заказ рассчитывается.
	OrderStatusInvalid    = OrdersStatus("INVALID")    // Система расчёта вознаграждений отказала в расчёте.
	OrderStatusProcessed  = OrdersStatus("PROCESSED")  // Данные по заказу проверены и информация о расчёте успешно получена.
)

type OrderNum string

type Order struct {
	ID         OrderId
	UserID     UserId
	Number     OrderNum     `json:"number"`
	Status     OrdersStatus `json:"status"`
	Accrual    int          `json:"accrual,omitempty"`
	UploadedAt time.Time    `json:"uploaded_at"`
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
	UserId
	OrderNum
}
