package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	repoAnalytics "github.com/asciiu/gomo/analytics-service/db/sql"
	protoAnalytics "github.com/asciiu/gomo/analytics-service/proto/analytics"
	protoBinance "github.com/asciiu/gomo/binance-service/proto/binance"
	constExch "github.com/asciiu/gomo/common/constants/exchange"
	constRes "github.com/asciiu/gomo/common/constants/response"
	evt "github.com/asciiu/gomo/common/proto/events"
	"github.com/lib/pq"
	micro "github.com/micro/go-micro"
)

// Processor will process orders
type AnalyticsService struct {
	sync.RWMutex
	DB               *sql.DB
	MarketClosePrice map[Market]float64

	BinanceClient protoBinance.BinanceServiceClient
	currencies    map[string]string
	//MarketCandles map[Market]string
	//ProcessQueue  map[Market]float64
	//CandlePub     micro.Publisher
	Directory map[string]*protoAnalytics.MarketInfo
}

func NewAnalyticsService(db *sql.DB, srv micro.Service) *AnalyticsService {
	service := AnalyticsService{
		DB:            db,
		Directory:     make(map[string]*protoAnalytics.MarketInfo),
		currencies:    make(map[string]string),
		BinanceClient: protoBinance.NewBinanceServiceClient("binance", srv.Client()),
	}

	currencies, err := repoAnalytics.GetCurrencyNames(db)
	switch {
	case err == sql.ErrNoRows:
		log.Println("Quaid, start the reactor!")
	case err != nil:
	default:
		for _, c := range currencies {
			service.currencies[c.TickerSymbol] = c.CurrencyName
		}
	}

	return &service
}

