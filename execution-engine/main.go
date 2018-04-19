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
	micro "github.com/micro/go-micro"
)

type TradeProcessor struct {
	DB *sql.DB
}

func (engine *TradeProcessor) ProcessEvent(ctx context.Context, event *evt.ExchangeEvent) error {
	fmt.Println("new event ", event)
	return nil
}

type BuyOrderProcessor struct {
	DB *sql.DB
}

func (engine *BuyOrderProcessor) ProcessEvent(ctx context.Context, event *evt.OrderEvent) error {
	fmt.Println("new order ", event)
	return nil
}

type SellOrderProcessor struct {
	DB *sql.DB
}

func (engine *SellOrderProcessor) ProcessEvent(ctx context.Context, event *evt.OrderEvent) error {
	fmt.Println("new order ", event)
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

	tradeProcess := TradeProcessor{gomoDB}
	buyProcess := BuyOrderProcessor{gomoDB}
	sellProcess := SellOrderProcessor{gomoDB}

	// subscribe to new key topic with a key validator
	micro.RegisterSubscriber(msg.TopicAggTrade, srv.Server(), &tradeProcess)
	micro.RegisterSubscriber(msg.TopicNewBuyOrder, srv.Server(), &buyProcess)
	micro.RegisterSubscriber(msg.TopicNewSellOrder, srv.Server(), &sellProcess)

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}

}
