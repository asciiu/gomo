package main

import (
	"log"

	protoBinance "github.com/asciiu/gomo/binance-service/proto/binance"
	constMessage "github.com/asciiu/gomo/common/constants/message"
	micro "github.com/micro/go-micro"
	k8s "github.com/micro/kubernetes/go/micro"
)

func main() {
	srv := k8s.NewService(
		micro.Name("fomo.binance"),
	)

	srv.Init()

	//verifiedPub := micro.NewPublisher(constMessage.TopicKeyVerified, srv.Client())
	//balancePub := micro.NewPublisher(constMessage.TopicBalanceUpdate, srv.Client())
	completedPub := micro.NewPublisher(constMessage.TopicCompletedOrder, srv.Client())

	//candleRetriever := CandleRetriever{}
	fulfiller := OrderFulfiller{
		CompletedPub: completedPub,
	}

	binanceService := new(BinanceService)

	// subscribe to new key topic with a key validator
	micro.RegisterSubscriber(constMessage.TopicTriggeredOrder, srv.Server(), &fulfiller)
	//micro.RegisterSubscriber(constMessage.TopicCandleDataRequest, srv.Server(), &candleRetriever)

	protoBinance.RegisterBinanceServiceHandler(srv.Server(), binanceService)

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
