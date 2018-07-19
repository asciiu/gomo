package main

import (
	"database/sql"
	"fmt"
	"time"

	evt "github.com/asciiu/gomo/common/proto/events"
)

// Processor will process orders
type Processor struct {
	DB          *sql.DB
	MarketQueue []string
}

func remove(s []string, i int) []string {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func (processor *Processor) Ticker() error {
	fmt.Println("ticker started")
	for {
		time.Sleep(1 * time.Second)
		length := len(processor.MarketQueue)

		if length > 0 {
			market := processor.MarketQueue[0]
			fmt.Println(market)
			remove(processor.MarketQueue, 0)
		}
	}
	return nil
}

// ProcessEvent will process ExchangeEvents. These events are published from the exchange sockets.
func (processor *Processor) ProcessEvents(payload *evt.TradeEvents) error {
	// every order check trade event price with order conditions
	for _, event := range payload.Events {
		marketName := event.MarketName
		fmt.Println(marketName)

		// found := false
		// for _, m := range processor.MarketQueue {
		// 	if m == marketName {
		// 		found = true
		// 	}
		// }
		// if !found {
		// 	processor.MarketQueue = append(processor.MarketQueue, marketName)
		// }
	}

	return nil
}
