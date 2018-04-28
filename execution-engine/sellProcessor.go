package main

import (
	"context"
	"database/sql"
	"log"
	"strconv"
	"strings"

	evt "github.com/asciiu/gomo/common/proto/events"
	"github.com/mattn/anko/vm"
	micro "github.com/micro/go-micro"
)

// SellProcessor will process and handle sell orders.
type SellProcessor struct {
	DB        *sql.DB
	Env       *vm.Env
	Receiver  *OrderReceiver
	Publisher micro.Publisher
}

// ProcessEvent will process ExchangeEvents. These events are published from the exchange sockets and
// are used as the order trigger.
func (process *SellProcessor) ProcessEvent(ctx context.Context, event *evt.ExchangeEvent) error {
	sellOrders := process.Receiver.Orders

	for i, sellOrder := range sellOrders {

		marketName := strings.Replace(sellOrder.EventOrigin.MarketName, "-", "", 1)
		if marketName != event.MarketName || sellOrder.EventOrigin.Exchange != event.Exchange {
			continue
		}

		conditions := sellOrder.Conditions
		for _, evaluateFunc := range conditions {
			f, _ := strconv.ParseFloat(event.Price, 64)

			if isValid, desc := evaluateFunc(f); isValid {
				process.Receiver.Orders = append(sellOrders[:i], sellOrders[i+1:]...)
				// if non simulated trigger buy event - exchange service subscribes to these events

				// if it is a simulated order trigger an update order event
				evt := sellOrder.EventOrigin
				evt.ExchangeOrderId = "paper"
				evt.ExchangeMarketName = "paper"
				evt.Status = "filled"
				evt.Condition = desc

				log.Printf("sell order triggered -- %+v\n", evt)

				if err := process.Publisher.Publish(ctx, evt); err != nil {
					log.Println("publish warning -- ", err)
				}
			}
		}
	}
	//fmt.Println("sell recv: ", event)

	return nil
}
