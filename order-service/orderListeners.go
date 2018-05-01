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
		childOrder, error := orderRepo.FindOrderWithParentID(receiver.DB, parentOrder.OrderID)

		switch {
		case error == sql.ErrNoRows:
			return nil

		case error != nil:
			log.Println("order filled error -- ", error.Error())

		case childOrder.Side == "buy":
			if err := receiver.Service.LoadBuyOrder(ctx, childOrder); err != nil {
				log.Println("order filled error -- ", err.Error())
			}

		case childOrder.Side == "sell":
			if err := receiver.Service.LoadSellOrder(ctx, childOrder); err != nil {
				log.Println("order filled error -- ", err.Error())
			}
		}

		return nil

	default:
		log.Println("order filled error -- ", error)
		return error
	}
}
