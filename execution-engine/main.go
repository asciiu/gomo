package main

import (
	"context"
	"fmt"
	"log"
	"os"

	constMessage "github.com/asciiu/gomo/common/constants/message"
	"github.com/asciiu/gomo/common/db"
	protoEvt "github.com/asciiu/gomo/common/proto/events"
	protoEngine "github.com/asciiu/gomo/execution-engine/proto/engine"
	"github.com/mattn/anko/core"
	"github.com/mattn/anko/vm"
	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/server"
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

	engine := Engine{
		DB:          gomoDB,
		Env:         env,
		Completed:   micro.NewPublisher(constMessage.TopicCompletedOrder, srv.Client()),
		FillBinance: micro.NewPublisher(constMessage.TopicFillBinanceOrder, srv.Client()),
		PriceLine:   make(map[string]float64),
		Plans:       make([]*Plan, 0),
	}
	protoEngine.RegisterExecutionEngineHandler(srv.Server(), &engine)

	// subscribe to trade events from the exchanges
	micro.RegisterSubscriber(constMessage.TopicAggTrade, srv.Server(), engine.HandleTradeEvents, server.SubscriberQueue("trade.event"))
	micro.RegisterSubscriber(constMessage.TopicAccountDeleted, srv.Server(), engine.HandleAccountDeleted, server.SubscriberQueue("account.deleted"))

	// fire this event on startup to tell the plan service to feed the engine active plans
	starter := micro.NewPublisher(constMessage.TopicEngineStart, srv.Client())
	// TODO when the engine is load balanced we need to know what the id of the engine is
	// as well as the plans that belong to the engine's id so we can reload the engine state
	starter.Publish(context.Background(), &protoEvt.EngineStartEvent{"replaceIDHERE"})

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
