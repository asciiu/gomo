package main

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	evt "github.com/asciiu/gomo/common/proto/events"
	micro "github.com/micro/go-micro"
)

// Processor will process orders
type Processor struct {
	sync.RWMutex
	DB               *sql.DB
	MarketClosePrice map[Market]float64
	CandlePub        micro.Publisher
}

type Market struct {
	Exchange string // market has a exchange
	Name     string // a name in the form of EOS-BTC
	Valid    bool   // true if market data is valid
}

func remove(s []string, i int) []string {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func (processor *Processor) Ticker() error {
	fmt.Println("ticker started")
	for {
		time.Sleep(1 * time.Second)

		processor.RLock()
		length := len(processor.MarketClosePrice)
		processor.RUnlock()

		fmt.Println(length)
		//if length > 0 {
		//	market := processor.MarketClosePrice[0]
		//	fmt.Println(market)
		//}
	}
	return nil
}

// ProcessEvent will process ExchangeEvents. These events are published from the exchange sockets.
func (processor *Processor) ProcessEvents(payload *evt.TradeEvents) error {
	// record close price for the market
	for _, event := range payload.Events {
		market := Market{
			Exchange: event.Exchange,
			Name:     event.MarketName,
			Valid:    true,
		}

		processor.Lock()
		processor.MarketClosePrice[market] = event.Price
		processor.Unlock()

		//found := false
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
