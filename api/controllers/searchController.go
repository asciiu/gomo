package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	asql "github.com/asciiu/gomo/api/db/sql"
	evt "github.com/asciiu/gomo/common/proto/events"
	"github.com/labstack/echo"
)

// A ResponseSearchSuccess will always contain a status of "successful".
// swagger:model responseSearchSuccess
type ResponseSearchSuccess struct {
	Status string          `json:"status"`
	Data   ResponseMarkets `json:"data"`
}

type ResponseMarkets struct {
	Markets []*Market `json:"markets"`
}

type Market struct {
	BaseCurrency       string `json:"baseCurrency"`
	BaseCurrencyLong   string `json:"baseCurrencyLong"`
	Exchange           string `json:"exchange"`
	ExchangeMarketName string `json:"exchangeMarketName"`
	MarketCurrency     string `json:"marketCurrency"`
	MarketCurrencyLong string `json:"marketCurrencyLong"`
	MarketName         string `json:"marketName"`
	Price              string `json:"price"`
	MinTradeSize       string `json:"minTradeSize"`
}

// This struct is used in the generated swagger docs,
// and it is not used anywhere.
// swagger:parameters searchMarkets
type SearchTerm struct {
	// Required: false
	// In: query
	Term string `json:"term"`
}

type SearchController struct {
	markets map[string]*Market
	mux     sync.Mutex
	// map of ticker symbol to full name
	currencies map[string]string
}

func NewSearchController(db *sql.DB) *SearchController {
	controller := SearchController{
		markets:    make(map[string]*Market),
		currencies: make(map[string]string),
	}

	currencies, err := asql.GetCurrencyNames(db)
	switch {
	case err == sql.ErrNoRows:
		log.Println("Quaid, start the reactor!")
	case err != nil:
	default:
		for _, c := range currencies {
			controller.currencies[c.TickerSymbol] = c.CurrencyName
		}
	}

	return &controller
}

// swagger:route GET /search search searchMarkets
//
// search markets (protected)
//
// Returns a list of active markets.
//
// responses:
//  200: responseSearchSuccess "data" will contain array of markets with "status": "success"
func (controller *SearchController) Search(c echo.Context) error {

	term := c.QueryParam("term")
	m := make([]*Market, 0)

	for k, v := range controller.markets {
		switch {
		case strings.Contains(strings.ToLower(k), strings.ToLower(term)):
			m = append(m, v)
		case strings.Contains(strings.ToLower(v.BaseCurrency), strings.ToLower(term)):
			m = append(m, v)
		case strings.Contains(strings.ToLower(v.MarketCurrency), strings.ToLower(term)):
			m = append(m, v)
		default:
		}
	}

	response := &ResponseSearchSuccess{
		Status: "success",
		Data:   ResponseMarkets{m},
	}

	return c.JSON(http.StatusOK, response)
}

// ProcessEvent will process ExchangeEvents. These events are published from the exchange sockets.
func (controller *SearchController) CacheEvents(tradeEvents *evt.TradeEvents) {
	for _, event := range tradeEvents.Events {
		names := strings.Split(event.MarketName, "-")
		baseCurrency := names[1]
		baseCurrencyName := controller.currencies[baseCurrency]
		marketCurrency := names[0]
		marketCurrencyName := controller.currencies[marketCurrency]

		// shorten trade event
		tevent := Market{
			BaseCurrency:       baseCurrency,
			BaseCurrencyLong:   baseCurrencyName,
			Exchange:           event.Exchange,
			ExchangeMarketName: names[0] + names[1],
			MarketCurrency:     marketCurrency,
			MarketCurrencyLong: marketCurrencyName,
			MarketName:         event.MarketName,
			Price:              fmt.Sprintf("%.8f", event.Price),
			MinTradeSize:       "0.0001",
		}

		key := fmt.Sprintf("%s-%s", event.Exchange, event.MarketName)
		//key = key
		controller.mux.Lock()
		controller.markets[key] = &tevent
		controller.mux.Unlock()
	}
}
