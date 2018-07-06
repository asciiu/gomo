package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	balances "github.com/asciiu/gomo/balance-service/proto/balance"
	"github.com/asciiu/gomo/common/constants/key"
	"github.com/asciiu/gomo/common/constants/plan"
	"github.com/asciiu/gomo/common/constants/response"
	evt "github.com/asciiu/gomo/common/proto/events"
	keys "github.com/asciiu/gomo/key-service/proto/key"
	planRepo "github.com/asciiu/gomo/plan-service/db/sql"
	protoPlan "github.com/asciiu/gomo/plan-service/proto/plan"
	micro "github.com/micro/go-micro"
)

// MinBalance needed to submit order
const MinBalance = 0.00001000

// PlanService ...
type PlanService struct {
	DB        *sql.DB
	Client    balances.BalanceServiceClient
	KeyClient keys.KeyServiceClient
	OrderPub  micro.Publisher
	//NewBuy    micro.Publisher
	//NewSell   micro.Publisher
}

// private: This is where the order events are published to the rest of the system
// this function should only be callable from within the PlanService. When a plan is
// published the first order of the plan will be emmitted as an ActiveOrderEvent to the
// system.
func (service *PlanService) publishPlan(ctx context.Context, plan *protoPlan.Plan) error {
	// TODO compute quantity as
	// Buy-limit: plan.baseBalance / planOrder.Price
	// Buy-market: plan.baseBalance / trigger.Price (can only determine this at trigger time)
	// Sell-limit: currencyBalance
	// Sell-market: currencyBalance
	nextOrders := make([]string, 0)
	for _, order := range plan.Orders {
		nextOrders = append(nextOrders, order.OrderID)
	}
	planOrder := plan.Orders[0]

	// convert order to order event
	activeOrder := evt.ActivateOrderEvent{
		Exchange:        plan.Exchange,
		OrderID:         planOrder.OrderID,
		PlanID:          plan.PlanID,
		UserID:          plan.UserID,
		BaseBalance:     plan.BaseBalance,
		CurrencyBalance: plan.CurrencyBalance,
		KeyID:           plan.KeyID,
		Key:             plan.Key,
		Secret:          plan.Secret,
		MarketName:      plan.MarketName,
		Side:            planOrder.Side,
		OrderType:       planOrder.OrderType,
		Price:           planOrder.Price,
		Conditions:      planOrder.Conditions,
		NextOrders:      nextOrders,
	}

	// //if err := publisher.Publish(context.Background(), &orderEvent); err != nil {
	if err := service.OrderPub.Publish(context.Background(), &activeOrder); err != nil {
		return fmt.Errorf("publish error: %s -- ActiveOrderEvent %+v", err, &activeOrder)
	}
	log.Printf("publish active order -- %+v\n", &activeOrder)
	return nil
}

// private: validateBalance
func (service *PlanService) validateBalance(ctx context.Context, currency string, balance float64, userID string, apikeyID string) error {
	balRequest := balances.GetUserBalanceRequest{
		UserID:   userID,
		KeyID:    apikeyID,
		Currency: currency,
	}

	balResponse, err := service.Client.GetUserBalance(ctx, &balRequest)
	if err != nil {
		return fmt.Errorf("ecountered error from GetUserBalance: %s", err.Error())
	}

	if balResponse.Data.Balance.Available < balance {
		return fmt.Errorf("insufficient %s balance, %.8f requested", currency, balance)
	}
	return nil
}

// LoadBuyPlan should not be invoked by the client. This function was designed to load an
// order after an order was filled.
func (service *PlanService) LoadPlan(ctx context.Context, order *protoPlan.Plan) error {

	// currencies := strings.Split(order.MarketName, "-")
	// // default market currency
	// currency := currencies[0]
	// if order.Side == side.Buy {
	// 	// buy uses base currency
	// 	currency = currencies[1]
	// }

	// if err := service.validateBalance(ctx, currency, order.UserID, order.KeyID); err != nil {
	// 	return err
	// }

	// if err := service.publishPlan(ctx, order); err != nil {
	// 	return err
	// }

	return nil
}

