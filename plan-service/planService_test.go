package main

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/asciiu/gomo/common/db"
	protoEngine "github.com/asciiu/gomo/execution-engine/proto/engine"
	testEngine "github.com/asciiu/gomo/execution-engine/test"
	constKey "github.com/asciiu/gomo/key-service/constants"
	repoKey "github.com/asciiu/gomo/key-service/db/sql"
	protoKey "github.com/asciiu/gomo/key-service/proto/key"
	testKey "github.com/asciiu/gomo/key-service/test"
	constPlan "github.com/asciiu/gomo/plan-service/constants"
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
		DB:           db,
		KeyClient:    testKey.MockKeyServiceClient(db),
		EngineClient: testEngine.MockEngineClient(db),
	}

	user := user.NewUser("first", "last", "test@email", "hash")
	_, error := repoUser.InsertUser(db, user)
	keyRequest := protoKey.KeyRequest{
		KeyID:       "92512e72-b7e8-49f4-bab5-271f4ba450d9",
		UserID:      user.ID,
		Exchange:    "test",
		Key:         "test_key",
		Secret:      "test_secret",
		Description: "testy",
	}

	key, error := repoKey.InsertKey(db, &keyRequest)
	checkErr(error)
	key.Status = constKey.Verified
	repoKey.UpdateKeyStatus(db, key)

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

	req := protoPlan.NewPlanRequest{
		UserID:          user.ID,
		Status:          "active",
		CloseOnComplete: false,
		PlanTemplateID:  "no_trigger_test",
		Orders: []*protoOrder.NewOrderRequest{
			&protoOrder.NewOrderRequest{
				OrderID:         "4d671984-d7dd-4dce-a20f-23f25d6daf7f",
				KeyID:           key.KeyID,
				OrderType:       "paper",
				OrderTemplateID: "mokie22",
				ParentOrderID:   "00000000-0000-0000-0000-000000000000",
				MarketName:      "BTC-USDT",
				Side:            "buy",
				InitialCurrencyBalance: 100,
			},
		},
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

	req := protoPlan.NewPlanRequest{
		UserID:          user.ID,
		Status:          "active",
		Title:           "Testing 123",
		CloseOnComplete: false,
		PlanTemplateID:  "bloody_test",
		Orders: []*protoOrder.NewOrderRequest{
			&protoOrder.NewOrderRequest{
				OrderID:         "4d671984-d7dd-4dce-a20f-23f25d6daf7f",
				KeyID:           key.KeyID,
				OrderType:       "paper",
				OrderTemplateID: "mokie22",
				ParentOrderID:   "00000000-0000-0000-0000-000000000000",
				MarketName:      "BTC-USDT",
				Side:            "buy",
				InitialCurrencyBalance: 70,
				Triggers: []*protoOrder.TriggerRequest{
					&protoOrder.TriggerRequest{
						TriggerID:         "ab4734f7-5ab7-46eb-9972-ed632ac752f8",
						Code:              "price <= 5",
						Index:             0,
						Name:              "price_trigger",
						Title:             "Testing update",
						TriggerTemplateID: "testtemplate",
						Actions:           []string{"placeOrder"},
					},
				},
			},
		},
	}

	res := protoPlan.PlanResponse{}
	service.NewPlan(context.Background(), &req, &res)

	assert.Equal(t, "success", res.Status, fmt.Sprintf("return status of inserting plan should be success got: %s", res.Message))
	assert.Equal(t, uint64(1), res.Data.Plan.UserPlanNumber, "the plan number should be 1")
	assert.Equal(t, "Testing 123", res.Data.Plan.Title, "titles don't match")

	repoUser.DeleteUserHard(service.DB, user.ID)
}

