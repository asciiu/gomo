package main

import (
	"context"
	"log"

	evt "github.com/asciiu/gomo/common/proto/events"
)

// EngineStartReceiver will listen to new engine instances and will load the engine instances
// with open orders. This is how the engine communicates to the order service that it needs
// to receive open orders.
type EngineStartReceiver struct {
	Service *OrderService
}

// ProcessEvent handles OrderEvents. These events are published by when an order was filled.
func (receiver *EngineStartReceiver) ProcessEvent(ctx context.Context, engine *evt.EngineStartEvent) error {

	log.Println("engine started load the open orders for: ", engine)
	return nil
}
