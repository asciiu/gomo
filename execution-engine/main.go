package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/asciiu/gomo/common/db"
	msg "github.com/asciiu/gomo/common/messages"
	evt "github.com/asciiu/gomo/common/proto/events"
	"github.com/mattn/anko/core"
	"github.com/mattn/anko/vm"
	micro "github.com/micro/go-micro"
)

type OrderReceiver struct {
	DB     *sql.DB
	Orders []*evt.OrderEvent
}

func (receiver *OrderReceiver) ProcessEvent(ctx context.Context, buy *evt.OrderEvent) error {
	receiver.Orders = append(receiver.Orders, buy)
	return nil
}

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

	buyReceiver := OrderReceiver{gomoDB, make([]*evt.OrderEvent, 0)}
	sellReceiver := OrderReceiver{gomoDB, make([]*evt.OrderEvent, 0)}
	buyProcessor := BuyProcessor{
		DB:       gomoDB,
		Env:      env,
		Receiver: &buyReceiver,
	}
	sellProcessor := SellProcessor{
		DB:       gomoDB,
		Env:      env,
		Receiver: &sellReceiver,
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
