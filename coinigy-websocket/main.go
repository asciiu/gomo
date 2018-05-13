package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Freeaqingme/go-socketcluster-client"
	micro "github.com/micro/go-micro"
)

type auth struct {
	ApiKey    string `json:"apiKey"`
	ApiSecret string `json:"apiSecret"`
}

type Exchange struct {
	ExID     float64           `json:"exch_id"`
	ExName   string            `json:"exch_name"`
	ExCode   string            `json:"exch_code"`
	ExFee    float64           `json:"exch_fee"`
	Channels []ExchangeChannel `json:"channels"`
}

type ExchangeChannel struct {
	ChannelName string `json:"channel"`
}

const wsUrl = "wss://sc-02.coinigy.com/socketcluster/"

var supportedExchanges = [...]string{
	"BTRX",
	"BINA",
	"KUCN",
}

func exchanges(client *scclient.Client) ([]Exchange, error) {
	if res, err := client.Emit("exchanges", nil); err != nil {
		return nil, err
	} else {
		//fmt.Println(string(res))

		var f []interface{}
		err := json.Unmarshal(res, &f)
		if err != nil {
			fmt.Println("error: ", err)
		}

		//ex := make([]Exchange, 0)

		exchanges := make([]Exchange, 0)
		jsonExchanges := f[0].([]interface{})
		//err = json.Unmarshal([]byte(exchanges), &ex)
		for _, jsonExchange := range jsonExchanges {
			var exchange Exchange
			ex := jsonExchange.(map[string]interface{})
			exchange.ExName = ex["exch_name"].(string)
			exchange.ExCode = ex["exch_code"].(string)
			exchange.ExFee = ex["exch_fee"].(float64)
			exchange.ExID = ex["exch_id"].(float64)

			chs, err := channels(client, exchange.ExCode)
			if err == nil {
				exchange.Channels = chs
			}
			exchanges = append(exchanges, exchange)

			if err != nil {
				panic(err)
			}
		}
		return exchanges, nil
	}
}

func channels(client *scclient.Client, exchange string) ([]ExchangeChannel, error) {
	if res, err := client.Emit("channels", exchange); err != nil {
		return nil, err
	} else {
		var f []interface{}
		err := json.Unmarshal(res, &f)
		if err != nil {
			fmt.Println("error: ", err)
		}

		chs := make([]ExchangeChannel, 0)
		jsonChannels := f[0].([]interface{})
		//err = json.Unmarshal([]byte(exchanges), &ex)
		for _, j := range jsonChannels {
			//fmt.Println(channel)
			var c ExchangeChannel
			m := j.(map[string]interface{})
			c.ChannelName = m["channel"].(string)
			chs = append(chs, c)
		}
		return chs, nil
	}
}

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

	// channel, err := client.Subscribe("TRADE-KRKN--XBT--EUR")
	// if err != nil {
	// 	panic(err)
	// }

	// go func() {
	// 	for msg := range channel {
	// 		fmt.Println("New kraken trade: " + string(msg))
	// 	}
	// }()
	exs, err := exchanges(client)
	if err != nil {
		panic(err)
	}

	channelNames := make([]string, 0)
	// loop through all exchanges
	for _, e := range exs {
		var flag = false

		for _, a := range supportedExchanges {
			if a == e.ExCode {
				flag = true
			}
		}

		if !flag {
			continue
		}

		for _, c := range e.Channels {
			if strings.Contains(c.ChannelName, "TRADE") {
				channelNames = append(channelNames, c.ChannelName)
			}
		}
	}

	for i, name := range channelNames {
		if i < 250 {
			channel, err := client.Subscribe(name)
			if err != nil {
				panic(err)
			}

			go func() {
				for msg := range channel {
					fmt.Println(string(msg))
				}
			}()

		}
	}

	srv := micro.NewService(
		micro.Name("micro.coinigy.websocket"),
	)
	srv.Init()
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
