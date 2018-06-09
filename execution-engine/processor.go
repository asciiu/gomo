package main

import (
	"context"
	"database/sql"
	"log"
	"strings"

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
	orders := processor.Receiver.Orders

	for i, order := range orders {

		marketName := strings.Replace(order.EventOrigin.MarketName, "-", "", 1)
		// market name and exchange must match
		if marketName != event.MarketName || order.EventOrigin.Exchange != event.Exchange {
			continue
		}

		conditions := order.Conditions
		// eval all conditions for this order
		for _, evaluateFunc := range conditions {

			// does condition of order eval to true?
			if isValid, desc := evaluateFunc(event.Price); isValid {
				// we have a condition that was true for this order, process this order now!
				log.Printf("order processed -- %+v\n", order)

				// remove this order from the processor
				processor.Receiver.Orders = append(orders[:i], orders[i+1:]...)

				evt := order.EventOrigin
				evt.Condition = desc

				// if non simulated trigger buy event - exchange service subscribes to these events
				if err := processor.Filler.Publish(ctx, evt); err != nil {
					log.Println("publish warning: ", err)
				}

				// if it is a simulated order trigger an update order event
				if order.EventOrigin.OrderType == types.VirtualOrder {
					evt.ExchangeOrderID = types.VirtualOrder
					evt.ExchangeMarketName = types.VirtualOrder
					evt.Status = status.Filled

					log.Printf("order processed -- %+v\n", evt)

					if err := processor.Filled.Publish(ctx, evt); err != nil {
						log.Println("publish warning: ", err, evt)
					}
				}
			}
		}
	}
	//fmt.Println("buy recv: ", event)

	return nil
}
