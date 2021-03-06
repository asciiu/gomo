package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	protoEvt "github.com/asciiu/gomo/common/proto/events"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

type WebsocketController struct {
	connections []*websocket.Conn
	buffer      []*protoEvt.TradeEvent
}

func NewWebsocketController() *WebsocketController {
	return &WebsocketController{
		buffer:      make([]*protoEvt.TradeEvent, 0),
		connections: make([]*websocket.Conn, 0),
	}
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
		log.Println(err)
	}

	// remove the connection from the connect pool
	controller.connections = append(controller.connections[:i], controller.connections[i+1:]...)
	return nil
}

func (controller *WebsocketController) Ticker() {
	for {
		time.Sleep(1 * time.Second)
		events := controller.buffer
		controller.buffer = nil

		// only send out events to clients when non nil
		if events != nil {
			// send events to all connected clients
			for _, conn := range controller.connections {
				json, err := json.Marshal(events)
				if err != nil {
					log.Println(err)
				}

				if err := conn.WriteMessage(websocket.TextMessage, json); err != nil {
					log.Println(err)
				}
			}
		}
	}
}

// ProcessEvent will process ExchangeEvents. These events are published from the exchange sockets.
func (controller *WebsocketController) CacheEvents(tradeEvents *protoEvt.TradeEvents) {
	for _, event := range tradeEvents.Events {
		// shorten trade event
		tevent := protoEvt.TradeEvent{
			Exchange:   event.Exchange,
			Type:       event.Type,
			MarketName: event.MarketName,
			Price:      event.Price,
		}

		found := false
		for _, e := range controller.buffer {
			if e.Exchange == tevent.Exchange && e.MarketName == tevent.MarketName {
				e.Type = tevent.Type
				e.Price = tevent.Price
				found = true
				break
			}
		}

		if !found {
			controller.buffer = append(controller.buffer, &tevent)
		}
	}
}
