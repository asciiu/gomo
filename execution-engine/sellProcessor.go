package main

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	evt "github.com/asciiu/gomo/common/proto/events"
	"github.com/mattn/anko/vm"
	micro "github.com/micro/go-micro"
)

type SellProcessor struct {
	DB        *sql.DB
	Env       *vm.Env
	Receiver  *OrderReceiver
	Publisher micro.Publisher
}

func (process *SellProcessor) ProcessEvent(ctx context.Context, event *evt.ExchangeEvent) error {
	sellOrders := process.Receiver.Orders

	for i, sellOrder := range sellOrders {

		marketName := strings.Replace(sellOrder.EventOrigin.MarketName, "-", "", 1)
		if marketName != event.MarketName || sellOrder.EventOrigin.Exchange != event.Exchange {
			continue
		}

		conditions := sellOrder.Conditions
		for _, cond := range conditions {
			f, _ := strconv.ParseFloat(event.Price, 64)

			if cond(f) {
				process.Receiver.Orders = append(sellOrders[:i], sellOrders[i+1:]...)
				fmt.Println("SELL NOW!!")
				fmt.Println(sellOrder)
			}
		}
	}
	//fmt.Println("sell recv: ", event)

	return nil
}
