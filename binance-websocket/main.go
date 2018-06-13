package main

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/asciiu/gomo/common/constants/exchange"
	msg "github.com/asciiu/gomo/common/constants/messages"
	evt "github.com/asciiu/gomo/common/proto/events"
	"github.com/gorilla/websocket"
	micro "github.com/micro/go-micro"
	k8s "github.com/micro/kubernetes/go/micro"
)

type BinanceConnection struct {
	group     *sync.WaitGroup
	Publisher micro.Publisher
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

func (bconn *BinanceConnection) Ticker() {
	url := "wss://stream.binance.com:9443/ws/!ticker@arr"
	log.Printf("connecting to %s", url)

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	// close the connection when this function returns
	defer conn.Close()

	conn.SetCloseHandler(func(code int, text string) error {
		log.Printf("closed connection %d %s\n", code, text)
		conn.Close()
		return nil
	})
	conn.SetPingHandler(func(appData string) error {
		log.Println("ping: ", appData)
		if err := conn.WriteMessage(websocket.PongMessage, []byte("pong")); err != nil {
			log.Println("ping error")
		}

		return nil
	})
	conn.SetPongHandler(func(appData string) error {
		log.Println("pong: ", appData)
		if err := conn.WriteMessage(websocket.PingMessage, []byte("ping")); err != nil {
			log.Println("pong error")
		}
		return nil
	})

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("websocket err: ", err)
			return
		}

		//aggTrade := BinanceAggTrade{}
		ticker := []*BinanceTicker{}
		if err = json.Unmarshal(message, &ticker); err != nil {
			log.Fatal("dial:", err)
		}

		for _, tick := range ticker {
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
				//EventTime:  tm.String(),
			}

			//fmt.Println(tickerEvent)

			if err := bconn.Publisher.Publish(context.Background(), &tickerEvent); err != nil {
				log.Println("publish warning: ", err, tickerEvent)
			}
		}
	}
	log.Println("websocket has stopped!")
}

func main() {
	srv := k8s.NewService(
		micro.Name("fomo.binance.websocket"),
	)

	srv.Init()
	tradePublisher := micro.NewPublisher(msg.TopicAggTrade, srv.Client())

	bconn := BinanceConnection{
		Publisher: tradePublisher,
	}

	go bconn.Ticker()

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
