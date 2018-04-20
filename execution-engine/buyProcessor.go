package main

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
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

		marketName := strings.Replace(buyOrder.EventOrigin.MarketName, "-", "", 1)
		if marketName != event.MarketName || buyOrder.EventOrigin.Exchange != event.Exchange {
			continue
		}

		conditions := buyOrder.Conditions
		for _, evaluate := range conditions {
			f, _ := strconv.ParseFloat(event.Price, 64)

			if evaluate(f) {
				process.Receiver.Orders = append(buyOrders[:i], buyOrders[i+1:]...)
				fmt.Println("BUY NOW!!")
				fmt.Println(buyOrder)
			}
		}
	}
	//fmt.Println("buy recv: ", event)

	return nil
}
