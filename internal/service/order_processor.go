package service

import (
	"context"
)

type orderProcessor struct {
	invoicesService *InvoicesService
	ordersService   *OrdersService
}

func NewOrdersProcessor(invoicesService *InvoicesService, ordersService *OrdersService) OrdersProcessor {
	return &orderProcessor{
		invoicesService: invoicesService,
		ordersService:   ordersService,
	}
}

func (p orderProcessor) RecalculateProcessedOrders(ctx context.Context) error {
	//ordersList, err := ordersService.GetNotCalculatedOrder(ctx, 1000)
	//if err != nil {
	//	return err
	//}
	//for _, order := range *ordersList {
	//	if *order.Accrual == 0 {
	//		continue
	//	}
	//
	//	account, err := invoicesService.
	//		invoicesService.UpdateBalance(ctx)
	//}

	return nil
}
