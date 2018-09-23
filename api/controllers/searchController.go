package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	protoAnalytics "github.com/asciiu/gomo/analytics-service/proto/analytics"
	repoToken "github.com/asciiu/gomo/api/db/sql"
	protoEvt "github.com/asciiu/gomo/common/proto/events"
	"github.com/labstack/echo"
	micro "github.com/micro/go-micro"
	"golang.org/x/net/context"
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
	BaseCurrencySymbol   string `json:"baseCurrencySymbol"`
	BaseCurrencyName     string `json:"baseCurrencyName"`
	BasePrecision        int32  `json:"basePrecision"`
	Exchange             string `json:"exchange"`
	ExchangeMarketName   string `json:"exchangeMarketName"`
	MarketCurrencySymbol string `json:"marketCurrencySymbol"`
	MarketCurrencyName   string `json:"marketCurrencyName"`
	MarketPrecision      int32  `json:"marketPrecision"`
	MarketName           string `json:"marketName"`
	MinTradeSize         string `json:"minTradeSize"`
	MaxTradeSize         string `json:"maxTradeSize"`
	TradeSizeStep        string `json:"tradeSizeStep"`
	MinMarketPrice       string `json:"minMarketPrice"`
	MaxMarketPrice       string `json:"maxMarketPrice"`
	MarketPriceStep      string `json:"marketPriceStep"`
	Price                string `json:"price"`
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
	currencies      map[string]string
	AnalyticsClient protoAnalytics.AnalyticsServiceClient
}

func NewSearchController(db *sql.DB, service micro.Service) *SearchController {
	controller := SearchController{
		markets:         make(map[string]*Market),
		currencies:      make(map[string]string),
		AnalyticsClient: protoAnalytics.NewAnalyticsServiceClient("analytics", service.Client()),
	}

	currencies, err := repoToken.GetCurrencyNames(db)
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

	response, _ := controller.AnalyticsClient.GetMarketInfo(context.Background(), &protoAnalytics.SearchMarketsRequest{term})
	//m := make([]*Market, 0)

	//for k, v := range controller.markets {
	//	switch {
	//	case strings.Contains(strings.ToLower(k), strings.ToLower(term)):
	//		m = append(m, v)
	//	case strings.Contains(strings.ToLower(v.BaseCurrencySymbol), strings.ToLower(term)):
	//		m = append(m, v)
	//	case strings.Contains(strings.ToLower(v.MarketCurrencySymbol), strings.ToLower(term)):
	//		m = append(m, v)
	//	default:
	//	}
	//}

	//response := &ResponseSearchSuccess{
	//	Status: constRes.Success,
	//	Data:   ResponseMarkets{m},
	//}

	return c.JSON(http.StatusOK, response)
}

// ProcessEvent will process ExchangeEvents. These events are published from the exchange sockets.
func (controller *SearchController) CacheEvents(tradeEvents *protoEvt.TradeEvents) {
	for _, event := range tradeEvents.Events {
		names := strings.Split(event.MarketName, "-")
		baseCurrency := names[1]
		baseCurrencyName := controller.currencies[baseCurrency]
		marketCurrency := names[0]
		marketCurrencyName := controller.currencies[marketCurrency]

		// shorten trade event
		tevent := Market{
			BaseCurrencySymbol:   baseCurrency,
			BaseCurrencyName:     baseCurrencyName,
			Exchange:             event.Exchange,
			ExchangeMarketName:   names[0] + names[1],
			MarketCurrencySymbol: marketCurrency,
			MarketCurrencyName:   marketCurrencyName,
			MarketName:           event.MarketName,
			Price:                fmt.Sprintf("%.8f", event.Price),
			BasePrecision:        8,
			MarketPrecision:      8,
			MarketPriceStep:      "0.00000001",
			MaxTradeSize:         "1000000000.0",
			MinTradeSize:         "0.00000001",
			TradeSizeStep:        "0.00000001"}

		key := fmt.Sprintf("%s-%s", event.Exchange, event.MarketName)
		//key = key
		controller.mux.Lock()
		controller.markets[key] = &tevent
		controller.mux.Unlock()
	}
}
