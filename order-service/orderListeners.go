package main

import (
	"context"
	"database/sql"
	"log"

	evt "github.com/asciiu/gomo/common/proto/events"
	orderRepo "github.com/asciiu/gomo/order-service/db/sql"
)

// OrderFilledReceiver handles order filled events
type OrderFilledReceiver struct {
	DB      *sql.DB
	Service *OrderService
}

// ProcessEvent handles OrderEvents. These events are published by when an order was filled.
func (receiver *OrderFilledReceiver) ProcessEvent(ctx context.Context, orderEvent *evt.OrderEvent) error {

	order, error := orderRepo.UpdateOrderStatus(receiver.DB, orderEvent)
	switch {
	case error == nil:
		order, error = orderRepo.FindOrderWithParentId(receiver.DB, order.OrderId)

		switch {
		case error != nil:
			log.Println("FindOrderWithParentId error ", error.Error())
		case order.Side == "buy":
			receiver.Service.LoadBuyOrder(ctx, order)
		case order.Side == "sell":
			receiver.Service.LoadSellOrder(ctx, order)
		}

		return nil
	default:
		log.Println("order fill: ", error)
		return error
	}
}
