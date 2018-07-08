package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	plan "github.com/asciiu/gomo/common/constants/plan"
	"github.com/asciiu/gomo/common/constants/status"
	evt "github.com/asciiu/gomo/common/proto/events"
	notifications "github.com/asciiu/gomo/notification-service/proto"
	planRepo "github.com/asciiu/gomo/plan-service/db/sql"
	micro "github.com/micro/go-micro"
)

// OrderFilledReceiver handles order filled events
type CompletedOrderReceiver struct {
	DB        *sql.DB
	Service   *PlanService
	NotifyPub micro.Publisher
}

// ProcessEvent handles OrderEvents. These events are published by when an order was filled.
func (receiver *CompletedOrderReceiver) ProcessEvent(ctx context.Context, completedOrderEvent *evt.CompletedOrderEvent) error {

	log.Printf("order completed -- %+v\n", completedOrderEvent)

	notification := notifications.Notification{
		UserID:      completedOrderEvent.UserID,
		Description: fmt.Sprintf("orderId: %s %s", completedOrderEvent.OrderID, completedOrderEvent.Details),
	}

	// notifiy the user of completed order
	if err := receiver.NotifyPub.Publish(context.Background(), &notification); err != nil {
		log.Println("could not publish notification: ", err)
	}

	order, error := planRepo.UpdateOrderStatus(receiver.DB, completedOrderEvent.OrderID, completedOrderEvent.Status)
	if _, err := planRepo.UpdatePlanStatus(receiver.DB, completedOrderEvent.PlanID, plan.Failed); err != nil {
		log.Println("completed order error trying to update the order status to filled -- ", err.Error())
		return nil
	}

	switch {
	case completedOrderEvent.Status == status.Failed:
		// failed orders should result in a failed plan
		if _, err := planRepo.UpdatePlanStatus(receiver.DB, completedOrderEvent.PlanID, plan.Failed); err != nil {
			log.Println("completed order error trying to update the plan status to completed -- ", err.Error())
		}

	case completedOrderEvent.Status == status.Filled:
		if err := planRepo.UpdatePlanBalances(receiver.DB, completedOrderEvent.PlanID, completedOrderEvent.BaseBalance, completedOrderEvent.CurrencyBalance); err != nil {
			log.Println("completed order error trying to update the plan balances -- ", err.Error())
			return nil
		}

		nextPlanOrder, error := planRepo.FindPlanWithOrderID(receiver.DB, order.NextOrderID)

		switch {
		case error == sql.ErrNoRows:
			// TODO if order status is failed the plan status should also be failed
			// set plan status to complete
			if _, err := planRepo.UpdatePlanStatus(receiver.DB, completedOrderEvent.PlanID, plan.Completed); err != nil {
				log.Println("completed order error trying to update the plan status to completed -- ", err.Error())
			}

		case error != nil:
			log.Println("completed order error on find plan -- ", error.Error())

		default:
			if err := receiver.Service.LoadPlanOrder(ctx, nextPlanOrder); err != nil {
				log.Println("completed order error load plan-- ", err.Error())
			}
			if _, err := planRepo.UpdateOrderStatus(receiver.DB, nextPlanOrder.Orders[0].OrderID, status.Active); err != nil {
				log.Println("completed order error trying to update order to active -- ", err.Error())
			}
		}

	default:
		log.Println("completed order error on update status -- ", error)
		return error
	}
	return nil
}
