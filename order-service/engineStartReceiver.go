package main

import (
	"context"
	"log"

	evt "github.com/asciiu/gomo/common/proto/events"
	orderRepo "github.com/asciiu/gomo/order-service/db/sql"
)

// EngineStartReceiver will listen to new engine instances and will load the engine instances
// with open orders. This is how the engine communicates to the order service that it needs
// to receive open orders.
type EngineStartReceiver struct {
	Service *OrderService
}

// ProcessEvent handles OrderEvents. These events are published by when an order was filled.
func (receiver *EngineStartReceiver) ProcessEvent(ctx context.Context, engine *evt.EngineStartEvent) error {

	orders, error := orderRepo.FindOpenOrders(receiver.Service.DB)
	if error != nil {
		log.Println("could not find open orders: ", error)
	}

	log.Println(orders)

	//log.Println("engine started load the open orders for: ", engine)
	return nil
}