// Test failure of order when incorrect order
func TestUnsortedPlan(t *testing.T) {
	service, user, key := setupService()

	defer service.DB.Close()

	req := protoPlan.NewPlanRequest{
		UserID:          user.ID,
		Status:          "active",
		CloseOnComplete: false,
		PlanTemplateID:  "bloody_test",
		Orders: []*protoOrder.NewOrderRequest{
			&protoOrder.NewOrderRequest{
				OrderID:         "da46997b-9d4f-45b2-9d2a-aaa404a213c5",
				KeyID:           key.KeyID,
				OrderType:       "paper",
				OrderTemplateID: "mokie22",
				ParentOrderID:   "4d671984-d7dd-4dce-a20f-23f25d6daf7f",
				MarketName:      "BTC-USDT",
				Side:            "buy",
				InitialCurrencyBalance: 70,
				Triggers: []*protoOrder.TriggerRequest{
					&protoOrder.TriggerRequest{
						TriggerID:         "ab4734f7-5ab7-46eb-9972-ed632ac752f8",
						Code:              "price <= 5500",
						Index:             0,
						Name:              "price_trigger",
						Title:             "Testing update",
						TriggerTemplateID: "testtemplate",
						Actions:           []string{"placeOrder"},
					},
				},
			},
			&protoOrder.NewOrderRequest{
				OrderID:         "4d671984-d7dd-4dce-a20f-23f25d6daf7f",
				KeyID:           key.KeyID,
				OrderType:       "paper",
				OrderTemplateID: "mokie22",
				ParentOrderID:   "00000000-0000-0000-0000-000000000000",
				MarketName:      "BTC-USDT",
				Side:            "buy",
				InitialCurrencyBalance: 70,
				Triggers: []*protoOrder.TriggerRequest{
					&protoOrder.TriggerRequest{
						TriggerID:         "f6aa71b0-8cac-11e8-9eb6-529269fb1450",
						Code:              "price <= 6000",
						Index:             0,
						Name:              "price_trigger",
						Title:             "Testing update",
						TriggerTemplateID: "testtemplate",
						Actions:           []string{"placeOrder"},
					},
				},
			},
		},
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

	//orders := make([]*protoOrder.NewOrderRequest, 0)
	req := protoPlan.NewPlanRequest{
		UserID:          user.ID,
		Status:          "active",
		CloseOnComplete: false,
		PlanTemplateID:  "update_test",
		Orders: []*protoOrder.NewOrderRequest{
			&protoOrder.NewOrderRequest{
				OrderID:         "4d671984-d7dd-4dce-a20f-23f25d6daf7f",
				KeyID:           key.KeyID,
				OrderType:       "paper",
				OrderTemplateID: "mokie",
				ParentOrderID:   "00000000-0000-0000-0000-000000000000",
				MarketName:      "BTC-USDT",
				Side:            "buy",
				InitialCurrencyBalance: 100,
				Triggers: []*protoOrder.TriggerRequest{
					&protoOrder.TriggerRequest{
						TriggerID:         "ab4734f7-5ab7-46eb-9972-ed632ac752f8",
						Code:              "test",
						Index:             0,
						Name:              "test_trigger",
						Title:             "Testing",
						TriggerTemplateID: "testtemplate",
						Actions:           []string{"placeOrder"},
					},
				},
			},
			&protoOrder.NewOrderRequest{
				OrderID:         "966fabac-8cad-11e8-9eb6-529269fb1459",
				KeyID:           key.KeyID,
				OrderType:       "paper",
				OrderTemplateID: "mokie",
				ParentOrderID:   "4d671984-d7dd-4dce-a20f-23f25d6daf7f",
				MarketName:      "EOS-BTC",
				Side:            "buy",
				Triggers: []*protoOrder.TriggerRequest{
					&protoOrder.TriggerRequest{
						TriggerID:         "ab4734f7-5ab7-46eb-9972-ed632ac75290",
						Code:              "test",
						Index:             0,
						Name:              "test_trigger",
						Title:             "Testing",
						TriggerTemplateID: "testtemplate",
						Actions:           []string{"placeOrder"},
					},
				},
			},
			&protoOrder.NewOrderRequest{
				OrderID:         "3c960c68-8cb0-11e8-9eb6-529269fb1459",
				KeyID:           key.KeyID,
				OrderType:       "paper",
				OrderTemplateID: "mokie",
				ParentOrderID:   "4d671984-d7dd-4dce-a20f-23f25d6daf7f",
				MarketName:      "BTC-USDT",
				Side:            "sell",
				Triggers: []*protoOrder.TriggerRequest{
					&protoOrder.TriggerRequest{
						TriggerID:         "ab4734f7-5ab7-46eb-9972-ed632ac75290",
						Code:              "test",
						Index:             0,
						Name:              "test_trigger",
						Title:             "Testing",
						TriggerTemplateID: "testtemplate",
						Actions:           []string{"placeOrder"},
					},
				},
			},
		},
	}
	res := protoPlan.PlanResponse{}
	service.NewPlan(context.Background(), &req, &res)

	assert.Equal(t, "success", res.Status, "return status of inserting plan should be success")

	req2 := protoPlan.UpdatePlanRequest{
		PlanID:          res.Data.Plan.PlanID,
		UserID:          user.ID,
		Status:          "active",
		CloseOnComplete: false,
		PlanTemplateID:  "thing2",
		Orders: []*protoOrder.NewOrderRequest{
			&protoOrder.NewOrderRequest{
				OrderID:         "4d671984-d7dd-4dce-a20f-23f25d6daf7f",
				KeyID:           key.KeyID,
				OrderType:       "paper",
				OrderTemplateID: "mokie22",
				ParentOrderID:   "00000000-0000-0000-0000-000000000000",
				MarketName:      "BTC-USDT",
				Side:            "buy",
				InitialCurrencyBalance: 70,
				Triggers: []*protoOrder.TriggerRequest{
					&protoOrder.TriggerRequest{
						TriggerID:         "ab4734f7-5ab7-46eb-9972-ed632ac752f8",
						Code:              "price <= 5",
						Index:             0,
						Name:              "price_trigger",
						Title:             "Testing update",
						TriggerTemplateID: "testtemplate",
						Actions:           []string{"placeOrder"},
					},
				},
			}, &protoOrder.NewOrderRequest{
				OrderID:         "4d671984-d7dd-4dce-a20f-23f25d6daf76",
				KeyID:           key.KeyID,
				OrderType:       "paper",
				OrderTemplateID: "mokie22",
				ParentOrderID:   "4d671984-d7dd-4dce-a20f-23f25d6daf7f",
				MarketName:      "ADA-BTC",
				Side:            "buy",
				Triggers: []*protoOrder.TriggerRequest{
					&protoOrder.TriggerRequest{
						TriggerID:         "ab4734f7-5ab7-46eb-9972-ed632ac752f8",
						Code:              "price <= 5",
						Index:             0,
						Name:              "price_trigger",
						Title:             "Testing update",
						TriggerTemplateID: "testtemplate",
						Actions:           []string{"placeOrder"},
					},
				},
			},
		},
	}

	res = protoPlan.PlanResponse{}
	service.UpdatePlan(context.Background(), &req2, &res)
	assert.Equal(t, "success", res.Status, res.Message)
	assert.Equal(t, 2, len(res.Data.Plan.Orders), "update should have yielded a single order")
	assert.Equal(t, 70.0, res.Data.Plan.Orders[0].InitialCurrencyBalance, "active currency balance after update order incorrect")
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
	assert.Equal(t, 70.0, res.Data.Plan.Orders[0].InitialCurrencyBalance, "active currency balance after update order incorrect")
	assert.Equal(t, "thing2", res.Data.Plan.PlanTemplateID, "the template ID should have changed")

	repoUser.DeleteUserHard(service.DB, user.ID)
}

// Plan update should fail if the client sends in an update request on executed order.
func TestOrderUpdatePlanFailure(t *testing.T) {
	service, user, key := setupService()

	defer service.DB.Close()

	req := protoPlan.NewPlanRequest{
		UserID:          user.ID,
		Status:          "active",
		CloseOnComplete: false,
		PlanTemplateID:  "thing",
		Orders: []*protoOrder.NewOrderRequest{
			&protoOrder.NewOrderRequest{
				OrderID:         "4d671984-d7dd-4dce-a20f-23f25d6daf7f",
				KeyID:           key.KeyID,
				OrderType:       "paper",
				OrderTemplateID: "mokie22",
				ParentOrderID:   "00000000-0000-0000-0000-000000000000",
				MarketName:      "BTC-USDT",
				Side:            "buy",
				InitialCurrencyBalance: 100,
				Triggers: []*protoOrder.TriggerRequest{
					&protoOrder.TriggerRequest{
						TriggerID:         "ab4734f7-5ab7-46eb-9972-ed632ac752f8",
						Code:              "price <= 5",
						Index:             0,
						Name:              "price_trigger",
						Title:             "Testing update",
						TriggerTemplateID: "testtemplate",
						Actions:           []string{"placeOrder"},
					},
				},
			},
		},
	}
	res := protoPlan.PlanResponse{}
	service.NewPlan(context.Background(), &req, &res)

	assert.Equal(t, "success", res.Status, "return status of inserting plan should be success")
	planID := res.Data.Plan.PlanID

	_, _, err := repoOrder.UpdateOrderStatus(service.DB, "4d671984-d7dd-4dce-a20f-23f25d6daf7f", constPlan.Filled)
	if err != nil {
		assert.Equal(t, "success", "failed", "update order status failed")
	}

	// this should fail because the order above was filled
	req2 := protoPlan.UpdatePlanRequest{
		PlanID:          res.Data.Plan.PlanID,
		UserID:          user.ID,
		Status:          "active",
		CloseOnComplete: false,
		PlanTemplateID:  "thing2",
		Orders: []*protoOrder.NewOrderRequest{
			&protoOrder.NewOrderRequest{
				OrderID:         "4d671984-d7dd-4dce-a20f-23f25d6daf7f",
				KeyID:           key.KeyID,
				OrderType:       "paper",
				OrderTemplateID: "mokie22",
				ParentOrderID:   "00000000-0000-0000-0000-000000000000",
				MarketName:      "ETH-USDT",
				Side:            "buy",
				InitialCurrencyBalance: 1000,
				Triggers: []*protoOrder.TriggerRequest{
					&protoOrder.TriggerRequest{
						TriggerID:         "ab4734f7-5ab7-46eb-9972-ed632ac752f8",
						Code:              "price <= 5",
						Index:             0,
						Name:              "price_trigger",
						Title:             "Testing update",
						TriggerTemplateID: "testtemplate",
						Actions:           []string{"placeOrder"},
					},
				},
			}, &protoOrder.NewOrderRequest{
				OrderID:         "4d671984-d7dd-4dce-a20f-23f25d6daf76",
				KeyID:           key.KeyID,
				OrderType:       "paper",
				OrderTemplateID: "mokie22",
				ParentOrderID:   "4d671984-d7dd-4dce-a20f-23f25d6daf7f",
				MarketName:      "ADA-ETH",
				Side:            "buy",
				Triggers: []*protoOrder.TriggerRequest{
					&protoOrder.TriggerRequest{
						TriggerID:         "ab4734f7-5ab7-46eb-9972-ed632ac752f8",
						Code:              "price <= 5",
						Index:             0,
						Name:              "price_trigger",
						Title:             "Testing update",
						TriggerTemplateID: "testtemplate",
						Actions:           []string{"placeOrder"},
					},
				},
			},
		},
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
	assert.Equal(t, 100.0, res.Data.Plan.Orders[0].InitialCurrencyBalance, "active currency balance after update order incorrect")
	assert.Equal(t, "BTC-USDT", res.Data.Plan.Orders[0].MarketName, "market should not have changed")
	assert.Equal(t, "thing", res.Data.Plan.PlanTemplateID, "the template ID should have changed")

	repoUser.DeleteUserHard(service.DB, user.ID)
}

// This should test an update when a plan order has executed
func TestOrderUpdatePlanFailureBecauseOfExecution(t *testing.T) {
	service, user, key := setupService()

	defer service.DB.Close()

	//orders := make([]*protoOrder.NewOrderRequest, 0)
	req := protoPlan.NewPlanRequest{
		UserID:          user.ID,
		Status:          "active",
		CloseOnComplete: false,
		PlanTemplateID:  "update_test",
		Orders: []*protoOrder.NewOrderRequest{
			&protoOrder.NewOrderRequest{
				OrderID:         "4d671984-d7dd-4dce-a20f-23f25d6daf7f",
				KeyID:           key.KeyID,
				OrderType:       "paper",
				OrderTemplateID: "mokie",
				ParentOrderID:   "00000000-0000-0000-0000-000000000000",
				MarketName:      "BTC-USDT",
				Side:            "buy",
				InitialCurrencyBalance: 100,
				Triggers: []*protoOrder.TriggerRequest{
					&protoOrder.TriggerRequest{
						TriggerID:         "ab4734f7-5ab7-46eb-9972-ed632ac752f8",
						Code:              "test",
						Index:             0,
						Name:              "test_trigger",
						Title:             "Testing",
						TriggerTemplateID: "testtemplate",
						Actions:           []string{"placeOrder"},
					},
				},
			},
			&protoOrder.NewOrderRequest{
				OrderID:         "966fabac-8cad-11e8-9eb6-529269fb1459",
				KeyID:           key.KeyID,
				OrderType:       "paper",
				OrderTemplateID: "mokie",
				ParentOrderID:   "4d671984-d7dd-4dce-a20f-23f25d6daf7f",
				MarketName:      "EOS-BTC",
				Side:            "buy",
				Triggers: []*protoOrder.TriggerRequest{
					&protoOrder.TriggerRequest{
						TriggerID:         "ab4734f7-5ab7-46eb-9972-ed632ac75290",
						Code:              "test",
						Index:             0,
						Name:              "test_trigger",
						Title:             "Testing",
						TriggerTemplateID: "testtemplate",
						Actions:           []string{"placeOrder"},
					},
				},
			},
			&protoOrder.NewOrderRequest{
				OrderID:         "3c960c68-8cb0-11e8-9eb6-529269fb1459",
				KeyID:           key.KeyID,
				OrderType:       "paper",
				OrderTemplateID: "mokie",
				ParentOrderID:   "4d671984-d7dd-4dce-a20f-23f25d6daf7f",
				MarketName:      "BTC-USDT",
				Side:            "sell",
				Triggers: []*protoOrder.TriggerRequest{
					&protoOrder.TriggerRequest{
						TriggerID:         "ab4734f7-5ab7-46eb-9972-ed632ac75290",
						Code:              "test",
						Index:             0,
						Name:              "test_trigger",
						Title:             "Testing",
						TriggerTemplateID: "testtemplate",
						Actions:           []string{"placeOrder"},
					},
				},
			},
		},
	}
	res := protoPlan.PlanResponse{}
	service.NewPlan(context.Background(), &req, &res)

	assert.Equal(t, "success", res.Status, "return status of inserting plan should be success")

	// next suppose the plan order executed already
	// simulate this by kill the plan in the mock
	killRequest := protoEngine.KillRequest{
		PlanID: res.Data.Plan.PlanID,
	}
	service.EngineClient.KillPlan(context.Background(), &killRequest)

	req2 := protoPlan.UpdatePlanRequest{
		PlanID:          res.Data.Plan.PlanID,
		UserID:          user.ID,
		Status:          "active",
		CloseOnComplete: false,
		PlanTemplateID:  "thing2",
		Orders: []*protoOrder.NewOrderRequest{
			&protoOrder.NewOrderRequest{
				OrderID:         "4d671984-d7dd-4dce-a20f-23f25d6daf7f",
				KeyID:           key.KeyID,
				OrderType:       "paper",
				OrderTemplateID: "mokie22",
				ParentOrderID:   "00000000-0000-0000-0000-000000000000",
				MarketName:      "BTC-USDT",
				Side:            "buy",
				InitialCurrencyBalance: 70,
				Triggers: []*protoOrder.TriggerRequest{
					&protoOrder.TriggerRequest{
						TriggerID:         "ab4734f7-5ab7-46eb-9972-ed632ac752f8",
						Code:              "price <= 5",
						Index:             0,
						Name:              "price_trigger",
						Title:             "Testing update",
						TriggerTemplateID: "testtemplate",
						Actions:           []string{"placeOrder"},
					},
				},
			}, &protoOrder.NewOrderRequest{
				OrderID:         "4d671984-d7dd-4dce-a20f-23f25d6daf76",
				KeyID:           key.KeyID,
				OrderType:       "paper",
				OrderTemplateID: "mokie22",
				ParentOrderID:   "4d671984-d7dd-4dce-a20f-23f25d6daf7f",
				MarketName:      "ADA-BTC",
				Side:            "buy",
				Triggers: []*protoOrder.TriggerRequest{
					&protoOrder.TriggerRequest{
						TriggerID:         "ab4734f7-5ab7-46eb-9972-ed632ac752f8",
						Code:              "price <= 5",
						Index:             0,
						Name:              "price_trigger",
						Title:             "Testing update",
						TriggerTemplateID: "testtemplate",
						Actions:           []string{"placeOrder"},
					},
				},
			},
		},
	}

	res = protoPlan.PlanResponse{}
	service.UpdatePlan(context.Background(), &req2, &res)
	assert.Equal(t, "fail", res.Status, "return status of updating an executed plan should be fail")

	repoUser.DeleteUserHard(service.DB, user.ID)
}

// Test delete plan
func TestDeletePlan(t *testing.T) {
	service, user, key := setupService()

	defer service.DB.Close()

	//orders := make([]*protoOrder.NewOrderRequest, 0)
	req := protoPlan.NewPlanRequest{
		UserID:          user.ID,
		Status:          "active",
		CloseOnComplete: false,
		PlanTemplateID:  "update_test",
		Orders: []*protoOrder.NewOrderRequest{
			&protoOrder.NewOrderRequest{
				OrderID:         "4d671984-d7dd-4dce-a20f-23f25d6daf7f",
				KeyID:           key.KeyID,
				OrderType:       "paper",
				OrderTemplateID: "mokie",
				ParentOrderID:   "00000000-0000-0000-0000-000000000000",
				MarketName:      "BTC-USDT",
				Side:            "buy",
				InitialCurrencyBalance: 100,
				Triggers: []*protoOrder.TriggerRequest{
					&protoOrder.TriggerRequest{
						TriggerID:         "ab4734f7-5ab7-46eb-9972-ed632ac752f8",
						Code:              "test",
						Index:             0,
						Name:              "test_trigger",
						Title:             "Testing",
						TriggerTemplateID: "testtemplate",
						Actions:           []string{"placeOrder"},
					},
				},
			},
			&protoOrder.NewOrderRequest{
				OrderID:         "966fabac-8cad-11e8-9eb6-529269fb1459",
				KeyID:           key.KeyID,
				OrderType:       "paper",
				OrderTemplateID: "mokie",
				ParentOrderID:   "4d671984-d7dd-4dce-a20f-23f25d6daf7f",
				MarketName:      "EOS-BTC",
				Side:            "buy",
				Triggers: []*protoOrder.TriggerRequest{
					&protoOrder.TriggerRequest{
						TriggerID:         "ab4734f7-5ab7-46eb-9972-ed632ac75290",
						Code:              "test",
						Index:             0,
						Name:              "test_trigger",
						Title:             "Testing",
						TriggerTemplateID: "testtemplate",
						Actions:           []string{"placeOrder"},
					},
				},
			},
			&protoOrder.NewOrderRequest{
				OrderID:         "3c960c68-8cb0-11e8-9eb6-529269fb1459",
				KeyID:           key.KeyID,
				OrderType:       "paper",
				OrderTemplateID: "mokie",
				ParentOrderID:   "4d671984-d7dd-4dce-a20f-23f25d6daf7f",
				MarketName:      "BTC-USDT",
				Side:            "sell",
				Triggers: []*protoOrder.TriggerRequest{
					&protoOrder.TriggerRequest{
						TriggerID:         "ab4734f7-5ab7-46eb-9972-ed632ac75290",
						Code:              "test",
						Index:             0,
						Name:              "test_trigger",
						Title:             "Testing",
						TriggerTemplateID: "testtemplate",
						Actions:           []string{"placeOrder"},
					},
				},
			},
		},
	}
	res := protoPlan.PlanResponse{}
	service.NewPlan(context.Background(), &req, &res)

	assert.Equal(t, "success", res.Status, "return status of inserting plan should be success")

	req2 := protoPlan.DeletePlanRequest{
		PlanID: res.Data.Plan.PlanID,
		UserID: user.ID,
	}

	res = protoPlan.PlanResponse{}
	service.DeletePlan(context.Background(), &req2, &res)
	assert.Equal(t, "success", res.Status, "return status of delete should be success")

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
	assert.Equal(t, "deleted", res.Data.Plan.Status, "status of plan should be deleted")

	repoUser.DeleteUserHard(service.DB, user.ID)
}
