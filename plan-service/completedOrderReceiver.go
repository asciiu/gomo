package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	protoActivity "github.com/asciiu/gomo/activity-bulletin/proto"
	constPlan "github.com/asciiu/gomo/common/constants/plan"
	"github.com/asciiu/gomo/common/constants/status"
	evt "github.com/asciiu/gomo/common/proto/events"
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

	notification := protoActivity.Activity{
		UserID:      completedOrderEvent.UserID,
		ObjectID:    completedOrderEvent.PlanID,
		Type:        "plan",
		Timestamp:   string(pq.FormatTimestamp(time.Now().UTC())),
		Title:       completedOrderEvent.MarketName,
		Subtitle:    completedOrderEvent.Side,
		Description: completedOrderEvent.Details,
		Details:     fmt.Sprintf("{orderID: %s}", completedOrderEvent.OrderID),
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
		now := string(pq.FormatTimestamp(time.Now().UTC()))
		if err := repoPlan.UpdateTriggerResults(receiver.DB,
			completedOrderEvent.TriggerID,
			completedOrderEvent.TriggeredPrice,
			completedOrderEvent.TriggeredCondition,
			now); err != nil {
			log.Println("completed order error trying to update the trigger -- ", err.Error())
			return nil
		}

		if err := repoPlan.UpdateOrderResults(receiver.DB,
			completedOrderEvent.OrderID,
			completedOrderEvent.InitialCurrencyTraded,
			completedOrderEvent.InitialCurrencyRemainder,
			completedOrderEvent.FinalCurrencyBalance,
			completedOrderEvent.FinalCurrencySymbol); err != nil {
			log.Println("completed order error trying to update the order -- ", err.Error())
			return nil
		}

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
