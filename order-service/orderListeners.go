package main

import (
	"context"
	"database/sql"
	"fmt"
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
		// load the next order in the chain if there is one
		fmt.Println("Load the next order if there is one ", order)

		return nil
	default:
		log.Println("order fill: ", error)
		return error
	}
}
