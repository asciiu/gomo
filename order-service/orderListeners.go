package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/asciiu/gomo/common/constants/side"
	evt "github.com/asciiu/gomo/common/proto/events"
	notifications "github.com/asciiu/gomo/notification-service/proto"
	orderRepo "github.com/asciiu/gomo/order-service/db/sql"
	micro "github.com/micro/go-micro"
)

// OrderFilledReceiver handles order filled events
type OrderFilledReceiver struct {
	DB        *sql.DB
	Service   *OrderService
	NotifyPub micro.Publisher
}

// ProcessEvent handles OrderEvents. These events are published by when an order was filled.
func (receiver *OrderFilledReceiver) ProcessEvent(ctx context.Context, orderEvent *evt.OrderEvent) error {

	log.Printf("order filled -- %+v\n", orderEvent)

	notification := notifications.Notification{
		UserID:      orderEvent.UserID,
		Description: fmt.Sprintf("orderId: %s status: %s", orderEvent.OrderID, orderEvent.Status),
	}

	// publish verify key event
	if err := receiver.NotifyPub.Publish(context.Background(), &notification); err != nil {
		log.Println("could not publish notification: ", err)
	}

	parentOrder, error := orderRepo.UpdateOrderStatus(receiver.DB, orderEvent)
	switch {
	case error == nil:
		childOrder, error := orderRepo.FindOrderWithParentID(receiver.DB, parentOrder.OrderID)

		switch {
		case error == sql.ErrNoRows:
			return nil

		case error != nil:
			log.Println("order filled error -- ", error.Error())

		case childOrder.Side == side.Buy:
			if err := receiver.Service.LoadBuyOrder(ctx, childOrder); err != nil {
				log.Println("order filled error -- ", err.Error())
			}

		case childOrder.Side == side.Sell:
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
