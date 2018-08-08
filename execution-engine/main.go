package main

import (
	"context"
	"fmt"
	"log"
	"os"

	msg "github.com/asciiu/gomo/common/constants/messages"
	"github.com/asciiu/gomo/common/db"
	evt "github.com/asciiu/gomo/common/proto/events"
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

	planReceiver := PlanReceiver{
		DB:      gomoDB,
		Plans:   make([]*Plan, 0),
		Env:     env,
		Aborted: micro.NewPublisher(msg.TopicAbortedOrder, srv.Client()),
	}
	processor := Processor{
		DB:        gomoDB,
		Receiver:  &planReceiver,
		Completed: micro.NewPublisher(msg.TopicCompletedOrder, srv.Client()),
		Triggered: micro.NewPublisher(msg.TopicTriggeredOrder, srv.Client()),
	}

	// subscribe to new key topic with a key validator
	micro.RegisterSubscriber(msg.TopicNewPlan, srv.Server(), &planReceiver)
	micro.RegisterSubscriber(msg.TopicAggTrade, srv.Server(), &processor)

	starter := micro.NewPublisher(msg.TopicEngineStart, srv.Client())
	starter.Publish(context.Background(), &evt.EngineStartEvent{"replaceIDHERE"})

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}

}
