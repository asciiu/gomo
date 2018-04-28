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

	log.Printf("order filled -- %+v\n", orderEvent)

	parentOrder, error := orderRepo.UpdateOrderStatus(receiver.DB, orderEvent)
	switch {
	case error == nil:
		childOrder, error := orderRepo.FindOrderWithParentId(receiver.DB, parentOrder.OrderId)

		switch {
		case error != nil:
			log.Println("FindOrderWithParentId error ", error.Error())
		case childOrder.Side == "buy":
			receiver.Service.LoadBuyOrder(ctx, childOrder)
		case childOrder.Side == "sell":
			receiver.Service.LoadSellOrder(ctx, childOrder)
		}

		return nil
	default:
		log.Println("order fill error: ", error)
		return error
	}
}
