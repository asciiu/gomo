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

// order has conditions
type Order struct {
	EventOrigin *evt.OrderEvent
	Conditions  []ConditionFunc
}

type OrderReceiver struct {
	DB     *sql.DB
	Orders []*Order
	Env    *vm.Env
}

func (receiver *OrderReceiver) ProcessEvent(ctx context.Context, buy *evt.OrderEvent) error {
	// convert OrderEvent to Order with conditions here
	strConditions := strings.Split(buy.Conditions, " or ")
	conditions := make([]ConditionFunc, 0)

	//var extractParams = """^.*?TrailingStop\((0\.\d{2,}),\s(\d+\.\d+).*?""".r
	//var rNum = regexp.MustCompile(`\d`)  // Has digit(s)
	//var rAbc = regexp.MustCompile(`abc`) // Contains "abc"

	for _, str := range strConditions {
		priceCond := PriceCondition{
			Env:       receiver.Env,
			Statement: str,
		}
		conditions = append(conditions, priceCond.evaluate)
	}

	order := Order{
		EventOrigin: buy,
		Conditions:  conditions,
	}
	receiver.Orders = append(receiver.Orders, &order)

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
		Env:      env,
		Receiver: &buyReceiver,
	}
	sellProcessor := SellProcessor{
		DB:       gomoDB,
		Env:      env,
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
