package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"time"

	exchanges "github.com/asciiu/gomo/common/exchanges"
	msg "github.com/asciiu/gomo/common/messages"
	evt "github.com/asciiu/gomo/common/proto/events"
	"github.com/gorilla/websocket"
	micro "github.com/micro/go-micro"
)

type BinanceConnection struct {
	channel   <-chan bool
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

func (bconn *BinanceConnection) Open(market string) {
	bconn.group.Add(1)
	//url := fmt.Sprintf("wss://stream.binance.com:9443/ws/%s@aggTrade", market)
	url := "wss://stream.binance.com:9443/ws/!ticker@arr"
	log.Printf("connecting to %s", url)

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	defer bconn.group.Done()

	// create channel for done
	done := make(chan struct{})

	go func() {
		defer close(done)
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
					Exchange:   exchanges.Binance,
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
	}()

	// loop indefinitely until one of the following happens
	for {
		select {
		case <-done:
			conn.Close()
			log.Println("unexpected close of market: ", market)
			return
		case <-bconn.channel:
			//Cleanly close the connection by sending a close message and then
			//waiting (with timeout) for the server to close the connection.
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}

			// wait on the above go routine to exit
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			conn.Close()
			return
		}
	}
}
func main() {
	srv := micro.NewService(
		micro.Name("micro.binance.websocket"),
	)

	srv.Init()
	tradePublisher := micro.NewPublisher(msg.TopicAggTrade, srv.Client())

	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	var wg sync.WaitGroup
	var channels []chan bool
	markets := [2]string{
		"adabtc",
		"bnbbtc",
	}

	for i := 0; i < len(markets); i++ {
		c := make(chan bool)
		bconn := BinanceConnection{
			channel:   c,
			group:     &wg,
			Publisher: tradePublisher,
		}
		channels = append(channels, c)
		go bconn.Open(markets[i])
	}

	for {
		select {
		case <-interrupt:

			for j := 0; j < len(channels); j++ {
				close(channels[j])
			}
			wg.Wait()
			log.Println("bye!")

			return
		}
	}
}
