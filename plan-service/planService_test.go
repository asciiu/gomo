package main

import (
	"context"
	"log"
	"testing"

	"github.com/asciiu/gomo/common/db"
	repoKey "github.com/asciiu/gomo/key-service/db/sql"
	protoKey "github.com/asciiu/gomo/key-service/proto/key"
	testKey "github.com/asciiu/gomo/key-service/test"
	protoOrder "github.com/asciiu/gomo/plan-service/proto/order"
	protoPlan "github.com/asciiu/gomo/plan-service/proto/plan"
	repoUser "github.com/asciiu/gomo/user-service/db/sql"
	user "github.com/asciiu/gomo/user-service/models"
	"github.com/stretchr/testify/assert"
)

func checkErr(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}

func setupService() (*PlanService, *user.User, *protoKey.Key) {
	dbUrl := "postgres://postgres@localhost:5432/gomo_test?&sslmode=disable"
	db, _ := db.NewDB(dbUrl)

	planService := PlanService{
		DB:        db,
		KeyClient: testKey.MockKeyServiceClient(),
	}

	user := user.NewUser("first", "last", "test@email", "hash")
	_, error := repoUser.InsertUser(db, user)
	keyRequest := protoKey.KeyRequest{
		KeyID:    "92512e72-b7e8-49f4-bab5-271f4ba450d9",
		UserID:   user.ID,
		Exchange: "test",
	}

	key, error := repoKey.InsertKey(db, &keyRequest)
	checkErr(error)

	return &planService, user, key
}

// You shouldn't be able to insert a plan with no orders. A new plan requires at least a single order.
func TestEmptyOrderPlan(t *testing.T) {
	service, user, _ := setupService()

	defer service.DB.Close()

	orders := make([]*protoOrder.NewOrderRequest, 0)
	req := protoPlan.NewPlanRequest{
		UserID:          user.ID,
		Status:          "active",
		CloseOnComplete: false,
		PlanTemplateID:  "thing",
		Orders:          orders,
	}
	res := protoPlan.PlanResponse{}
	service.NewPlan(context.Background(), &req, &res)

	assert.Equal(t, "fail", res.Status, "return status of inserting an empty order plan should be fail")

	repoUser.DeleteUserHard(service.DB, user.ID)
}

// Orders must have a trigger
func TestOrdersMustHaveTriggers(t *testing.T) {
	service, user, key := setupService()

	defer service.DB.Close()

	orders := make([]*protoOrder.NewOrderRequest, 0)
	order := protoOrder.NewOrderRequest{
		OrderID:         "4d671984-d7dd-4dce-a20f-23f25d6daf7f",
		KeyID:           key.KeyID,
		OrderType:       "paper",
		OrderTemplateID: "mokie",
		ParentOrderID:   "00000000-0000-0000-0000-000000000000",
		MarketName:      "ADA-BTC",
		Side:            "buy",
		ActiveCurrencyBalance: 100,
	}
	orders = append(orders, &order)
	req := protoPlan.NewPlanRequest{
		UserID:          user.ID,
		Status:          "active",
		CloseOnComplete: false,
		PlanTemplateID:  "thing",
		Orders:          orders,
	}
	res := protoPlan.PlanResponse{}
	service.NewPlan(context.Background(), &req, &res)

	assert.Equal(t, "fail", res.Status, "return status should be fail when submitting orders without triggers")

	repoUser.DeleteUserHard(service.DB, user.ID)
}

// Successfully inserting a new plan
func TestSuccessfulOrderPlan(t *testing.T) {
	service, user, key := setupService()

	defer service.DB.Close()

	triggers := make([]*protoOrder.TriggerRequest, 0)
	trigger := protoOrder.TriggerRequest{
		TriggerID:         "ab4734f7-5ab7-46eb-9972-ed632ac752f8",
		Code:              "test",
		Index:             0,
		Name:              "test_trigger",
		Title:             "Testing",
		TriggerTemplateID: "testtemplate",
		Actions:           []string{"placeOrder"},
	}
	triggers = append(triggers, &trigger)

	orders := make([]*protoOrder.NewOrderRequest, 0)
	order := protoOrder.NewOrderRequest{
		OrderID:         "4d671984-d7dd-4dce-a20f-23f25d6daf7f",
		KeyID:           key.KeyID,
		OrderType:       "paper",
		OrderTemplateID: "mokie",
		ParentOrderID:   "00000000-0000-0000-0000-000000000000",
		MarketName:      "ADA-BTC",
		Side:            "buy",
		ActiveCurrencyBalance: 100,
		Triggers:              triggers,
	}
	orders = append(orders, &order)
	req := protoPlan.NewPlanRequest{
		UserID:          user.ID,
		Status:          "active",
		CloseOnComplete: false,
		PlanTemplateID:  "thing",
		Orders:          orders,
	}
	res := protoPlan.PlanResponse{}
	service.NewPlan(context.Background(), &req, &res)

	assert.Equal(t, "success", res.Status, "return status of inserting plan should be success")

	repoUser.DeleteUserHard(service.DB, user.ID)
}

func TestOrderUpdatePlan(t *testing.T) {
	service, user, key := setupService()

	defer service.DB.Close()

	triggers := make([]*protoOrder.TriggerRequest, 0)
	trigger := protoOrder.TriggerRequest{
		TriggerID:         "ab4734f7-5ab7-46eb-9972-ed632ac752f8",
		Code:              "test",
		Index:             0,
		Name:              "test_trigger",
		Title:             "Testing",
		TriggerTemplateID: "testtemplate",
		Actions:           []string{"placeOrder"},
	}
	triggers = append(triggers, &trigger)

	orders := make([]*protoOrder.NewOrderRequest, 0)
	order := protoOrder.NewOrderRequest{
		OrderID:         "4d671984-d7dd-4dce-a20f-23f25d6daf7f",
		KeyID:           key.KeyID,
		OrderType:       "paper",
		OrderTemplateID: "mokie",
		ParentOrderID:   "00000000-0000-0000-0000-000000000000",
		MarketName:      "ADA-BTC",
		Side:            "buy",
		ActiveCurrencyBalance: 100,
		Triggers:              triggers,
	}
	orders = append(orders, &order)
	req := protoPlan.NewPlanRequest{
		UserID:          user.ID,
		Status:          "active",
		CloseOnComplete: false,
		PlanTemplateID:  "thing",
		Orders:          orders,
	}
	res := protoPlan.PlanResponse{}
	service.NewPlan(context.Background(), &req, &res)

	assert.Equal(t, "success", res.Status, "return status of inserting plan should be success")

	orders = make([]*protoOrder.NewOrderRequest, 0)
	order = protoOrder.NewOrderRequest{
		OrderID:         "4d671984-d7dd-4dce-a20f-23f25d6daf7f",
		KeyID:           key.KeyID,
		OrderType:       "paper",
		OrderTemplateID: "mokie",
		ParentOrderID:   "00000000-0000-0000-0000-000000000000",
		MarketName:      "ADA-BTC",
		Side:            "buy",
		ActiveCurrencyBalance: 100,
	}

	orders = append(orders, &order)
	req2 := protoPlan.UpdatePlanRequest{
		PlanID:          res.Data.Plan.PlanID,
		UserID:          user.ID,
		Status:          "active",
		CloseOnComplete: false,
		PlanTemplateID:  "thing",
		Orders:          orders,
	}
	res = protoPlan.PlanResponse{}
	service.UpdatePlan(context.Background(), &req2, &res)

	assert.Equal(t, "fail", res.Status, "return status of updating an invalid plan should be fail")

	repoUser.DeleteUserHard(service.DB, user.ID)
}
