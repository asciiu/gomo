package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	evt "github.com/asciiu/gomo/common/proto/events"
	"github.com/mattn/anko/vm"
)

type SellProcessor struct {
	DB       *sql.DB
	Env      *vm.Env
	Receiver *OrderReceiver
}

func (process *SellProcessor) ProcessEvent(ctx context.Context, event *evt.ExchangeEvent) error {
	sellOrders := process.Receiver.Orders

	for i, sellOrder := range sellOrders {

		marketName := strings.Replace(sellOrder.MarketName, "-", "", 1)
		if marketName != event.MarketName || sellOrder.Exchange != event.Exchange {
			continue
		}

		conditions := strings.Replace(sellOrder.Conditions, "price", event.Price, -1)

		result, err := process.Env.Execute(conditions)
		if err != nil {
			panic(err)
		}

		if result == true {
			// remove order
			process.Receiver.Orders = append(sellOrders[:i], sellOrders[i+1:]...)
			fmt.Println("Sell NOW!!")
		}
	}
	//fmt.Println("sell recv: ", event)

	return nil
}
