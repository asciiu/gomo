package main

import (
	"context"
	"log"
	"testing"

	constStatus "github.com/asciiu/gomo/common/constants/status"
	"github.com/asciiu/gomo/common/db"
	repoKey "github.com/asciiu/gomo/key-service/db/sql"
	protoKey "github.com/asciiu/gomo/key-service/proto/key"
	testKey "github.com/asciiu/gomo/key-service/test"
	repoOrder "github.com/asciiu/gomo/plan-service/db/sql"
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

// Test failure of order when incorrect order
func TestUnsortedPlan(t *testing.T) {
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
	order1 := protoOrder.NewOrderRequest{
		OrderID:         "4d671984-d7dd-4dce-a20f-23f25d6daf7f",
		KeyID:           key.KeyID,
		OrderType:       "paper",
		OrderTemplateID: "mokie22",
		ParentOrderID:   "4d671984-d7dd-4dce-a20f-23f25d6daf7f",
		MarketName:      "ADA-BTC",
		Side:            "buy",
		ActiveCurrencyBalance: 70,
		Triggers:              triggers,
	}
	order2 := protoOrder.NewOrderRequest{
		OrderID:         "4d671984-d7dd-4dce-a20f-23f25d6daf76",
		KeyID:           key.KeyID,
		OrderType:       "paper",
		OrderTemplateID: "mokie22",
		ParentOrderID:   "00000000-0000-0000-0000-000000000000",
		MarketName:      "BTC-USDT",
		Side:            "buy",
		ActiveCurrencyBalance: 70,
		Triggers:              triggers,
	}
	orders = append(orders, &order1, &order2)

	req := protoPlan.NewPlanRequest{
		UserID:          user.ID,
		Status:          "active",
		CloseOnComplete: false,
		PlanTemplateID:  "thing",
		Orders:          orders,
	}
	res := protoPlan.PlanResponse{}
	service.NewPlan(context.Background(), &req, &res)

	assert.Equal(t, "fail", res.Status, "an unsorted plan should have failed")

	repoUser.DeleteUserHard(service.DB, user.ID)
}

// Test updating a plan
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

	triggers = make([]*protoOrder.TriggerRequest, 0)
	trigger = protoOrder.TriggerRequest{
		TriggerID:         "ab4734f7-5ab7-46eb-9972-ed632ac752f8",
		Code:              "price <= 5",
		Index:             0,
		Name:              "price_trigger",
		Title:             "Testing update",
		TriggerTemplateID: "testtemplate",
		Actions:           []string{"placeOrder"},
	}
	triggers = append(triggers, &trigger)

	orders = make([]*protoOrder.NewOrderRequest, 0)
	order1 := protoOrder.NewOrderRequest{
		OrderID:         "4d671984-d7dd-4dce-a20f-23f25d6daf7f",
		KeyID:           key.KeyID,
		OrderType:       "paper",
		OrderTemplateID: "mokie22",
		ParentOrderID:   "00000000-0000-0000-0000-000000000000",
		MarketName:      "BTC-USDT",
		Side:            "buy",
		ActiveCurrencyBalance: 70,
		Triggers:              triggers,
	}
	order2 := protoOrder.NewOrderRequest{
		OrderID:         "4d671984-d7dd-4dce-a20f-23f25d6daf76",
		KeyID:           key.KeyID,
		OrderType:       "paper",
		OrderTemplateID: "mokie22",
		ParentOrderID:   "4d671984-d7dd-4dce-a20f-23f25d6daf7f",
		MarketName:      "ADA-BTC",
		Side:            "buy",
		Triggers:        triggers,
	}

	orders = append(orders, &order1, &order2)
	req2 := protoPlan.UpdatePlanRequest{
		PlanID:          res.Data.Plan.PlanID,
		UserID:          user.ID,
		Status:          "active",
		CloseOnComplete: false,
		PlanTemplateID:  "thing2",
		Orders:          orders,
	}
	res = protoPlan.PlanResponse{}
	service.UpdatePlan(context.Background(), &req2, &res)
	assert.Equal(t, "success", res.Status, "return status of updating an invalid plan should be fail")
	assert.Equal(t, 2, len(res.Data.Plan.Orders), "update should have yielded a single order")
	assert.Equal(t, 70.0, res.Data.Plan.Orders[0].ActiveCurrencyBalance, "active currency balance after update order incorrect")
	assert.Equal(t, "thing2", res.Data.Plan.PlanTemplateID, "the template ID should have changed")

	// Next verify that when we read the plan again that things indeed change
	req3 := protoPlan.GetUserPlanRequest{
		PlanID:     res.Data.Plan.PlanID,
		UserID:     user.ID,
		PlanDepth:  0,
		PlanLength: 10,
	}
	res = protoPlan.PlanResponse{}
	service.GetUserPlan(context.Background(), &req3, &res)
	assert.Equal(t, "success", res.Status, "return status for get plan should be success")
	assert.Equal(t, 2, len(res.Data.Plan.Orders), "update should have yielded a single order")
	assert.Equal(t, 70.0, res.Data.Plan.Orders[0].ActiveCurrencyBalance, "active currency balance after update order incorrect")
	assert.Equal(t, "thing2", res.Data.Plan.PlanTemplateID, "the template ID should have changed")

	repoUser.DeleteUserHard(service.DB, user.ID)
}

