package controllers

import (
	"fmt"

	evt "github.com/asciiu/gomo/common/proto/events"
	"github.com/labstack/echo"
	"golang.org/x/net/context"
)

type SearchController struct {
	buffer map[string]*evt.TradeEvent
}

func NewSearchController() *SearchController {
	return &SearchController{
		buffer: make(map[string]*evt.TradeEvent),
	}
}

func (controller *SearchController) Search(c echo.Context) error {
	return nil
}

// ProcessEvent will process ExchangeEvents. These events are published from the exchange sockets.
func (controller *SearchController) ProcessEvent(ctx context.Context, event *evt.TradeEvent) error {

	key := fmt.Sprintf("%s-%s", event.Exchange, event.MarketName)
	fmt.Println(key)
	controller.buffer[key] = event

	return nil
}
