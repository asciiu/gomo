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
	"github.com/asciiu/gomo/common/constants/side"
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

	// the first plan order will always be the active one
	planOrder := plan.Orders[0]

	// convert order to order event
	activeOrder := evt.ActiveOrderEvent{
		Exchange:        plan.Exchange,
		OrderID:         planOrder.OrderID,
		PlanID:          plan.PlanID,
		UserID:          plan.UserID,
		BaseBalance:     plan.BaseBalance,
		BasePercent:     planOrder.BasePercent,
		CurrencyBalance: plan.CurrencyBalance,
		CurrencyPercent: planOrder.CurrencyPercent,
		KeyID:           plan.KeyID,
		Key:             plan.Key,
		Secret:          plan.KeySecret,
		MarketName:      plan.MarketName,
		Side:            planOrder.Side,
		OrderType:       planOrder.OrderType,
		Price:           planOrder.Price,
		Conditions:      planOrder.Conditions,
		NextOrderID:     planOrder.NextOrderID,
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
func (service *PlanService) LoadPlanOrder(ctx context.Context, plan *protoPlan.Plan) error {

	planOrder := plan.Orders[0]
	currencies := strings.Split(plan.MarketName, "-")
	// default market currency
	currency := currencies[0]
	balance := plan.CurrencyBalance
	if planOrder.Side == side.Buy {
		// buy uses base currency
		currency = currencies[1]
		balance = plan.BaseBalance
	}

	if err := service.validateBalance(ctx, currency, balance, plan.UserID, plan.KeyID); err != nil {
		return err
	}

	if err := service.publishPlan(ctx, plan); err != nil {
		return err
	}

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
		pln.KeySecret = ky.Secret

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
func (service *PlanService) GetUserPlan(ctx context.Context, req *protoPlan.GetUserPlanRequest, res *protoPlan.PlanWithPagedOrdersResponse) error {
	pagedPlan, error := planRepo.FindPlanWithPagedOrders(service.DB, req)

	switch {
	case error == nil:
		res.Status = response.Success
		res.Data = pagedPlan
	case pagedPlan.OrdersPage.Total < (req.PageSize * req.Page):
		res.Status = response.Nonentity
		res.Message = "page index out of bounds"
	default:
		res.Status = response.Error
		res.Message = error.Error()
	}
	return nil
}

// GetUserPlans returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object.
func (service *PlanService) GetUserPlans(ctx context.Context, req *protoPlan.GetUserPlansRequest, res *protoPlan.PlansPageResponse) error {

	var page *protoPlan.PlansPage
	var err error

	switch {
	case req.MarketName == "" && req.Exchange != "":
		// search by userID, exchange, status when no marketName
		page, err = planRepo.FindUserExchangePlansWithStatus(service.DB, req.UserID, req.Status, req.Exchange, req.Page, req.PageSize)
	case req.MarketName != "" && req.Exchange != "":
		// search by userID, exchange, marketName, status
		page, err = planRepo.FindUserExchangePlansWithStatus(service.DB, req.UserID, req.Status, req.Exchange, req.Page, req.PageSize)
	default:
		// search by userID and status
		page, err = planRepo.FindUserPlansWithStatus(service.DB, req.UserID, req.Status, req.Page, req.PageSize)
	}

	switch {
	case err == nil:
		res.Status = response.Success
		res.Data = page
	default:
		res.Status = response.Error
		res.Message = err.Error()
	}

	return nil
}

// RemovePlan returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object.
//func (service *PlanService) DeletePlan(ctx context.Context, req *protoPlan.DeletePlanRequest, res *protoPlan.PlanResponse) error {
// error := orderRepo.DeletePlan(service.DB, req.PlanID)
// switch {
// case error == nil:
// 	res.Status = response.Success
// default:
// 	res.Status = response.Error
// 	res.Message = error.Error()
// }
//	return nil
//}

// UpdatePlan returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object.
//func (service *PlanService) UpdatePlan(ctx context.Context, req *protoPlan.PlanRequest, res *protoPlan.PlanResponse) error {
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
//	return nil
//}
