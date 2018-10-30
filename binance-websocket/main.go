package main

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"

	constExch "github.com/asciiu/gomo/common/constants/exchange"
	constMessage "github.com/asciiu/gomo/common/constants/message"
	protoEvt "github.com/asciiu/gomo/common/proto/events"
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
	// these may come in as -1
	FirstTradeID int64 `json:"F"`
	LastTradeID  int64 `json:"L"`
	TotalTrades  int64 `json:"n"`
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 10 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 7) / 10

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
				log.Println("failed to ping:", err)
				break
			}
		}
	}
}

func (c *BinanceClient) readPump() {
	defer c.ws.Close()

	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error {
		//log.Println("pong: ", time.Now().UTC())
		c.ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			log.Println("read error: ", err)
			return
		}

		binanceTickers := []*BinanceTicker{}
		if err = json.Unmarshal(message, &binanceTickers); err != nil {
			log.Println("unmarshall error:", err)
			return
		}

		marketTickers := make([]*protoEvt.TradeEvent, 0)
		for _, tick := range binanceTickers {
			//tm := time.Unix(int64(tick.EventTime), 0)
			p, err := strconv.ParseFloat(tick.ClosePrice, 64)
			if err != nil {
				log.Println("failed to parse float from close price: ", tick.ClosePrice)
			}

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

			tickerEvent := protoEvt.TradeEvent{
				Exchange:   constExch.Binance,
				MarketName: symbol,
				Price:      p,
			}
			marketTickers = append(marketTickers, &tickerEvent)
		}
		payload := protoEvt.TradeEvents{
			Events: marketTickers,
		}

		if err := c.Publisher.Publish(context.Background(), &payload); err != nil {
			log.Println("publish warning: ", err)
		}
	}
}

func (c *BinanceClient) run() {
	// loop indefinitely
	for {
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

		// in theory when the readPump returns the defer statement in that function
		// should close the socket connection gracefully
		log.Println("...binance websocket has stopped, reconnecting in 5 seconds")

		// wait for 10 seconds to pull the trade results
		time.Sleep(5 * time.Second)
	}
}

func main() {
	srv := k8s.NewService(
		micro.Name("fomo.binance.websocket"),
	)

	srv.Init()
	tradePublisher := micro.NewPublisher(constMessage.TopicAggTrade, srv.Client())

	client := BinanceClient{
		Publisher:   tradePublisher,
		lastReceive: time.Now().UTC(),
	}

	go client.run()

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
