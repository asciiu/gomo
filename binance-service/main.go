package main

import (
	"log"

	protoBinance "github.com/asciiu/gomo/binance-service/proto/binance"
	constMessage "github.com/asciiu/gomo/common/constants/message"
	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/server"
	k8s "github.com/micro/kubernetes/go/micro"
)

func main() {
	srv := k8s.NewService(
		micro.Name("fomo.binance"),
	)
	srv.Init()

	//candleRetriever := CandleRetriever{}
	bservice := BinanceService{
		CompletedPub: micro.NewPublisher(constMessage.TopicCompletedOrder, srv.Client()),
	}

	// subscribe to new key topic with a key validator
	micro.RegisterSubscriber(constMessage.TopicFillBinanceOrder, srv.Server(), bservice.HandleFillOrder, server.SubscriberQueue("fill.order"))
	//micro.RegisterSubscriber(constMessage.TopicCandleDataRequest, srv.Server(), &candleRetriever)

	protoBinance.RegisterBinanceServiceHandler(srv.Server(), &bservice)

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
