package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Freeaqingme/go-socketcluster-client"
	micro "github.com/micro/go-micro"
)

type auth struct {
	ApiKey    string `json:"apiKey"`
	ApiSecret string `json:"apiSecret"`
}

const wsUrl = "wss://sc-02.coinigy.com/socketcluster/"

func main() {
	client := scclient.New(wsUrl)

	// Supply a callback for any events that need to be performed upon every reconnnect
	auth := auth{
		ApiKey:    fmt.Sprintf("%s", os.Getenv("API_KEY")),
		ApiSecret: fmt.Sprintf("%s", os.Getenv("API_SECRET")),
	}
	client.ConnectCallback = func() error {
		_, err := client.Emit("auth", &auth)
		return err
	}

	if err := client.Connect(); err != nil {
		panic(err)
	}

	channel, err := client.Subscribe("TRADE-KRKN--XBT--EUR")
	if err != nil {
		panic(err)
	}

	go func() {
		for msg := range channel {
			fmt.Println("New kraken trade: " + string(msg))
		}
	}()

	if res, err := client.Emit("exchanges", nil); err != nil {
		panic(err)
	} else {
		fmt.Println("exchanges: " + string(res))
	}

	srv := micro.NewService(
		micro.Name("micro.coinigy.websocket"),
	)
	srv.Init()
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
