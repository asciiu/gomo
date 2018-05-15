package controllers

import (
	"fmt"
	"net/http"
	"sync"

	evt "github.com/asciiu/gomo/common/proto/events"
	"github.com/labstack/echo"
	"golang.org/x/net/context"
)

type SearchController struct {
	markets map[string]*evt.TradeEvent
	mux     sync.Mutex
}

func NewSearchController() *SearchController {
	return &SearchController{
		markets: make(map[string]*evt.TradeEvent),
	}
}

func (controller *SearchController) Search(c echo.Context) error {

	name := c.QueryParam("name")
	return c.String(http.StatusOK, name)
}

// ProcessEvent will process ExchangeEvents. These events are published from the exchange sockets.
func (controller *SearchController) ProcessEvent(ctx context.Context, event *evt.TradeEvent) error {

	key := fmt.Sprintf("%s-%s", event.Exchange, event.MarketName)
	//key = key
	controller.mux.Lock()
	controller.markets[key] = event
	controller.mux.Unlock()

	return nil
}
