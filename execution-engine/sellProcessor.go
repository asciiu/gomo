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

// SellProcessor will process and handle sell orders.
type SellProcessor struct {
	DB            *sql.DB
	Env           *vm.Env
	Receiver      *OrderReceiver
	Publisher     micro.Publisher
	FillPublisher micro.Publisher
}

// ProcessEvent will process ExchangeEvents. These events are published from the exchange sockets and
// are used as the order trigger.
func (process *SellProcessor) ProcessEvent(ctx context.Context, event *evt.TradeEvent) error {
	sellOrders := process.Receiver.Orders

	for i, sellOrder := range sellOrders {

		marketName := strings.Replace(sellOrder.EventOrigin.MarketName, "-", "", 1)
		if marketName != event.MarketName || sellOrder.EventOrigin.Exchange != event.Exchange {
			continue
		}

		conditions := sellOrder.Conditions
		for _, evaluateFunc := range conditions {

			if isValid, desc := evaluateFunc(event.Price); isValid {
				process.Receiver.Orders = append(sellOrders[:i], sellOrders[i+1:]...)
				// if non simulated trigger buy event - exchange service subscribes to these events
				evt := sellOrder.EventOrigin
				evt.Condition = desc

				// if non simulated trigger buy event - exchange service subscribes to these events
				if err := process.FillPublisher.Publish(ctx, evt); err != nil {
					log.Println("publish warning: ", err)
				}

				// if it is a simulated order trigger an update order event
				if sellOrder.EventOrigin.OrderType == types.VirtualOrder {
					evt.ExchangeOrderID = types.VirtualOrder
					evt.ExchangeMarketName = types.VirtualOrder
					evt.Status = status.Filled

					log.Printf("sell order triggered -- %+v\n", evt)

					if err := process.Publisher.Publish(ctx, evt); err != nil {
						log.Println("publish warning -- ", err)
					}
				}
			}
		}
	}
	//fmt.Println("sell recv: ", event)

	return nil
}