// Plan update should fail if the client sends in an update request on executed order.
func TestOrderUpdatePlanFailure(t *testing.T) {
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
	planID := res.Data.Plan.PlanID

	_, err := repoOrder.UpdateOrderStatus(service.DB, order.OrderID, constStatus.Filled)
	if err != nil {
		assert.Equal(t, "success", "failed", "update order status failed")
	}

	triggers = make([]*protoOrder.TriggerRequest, 0)
	trigger = protoOrder.TriggerRequest{
		TriggerID:         "ab4734f7-5ab7-46eb-9972-ed632ac752f8",
		Code:              "price <= 5",
		Index:             0,
		Name:              "price_trigger",
		Title:             "Testing update",
		TriggerTemplateID: "testtemplate",
		Actions:           []string{"placeOrder"},
	}
	triggers = append(triggers, &trigger)

	orders = make([]*protoOrder.NewOrderRequest, 0)
	order1 := protoOrder.NewOrderRequest{
		OrderID:         "4d671984-d7dd-4dce-a20f-23f25d6daf7f",
		KeyID:           key.KeyID,
		OrderType:       "paper",
		OrderTemplateID: "mokie22",
		ParentOrderID:   "00000000-0000-0000-0000-000000000000",
		MarketName:      "BTC-USDT",
		Side:            "buy",
		ActiveCurrencyBalance: 70,
		Triggers:              triggers,
	}
	order2 := protoOrder.NewOrderRequest{
		OrderID:         "4d671984-d7dd-4dce-a20f-23f25d6daf76",
		KeyID:           key.KeyID,
		OrderType:       "paper",
		OrderTemplateID: "mokie22",
		ParentOrderID:   "4d671984-d7dd-4dce-a20f-23f25d6daf7f",
		MarketName:      "ADA-BTC",
		Side:            "buy",
		Triggers:        triggers,
	}

	orders = append(orders, &order1, &order2)
	req2 := protoPlan.UpdatePlanRequest{
		PlanID:          res.Data.Plan.PlanID,
		UserID:          user.ID,
		Status:          "active",
		CloseOnComplete: false,
		PlanTemplateID:  "thing2",
		Orders:          orders,
	}
	res = protoPlan.PlanResponse{}
	service.UpdatePlan(context.Background(), &req2, &res)
	assert.Equal(t, "fail", res.Status, "return status of updating an invalid plan should be fail")

	// Next verify that when we read the plan again that things indeed change
	req3 := protoPlan.GetUserPlanRequest{
		PlanID:     planID,
		UserID:     user.ID,
		PlanDepth:  0,
		PlanLength: 10,
	}
	res = protoPlan.PlanResponse{}
	service.GetUserPlan(context.Background(), &req3, &res)
	assert.Equal(t, "success", res.Status, "return status for get plan should be success")
	assert.Equal(t, 1, len(res.Data.Plan.Orders), "update should have yielded a single order")
	assert.Equal(t, 100.0, res.Data.Plan.Orders[0].ActiveCurrencyBalance, "active currency balance after update order incorrect")
	assert.Equal(t, "thing", res.Data.Plan.PlanTemplateID, "the template ID should have changed")

	repoUser.DeleteUserHard(service.DB, user.ID)
}
