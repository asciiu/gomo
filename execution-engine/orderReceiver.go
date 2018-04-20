package main

import (
	"context"
	"database/sql"
	"strings"

	evt "github.com/asciiu/gomo/common/proto/events"
	"github.com/mattn/anko/vm"
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
