package main

import (
	"context"
	"log"
	"time"

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

	orders, error := orderRepo.FindActiveOrders(receiver.Service.DB)
	if error != nil {
		log.Println("could not find open orders -- ", error)
	}

	time.Sleep(2 * time.Second)

	// TODO we need to explore a different approach here that is more efficient.
	for _, order := range orders {
		if error = receiver.Service.LoadOrder(ctx, order); error != nil {
			log.Println("load order error -- ", error)
		} else {
			log.Printf("loaded order -- %+v\n", order)
		}
	}

	return nil
}
