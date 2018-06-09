package main

import (
	"context"
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

	orderReceiver := OrderReceiver{
		DB:     gomoDB,
		Orders: make([]*Order, 0),
		Env:    env,
	}
	processor := Processor{
		DB:       gomoDB,
		Receiver: &orderReceiver,
		Filled:   micro.NewPublisher(msg.TopicOrderFilled, srv.Client()),
		Filler:   micro.NewPublisher(msg.TopicFillOrder, srv.Client()),
	}

	// subscribe to new key topic with a key validator
	micro.RegisterSubscriber(msg.TopicNewOrder, srv.Server(), &orderReceiver)
	micro.RegisterSubscriber(msg.TopicAggTrade, srv.Server(), &processor)

	starter := micro.NewPublisher(msg.TopicEngineStart, srv.Client())
	starter.Publish(context.Background(), &processor)

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}

}
