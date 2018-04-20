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
)

func main() {
	srv := micro.NewService(
		micro.Name("micro.execution.engine"),
	)

	srv.Init()

	dbUrl := fmt.Sprintf("%s", os.Getenv("DB_URL"))
	gomoDB, err := db.NewDB(dbUrl)
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
		DB:       gomoDB,
		Receiver: &buyReceiver,
	}
	sellProcessor := SellProcessor{
		DB:       gomoDB,
		Receiver: &sellReceiver,
	}

	DeclareConditions(env)

	// subscribe to new key topic with a key validator
	micro.RegisterSubscriber(msg.TopicNewBuyOrder, srv.Server(), &buyReceiver)
	micro.RegisterSubscriber(msg.TopicNewSellOrder, srv.Server(), &sellReceiver)
	micro.RegisterSubscriber(msg.TopicAggTrade, srv.Server(), &buyProcessor)
	micro.RegisterSubscriber(msg.TopicAggTrade, srv.Server(), &sellProcessor)

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}

}
