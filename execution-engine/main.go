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

type TradeEngine struct {
	DB *sql.DB
}

func (engine *TradeEngine) Process(ctx context.Context, event *evt.ExchangeEvent) error {
	fmt.Println("new event ", event)
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

	// subscribe to new key topic with a key validator
	micro.RegisterSubscriber(msg.TopicAggTrade, srv.Server(), &TradeEngine{gomoDB})
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}

}
