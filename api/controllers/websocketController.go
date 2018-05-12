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

// Connect handles websocket connections
func (controller *WebsocketController) Connect(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	if err := ws.WriteMessage(websocket.TextMessage, []byte("ready!")); err != nil {
		log.Println("write:", err)
		return err
	}
	i := len(controller.connections)
	controller.connections = append(controller.connections, ws)

	// block until client closes
	if _, _, err := ws.ReadMessage(); err != nil {
		// client closes this will read: websocket: close 1005 (no status)
		log.Println("read:", err)
	}

	// remove the connection from the connect pool
	controller.connections = append(controller.connections[:i], controller.connections[i+1:]...)
	return nil
}

// ProcessEvent will process ExchangeEvents. These events are published from the exchange sockets.
func (controller *WebsocketController) ProcessEvent(ctx context.Context, event *evt.ExchangeEvent) error {
	for _, conn := range controller.connections {
		json, err := json.Marshal(event)
		if err != nil {
			log.Println(err)
			return err
		}

		if err := conn.WriteMessage(websocket.TextMessage, json); err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}
