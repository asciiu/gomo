package main

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"testing"

	"github.com/asciiu/gomo/common/db"
	protoEvt "github.com/asciiu/gomo/common/proto/events"
	protoEngine "github.com/asciiu/gomo/execution-engine/proto/engine"
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

func setupEngine() *Engine {
	dbUrl := "postgres://postgres@localhost:5432/gomo_test?&sslmode=disable"
	db, _ := db.NewDB(dbUrl)

	env := vm.NewEnv()
	core.Import(env)

	engine := Engine{
		DB:    db,
		Plans: make([]*Plan, 0),
		Env:   env,
	}

	return &engine
}

// You shouldn't be able to insert a plan with no orders. A new plan requires at least a single order.
func TestAddPlan(t *testing.T) {
	engine := setupEngine()

	defer engine.DB.Close()

	req := protoEngine.NewPlanRequest{
		PlanID:                  "1aa6ae7f-76bc-49ba-a58b-5cf431e0337c",
		UserID:                  "30c30397-066c-4a48-9b9a-b778ef11f291",
		CommittedCurrencySymbol: "USDT",
		CommittedCurrencyAmount: 100.00,
		Orders: []*protoEngine.Order{
			&protoEngine.Order{
				OrderID:     "4d67dcac-0e46-49f5-9258-234fdba373ae",
				Exchange:    "testex",
				MarketName:  "BTC-USDT",
				Side:        "buy",
				LimitPrice:  4950.00,
				OrderType:   "limit",
				OrderStatus: "active",
				Triggers: []*protoEngine.Trigger{
					&protoEngine.Trigger{
						TriggerID: "76fa5d0c-34c7-4c67-ad64-389f92ce82a5",
						OrderID:   "4d67dcac-0e46-49f5-9258-234fdba373ae",
						Name:      "My Price!",
						Code:      "price <= 5000",
						Triggered: false,
						Actions:   []string{"placeOrder"},
					},
				},
			},
		},
	}
	res := protoEngine.PlanResponse{}
	engine.AddPlan(context.Background(), &req, &res)

	assert.Equal(t, 1, len(engine.Plans), "the engine should have a single plan")
	plan := engine.Plans[0]
	order := plan.Orders[0]
	trigger := order.TriggerExs[0]
	result, desc := trigger.Evaluate(4999.00)

	assert.Equal(t, true, result, "the trigger statement should be true")
	assert.Equal(t, "4999.00000000 <= 5000", desc, "the description for the trigger did not match")
}

// When an account is deleted any orders related to the account should be removed from execution
func TestDeletedAccount(t *testing.T) {
	engine := setupEngine()

	defer engine.DB.Close()

	req := protoEngine.NewPlanRequest{
		PlanID:                  "1aa6ae7f-76bc-49ba-a58b-5cf431e0337c",
		UserID:                  "30c30397-066c-4a48-9b9a-b778ef11f291",
		CommittedCurrencySymbol: "USDT",
		CommittedCurrencyAmount: 100.00,
		Orders: []*protoEngine.Order{
			&protoEngine.Order{
				AccountID:   "188077aa-2d7a-4b18-8011-1b3b32340e79",
				OrderID:     "4d67dcac-0e46-49f5-9258-234fdba373ae",
				Exchange:    "testex",
				MarketName:  "BTC-USDT",
				Side:        "buy",
				LimitPrice:  4950.00,
				OrderType:   "limit",
				OrderStatus: "active",
				Triggers: []*protoEngine.Trigger{
					&protoEngine.Trigger{
						TriggerID: "76fa5d0c-34c7-4c67-ad64-389f92ce82a5",
						OrderID:   "4d67dcac-0e46-49f5-9258-234fdba373ae",
						Name:      "My Price!",
						Code:      "price <= 5000",
						Triggered: false,
						Actions:   []string{"placeOrder"},
					},
				},
			},
		},
	}
	res := protoEngine.PlanResponse{}
	engine.AddPlan(context.Background(), &req, &res)
	engine.HandleAccountDeleted(context.Background(), &protoEvt.DeletedAccountEvent{"188077aa-2d7a-4b18-8011-1b3b32340e79"})

	assert.Equal(t, 0, len(engine.Plans), "the engine should have a no plans")
}

func TestRegex(t *testing.T) {
	stopLoss := regexp.MustCompile(`^.*?StopLoss\((0\.\d{2,}).*?`)
	str := "StopLoss(0.10)"
	switch {
	case stopLoss.MatchString(str):
		rs := stopLoss.FindStringSubmatch(str)
		fmt.Println(rs[1])
	}
	fmt.Println("done")
}