func (service *PlanService) fetchKey(keyID, userID string) (*keys.Key, error) {
	getRequest := keys.GetUserKeyRequest{
		KeyID:  keyID,
		UserID: userID,
	}

	r, _ := service.KeyClient.GetUserKey(context.Background(), &getRequest)
	if r.Status != response.Success {
		if r.Status == response.Fail {
			return nil, fmt.Errorf(r.Message)
		}
		if r.Status == response.Error {
			return nil, fmt.Errorf(r.Message)
		}
		if r.Status == response.Nonentity {
			return nil, fmt.Errorf("invalid key")
		}
	}

	// key must be verified status
	if r.Data.Key.Status != key.Verified {
		return nil, fmt.Errorf("key must be verified")
	}
	return r.Data.Key, nil
}

// AddPlans returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object.
func (service *PlanService) AddPlan(ctx context.Context, req *protoPlan.PlanRequest, res *protoPlan.PlanResponse) error {
	currencies := strings.Split(req.MarketName, "-")
	var currency string
	var balance float64

	switch {
	case req.BaseBalance > 0:
		// base currency will be second
		currency = currencies[1]
		balance = req.BaseBalance
	case req.CurrencyBalance > 0:
		currency = currencies[0]
		balance = req.CurrencyBalance
	default:
		res.Status = response.Fail
		res.Message = "baseBalance and currencyBalance are 0"
		return nil
	}

	if err := service.validateBalance(ctx, currency, balance, req.UserID, req.KeyID); err != nil {
		res.Status = response.Fail
		res.Message = err.Error()
		return nil
	}

	// validate plan key
	ky, err := service.fetchKey(req.KeyID, req.UserID)
	if err != nil {
		res.Status = response.Fail
		res.Message = err.Error()
		return nil
	}

	// insert the exchange name from the key
	req.Exchange = ky.Exchange

	pln, error := planRepo.InsertPlan(service.DB, req)
	if error != nil {
		res.Status = response.Error
		res.Message = "AddPlan error: " + error.Error()
		return nil
	}

	// activate first plan order if plan is active
	if pln.Status == plan.Active {
		// send key and secret with plan
		pln.Key = ky.Key
		pln.Secret = ky.Secret

		if err := service.publishPlan(ctx, pln); err != nil {
			// TODO return a warning here
			res.Status = response.Error
			res.Message = "could not publish first order: " + err.Error()
			return nil
		}
	}

	res.Status = response.Success
	res.Data = &protoPlan.PlanData{
		Plan: pln,
	}
	return nil
}

// GetUserPlan returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object.
func (service *PlanService) GetUserPlan(ctx context.Context, req *protoPlan.GetUserPlanRequest, res *protoPlan.PlanResponse) error {
	// order, error := orderRepo.FindPlanByID(service.DB, req)

	// switch {
	// case error == nil:
	// 	res.Status = response.Success
	// 	res.Data = &plans.UserPlanData{
	// 		Plan: order,
	// 	}
	// default:
	// 	res.Status = response.Error
	// 	res.Message = error.Error()
	// }
	return nil
}

// GetUserPlans returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object.
func (service *PlanService) GetUserPlans(ctx context.Context, req *protoPlan.GetUserPlansRequest, res *protoPlan.PlansResponse) error {
	// ordrs, error := orderRepo.FindPlansByUserID(service.DB, req)

	// switch {
	// case error == nil:
	// 	res.Status = response.Success
	// 	res.Data = &plans.UserPlansData{
	// 		Plans: ordrs,
	// 	}
	// default:
	// 	res.Status = response.Error
	// 	res.Message = error.Error()
	// }
	return nil
}

// RemovePlan returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object.
func (service *PlanService) DeletePlan(ctx context.Context, req *protoPlan.DeletePlanRequest, res *protoPlan.PlanResponse) error {
	// error := orderRepo.DeletePlan(service.DB, req.PlanID)
	// switch {
	// case error == nil:
	// 	res.Status = response.Success
	// default:
	// 	res.Status = response.Error
	// 	res.Message = error.Error()
	// }
	return nil
}

// UpdatePlan returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object.
func (service *PlanService) UpdatePlan(ctx context.Context, req *protoPlan.PlanRequest, res *protoPlan.PlanResponse) error {
	// order, error := orderRepo.UpdatePlan(service.DB, req)
	// switch {
	// case error == nil:
	// 	res.Status = response.Success
	// 	res.Data = &plans.UserPlanData{
	// 		Plan: order,
	// 	}
	// default:
	// 	res.Status = response.Error
	// 	res.Message = error.Error()
	// }
	return nil
}
