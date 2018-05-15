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
	markets map[string]*evt.TradeEvent
	mux     sync.Mutex
}

// A ResponseSearchSuccess will always contain a status of "successful".
// swagger:model responseDeviceSuccess
type ResponseSearchSuccess struct {
	Status string            `json:"status"`
	Data   []*evt.TradeEvent `json:"data"`
}

func NewSearchController() *SearchController {
	return &SearchController{
		markets: make(map[string]*evt.TradeEvent),
	}
}

func (controller *SearchController) Search(c echo.Context) error {

	term := c.QueryParam("term")
	m := make([]*evt.TradeEvent, 0)

	for k, v := range controller.markets {
		if strings.Contains(strings.ToLower(k), term) {
			m = append(m, v)
		}
	}

	response := &ResponseSearchSuccess{
		Status: "success",
		Data:   m,
	}

	return c.JSON(http.StatusOK, response)
}

// ProcessEvent will process ExchangeEvents. These events are published from the exchange sockets.
func (controller *SearchController) ProcessEvent(ctx context.Context, event *evt.TradeEvent) error {
	// shorten trade event
	tevent := evt.TradeEvent{
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
