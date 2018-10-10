package main

import (
	"context"
	"testing"
	"time"

	protoEvt "github.com/asciiu/gomo/common/proto/events"
	protoOrder "github.com/asciiu/gomo/plan-service/proto/order"
	protoPlan "github.com/asciiu/gomo/plan-service/proto/plan"
	repoUser "github.com/asciiu/gomo/user-service/db/sql"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

// This should test an update when a plan order has executed
func TestFilledOrder(t *testing.T) {
	service, user, acc := setupService()

	defer service.DB.Close()

	//orders := make([]*protoOrder.NewOrderRequest, 0)
	req := protoPlan.NewPlanRequest{
		UserID:                  user.ID,
		Status:                  "active",
		UserCurrencySymbol:      "USDT",
		CloseOnComplete:         false,
		PlanTemplateID:          "update_test",
		CommittedCurrencySymbol: "USDT",
		CommittedCurrencyAmount: 100.0,
		Orders: []*protoOrder.NewOrderRequest{
			&protoOrder.NewOrderRequest{
				OrderID:         "4d671984-d7dd-4dce-a20f-23f25d6daf7f",
				AccountID:       acc.AccountID,
				OrderType:       "market",
				OrderTemplateID: "mokie",
				ParentOrderID:   "00000000-0000-0000-0000-000000000000",
				MarketName:      "BTC-USDT",
				Side:            "buy",
				InitialCurrencyBalance: 0,
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
		},
	}
	res := protoPlan.PlanResponse{}
	service.NewPlan(context.Background(), &req, &res)
	assert.Equal(t, "success", res.Status, "return status of inserting plan should be success got: "+res.Message)

	now := string(pq.FormatTimestamp(time.Now().UTC()))

	completedEvent := protoEvt.CompletedOrderEvent{
		UserID:                   user.ID,
		PlanID:                   res.Data.Plan.PlanID,
		OrderID:                  "4d671984-d7dd-4dce-a20f-23f25d6daf7f",
		Exchange:                 acc.Exchange,
		MarketName:               "BTC-USDT",
		Side:                     "buy",
		AccountID:                acc.AccountID,
		TriggerID:                "ab4734f7-5ab7-46eb-9972-ed632ac752f8",
		TriggeredPrice:           1001,
		TriggeredCondition:       "test",
		TriggeredTime:            now,
		CloseOnComplete:          false,
		InitialCurrencyBalance:   100,
		InitialCurrencySymbol:    "USDT",
		InitialCurrencyTraded:    98,
		InitialCurrencyRemainder: 2,
		Status:               "filled",
		ExchangeTime:         now,
		ExchangePrice:        998,
		ExchangeOrderID:      "123",
		FeeCurrencySymbol:    "BNB",
		FeeCurrencyAmount:    10,
		FinalCurrencyBalance: 0.1,
		FinalCurrencySymbol:  "BTC",
	}
	service.HandleCompletedOrder(context.Background(), &completedEvent)

	req3 := protoPlan.GetUserPlanRequest{
		PlanID:     res.Data.Plan.PlanID,
		UserID:     user.ID,
		PlanDepth:  0,
		PlanLength: 10,
	}
	res = protoPlan.PlanResponse{}
	service.GetUserPlan(context.Background(), &req3, &res)

	assert.Equal(t, "success", res.Status, "return status for get plan should be success got: "+res.Message)
	assert.Equal(t, 1, len(res.Data.Plan.Orders), "update should have yielded a single order")
	assert.Equal(t, 100.0, res.Data.Plan.Orders[0].InitialCurrencyBalance, "active currency balance after update order incorrect")
	assert.Equal(t, 10.0, res.Data.Plan.Orders[0].FeeCurrencyAmount, "fee amount off")
	assert.Equal(t, "BNB", res.Data.Plan.Orders[0].FeeCurrencySymbol, "fee symbol")
	assert.Equal(t, "BTC", res.Data.Plan.Orders[0].FinalCurrencySymbol, "final should be BTC")
	assert.Equal(t, 0.1, res.Data.Plan.Orders[0].FinalCurrencyBalance, "final balance off")
	assert.Equal(t, 998.0, res.Data.Plan.Orders[0].ExchangePrice, "exchange price off")
	assert.Equal(t, 100.0, res.Data.Plan.UserCurrencyBalanceAtInit, "assert value at init incorrect")

	repoUser.DeleteUserHard(service.DB, user.ID)
}

func TestFailedOrder(t *testing.T) {
	service, user, acc := setupService()

	defer service.DB.Close()

	//orders := make([]*protoOrder.NewOrderRequest, 0)
	req := protoPlan.NewPlanRequest{
		UserID:                  user.ID,
		Status:                  "active",
		CloseOnComplete:         false,
		PlanTemplateID:          "update_test",
		UserCurrencySymbol:      "USDT",
		CommittedCurrencySymbol: "USDT",
		CommittedCurrencyAmount: 100.0,
		Orders: []*protoOrder.NewOrderRequest{
			&protoOrder.NewOrderRequest{
				OrderID:         "4d671984-d7dd-4dce-a20f-23f25d6daf7f",
				AccountID:       acc.AccountID,
				OrderType:       "market",
				OrderTemplateID: "mokie",
				ParentOrderID:   "00000000-0000-0000-0000-000000000000",
				MarketName:      "BTC-USDT",
				Side:            "buy",
				InitialCurrencyBalance: 0,
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
		},
	}
	res := protoPlan.PlanResponse{}
	service.NewPlan(context.Background(), &req, &res)
	assert.Equal(t, "success", res.Status, "return status of inserting plan should be success got: "+res.Message)

	now := string(pq.FormatTimestamp(time.Now().UTC()))
	completedEvent := protoEvt.CompletedOrderEvent{
		UserID:                   user.ID,
		PlanID:                   res.Data.Plan.PlanID,
		OrderID:                  "4d671984-d7dd-4dce-a20f-23f25d6daf7f",
		Exchange:                 acc.Exchange,
		MarketName:               "BTC-USDT",
		Side:                     "buy",
		AccountID:                acc.AccountID,
		TriggerID:                "ab4734f7-5ab7-46eb-9972-ed632ac752f8",
		TriggeredPrice:           1001,
		TriggeredCondition:       "test",
		TriggeredTime:            now,
		CloseOnComplete:          false,
		InitialCurrencyBalance:   100,
		InitialCurrencySymbol:    "USDT",
		InitialCurrencyTraded:    98,
		InitialCurrencyRemainder: 2,
		Status:  "failed",
		Details: "nope",
	}
	service.HandleCompletedOrder(context.Background(), &completedEvent)

	req3 := protoPlan.GetUserPlanRequest{
		PlanID:     res.Data.Plan.PlanID,
		UserID:     user.ID,
		PlanDepth:  0,
		PlanLength: 10,
	}
	res = protoPlan.PlanResponse{}
	service.GetUserPlan(context.Background(), &req3, &res)

	assert.Equal(t, 0.0, res.Data.Plan.Orders[0].FeeCurrencyAmount, "fee amount off")
	assert.Equal(t, "nope", res.Data.Plan.Orders[0].Errors, "errors incorrect")
	assert.Equal(t, "", res.Data.Plan.Orders[0].FeeCurrencySymbol, "fee symbol")
	assert.Equal(t, "", res.Data.Plan.Orders[0].FinalCurrencySymbol, "final should be BTC")
	assert.Equal(t, 0.0, res.Data.Plan.Orders[0].FinalCurrencyBalance, "final balance off")
	assert.Equal(t, 0.0, res.Data.Plan.Orders[0].ExchangePrice, "exchange price off")

	repoUser.DeleteUserHard(service.DB, user.ID)
}
