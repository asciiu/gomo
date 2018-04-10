// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type BinanceConnection struct {
	channel <-chan bool
	group   *sync.WaitGroup
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
				log.Printf("%s websocket client read: %s", market, err)
				return
			}
			log.Printf("recv: %s", message)
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
			channel: c,
			group:   &wg,
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
