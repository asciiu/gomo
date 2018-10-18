package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	// "github.com/asciiu/gomo/analytics-service/fomolang/repl"
	protoAnalytics "github.com/asciiu/gomo/analytics-service/proto/analytics"
	constMessage "github.com/asciiu/gomo/common/constants/message"
	"github.com/asciiu/gomo/common/db"
	protoEvt "github.com/asciiu/gomo/common/proto/events"
	micro "github.com/micro/go-micro"
	k8s "github.com/micro/kubernetes/go/micro"
)

func main() {

	// user, err := user.Current()
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Printf("Hello %s! This is the Fomo programming language!\n",
	// 	user.Username)
	// fmt.Printf("Feel free to type in commands\n")

	// repl.Start(os.Stdin, os.Stdout)
	srv := k8s.NewService(
		micro.Name("analytics"),
		micro.Version("latest"),
	)

	srv.Init()

	dbURL := fmt.Sprintf("%s", os.Getenv("DB_URL"))
	gomoDB, err := db.NewDB(dbURL)
	if err != nil {
		log.Fatalf(err.Error())
	}

	priceTicker := PriceTicker{
		DB:           gomoDB,
		MarketPrices: make(map[Market]float64),
		TimePeriod:   time.Duration(5) * time.Minute, // 5 minute period
	}

	service := NewAnalyticsService(gomoDB, srv)
	// subscribe to the exchange events here
	micro.RegisterSubscriber(constMessage.TopicAggTrade, srv.Server(), func(ctx context.Context, tradeEvents *protoEvt.TradeEvents) error {
		priceTicker.HandleExchangeEvent(tradeEvents)
		service.HandleExchangeEvent(tradeEvents)
		return nil
	})
	//micro.RegisterSubscriber(constMessage.TopicStartIndicator, srv.Server(), service.HandleStartIndicator, server.SubscriberQueue("subscribe.indicator"))

	go priceTicker.Ticker()

	protoAnalytics.RegisterAnalyticsServiceHandler(srv.Server(), service)

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
