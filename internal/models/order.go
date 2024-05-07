package models

import "time"

type OrderId int

type OrdersStatus string

const (
	New        = OrdersStatus("NEW")        // Заказ загружен в систему, но не попал в обработку.
	Processing = OrdersStatus("PROCESSING") // Вознаграждение за заказ рассчитывается.
	Invalid    = OrdersStatus("INVALID")    // Система расчёта вознаграждений отказала в расчёте.
	Processed  = OrdersStatus("PROCESSED")  // Данные по заказу проверены и информация о расчёте успешно получена.
)

type OrderNum string

type Order struct {
	ID         OrderId
	UserID     UserId
	Number     OrderNum     `json:"number"`
	Status     OrdersStatus `json:"status"`
	Accrual    int          `json:"accrual"`
	UploadedAt *time.Time   `json:"uploaded_at"`
}

type OrdersList = []Order

type OrderInfo struct {
	Number OrderNum
}
