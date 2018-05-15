package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	evt "github.com/asciiu/gomo/common/proto/events"
	"github.com/labstack/echo"
	"golang.org/x/net/context"
)

type SearchController struct {
	markets map[string]*Market
	mux     sync.Mutex
}

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
	Exchange   string  `json:"exchange,omitempty"`
	Type       string  `json:"type,omitempty"`
	MarketName string  `json:"marketName,omitempty"`
	Price      float64 `json:"price,omitempty"`
}

func NewSearchController() *SearchController {
	return &SearchController{
		markets: make(map[string]*Market),
	}
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
		if strings.Contains(strings.ToLower(k), strings.ToLower(term)) {
			m = append(m, v)
		}
	}

	response := &ResponseSearchSuccess{
		Status: "success",
		Data:   ResponseMarkets{m},
	}

	return c.JSON(http.StatusOK, response)
}

// ProcessEvent will process ExchangeEvents. These events are published from the exchange sockets.
func (controller *SearchController) ProcessEvent(ctx context.Context, event *evt.TradeEvent) error {
	// shorten trade event
	tevent := Market{
		Exchange:   event.Exchange,
		Type:       event.Type,
		MarketName: event.MarketName,
		Price:      event.Price,
	}

	key := fmt.Sprintf("%s-%s", event.Exchange, event.MarketName)
	//key = key
	controller.mux.Lock()
	controller.markets[key] = &tevent
	controller.mux.Unlock()

	return nil
}
