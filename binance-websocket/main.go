package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	msg "github.com/asciiu/gomo/common/messages"
	ep "github.com/asciiu/gomo/common/proto/events"
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
	TradeId      uint64 `json:"a"`
	Price        string `json:"p"`
	Quantity     string `json:"q"`
	FirstTradeId uint64 `json:"f"`
	LastTradeId  uint64 `json:"l"`
	TradeTime    uint64 `json:"T"`
	IsBuyerMaker bool   `json:"m"`
	Ignore       bool   `json:"M"`
}

func (bconn *BinanceConnection) Open(market string) {
	bconn.group.Add(1)
	url := fmt.Sprintf("wss://stream.binance.com:9443/ws/%s@aggTrade", market)
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

			aggTrade := BinanceAggTrade{}
			if err = json.Unmarshal(message, &aggTrade); err != nil {
				log.Println("nope")
			}

			exchangeEvent := ep.ExchangeEvent{
				Exchange:   "Binance",
				Type:       aggTrade.Type,
				EventTime:  aggTrade.EventTime,
				MarketName: aggTrade.Symbol,
				TradeId:    aggTrade.TradeId,
				Price:      aggTrade.Price,
				Quantity:   aggTrade.Quantity,
				TradeTime:  aggTrade.TradeTime,
			}

			if err := bconn.Publisher.Publish(context.Background(), &exchangeEvent); err != nil {
				log.Println("could not publish binance trade event: ", err)
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
		micro.Name("go.micro.srv.binance.websocket"),
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
