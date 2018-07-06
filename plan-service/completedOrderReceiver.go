package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	evt "github.com/asciiu/gomo/common/proto/events"
	notifications "github.com/asciiu/gomo/notification-service/proto"
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
	return nil

	// // publish notification about order fill
	// if err := receiver.NotifyPub.Publish(context.Background(), &notification); err != nil {
	// 	log.Println("could not publish notification: ", err)
	// }

	// parentOrder, error := orderRepo.UpdateOrderStatus(receiver.DB, orderEvent)
	// switch {
	// case error == nil:
	// 	childOrder, error := orderRepo.FindOrderWithParentID(receiver.DB, parentOrder.OrderID)

	// 	switch {
	// 	case error == sql.ErrNoRows:
	// 		return nil

	// 	case error != nil:
	// 		log.Println("order filled error -- ", error.Error())

	// 	default:
	// 		if err := receiver.Service.LoadOrder(ctx, childOrder); err != nil {
	// 			log.Println("order filled error -- ", err.Error())
	// 		}
	// 	}

	// 	return nil

	// default:
	// 	log.Println("order filled error -- ", error)
	// 	return error
	// }
}
