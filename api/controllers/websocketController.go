package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	evt "github.com/asciiu/gomo/common/proto/events"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"golang.org/x/net/context"
)

type WebsocketController struct {
	connections []*websocket.Conn
}

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
	//defer ws.Close()

	controller.connections = append(controller.connections, ws)
	if err := ws.WriteMessage(websocket.TextMessage, []byte("what's the frequency!")); err != nil {
		c.Logger().Error(err)
	}
	return nil
}

// ProcessEvent will process ExchangeEvents. These events are published from the exchange sockets.
func (controller *WebsocketController) ProcessEvent(ctx context.Context, event *evt.ExchangeEvent) error {
	//log.Println(event)

	for _, conn := range controller.connections {
		json, err := json.Marshal(event)
		if err != nil {
			log.Println(err)
			return nil
		}

		if err := conn.WriteMessage(websocket.TextMessage, json); err != nil {
			log.Println(err)
		}
	}

	return nil
}
