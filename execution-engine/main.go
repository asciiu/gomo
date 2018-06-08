package main

import (
	"fmt"
	"log"
	"os"

	"github.com/asciiu/gomo/common/db"
	msg "github.com/asciiu/gomo/common/messages"
	"github.com/mattn/anko/core"
	"github.com/mattn/anko/vm"
	micro "github.com/micro/go-micro"
	k8s "github.com/micro/kubernetes/go/micro"
)

func main() {
	srv := k8s.NewService(
		micro.Name("fomo.execution"),
		micro.Version("latest"),
	)

	srv.Init()

	dbURL := fmt.Sprintf("%s", os.Getenv("DB_URL"))
	gomoDB, err := db.NewDB(dbURL)
	if err != nil {
		log.Fatalf(err.Error())
	}

	env := vm.NewEnv()
	core.Import(env)

	buyReceiver := OrderReceiver{
		DB:     gomoDB,
		Orders: make([]*Order, 0),
		Env:    env,
	}
	sellReceiver := OrderReceiver{
		DB:     gomoDB,
		Orders: make([]*Order, 0),
		Env:    env,
	}
	buyProcessor := BuyProcessor{
		DB:            gomoDB,
		Receiver:      &buyReceiver,
		Publisher:     micro.NewPublisher(msg.TopicOrderFilled, srv.Client()),
		FillPublisher: micro.NewPublisher(msg.TopicFillOrder, srv.Client()),
	}
	sellProcessor := SellProcessor{
		DB:            gomoDB,
		Receiver:      &sellReceiver,
		Publisher:     micro.NewPublisher(msg.TopicOrderFilled, srv.Client()),
		FillPublisher: micro.NewPublisher(msg.TopicFillOrder, srv.Client()),
	}

	// subscribe to new key topic with a key validator
	micro.RegisterSubscriber(msg.TopicNewBuyOrder, srv.Server(), &buyReceiver)
	micro.RegisterSubscriber(msg.TopicNewSellOrder, srv.Server(), &sellReceiver)
	micro.RegisterSubscriber(msg.TopicAggTrade, srv.Server(), &buyProcessor)
	micro.RegisterSubscriber(msg.TopicAggTrade, srv.Server(), &sellProcessor)

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}

}
