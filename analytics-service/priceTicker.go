package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	repoAnalytics "github.com/asciiu/gomo/analytics-service/db/sql"
	protoAnalytics "github.com/asciiu/gomo/analytics-service/proto/analytics"
	protoEvt "github.com/asciiu/gomo/common/proto/events"
	"github.com/lib/pq"
)

// Processor will process orders
type PriceTicker struct {
	DB            *sql.DB
	MarketPrices  map[Market]float64
	CurrentPeriod string
	TimePeriod    time.Duration
}

func (service *PriceTicker) Ticker() error {
	fmt.Println("ticker started: ", time.Now())
	truncTime := time.Now().UTC().Truncate(service.TimePeriod)
	truncTimeStr := string(pq.FormatTimestamp(truncTime))
	service.CurrentPeriod = truncTimeStr

	for {
		time.Sleep(1 * time.Second)

		truncTime = time.Now().UTC().Truncate(service.TimePeriod)
		truncTimeStr = string(pq.FormatTimestamp(truncTime))

		if service.CurrentPeriod != truncTimeStr {
			rates := make([]*protoAnalytics.MarketPrice, 0)
			for market, price := range service.MarketPrices {
				marketPrice := protoAnalytics.MarketPrice{
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
func (service *PriceTicker) HandleExchangeEvent(payload *protoEvt.TradeEvents) error {
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

// from and to params are symbols: e.g. from: BTC to: USDT
func AmigoniSpecial(exchange, from, to, atTime string, fromAmount float64) float64 {
	var rate, reverse, fromRate, toRate float64

	//Case in which to and from are the same i.e. BTCBTC
	if from == to {
		return fromAmount
	}

	// find all prices here for given time
	// repoAnalytics.FindPricesAtTime(exchange, time)
	markets := make([]*protoAnalytics.MarketPrice, 0)

	for _, market := range markets {
		if market.MarketName == from+"-"+to {
			//Simple Case where the rate exists i.e. ADA-BTC
			rate = market.ClosedAtPrice
			break
		}
		if market.MarketName == to+"-"+from {
			// reverse case exists BTC-ADA
			reverse = 1 / market.ClosedAtPrice
		}
		if market.MarketName == from+"-BTC" {
			// indirect from rate
			fromRate = market.ClosedAtPrice
		}
		if market.MarketName == to+"-BTC" {
			// indirect to rate
			toRate = market.ClosedAtPrice
		}
	}

	switch {
	case rate == 0 && reverse != 0:
		rate = reverse
	case rate == 0 && reverse == 0:
		// direct rate doesn't exist so going through BTC to convert i.e ADAXVG
		rate = fromRate / toRate
	}

	//Ti.API.trace("Convert: "+from+" to "+to+" = "+rate+" "+fromExchange);
	return rate * fromAmount
}
