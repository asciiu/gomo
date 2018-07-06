package main

import (
	"context"
	"log"
	"time"

	evt "github.com/asciiu/gomo/common/proto/events"
	planRepo "github.com/asciiu/gomo/plan-service/db/sql"
)

// EngineStartReceiver will listen to new engine instances and will load the engine instances
// with open orders. This is how the engine communicates to the order service that it needs
// to receive open orders.
type EngineStartReceiver struct {
	Service *PlanService
}

// ProcessEvent handles OrderEvents. These events are published by when an order was filled.
func (receiver *EngineStartReceiver) ProcessEvent(ctx context.Context, engine *evt.EngineStartEvent) error {

	plans, error := planRepo.FindActivePlans(receiver.Service.DB)
	if error != nil {
		log.Println("could not find acive plans -- ", error)
	}

	// must sleep before sending off to execution engine
	// because engine might not have fully started yet
	time.Sleep(2 * time.Second)

	log.Println(plans)

	// // TODO we need to explore a different approach here that is more efficient.
	// for _, order := range orders {
	// 	if error = receiver.Service.LoadOrder(ctx, order); error != nil {
	// 		log.Println("load order error -- ", error)
	// 	} else {
	// 		log.Printf("loaded order -- %+v\n", order)
	// 	}
	// }

	return nil
}
