package main

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/asciiu/gomo/common/constants/exchange"
	msg "github.com/asciiu/gomo/common/constants/messages"
	evt "github.com/asciiu/gomo/common/proto/events"
	"github.com/gorilla/websocket"
	micro "github.com/micro/go-micro"
	k8s "github.com/micro/kubernetes/go/micro"
)

type BinanceClient struct {
	ws          *websocket.Conn
	Publisher   micro.Publisher
	lastReceive time.Time
}

type BinanceAggTrade struct {
	Type         string `json:"e"`
	EventTime    uint64 `json:"E"`
	Symbol       string `json:"s"`
	TradeID      uint64 `json:"a"`
	Price        string `json:"p"`
	Quantity     string `json:"q"`
	FirstTradeID uint64 `json:"f"`
	LastTradeID  uint64 `json:"l"`
	TradeTime    uint64 `json:"T"`
	IsBuyerMaker bool   `json:"m"`
	Ignore       bool   `json:"M"`
}

type BinanceTicker struct {
	Type                string `json:"e"`
	EventTime           uint64 `json:"E"`
	Symbol              string `json:"s"`
	PriceChange         string `json:"p"`
	PriceChangePercent  string `json:"P"`
	WeightedAvgPrice    string `json:"w"`
	PreviousDayClose    string `json:"x"`
	ClosePrice          string `json:"c"`
	CloseQuantity       string `json:"Q"`
	BestBidPrice        string `json:"b"`
	BestBidQuantity     string `json:"B"`
	BestAskPrice        string `json:"a"`
	BestAskQuantity     string `json:"A"`
	OpenPrice           string `json:"o"`
	HighPrice           string `json:"h"`
	LowPrice            string `json:"l"`
	TotalTradedBaseVol  string `json:"v"`
	TotalTradedAssetVol string `json:"q"`
	OpenTime            uint64 `json:"O"`
	CloseTime           uint64 `json:"C"`
	FirstTradeID        uint64 `json:"F"`
	LastTradeID         uint64 `json:"L"`
	TotalTrades         uint64 `json:"n"`
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 30 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1024 * 1024
)

func (c *BinanceClient) writePump() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()

	for {
		select {
		case <-ticker.C:
			c.ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}

			if c.lastReceive.Add(pongWait).Before(time.Now().UTC()) {
				log.Println("pong time elapsed for receive")
				return
			}
			log.Println("ping: ", time.Now().UTC())
		}
	}
}

func (c *BinanceClient) readPump() {
	defer func() {
		c.ws.Close()
	}()

	c.ws.SetReadLimit(maxMessageSize)
	//c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error {
		log.Println("pong: ", time.Now().UTC())
		//c.ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		c.ws.SetReadDeadline(time.Now().Add(pongWait))
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			log.Println("websocket err: ", err)
			return
		}

		tickers := []*BinanceTicker{}
		if err = json.Unmarshal(message, &tickers); err != nil {
			log.Println("unmarshall error:", err)
			return
		}

		for _, tick := range tickers {
			//tm := time.Unix(int64(tick.EventTime), 0)
			p, _ := strconv.ParseFloat(tick.ClosePrice, 64)

			// marketname must include hyphen
			symbol := tick.Symbol
			switch {
			case strings.HasSuffix(symbol, "BTC"):
				symbol = strings.Replace(symbol, "BTC", "", 1)
				symbol = symbol + "-BTC"
			case strings.HasSuffix(symbol, "USDT"):
				symbol = strings.Replace(symbol, "USDT", "", 1)
				symbol = symbol + "-USDT"
			case strings.HasSuffix(symbol, "ETH"):
				symbol = strings.Replace(symbol, "ETH", "", 1)
				symbol = symbol + "-ETH"
			case strings.HasSuffix(symbol, "BNB"):
				symbol = strings.Replace(symbol, "BNB", "", 1)
				symbol = symbol + "-BNB"
			}

			tickerEvent := evt.TradeEvent{
				Exchange:   exchange.Binance,
				MarketName: symbol,
				Price:      p,
			}

			//fmt.Println(tickerEvent)

			if err := c.Publisher.Publish(context.Background(), &tickerEvent); err != nil {
				log.Println("publish warning: ", err, tickerEvent)
			}
			c.lastReceive = time.Now().UTC()
		}
	}
}

func (c *BinanceClient) run() {
	url := "wss://stream.binance.com:9443/ws/!ticker@arr"

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	c.ws = conn
	log.Printf("connected to %s", url)

	conn.SetCloseHandler(func(code int, text string) error {
		log.Printf("closed connection %d %s\n", code, text)
		return nil
	})

	go c.writePump()
	c.readPump()
	log.Println("websocket has stopped!")
}

func main() {
	srv := k8s.NewService(
		micro.Name("fomo.binance.websocket"),
	)

	srv.Init()
	tradePublisher := micro.NewPublisher(msg.TopicAggTrade, srv.Client())

	client := BinanceClient{
		Publisher:   tradePublisher,
		lastReceive: time.Now().UTC(),
	}

	go client.run()

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
