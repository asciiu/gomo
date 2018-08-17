package main

import (
	"context"
	"fmt"
	"log"
	"os"

	msg "github.com/asciiu/gomo/common/constants/messages"
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
		DB:        gomoDB,
		Env:       env,
		Aborted:   micro.NewPublisher(msg.TopicAbortedOrder, srv.Client()),
		Completed: micro.NewPublisher(msg.TopicCompletedOrder, srv.Client()),
		Triggered: micro.NewPublisher(msg.TopicTriggeredOrder, srv.Client()),
		PriceLine: make(map[string]float64),
		Plans:     make([]*Plan, 0),
	}
	protoEngine.RegisterExecutionEngineHandler(srv.Server(), &engine)

	// subscribe to trade events from the exchanges
	micro.RegisterSubscriber(msg.TopicAggTrade, srv.Server(), engine.ProcessTradeEvents, server.SubscriberQueue("trade.event"))

	// fire this event on startup to tell the plan service to feed the engine active plans
	starter := micro.NewPublisher(msg.TopicEngineStart, srv.Client())
	starter.Publish(context.Background(), &protoEvt.EngineStartEvent{"replaceIDHERE"})

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
