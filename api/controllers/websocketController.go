package controllers

import (
	"fmt"
	"log"
	"net/http"

	evt "github.com/asciiu/gomo/common/proto/events"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"golang.org/x/net/context"
)

type WebsocketController struct{}

func NewWebsocketController() *WebsocketController {
	return &WebsocketController{}
}

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func (controller *WebsocketController) Connect(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		// Write
		err := ws.WriteMessage(websocket.TextMessage, []byte("Hello, Client!"))
		if err != nil {
			c.Logger().Error(err)
		}

		//// Read
		_, msg, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
		}
		fmt.Printf("%s\n", msg)
	}
}

// ProcessEvent will process ExchangeEvents. These events are published from the exchange sockets.
func (controller *WebsocketController) ProcessEvent(ctx context.Context, event *evt.ExchangeEvent) error {
	log.Println(event)

	return nil
}
