package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	repoAnalytics "github.com/asciiu/gomo/analytics-service/db/sql"
	protoPrice "github.com/asciiu/gomo/analytics-service/proto/price"
	protoEvt "github.com/asciiu/gomo/common/proto/events"
	"github.com/lib/pq"
)

// Processor will process orders
type PriceService struct {
	DB            *sql.DB
	MarketPrices  map[Market]float64
	CurrentPeriod string
	TimePeriod    time.Duration
}

func (service *PriceService) Ticker() error {
	fmt.Println("ticker started: ", time.Now())
	truncTime := time.Now().UTC().Truncate(service.TimePeriod)
	truncTimeStr := string(pq.FormatTimestamp(truncTime))
	service.CurrentPeriod = truncTimeStr

	for {
		time.Sleep(1 * time.Second)

		truncTime = time.Now().UTC().Truncate(service.TimePeriod)
		truncTimeStr = string(pq.FormatTimestamp(truncTime))

		if service.CurrentPeriod != truncTimeStr {
			rates := make([]*protoPrice.MarketPrice, 0)
			for market, price := range service.MarketPrices {
				marketPrice := protoPrice.MarketPrice{
					Exchange:      market.Exchange,
					MarketName:    market.Name,
					ClosedAtPrice: price,
					ClosedAtTime:  truncTimeStr,
				}
				rates = append(rates, &marketPrice)
			}

			// archive the market rates
			if err := repoAnalytics.InsertPrices(service.DB, rates); err != nil {
				log.Println("error on insert prices: ", err.Error())
			}

			service.CurrentPeriod = truncTimeStr
		}
	}
	return nil
}

// These events are published from the exchange socket services.
func (service *PriceService) HandleExchangeEvent(payload *protoEvt.TradeEvents) error {
	// record close price for the market
	for _, event := range payload.Events {
		market := Market{
			Exchange: event.Exchange,
			Name:     event.MarketName,
		}

		service.MarketPrices[market] = event.Price
	}

	return nil
}
