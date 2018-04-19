package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	evt "github.com/asciiu/gomo/common/proto/events"
	"github.com/mattn/anko/vm"
)

type BuyProcessor struct {
	DB       *sql.DB
	Env      *vm.Env
	Receiver *OrderReceiver
}

func (process *BuyProcessor) ProcessEvent(ctx context.Context, event *evt.ExchangeEvent) error {
	buyOrders := process.Receiver.Orders

	for i, buyOrder := range buyOrders {

		marketName := strings.Replace(buyOrder.MarketName, "-", "", 1)
		if marketName != event.MarketName || buyOrder.Exchange != event.Exchange {
			continue
		}

		conditions := strings.Replace(buyOrder.Conditions, "price", event.Price, -1)

		result, err := process.Env.Execute(conditions)
		if err != nil {
			panic(err)
		}

		if result == true {
			// remove order
			process.Receiver.Orders = append(buyOrders[:i], buyOrders[i+1:]...)
			fmt.Println("BUY NOW!!")
			fmt.Println(buyOrder)
		}
	}
	//fmt.Println("buy recv: ", event)

	return nil
}
