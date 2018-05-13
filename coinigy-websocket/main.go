package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Freeaqingme/go-socketcluster-client"
	msg "github.com/asciiu/gomo/common/messages"
	evt "github.com/asciiu/gomo/common/proto/events"
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

// {"market_history_id":160554809190,"exchange":"BTRX","marketid":0,"label":"KMD/BTC",
//"tradeid":"29151260","time":"2018-05-13T05:12:27","price":0.000353,"quantity":300,
//"total":0.1059,"timestamp":"2018-05-13T05:13:10","time_local":"2018-05-13 05:12:27",
//"type":"BUY","exchId":15,"channel":"TRADE-BTRX--KMD--BTC"}
type CoinigyTrade struct {
	MarketHistoryID int     `json:"market_history_id"`
	ExCode          string  `json:"exchange"`
	MarketID        int     `json:"marketid"`
	Label           string  `json:"label"`
	TradeID         string  `json:"tradeid"`
	Time            string  `json:"time"`
	Price           float64 `json:"price"`
	Quantity        float64 `json:"quantity"`
	Total           float64 `json:"total"`
	TimeStamp       string  `json:"timestamp"`
	TimeLocal       string  `json:"time_local"`
	Type            string  `json:"type"`
	ExchangeID      int     `json:"exchId"`
	Channel         string  `json:"channel"`
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
		ApiKey:    os.Getenv("API_KEY"),
		ApiSecret: os.Getenv("API_SECRET"),
	}
	client.ConnectCallback = func() error {
		_, err := client.Emit("auth", &auth)
		return err
	}

	if err := client.Connect(); err != nil {
		panic(err)
	}

	exs, err := exchanges(client)
	if err != nil {
		panic(err)
	}

	channelNames := make([]string, 0)
	exchangeNames := make(map[string]string)
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

		exchangeNames[e.ExCode] = e.ExName

		for _, c := range e.Channels {
			if strings.Contains(c.ChannelName, "TRADE") {
				channelNames = append(channelNames, c.ChannelName)
			}
		}
	}

	//fmt.Println(exchangeNames)

	srv := micro.NewService(
		micro.Name("micro.coinigy.websocket"),
	)
	srv.Init()
	tradePublisher := micro.NewPublisher(msg.TopicAggTrade, srv.Client())

	for i, name := range channelNames {
		if i < 250 {
			channel, err := client.Subscribe(name)
			if err != nil {
				panic(err)
			}

			go func() {
				for msg := range channel {
					var md CoinigyTrade
					err := json.Unmarshal(msg, &md)
					if err != nil {
						log.Println(err)
					}

					tradeEvent := evt.TradeEvent{
						Exchange:   exchangeNames[md.ExCode],
						Type:       md.Type,
						EventTime:  md.TimeStamp,
						MarketName: md.Label,
						TradeID:    md.TradeID,
						Price:      md.Price,
						Quantity:   md.Quantity,
						Total:      md.Total,
					}

					if err := tradePublisher.Publish(context.Background(), &tradeEvent); err != nil {
						log.Println("publish warning: ", err, tradeEvent)
					}

					//fmt.Println(string(msg))
					//fmt.Println(md)
				}
			}()

		}
	}

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
