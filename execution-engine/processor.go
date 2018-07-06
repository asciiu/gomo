package main

import (
	"context"
	"database/sql"
	"log"

	types "github.com/asciiu/gomo/common/constants/order"
	"github.com/asciiu/gomo/common/constants/status"
	evt "github.com/asciiu/gomo/common/proto/events"
	"github.com/mattn/anko/vm"
	micro "github.com/micro/go-micro"
)

// Processor will process orders
type Processor struct {
	DB        *sql.DB
	Env       *vm.Env
	Receiver  *OrderReceiver
	Completed micro.Publisher
	Triggered micro.Publisher
}

// ProcessEvent will process ExchangeEvents. These events are published from the exchange sockets.
func (processor *Processor) ProcessEvent(ctx context.Context, payload *evt.TradeEvents) error {
	//log.Println("buy recv: ", event)
	orders := processor.Receiver.Orders

	// every order check trade event price with order conditions
	for i, order := range orders {
		for _, event := range payload.Events {
			marketName := order.EventOrigin.MarketName
			//marketName := strings.Replace(order.EventOrigin.MarketName, "-", "", 1)
			// market name and exchange must match
			if marketName != event.MarketName || order.EventOrigin.Exchange != event.Exchange {
				continue
			}

			conditions := order.Conditions
			// eval all conditions for this order
			for _, evaluateFunc := range conditions {

				// does condition of order eval to true?
				if isValid, desc := evaluateFunc(event.Price); isValid {
					// remove this order from the processor
					processor.Receiver.Orders = append(orders[:i], orders[i+1:]...)

					switch {
					case order.EventOrigin.OrderType == types.VirtualOrder:
						completedEvent := evt.CompletedOrderEvent{
							UserID:             order.EventOrigin.UserID,
							OrderID:            order.EventOrigin.OrderID,
							Side:               order.EventOrigin.Side,
							TriggeredPrice:     event.Price,
							TriggeredCondition: desc,
							ExchangeOrderID:    types.VirtualOrder,
							ExchangeMarketName: types.VirtualOrder,
							Status:             status.Filled,
						}

						// Never log the secrets contained in the event
						log.Printf("virtual order processed -- orderID: %s, market: %s\n", order.EventOrigin.OrderID, order.EventOrigin.MarketName)

						if err := processor.Completed.Publish(ctx, &completedEvent); err != nil {
							log.Println("publish warning: ", err, completedEvent)
						}
					default:
						// convert this active order event to a triggered order event
						triggeredEvent := evt.TriggeredOrderEvent{
							Exchange:           order.EventOrigin.Exchange,
							OrderID:            order.EventOrigin.OrderID,
							PlanID:             order.EventOrigin.PlanID,
							UserID:             order.EventOrigin.UserID,
							BaseBalance:        order.EventOrigin.BaseBalance,
							BasePercent:        order.EventOrigin.BasePercent,
							CurrencyBalance:    order.EventOrigin.CurrencyBalance,
							CurrencyPercent:    order.EventOrigin.CurrencyPercent,
							KeyID:              order.EventOrigin.KeyID,
							Key:                order.EventOrigin.Key,
							Secret:             order.EventOrigin.Secret,
							MarketName:         order.EventOrigin.MarketName,
							Side:               order.EventOrigin.Side,
							OrderType:          order.EventOrigin.OrderType,
							Price:              order.EventOrigin.Price,
							TriggeredPrice:     event.Price,
							TriggeredCondition: desc,
							NextOrders:         order.EventOrigin.NextOrders,
						}

						// Never log the secrets contained in the event
						log.Printf("triggered order -- orderID: %s, market: %s\n", order.EventOrigin.OrderID, order.EventOrigin.MarketName)
						// if non simulated trigger buy event - exchange service subscribes to these events
						if err := processor.Triggered.Publish(ctx, &triggeredEvent); err != nil {
							log.Println("publish warning: ", err)
						}
					}
				}
			}
		}
	}

	return nil
}
