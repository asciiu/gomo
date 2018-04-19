package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/asciiu/gomo/common/db"
	msg "github.com/asciiu/gomo/common/messages"
	evt "github.com/asciiu/gomo/common/proto/events"
	"github.com/mattn/anko/core"
	"github.com/mattn/anko/vm"
	micro "github.com/micro/go-micro"
)

type TradeProcessor struct {
	DB           *sql.DB
	Env          *vm.Env
	BuyReceiver  *BuyOrderReceiver
	SellReceiver *SellOrderReceiver
}

func (engine *TradeProcessor) ProcessEvent(ctx context.Context, event *evt.ExchangeEvent) error {
	buyOrders := engine.BuyReceiver.Orders

	for i, buyOrder := range buyOrders {

		marketName := strings.Replace(buyOrder.MarketName, "-", "", 1)
		if marketName != event.MarketName {
			continue
		}

		conditions := strings.Replace(buyOrder.Conditions, "price", event.Price, -1)
		fmt.Println(conditions)

		result, err := engine.Env.Execute(conditions)
		if err != nil {
			panic(err)
		}

		if result == true {
			// remove order
			engine.BuyReceiver.Orders = append(buyOrders[:i], buyOrders[i+1:]...)
			fmt.Println("BUY NOW!!")
		}
	}

	for _, sellOrder := range engine.SellReceiver.Orders {
		fmt.Println(sellOrder)

		conditions := strings.Replace(sellOrder.Conditions, "price", event.Price, -1)

		fmt.Println(conditions)

		v, err := engine.Env.Execute(conditions)

		if err != nil {
			panic(err)
		}

		if v == true {
			fmt.Println("Sell Order Execute NOW!!")
		}
	}

	fmt.Println("trade event ", event)
	return nil
}

type BuyOrderReceiver struct {
	DB     *sql.DB
	Orders []*evt.OrderEvent
}

func (receiver *BuyOrderReceiver) ProcessEvent(ctx context.Context, buy *evt.OrderEvent) error {
	receiver.Orders = append(receiver.Orders, buy)
	fmt.Println("received new buy ", buy)
	return nil
}

type SellOrderReceiver struct {
	DB     *sql.DB
	Orders []*evt.OrderEvent
}

func (receiver *SellOrderReceiver) ProcessEvent(ctx context.Context, sell *evt.OrderEvent) error {
	receiver.Orders = append(receiver.Orders, sell)
	fmt.Println("received new sell ", sell)
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

	buyReceiver := BuyOrderReceiver{gomoDB, make([]*evt.OrderEvent, 0)}
	sellReceiver := SellOrderReceiver{gomoDB, make([]*evt.OrderEvent, 0)}
	tradeProcess := TradeProcessor{
		DB:           gomoDB,
		Env:          env,
		BuyReceiver:  &buyReceiver,
		SellReceiver: &sellReceiver,
	}

	// subscribe to new key topic with a key validator
	micro.RegisterSubscriber(msg.TopicAggTrade, srv.Server(), &tradeProcess)
	micro.RegisterSubscriber(msg.TopicNewBuyOrder, srv.Server(), &buyReceiver)
	micro.RegisterSubscriber(msg.TopicNewSellOrder, srv.Server(), &sellReceiver)

	//v, err := env.Execute(`4 < 10`)
	//if err != nil {
	//	panic(err)
	//}

	//fmt.Println(v)

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}

}
