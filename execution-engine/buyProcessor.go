package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	evt "github.com/asciiu/gomo/common/proto/events"
	"github.com/mattn/anko/vm"
	micro "github.com/micro/go-micro"
)

// BuyProcessor will handle all buys
type BuyProcessor struct {
	DB        *sql.DB
	Env       *vm.Env
	Receiver  *OrderReceiver
	Publisher micro.Publisher
}

// ProcessEvent will process ExchangeEvents. These events are published from the exchange sockets.
func (process *BuyProcessor) ProcessEvent(ctx context.Context, event *evt.TradeEvent) error {
	buyOrders := process.Receiver.Orders

	fmt.Println(event)

	for i, buyOrder := range buyOrders {

		marketName := strings.Replace(buyOrder.EventOrigin.MarketName, "-", "", 1)
		// market name and exchange must match
		if marketName != event.MarketName || buyOrder.EventOrigin.Exchange != event.Exchange {
			continue
		}

		conditions := buyOrder.Conditions
		// eval all conditions for this order
		for _, evaluateFunc := range conditions {

			// does condition of order eval to true?
			if isValid, desc := evaluateFunc(event.Price); isValid {
				// remove this order from the process
				process.Receiver.Orders = append(buyOrders[:i], buyOrders[i+1:]...)

				// if non simulated trigger buy event - exchange service subscribes to these events

				// if it is a simulated order trigger an update order event
				evt := buyOrder.EventOrigin
				evt.ExchangeOrderID = "paper"
				evt.ExchangeMarketName = "paper"
				evt.Status = "filled"
				evt.Condition = desc

				log.Printf("buy order triggered -- %+v\n", evt)

				if err := process.Publisher.Publish(ctx, evt); err != nil {
					log.Println("publish warning: ", err, evt)
				}
			}
		}
	}
	//fmt.Println("buy recv: ", event)

	return nil
}
