package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Symbol struct {
	Symbol             string                    `json:"symbol"`
	Status             string                    `json:"status"`
	BaseAsset          string                    `json:"baseAsset"`
	BaseAssetPrecision int16                     `json:"baseAssetPrecision"`
	QuoteAsset         string                    `json:"quoteAsset"`
	QuotePrecision     int16                     `json:"quotePrecision"`
	OrderTypes         []string                  `json:"orderTypes"`
	IcebergAllowed     bool                      `json:"icebergAllowed"`
	Filters            []*map[string]interface{} `json:"filters"`
}

type RateLimit struct {
	RateLimitType string `json:"rateLimitType`
	Interval      string `json:"interval"`
	Limit         int32  `json:"limit"`
}

type ExchangeInfo struct {
	TimeZone   string       `json:"timezone"`
	ServerTime int64        `json:"serverTime"`
	RateLimits []*RateLimit `json:"rateLimits"`
	Symbols    []*Symbol    `json:"symbols"`
}

type BinanceExchangeInfo struct {
	// market name example: ADA-BTC
	Markets map[string]*Symbol
}

func NewBinanceExchangeInfo() *BinanceExchangeInfo {
	// TODO this data should be updated periodically
	bexinfo := BinanceExchangeInfo{
		Markets: make(map[string]*Symbol),
	}
	resp, err := http.Get("https://api.binance.com/api/v1/exchangeInfo")
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var info ExchangeInfo
	if err := json.Unmarshal(body, &info); err != nil {
		panic(err)
	}

	for _, symbol := range info.Symbols {
		marketName := symbol.BaseAsset + "-" + symbol.QuoteAsset
		bexinfo.Markets[marketName] = symbol
	}
	return &bexinfo
}

type LotSize struct {
	MinQty   float64
	MaxQty   float64
	StepSize float64
}

func (binance *BinanceExchangeInfo) LotSize(marketName string) *LotSize {
	symbol := binance.Markets[marketName]

	for _, filter := range symbol.Filters {
		if (*filter)["filterType"] == "LOT_SIZE" {
			minQty, _ := strconv.ParseFloat((*filter)["minQty"].(string), 64)
			maxQty, _ := strconv.ParseFloat((*filter)["maxQty"].(string), 64)
			stepSize, _ := strconv.ParseFloat((*filter)["stepSize"].(string), 64)
			return &LotSize{
				MinQty:   minQty,
				MaxQty:   maxQty,
				StepSize: stepSize,
			}
		}
	}
	return nil
}

type PriceFilter struct {
	Min      float64
	Max      float64
	TickSize float64
}

func (binance *BinanceExchangeInfo) PriceFilter(marketName string) *PriceFilter {
	symbol := binance.Markets[marketName]

	for _, filter := range symbol.Filters {
		if (*filter)["filterType"] == "PRICE_FILTER" {
			min, _ := strconv.ParseFloat((*filter)["minPrice"].(string), 64)
			max, _ := strconv.ParseFloat((*filter)["maxPrice"].(string), 64)
			tick, _ := strconv.ParseFloat((*filter)["tickSize"].(string), 64)
			return &PriceFilter{
				Min:      min,
				Max:      max,
				TickSize: tick,
			}
		}
	}
	return nil
}

type MinNotional struct {
	MinNotional float64
}

func (binance *BinanceExchangeInfo) MinNotional(marketName string) *MinNotional {
	symbol := binance.Markets[marketName]

	for _, filter := range symbol.Filters {
		if (*filter)["filterType"] == "MIN_NOTIONAL" {
			minNote, _ := strconv.ParseFloat((*filter)["minNotional"].(string), 64)
			return &MinNotional{
				MinNotional: minNote,
			}
		}
	}
	return nil
}

type IcebergParts struct {
	Limit float64
}

func (binance *BinanceExchangeInfo) IcebergParts(marketName string) *IcebergParts {
	symbol := binance.Markets[marketName]

	for _, filter := range symbol.Filters {
		if (*filter)["filterType"] == "ICEBERG_PARTS" {
			limit, _ := (*filter)["limit"].(float64)
			return &IcebergParts{
				Limit: limit,
			}
		}
	}
	return nil
}

type MaxAlgoOrders struct {
	MaxNum float64
}

func (binance *BinanceExchangeInfo) MaxAlgoOrders(marketName string) *MaxAlgoOrders {
	symbol := binance.Markets[marketName]

	for _, filter := range symbol.Filters {
		if (*filter)["filterType"] == "MAX_NUM_ALGO_ORDERS" {
			max, _ := (*filter)["maxNumAlgoOrders"].(float64)
			return &MaxAlgoOrders{
				MaxNum: max,
			}
		}
	}
	return nil
}
