package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"e1m0re/loyalty-srv/internal/models"
)

type orderProcessor struct {
	accrualSystemAddress string
	httpClient           *http.Client
	invoicesService      InvoicesService
	ordersService        OrdersService
}

func NewOrdersProcessor(invoicesService InvoicesService, ordersService OrdersService, accrualSystemAddress string) OrdersProcessor {
	return &orderProcessor{
		accrualSystemAddress: accrualSystemAddress,
		httpClient:           &http.Client{},
		invoicesService:      invoicesService,
		ordersService:        ordersService,
	}
}

func (p orderProcessor) RecalculateProcessedOrders(ctx context.Context) error {
	order, err := p.ordersService.GetNotCalculatedOrder(ctx)
	if err != nil {
		return err
	}

	if order == nil || order.Accrual == nil || *order.Accrual == 0 {
		return nil
	}

	invoice, err := p.invoicesService.GetInvoiceByUserID(ctx, order.UserID)
	if err != nil {
		slog.Warn("Recalculate processed orders",
			slog.String("step", "getting invoice"),
			slog.String("order", strconv.Itoa(int(order.UserID))),
			slog.String("error", err.Error()),
		)

		return err
	}

	_, err = p.invoicesService.UpdateBalance(ctx, *invoice, *order.Accrual, order.Number)
	if err != nil {
		slog.Warn("Recalculate processed orders",
			slog.String("step", "update balance of invoice"),
			slog.String("invoice", strconv.Itoa(int(invoice.ID))),
			slog.String("order", strconv.Itoa(int(order.UserID))),
			slog.String("error", err.Error()),
		)

		return err
	}

	_, err = p.ordersService.UpdateOrdersCalculated(ctx, *order, true)
	if err != nil {
		slog.Warn("Recalculate processed orders",
			slog.String("step", "mark order as calculated"),
			slog.String("invoice", strconv.Itoa(int(invoice.ID))),
			slog.String("order", strconv.Itoa(int(order.UserID))),
			slog.String("error", err.Error()),
		)

		return err
	}

	return nil
}

func (p orderProcessor) CheckProcessingOrders(ctx context.Context) error {
	order, err := p.ordersService.GetNotProcessedOrder(ctx)
	if err != nil {
		return err
	}

	if order == nil {
		return nil
	}

	si, err := p.RequestOrdersStatus(ctx, order.Number)
	if err != nil {
		return err
	}

	if si == nil {
		return nil
	}

	_, err = p.ordersService.UpdateOrdersStatus(ctx, *order, si.Status, si.Accrual)
	if err != nil {
		return err
	}

	return nil
}

func (p orderProcessor) RequestOrdersStatus(ctx context.Context, orderNum models.OrderNum) (*models.OrdersStatusInfo, error) {
	url := fmt.Sprintf("%s/api/orders/%s", p.accrualSystemAddress, orderNum)
	response, err := p.httpClient.Get(url)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	osi := &models.OrdersStatusInfo{}
	err = json.NewDecoder(response.Body).Decode(osi)
	if err != nil {
		return nil, err
	}

	return osi, nil
}
