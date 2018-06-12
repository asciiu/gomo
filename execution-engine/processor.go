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
	DB       *sql.DB
	Env      *vm.Env
	Receiver *OrderReceiver
	Filled   micro.Publisher
	Filler   micro.Publisher
}

// ProcessEvent will process ExchangeEvents. These events are published from the exchange sockets.
func (processor *Processor) ProcessEvent(ctx context.Context, event *evt.TradeEvent) error {
	//log.Println("buy recv: ", event)
	orders := processor.Receiver.Orders

	for i, order := range orders {

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

				event := order.EventOrigin
				event.Condition = desc

				switch {
				case order.EventOrigin.OrderType == types.VirtualOrder:
					event.ExchangeOrderID = types.VirtualOrder
					event.ExchangeMarketName = types.VirtualOrder
					event.Status = status.Filled

					// Never log the secrets contained in the event
					log.Printf("virtual order processed -- orderID: %s, market: %s\n", event.OrderID, event.MarketName)

					if err := processor.Filled.Publish(ctx, event); err != nil {
						log.Println("publish warning: ", err, event)
					}
				default:
					// Never log the secrets contained in the event
					log.Printf("order processed -- orderID: %s, market: %s\n", event.OrderID, event.MarketName)
					// if non simulated trigger buy event - exchange service subscribes to these events
					if err := processor.Filler.Publish(ctx, event); err != nil {
						log.Println("publish warning: ", err)
					}
				}
			}
		}
	}

	return nil
}
