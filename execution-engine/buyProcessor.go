package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
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
func (process *BuyProcessor) ProcessEvent(ctx context.Context, event *evt.ExchangeEvent) error {
	buyOrders := process.Receiver.Orders

	for i, buyOrder := range buyOrders {

		marketName := strings.Replace(buyOrder.EventOrigin.MarketName, "-", "", 1)
		// market name and exchange must match
		if marketName != event.MarketName || buyOrder.EventOrigin.Exchange != event.Exchange {
			continue
		}

		conditions := buyOrder.Conditions
		// eval all conditions for this order
		for _, evaluateFunc := range conditions {
			f, _ := strconv.ParseFloat(event.Price, 64)

			// does condition of order eval to true?
			if isValid, desc := evaluateFunc(f); isValid {
				// remove this order from the process
				process.Receiver.Orders = append(buyOrders[:i], buyOrders[i+1:]...)

				// trigger buy order event
				fmt.Println("BUY NOW!! ", buyOrder)
				// if non simulated trigger buy event - exchange service subscribes to these events

				// if it is a simulated order trigger an update order event
				evt := buyOrder.EventOrigin
				evt.ExchangeOrderId = "simulated"
				evt.ExchangeMarketName = evt.MarketName
				evt.Status = "filled"
				evt.Condition = desc

				if err := process.Publisher.Publish(ctx, &evt); err != nil {
					log.Println("publish warning: ", err, evt)
				}
			}
		}
	}
	//fmt.Println("buy recv: ", event)

	return nil
}
