package main

import (
	"fmt"
	"log"
	"os"

	"github.com/asciiu/socketcluster-client-go/scclient"
	micro "github.com/micro/go-micro"
)

type auth struct {
	ApiKey    string `json:"apiKey"`
	ApiSecret string `json:"apiSecret"`
}

func onConnect(client scclient.Client) {
	fmt.Println("Connected to server")
}

func onDisconnect(client scclient.Client, err error) {
	fmt.Printf("Error: %s\n", err.Error())
}

func onConnectError(client scclient.Client, err error) {
	fmt.Printf("Error: %s\n", err.Error())
}

func onSetAuthentication(client scclient.Client, token string) {
	fmt.Println("Auth token received :", token)
}

func onAuthentication(client scclient.Client, isAuthenticated bool) {
	fmt.Println("Client authenticated :", isAuthenticated)
	auth := auth{
		ApiKey:    fmt.Sprintf("%s", os.Getenv("API_KEY")),
		ApiSecret: fmt.Sprintf("%s", os.Getenv("API_SECRET")),
	}

	client.EmitAck("auth", auth, func(eventName string, error interface{}, data interface{}) {
		if error == nil {
			go startCode(client)
		}
	})
}

func main() {
	client := scclient.New("wss://sc-02.coinigy.com/socketcluster/")
	client.SetBasicListener(onConnect, onConnectError, onDisconnect)
	client.SetAuthenticationListener(onSetAuthentication, onAuthentication)
	go client.Connect()

	srv := micro.NewService(
		micro.Name("micro.coinigy.websocket"),
	)
	srv.Init()
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}

func startCode(client scclient.Client) {
	//client.Emit("channels", "BTRX")
	client.SubscribeAck("TRADE-BTRX--ADA--BTC", func(channelName string, error interface{}, data interface{}) {
		if error == nil {
			fmt.Println("Subscribed to channel ", channelName, "successfully")
		}
	})

	client.OnChannel("TRADE-BTRX--ADA--BTC", func(channelName string, data interface{}) {
		fmt.Println("Got data ", data, " for channel ", channelName)
	})
}
