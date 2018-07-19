package main

import (
	"context"
	"fmt"
	"log"
	"os"

	msg "github.com/asciiu/gomo/common/constants/messages"
	"github.com/asciiu/gomo/common/db"
	evt "github.com/asciiu/gomo/common/proto/events"
	micro "github.com/micro/go-micro"
	k8s "github.com/micro/kubernetes/go/micro"
)

func main() {
	srv := k8s.NewService(
		micro.Name("fomo.analytics"),
		micro.Version("latest"),
	)

	srv.Init()

	dbURL := fmt.Sprintf("%s", os.Getenv("DB_URL"))
	gomoDB, err := db.NewDB(dbURL)
	if err != nil {
		log.Fatalf(err.Error())
	}

	processor := Processor{
		DB:               gomoDB,
		MarketClosePrice: make(map[Market]float64),
		CandlePub:        micro.NewPublisher(msg.TopicCandleDataRequest, srv.Client()),
	}

	// subscribe to new key topic with a key validator
	micro.RegisterSubscriber(msg.TopicAggTrade, srv.Server(), func(ctx context.Context, tradeEvents *evt.TradeEvents) error {
		processor.ProcessEvents(tradeEvents)
		return nil
	})

	go processor.Ticker()

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
