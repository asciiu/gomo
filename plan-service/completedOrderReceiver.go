package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

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
		Description: fmt.Sprintf("orderId: %s status: %s", completedOrderEvent.OrderID, completedOrderEvent.Status),
	}

	// notifiy the user of completed order
	if err := receiver.NotifyPub.Publish(context.Background(), &notification); err != nil {
		log.Println("could not publish notification: ", err)
	}

	order, error := planRepo.UpdateOrderStatus(receiver.DB, completedOrderEvent.OrderID, completedOrderEvent.Status)
	switch {
	case error == nil:
		next, error := planRepo.FindPlanWithOrderID(receiver.DB, order.NextOrderID)

		switch {
		case error == sql.ErrNoRows:
			return nil

		case error != nil:
			log.Println("completed order error on find plan -- ", error.Error())

		default:
			if err := receiver.Service.LoadPlanOrder(ctx, next); err != nil {
				log.Println("completed order error load plan-- ", err.Error())
			}
			if _, err := planRepo.UpdateOrderStatus(receiver.DB, next.Orders[0].OrderID, status.Active); err != nil {
				log.Println("completed order error trying to update order to active -- ", err.Error())
			}
		}

		return nil

	default:
		log.Println("completed order error on update status -- ", error)
		return error
	}
}
