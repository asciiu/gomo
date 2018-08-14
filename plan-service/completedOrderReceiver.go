package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	constPlan "github.com/asciiu/gomo/common/constants/plan"
	"github.com/asciiu/gomo/common/constants/status"
	evt "github.com/asciiu/gomo/common/proto/events"
	notifications "github.com/asciiu/gomo/notification-service/proto"
	repoPlan "github.com/asciiu/gomo/plan-service/db/sql"
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

	planID, depth, err := repoPlan.UpdateOrderStatus(receiver.DB, completedOrderEvent.OrderID, completedOrderEvent.Status)
	if err != nil {
		log.Println("could not update order status -- ", err.Error())
		return nil
	}

	if completedOrderEvent.Status == status.Filled {
		if err := repoPlan.UpdatePlanContext(receiver.DB,
			planID,
			completedOrderEvent.OrderID,
			completedOrderEvent.Exchange,
			completedOrderEvent.MarketName,
			completedOrderEvent.FinalCurrencySymbol,
			completedOrderEvent.FinalCurrencyBalance,
			depth); err != nil {
			log.Println("completed order error trying to update the plan context -- ", err.Error())
			return nil
		}

		// load the child orders of this completed order
		nextPlanOrders, err := repoPlan.FindChildOrders(receiver.DB, planID, completedOrderEvent.OrderID)

		switch {
		case err == sql.ErrNoRows:
			if completedOrderEvent.CloseOnComplete && repoPlan.UpdatePlanStatus(receiver.DB, planID, constPlan.Closed) != nil {
				log.Printf("could not close plan -- %s\n", planID)
			}

		case err != nil:
			log.Println("completed order error on find plan -- ", err.Error())

		default:
			// load new plan order with false - it is not a revision of an active order
			if err := receiver.Service.LoadPlanOrder(ctx, nextPlanOrders, false); err != nil {
				log.Println("could not load the plan orders -- ", err.Error())
			}
		}
	}

	return nil
}