func remove(s []string, i int) []string {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func (service *AnalyticsService) Ticker() error {
	fmt.Println("ticker started: ", time.Now())
	for {
		time.Sleep(2 * time.Second)

		// for k := range service.ProcessQueue {
		// 	service.Lock()

		// 	service.MarketCandles[k] = "ding"
		// 	delete(service.ProcessQueue, k)

		// 	service.Unlock()
		// 	fmt.Println(k)
		// 	break
		// }

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
func (service *AnalyticsService) HandleExchangeEvent(payload *evt.TradeEvents) error {
	// record close price for the market
	for _, event := range payload.Events {
		// market := Market{
		// 	Exchange: event.Exchange,
		// 	Name:     event.MarketName,
		// }

		// service.RLock()
		// _, ok1 := service.MarketCandles[market]
		// _, ok2 := service.ProcessQueue[market]
		// service.RUnlock()

		// if !ok1 && !ok2 {
		// 	service.Lock()
		// 	service.ProcessQueue[market] = event.Price
		// 	service.Unlock()
		// }

		names := strings.Split(event.MarketName, "-")
		baseCurrency := names[1]
		baseCurrencyName := service.currencies[baseCurrency]
		marketCurrency := names[0]
		marketCurrencyName := service.currencies[marketCurrency]

		market := protoAnalytics.MarketInfo{
			BaseCurrencySymbol:   baseCurrency,
			BaseCurrencyName:     baseCurrencyName,
			Exchange:             event.Exchange,
			ExchangeMarketName:   names[0] + names[1],
			MarketCurrencySymbol: marketCurrency,
			MarketCurrencyName:   marketCurrencyName,
			MarketName:           event.MarketName,
			Price:                fmt.Sprintf("%.8f", event.Price)}

		key := fmt.Sprintf("%s-%s", event.Exchange, event.MarketName)

		service.Lock()
		if m, ok := service.Directory[key]; ok {
			// update the price only
			m.Price = fmt.Sprintf("%.8f", event.Price)
		} else {
			// grab exchange rules for this market here
			switch event.Exchange {
			case constExch.Binance:
				rules, _ := service.BinanceClient.GetMarketRestrictions(context.Background(), &protoBinance.MarketRestrictionRequest{MarketName: event.MarketName})
				if rules.Status == constRes.Success {
					market.MinTradeSize = fmt.Sprintf("%.8f", rules.Data.Restrictions.MinTradeSize)
					market.MaxTradeSize = fmt.Sprintf("%.8f", rules.Data.Restrictions.MaxTradeSize)
					market.TradeSizeStep = fmt.Sprintf("%.8f", rules.Data.Restrictions.TradeSizeStep)
					market.MinMarketPrice = fmt.Sprintf("%.8f", rules.Data.Restrictions.MinMarketPrice)
					market.MaxMarketPrice = fmt.Sprintf("%.8f", rules.Data.Restrictions.MaxMarketPrice)
					market.MarketPriceStep = fmt.Sprintf("%.8f", rules.Data.Restrictions.MarketPriceStep)
					market.BasePrecision = rules.Data.Restrictions.BasePrecision
					market.MarketPrecision = rules.Data.Restrictions.MarketPrecision
				} else {
					log.Println("could not get rules for ", event.MarketName)
				}
			}
			service.Directory[key] = &market

		}
		service.Unlock()

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

// This function was formly known as the Amigoni special. It has been refined by yours truely - Axl Codes.
func (service *AnalyticsService) ConvertCurrency(ctx context.Context, req *protoAnalytics.ConversionRequest, res *protoAnalytics.ConversionResponse) error {
	var rate, reverse, fromRate, toRate float64
	from := req.From
	to := req.To
	atTime, _ := time.Parse(time.RFC3339, req.AtTimestamp)
	trunctTime := atTime.Truncate(time.Duration(5) * time.Minute)
	trunctTimeStr := string(pq.FormatTimestamp(trunctTime))

	//Case in which to and from are the same i.e. BTCBTC
	if from == to {
		res.Status = constRes.Success
		res.Data = &protoAnalytics.ConversionAmount{
			ConvertedAmount: req.FromAmount,
		}
		return nil
	}

	// find all prices here for given time
	markets, err := repoAnalytics.FindExchangeRates(service.DB, req.Exchange, trunctTimeStr)
	if err != nil {
		res.Status = constRes.Error
		res.Message = err.Error()
		return nil
	}

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
		//if market.MarketName == from+"-BTC" {
		//	// indirect from rate
		//	fromRate = market.ClosedAtPrice
		//}
		//if market.MarketName == to+"-BTC" {
		//	// indirect to rate
		//	toRate = market.ClosedAtPrice
		//}
		if market.MarketName == from+"-BTC" {
			// indirect from rate
			fromRate = market.ClosedAtPrice
		}
		if market.MarketName == "BTC-"+to {
			// indirect to rate
			toRate = market.ClosedAtPrice
		}
	}

	switch {
	case rate == 0 && reverse != 0:
		rate = reverse
	case rate == 0 && reverse == 0:
		// direct rate doesn't exist so going through BTC to convert i.e ADAXVG
		//rate = fromRate / toRate
		rate = fromRate * toRate
	}

	//Ti.API.trace("Convert: "+from+" to "+to+" = "+rate+" "+fromExchange);
	res.Status = constRes.Success
	res.Data = &protoAnalytics.ConversionAmount{
		ConvertedAmount: rate * req.FromAmount,
	}
	return nil
}

func (service *AnalyticsService) GetMarketInfo(ctx context.Context, req *protoAnalytics.SearchMarketsRequest, res *protoAnalytics.MarketsResponse) error {
	m := make([]*protoAnalytics.MarketInfo, 0)
	term := req.Term

	for k, v := range service.Directory {
		// if the key contains the term or base or market currency
		// append to results
		if strings.Contains(strings.ToLower(k), strings.ToLower(term)) ||
			strings.Contains(strings.ToLower(v.BaseCurrencySymbol), strings.ToLower(term)) ||
			strings.Contains(strings.ToLower(v.MarketCurrencySymbol), strings.ToLower(term)) {
			m = append(m, v)
		}
	}

	res.Status = constRes.Success
	res.Data = &protoAnalytics.MarketInfoResponse{
		Markets: m,
	}
	return nil
}
