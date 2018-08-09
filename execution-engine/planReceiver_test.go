package main

import (
	"context"
	"log"
	"testing"

	"github.com/asciiu/gomo/common/db"
	protoEvents "github.com/asciiu/gomo/common/proto/events"
	"github.com/mattn/anko/core"
	"github.com/mattn/anko/vm"
	"github.com/stretchr/testify/assert"
)

func checkErr(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}

func setupReceiver() *PlanReceiver {
	dbUrl := "postgres://postgres@localhost:5432/gomo_test?&sslmode=disable"
	db, _ := db.NewDB(dbUrl)

	env := vm.NewEnv()
	core.Import(env)

	planReceiver := PlanReceiver{
		DB:    db,
		Plans: make([]*Plan, 0),
		Env:   env,
	}

	return &planReceiver
}

// You shouldn't be able to insert a plan with no orders. A new plan requires at least a single order.
func TestAddPlan(t *testing.T) {
	receiver := setupReceiver()

	defer receiver.DB.Close()

	newPlanEvt := protoEvents.NewPlanEvent{
		PlanID: "1aa6ae7f-76bc-49ba-a58b-5cf431e0337c",
		UserID: "30c30397-066c-4a48-9b9a-b778ef11f291",
		Orders: []*protoEvents.Order{
			&protoEvents.Order{
				OrderID:               "4d67dcac-0e46-49f5-9258-234fdba373ae",
				Exchange:              "testex",
				MarketName:            "BTC-USDT",
				Side:                  "buy",
				LimitPrice:            1.00,
				OrderType:             "paper",
				OrderStatus:           "active",
				ActiveCurrencySymbol:  "USDT",
				ActiveCurrencyBalance: 100.00,
				Triggers: []*protoEvents.Trigger{
					&protoEvents.Trigger{
						TriggerID: "",
						OrderID:   "4d67dcac-0e46-49f5-9258-234fdba373ae",
						Name:      "triggex",
						Code:      "price <= 5000",
						Triggered: false,
						Actions:   []string{"placeOrder"},
					},
				},
			},
		},
	}
	receiver.AddPlan(context.Background(), &newPlanEvt)

	assert.Equal(t, 1, len(receiver.Plans), "the receiver should have a single plan")
	plan := receiver.Plans[0]
	order := plan.Orders[0]
	trigger := order.TriggerExs[0]
	result, desc := trigger.Evaluate(4999.00)

	assert.Equal(t, true, result, "the trigger statement should be true")
	assert.Equal(t, "4999.00000000 <= 5000", desc, "what?")
}
