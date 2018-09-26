package controllers

import (
	"database/sql"
	"net/http"

	protoAnalytics "github.com/asciiu/gomo/analytics-service/proto/analytics"
	constRes "github.com/asciiu/gomo/common/constants/response"
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
	AnalyticsClient protoAnalytics.AnalyticsServiceClient
}

func NewSearchController(db *sql.DB, service micro.Service) *SearchController {
	controller := SearchController{
		AnalyticsClient: protoAnalytics.NewAnalyticsServiceClient("analytics", service.Client()),
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
	markets := make([]*Market, 0)

	for _, market := range response.Data.Markets {
		markets = append(markets, &Market{
			Exchange:             market.Exchange,
			ExchangeMarketName:   market.ExchangeMarketName,
			BaseCurrencySymbol:   market.BaseCurrencySymbol,
			BaseCurrencyName:     market.BaseCurrencyName,
			BasePrecision:        market.BasePrecision,
			MarketCurrencySymbol: market.MarketCurrencySymbol,
			MarketCurrencyName:   market.MarketCurrencyName,
			MarketPrecision:      market.MarketPrecision,
			MarketName:           market.MarketName,
			MinTradeSize:         market.MinTradeSize,
			MaxTradeSize:         market.MaxTradeSize,
			TradeSizeStep:        market.TradeSizeStep,
			MinMarketPrice:       market.MinMarketPrice,
			MaxMarketPrice:       market.MaxMarketPrice,
			MarketPriceStep:      market.MarketPriceStep,
			Price:                market.Price,
		})
	}

	res := &ResponseSearchSuccess{
		Status: constRes.Success,
		Data:   ResponseMarkets{markets},
	}

	return c.JSON(http.StatusOK, res)
}
