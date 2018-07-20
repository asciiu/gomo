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

	MarketCandles map[Market]string
	ProcessQueue  map[Market]float64
	CandlePub     micro.Publisher
}

func remove(s []string, i int) []string {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func (processor *Processor) Ticker() error {
	fmt.Println("ticker started: ", time.Now())
	for {
		time.Sleep(2 * time.Second)

		for k := range processor.ProcessQueue {
			processor.Lock()

			processor.MarketCandles[k] = "ding"
			delete(processor.ProcessQueue, k)

			processor.Unlock()
			fmt.Println(k)
			break
		}

		// get market from collection
		// send out candle request if no candles
		// otherwise ignore
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
		}

		processor.RLock()
		_, ok1 := processor.MarketCandles[market]
		_, ok2 := processor.ProcessQueue[market]
		processor.RUnlock()

		if !ok1 && !ok2 {
			processor.Lock()
			processor.ProcessQueue[market] = event.Price
			processor.Unlock()
		}

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
