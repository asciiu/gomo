package main

import (
	"context"
	"database/sql"
	"log"
	"regexp"
	"strconv"

	evt "github.com/asciiu/gomo/common/proto/events"
	"github.com/mattn/anko/vm"
	micro "github.com/micro/go-micro"
)

// Order has conditions
type Plan struct {
	PlanID                string
	UserID                string
	ActiveCurrencySymbol  string
	ActiveCurrencyBalance float64
	Orders                []*Order
}

type Order struct {
	OrderID     string
	Exchange    string
	MarketName  string
	Side        string
	LimitPrice  float64
	OrderType   string
	OrderStatus string
	KeyID       string
	KeyPublic   string
	KeySecret   string
	TriggerExs  []*TriggerEx
}

type TriggerEx struct {
	TriggerID string
	OrderID   string
	Name      string
	Triggered bool
	Actions   []string
	Evaluate  Expression
}

// This approach may not be feasible because there may be
// multiple orders with different market names in a single plan
type MarketPlans struct {
	// maps market name to plans for that market
	// example: ADA-BTC -> plans for ADA-BTC
	MarketPlans map[string][]*Plan
}

// OrderReceiver will receive and prep an order conditions
type PlanReceiver struct {
	DB    *sql.DB
	Plans []*Plan

	// example: binance -> (binance only market plans)
	//ExchangeMarkets map[string]MarketPlans
	Env     *vm.Env
	Aborted micro.Publisher
}

// ProcessEvent handles ActiveOrderEvents. These events are published by the plan service.
func (receiver *PlanReceiver) AddPlan(ctx context.Context, plan *evt.NewPlanEvent) error {
	log.Printf("received new plan event for plan ID: %s", plan.PlanID)

	// convert plan orders to have trigger expressions
	orders := make([]*Order, 0)

	trailingPoint := regexp.MustCompile(`^.*?TrailingStopPoint\((0\.\d{2,}),\s(\d+\.\d+).*?`)
	trailingPercent := regexp.MustCompile(`^.*?TrailingStopPercent\((0\.\d{2,}),\s(\d+\.\d+).*?`)
	for _, order := range plan.Orders {

		// convert trigger to trigger expression
		expressions := make([]*TriggerEx, 0)

		for _, trigger := range order.Triggers {
			str := trigger.Code
			var expression Expression
			switch {
			case trailingPoint.MatchString(str):
				rs := trailingPoint.FindStringSubmatch(str)
				top, _ := strconv.ParseFloat(rs[1], 64)
				points, _ := strconv.ParseFloat(rs[2], 64)

				ts := TrailingStopPoint{
					Top:    top,
					Points: points,
				}
				expression = (&ts).evaluate
			case trailingPercent.MatchString(str):
				rs := trailingPercent.FindStringSubmatch(str)
				top, _ := strconv.ParseFloat(rs[1], 64)
				percent, _ := strconv.ParseFloat(rs[2], 64)

				ts := TrailingStopPercent{
					Top:     top,
					Percent: percent,
				}
				expression = (&ts).evaluate
			default:
				priceCond := PriceCondition{
					Env:       receiver.Env,
					Statement: str,
				}
				expression = (&priceCond).evaluate
			}

			expressions = append(expressions, &TriggerEx{
				TriggerID: trigger.TriggerID,
				OrderID:   trigger.OrderID,
				Name:      trigger.Name,
				Triggered: trigger.Triggered,
				Actions:   trigger.Actions,
				Evaluate:  expression,
			})
		}
		orders = append(orders, &Order{
			OrderID:     order.OrderID,
			Exchange:    order.Exchange,
			MarketName:  order.MarketName,
			Side:        order.Side,
			LimitPrice:  order.LimitPrice,
			OrderType:   order.OrderType,
			OrderStatus: order.OrderStatus,
			KeyID:       order.KeyID,
			KeyPublic:   order.KeyPublic,
			KeySecret:   order.KeySecret,
			TriggerExs:  expressions,
		})
	}

	receiver.Plans = append(receiver.Plans, &Plan{
		PlanID:                plan.PlanID,
		UserID:                plan.UserID,
		ActiveCurrencySymbol:  plan.ActiveCurrencySymbol,
		ActiveCurrencyBalance: plan.ActiveCurrencyBalance,
		Orders:                orders,
	})

	// TODO must replace existing plan if this plan ID already exists
	// // if we have an order revision we must remove the original active order from memory
	// if activeOrder.Revision {
	// 	orders := receiver.Orders
	// 	for i, order := range orders {
	// 		if order.EventOrigin.OrderID == activeOrder.OrderID {
	// 			// remove this order from the processor
	// 			receiver.Orders = append(orders[:i], orders[i+1:]...)

	// 			if activeOrder.OrderStatus == status.Aborted {
	// 				abortedEvent := evt.AbortedOrderEvent{
	// 					OrderID: order.EventOrigin.OrderID,
	// 					PlanID:  order.EventOrigin.PlanID,
	// 				}

	// 				// Never log the secrets contained in the event
	// 				log.Printf("aborted order -- orderID: %s, market: %s\n", order.EventOrigin.OrderID, order.EventOrigin.MarketName)
	// 				// notify the plan service of abort success
	// 				if err := receiver.Aborted.Publish(ctx, &abortedEvent); err != nil {
	// 					log.Println("publish abort warning: ", err)
	// 				}
	// 				return nil
	// 			}
	// 		}
	// 	}
	// }

	// receiver.Orders = append(receiver.Orders, &order)

	// if activeOrder.Revision {
	// 	log.Printf("revision to order received -- orderID: %s, market: %s\n", order.EventOrigin.OrderID, order.EventOrigin.MarketName)
	// } else {
	// 	log.Printf("new order received -- orderID: %s, market: %s\n", order.EventOrigin.OrderID, order.EventOrigin.MarketName)
	// }

	return nil
}
