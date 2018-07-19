package main

import (
	"context"
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

func (processor *Processor) Ticker() error {
	for {
		time.Sleep(1 * time.Second)

		if len(processor.MarketQueue) > 0 {
			market := processor.MarketQueue[0]
			fmt.Println(market)
		}
	}
}

// ProcessEvent will process ExchangeEvents. These events are published from the exchange sockets.
func (processor *Processor) ProcessEvent(ctx context.Context, payload *evt.TradeEvents) error {
	// every order check trade event price with order conditions
	for _, event := range payload.Events {
		marketName := event.MarketName

		found := false
		for _, m := range processor.MarketQueue {
			if m == marketName {
				found = true
			}
		}
		if !found {
			processor.MarketQueue = append(processor.MarketQueue, marketName)
		}
	}

	return nil
}
