package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/asciiu/gomo/common/constants/plan"
	"github.com/asciiu/gomo/common/constants/status"
	evt "github.com/asciiu/gomo/common/proto/events"
	planRepo "github.com/asciiu/gomo/plan-service/db/sql"
)

// AbortedOrderReceiver handles aborted events
type AbortedOrderReceiver struct {
	DB      *sql.DB
	Service *PlanService
}

// ProcessEvent handles AbortedOrderEvent. These events are published by when an order was filled.
func (receiver *AbortedOrderReceiver) ProcessAbortedEvent(ctx context.Context, abortedOrderEvent *evt.AbortedOrderEvent) error {

	log.Printf("aborted event received -- %+v\n", abortedOrderEvent)

	pln, error := planRepo.FindPlanWithOrderID(receiver.DB, abortedOrderEvent.OrderID)
	if error != nil {
		log.Printf("error in ProcessAbortedEvent when querying for plan: %s\n", error.Error())
	}

	if pln.ActiveOrderNumber == 0 {
		// the active order that was aborted was the first order - safe to delete plan
		error = planRepo.DeletePlan(receiver.DB, abortedOrderEvent.PlanID)
	} else {
		_, error = planRepo.UpdatePlanStatus(receiver.DB, abortedOrderEvent.PlanID, plan.Aborted)
		_, error = planRepo.UpdatePlanOrder(receiver.DB, abortedOrderEvent.OrderID, status.Aborted)
	}

	return nil
}
