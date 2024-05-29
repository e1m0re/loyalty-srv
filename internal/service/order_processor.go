package service

import (
	"context"
	"log/slog"
	"strconv"
)

type orderProcessor struct {
	invoicesService InvoicesService
	ordersService   OrdersService
}

func NewOrdersProcessor(invoicesService InvoicesService, ordersService OrdersService) OrdersProcessor {
	return &orderProcessor{
		invoicesService: invoicesService,
		ordersService:   ordersService,
	}
}

func (p orderProcessor) RecalculateProcessedOrders(ctx context.Context) error {
	ordersList, err := p.ordersService.GetNotCalculatedOrder(ctx, 1000)
	if err != nil {
		return err
	}

	for _, order := range *ordersList {
		if order.Accrual == nil || *order.Accrual == 0 {
			continue
		}

		invoice, err := p.invoicesService.GetInvoiceByUserID(ctx, order.UserID)
		if err != nil {
			slog.Warn("Recalculate processed orders",
				slog.String("step", "getting invoice"),
				slog.String("order", strconv.Itoa(int(order.UserID))),
				slog.String("error", err.Error()),
			)
		}

		_, err = p.invoicesService.UpdateBalance(ctx, *invoice, float64(*order.Accrual), order.Number)
		if err != nil {
			slog.Warn("Recalculate processed orders",
				slog.String("step", "update balance of invoice"),
				slog.String("invoice", strconv.Itoa(int(invoice.ID))),
				slog.String("order", strconv.Itoa(int(order.UserID))),
				slog.String("error", err.Error()),
			)
			continue
		}

		_, err = p.ordersService.UpdateOrdersCalculated(ctx, order, true)
		if err != nil {
			slog.Warn("Recalculate processed orders",
				slog.String("step", "mark order as calculated"),
				slog.String("invoice", strconv.Itoa(int(invoice.ID))),
				slog.String("order", strconv.Itoa(int(order.UserID))),
				slog.String("error", err.Error()),
			)
		}
	}

	return nil
}

func (p orderProcessor) CheckProcessingOrders(ctx context.Context) error {
	return nil
}
