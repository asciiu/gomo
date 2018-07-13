package main

import (
	"context"
	"database/sql"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/asciiu/gomo/common/constants/status"
	evt "github.com/asciiu/gomo/common/proto/events"
	"github.com/mattn/anko/vm"
	micro "github.com/micro/go-micro"
)

// Order has conditions
type Order struct {
	EventOrigin *evt.ActiveOrderEvent
	Conditions  []ConditionFunc
}

// OrderReceiver will receive and prep an order conditions
type OrderReceiver struct {
	DB      *sql.DB
	Orders  []*Order
	Env     *vm.Env
	Aborted micro.Publisher
}

// ProcessEvent handles ActiveOrderEvents. These events are published by the plan service.
func (receiver *OrderReceiver) ProcessEvent(ctx context.Context, activeOrder *evt.ActiveOrderEvent) error {
	// convert OrderEvent to Order with conditions here
	strConditions := strings.Split(activeOrder.Conditions, " or ")
	conditions := make([]ConditionFunc, 0)

	trailingPoint := regexp.MustCompile(`^.*?TrailingStopPoint\((0\.\d{2,}),\s(\d+\.\d+).*?`)
	trailingPercent := regexp.MustCompile(`^.*?TrailingStopPercent\((0\.\d{2,}),\s(\d+\.\d+).*?`)

	for _, str := range strConditions {
		switch {
		case trailingPoint.MatchString(str):
			rs := trailingPoint.FindStringSubmatch(str)
			top, _ := strconv.ParseFloat(rs[1], 64)
			points, _ := strconv.ParseFloat(rs[2], 64)

			ts := TrailingStopPoint{
				Top:    top,
				Points: points,
			}
			conditions = append(conditions, (&ts).evaluate)

		case trailingPercent.MatchString(str):
			rs := trailingPercent.FindStringSubmatch(str)
			top, _ := strconv.ParseFloat(rs[1], 64)
			percent, _ := strconv.ParseFloat(rs[2], 64)

			ts := TrailingStopPercent{
				Top:     top,
				Percent: percent,
			}
			conditions = append(conditions, (&ts).evaluate)

		default:

			priceCond := PriceCondition{
				Env:       receiver.Env,
				Statement: str,
			}
			conditions = append(conditions, (&priceCond).evaluate)
		}
	}

	order := Order{
		EventOrigin: activeOrder,
		Conditions:  conditions,
	}

	// if we have an order revision we must remove the original active order from memory
	if activeOrder.Revision {
		orders := receiver.Orders
		for i, order := range orders {
			if order.EventOrigin.OrderID == activeOrder.OrderID {
				// remove this order from the processor
				receiver.Orders = append(orders[:i], orders[i+1:]...)

				if activeOrder.OrderStatus == status.Aborted {
					abortedEvent := evt.AbortedOrderEvent{
						OrderID: order.EventOrigin.OrderID,
						PlanID:  order.EventOrigin.PlanID,
					}

					// Never log the secrets contained in the event
					log.Printf("aborted order -- orderID: %s, market: %s\n", order.EventOrigin.OrderID, order.EventOrigin.MarketName)
					// notify the plan service of abort success
					if err := receiver.Aborted.Publish(ctx, &abortedEvent); err != nil {
						log.Println("publish abort warning: ", err)
					}
					return nil
				}
			}
		}
	}

	receiver.Orders = append(receiver.Orders, &order)

	if activeOrder.Revision {
		log.Printf("revision to order received -- orderID: %s, market: %s\n", order.EventOrigin.OrderID, order.EventOrigin.MarketName)
	} else {
		log.Printf("new order received -- orderID: %s, market: %s\n", order.EventOrigin.OrderID, order.EventOrigin.MarketName)
	}

	return nil
}
