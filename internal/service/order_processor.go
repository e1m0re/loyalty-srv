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

func (p orderProcessor) CheckProcessingOrders(ctx context.Context) (timeout int64, err error) {
	order, err := p.ordersService.GetNotProcessedOrder(ctx)
	if err != nil {

		return 0, err
	}

	if order == nil {

		return 0, nil
	}

	si, timeout, err := p.RequestOrdersStatus(ctx, order.Number)
	if err != nil {

		return timeout, err
	}

	if si == nil {

		return 0, nil
	}

	_, err = p.ordersService.UpdateOrdersStatus(ctx, *order, si.Status, si.Accrual)
	if err != nil {

		return 0, err
	}

	return 0, nil
}

func (p orderProcessor) RequestOrdersStatus(ctx context.Context, orderNum models.OrderNum) (osi *models.OrdersStatusInfo, timeout int64, err error) {
	url := fmt.Sprintf("%s/api/orders/%s", p.accrualSystemAddress, orderNum)
	response, err := p.httpClient.Get(url)
	if err != nil {

		return nil, 0, err
	}

	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusTooManyRequests:
		if s, ok := response.Header["Retry-After"]; ok {
			if timeout, err = strconv.ParseInt(s[0], 10, 32); err != nil {
				return nil, 0, nil
			}
		}

		return nil, timeout, nil
	case http.StatusNoContent:
		return nil, 0, nil
	case http.StatusOK:
		osi = &models.OrdersStatusInfo{}
		err = json.NewDecoder(response.Body).Decode(osi)
		if err != nil {
			return nil, 0, err
		}

		return osi, 0, nil
	case http.StatusInternalServerError:
		var errMsg []byte
		response.Body.Read(errMsg)

		return nil, 0, fmt.Errorf("internal server error: %s", errMsg)
	default:

		return nil, 0, fmt.Errorf("unknow resonse type")
	}
}
