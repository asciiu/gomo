package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	evt "github.com/asciiu/gomo/common/proto/events"
	notifications "github.com/asciiu/gomo/notification-service/proto"
	"github.com/lib/pq"
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

	notification := notifications.Notification{
		UserID:           completedOrderEvent.UserID,
		ObjectID:         completedOrderEvent.OrderID,
		NotificationType: "order",
		Timestamp:        string(pq.FormatTimestamp(time.Now().UTC())),
		Title:            fmt.Sprintf("%s", completedOrderEvent.MarketName),
		Description:      completedOrderEvent.Details,
	}

	log.Printf("%+v\n", notification)

	// notify the user of completed order
	if err := receiver.NotifyPub.Publish(context.Background(), &notification); err != nil {
		log.Println("could not publish notification: ", err)
	}

	// order, error := planRepo.UpdatePlanOrder(receiver.DB, completedOrderEvent.OrderID, completedOrderEvent.Status)
	// if _, err := planRepo.UpdatePlanStatus(receiver.DB, completedOrderEvent.PlanID, plan.Failed); err != nil {
	// 	log.Println("completed order error trying to update the order status to filled -- ", err.Error())
	// 	return nil
	// }

	// switch {
	// case completedOrderEvent.Status == status.Failed:
	// 	// failed orders should result in a failed plan
	// 	if _, err := planRepo.UpdatePlanStatus(receiver.DB, completedOrderEvent.PlanID, plan.Failed); err != nil {
	// 		log.Println("completed order error trying to update the plan status to completed -- ", err.Error())
	// 	}

	// case completedOrderEvent.Status == status.Filled:
	// 	if err := planRepo.UpdatePlanBalances(receiver.DB, completedOrderEvent.PlanID, completedOrderEvent.BaseBalance, completedOrderEvent.CurrencyBalance); err != nil {
	// 		log.Println("completed order error trying to update the plan balances -- ", err.Error())
	// 		return nil
	// 	}

	// 	nextPlanOrder, error := planRepo.FindChildOrders(receiver.DB, order.PlanID, order.OrderNumber)

	// 	switch {
	// 	case error == sql.ErrNoRows:
	// 		// TODO if order status is failed the plan status should also be failed
	// 		// set plan status to complete
	// 		if _, err := planRepo.UpdatePlanStatus(receiver.DB, completedOrderEvent.PlanID, plan.Completed); err != nil {
	// 			log.Println("completed order error trying to update the plan status to completed -- ", err.Error())
	// 		}

	// 	case error != nil:
	// 		log.Println("completed order error on find plan -- ", error.Error())

	// 	default:
	// 		// load new plan order with false - it is not a revision of an active order
	// 		if err := receiver.Service.LoadPlanOrder(ctx, nextPlanOrder, false); err != nil {
	// 			log.Println("completed order error load plan-- ", err.Error())
	// 		}
	// 		if _, err := planRepo.UpdatePlanOrder(receiver.DB, nextPlanOrder.Orders[0].OrderID, status.Active); err != nil {
	// 			log.Println("completed order error trying to update order to active -- ", err.Error())
	// 		}
	// 	}

	// default:
	// 	log.Println("completed order error on update status -- ", error)
	// 	return error
	// }
	return nil
}
