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
type AnalyticsService struct {
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

func (service *AnalyticsService) Ticker() error {
	fmt.Println("ticker started: ", time.Now())
	for {
		time.Sleep(2 * time.Second)

		for k := range service.ProcessQueue {
			service.Lock()

			service.MarketCandles[k] = "ding"
			delete(service.ProcessQueue, k)

			service.Unlock()
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
func (service *AnalyticsService) HandleTradeEvent(payload *evt.TradeEvents) error {
	// record close price for the market
	for _, event := range payload.Events {
		market := Market{
			Exchange: event.Exchange,
			Name:     event.MarketName,
		}

		service.RLock()
		_, ok1 := service.MarketCandles[market]
		_, ok2 := service.ProcessQueue[market]
		service.RUnlock()

		if !ok1 && !ok2 {
			service.Lock()
			service.ProcessQueue[market] = event.Price
			service.Unlock()
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
