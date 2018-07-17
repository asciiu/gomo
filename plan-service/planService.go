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
	"github.com/asciiu/gomo/common/constants/status"
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
}

// private: This is where the order events are published to the rest of the system
// this function should only be callable from within the PlanService. When a plan is
// published the first order of the plan will be emmitted as an ActiveOrderEvent to the
// system.
//
// VERY IMPORTANT: Only send Plans where the first order plan's orders is the next order to active.
// That is to say. DO NOT load a plan where the first order in the orders array has been filled. Fuck
// it, I'm going to implement a check here to ensure this never happens.
func (service *PlanService) publishPlan(ctx context.Context, plan *protoPlan.Plan, isRevision bool) error {
	// the first plan order will always be the active one
	planOrder := plan.Orders[0]

	// only pub plan if the next plan order is active or inactive
	// we do not pub plan orders that have been filled, failed, or aborted
	// reexecuting those plan orders would be very bad!
	if planOrder.Status != status.Active && planOrder.Status != status.Inactive {
		return nil
	}
	triggers := make([]*evt.Trigger, 0)
	for _, t := range planOrder.Triggers {
		trig := evt.Trigger{
			TriggerID: t.TriggerID,
			OrderID:   t.OrderID,
			Name:      t.Name,
			Code:      t.Code,
			Triggered: t.Triggered,
			Actions:   t.Actions,
		}
		triggers = append(triggers, &trig)
	}

	// convert order to order event
	activeOrder := evt.ActiveOrderEvent{
		Exchange:        plan.Exchange,
		OrderID:         planOrder.OrderID,
		PlanID:          plan.PlanID,
		UserID:          plan.UserID,
		BaseBalance:     plan.BaseBalance,
		CurrencyBalance: plan.CurrencyBalance,
		BalancePercent:  planOrder.BalancePercent,
		KeyID:           plan.KeyID,
		Key:             plan.Key,
		Secret:          plan.KeySecret,
		MarketName:      plan.MarketName,
		Side:            planOrder.Side,
		OrderType:       planOrder.OrderType,
		Price:           planOrder.LimitPrice,
		NextOrderID:     planOrder.NextOrderID,
		Revision:        isRevision,
		OrderStatus:     planOrder.Status,
		Triggers:        triggers,
	}

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

// LoadPlanOrder will activate an order (i.e. send a plan order) to the execution engine to process.
func (service *PlanService) LoadPlanOrder(ctx context.Context, plan *protoPlan.Plan, isRevision bool) error {

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

	if err := service.publishPlan(ctx, plan, isRevision); err != nil {
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

		// this is a new plan
		if err := service.publishPlan(ctx, pln, false); err != nil {
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
	case error == sql.ErrNoRows:
		res.Status = response.Nonentity
		res.Message = fmt.Sprintf("planID not found %s", req.PlanID)
	case error != nil:
		res.Status = response.Error
		res.Message = error.Error()
	case pagedPlan.OrdersPage.Total < (req.PageSize * req.Page):
		res.Status = response.Nonentity
		res.Message = "page index out of bounds"
	case error == nil:
		res.Status = response.Success
		res.Data = pagedPlan
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

// We can delete plans that have no filled orders and that are inactive. This becomes an abort plan
// if the plan status is active.
func (service *PlanService) DeletePlan(ctx context.Context, req *protoPlan.DeletePlanRequest, res *protoPlan.PlanResponse) error {
	pln, err := planRepo.FindPlanSummary(service.DB, req.PlanID)
	switch {
	case err == sql.ErrNoRows:
		res.Status = response.Nonentity
		res.Message = fmt.Sprintf("planID not found %s", req.PlanID)
		return nil

	case err != nil:
		res.Status = response.Error
		res.Message = fmt.Sprintf("unexpected error in DeletePlan: %s", err.Error())
		return nil

	default:

		switch {
		case pln.ActiveOrderNumber == 0 && pln.Status != plan.Active:
			// we can safely delete this plan from the system because the plan is not in memory
			// (i.e. not active) and the first order of the plan has not been executed
			err = planRepo.DeletePlan(service.DB, req.PlanID)
			if err != nil {
				res.Status = response.Error
				res.Message = err.Error()
			} else {
				pln.Status = plan.Deleted
				res.Status = response.Success
				res.Data = &protoPlan.PlanData{
					Plan: pln,
				}
			}

		case pln.Status == plan.Active:
			pln.Status = plan.PendingAbort
			_, err = planRepo.UpdatePlanStatus(service.DB, req.PlanID, pln.Status)
			if err != nil {
				res.Status = response.Error
				res.Message = err.Error()
				return nil
			}

			// set the plan order status to aborted we are going to use
			// this status in the execution engine to remove order from memory
			pln.Orders[0].Status = status.Aborted
			// publish this revision to the system so the plan order can be removed from execution
			if err := service.publishPlan(ctx, pln, true); err != nil {
				res.Status = response.Error
				res.Message = fmt.Sprintf("failed to remove active plan order from execution: %s", err.Error())
				return nil
			}

			res.Status = response.Success
			res.Data = &protoPlan.PlanData{
				Plan: pln,
			}

		default:
			// what's this?
		}
	}
	return nil
}

// UpdatePlan returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object.
func (service *PlanService) UpdatePlan(ctx context.Context, req *protoPlan.UpdatePlanRequest, res *protoPlan.PlanResponse) error {
	pln, err := planRepo.FindPlanSummary(service.DB, req.PlanID)
	switch {
	case err == sql.ErrNoRows:
		res.Status = response.Nonentity
		res.Message = fmt.Sprintf("planID not found %s", req.PlanID)
		return nil

	case err != nil:
		res.Status = response.Error
		res.Message = err.Error()
		return nil
	default:

		switch {
		// can't set base balance to 0 if first is buy
		case pln.ActiveOrderNumber == 0 && pln.Orders[0].Side == side.Buy && req.BaseBalance == 0:
			res.Status = response.Fail
			res.Message = fmt.Sprintf("base balance for buy plan cannot be 0")
			return nil

		// can't set currency balance to 0 if first is sell
		case pln.ActiveOrderNumber == 0 && pln.Orders[0].Side == side.Sell && req.CurrencyBalance == 0:
			res.Status = response.Fail
			res.Message = fmt.Sprintf("currency balance for sell plan cannot be 0")
			return nil

		// must specify a valid status if a status was specified
		case req.Status != "" && !plan.ValidateUpdatePlanStatus(req.Status):
			res.Status = response.Fail
			res.Message = fmt.Sprintf("invalid status for update plan")
			return nil

		// when active order is not first order there must be a status param
		case pln.ActiveOrderNumber != 0 && req.Status == "":
			res.Status = response.Fail
			res.Message = fmt.Sprintf("must specify non empty status for update plan")
			return nil

		case pln.ActiveOrderNumber == 0:
			if req.CurrencyBalance >= 0 {
				pln, err = planRepo.UpdatePlanCurrencyBalance(service.DB, req.PlanID, req.CurrencyBalance)
				if err != nil {
					res.Status = response.Error
					res.Message = err.Error()
					return nil
				}
			}
			if req.BaseBalance >= 0 {
				pln, err = planRepo.UpdatePlanBaseBalance(service.DB, req.PlanID, req.BaseBalance)
				if err != nil {
					res.Status = response.Error
					res.Message = err.Error()
					return nil
				}
			}

		default:
		}

		isActive := pln.Status
		if req.Status != "" {
			pln, err = planRepo.UpdatePlanStatus(service.DB, req.PlanID, req.Status)
			if err != nil {
				res.Status = response.Error
				res.Message = err.Error()
				return nil
			}
		}

		if isActive == plan.Active {
			// publish the revised plan to the system
			if err := service.publishPlan(ctx, pln, true); err != nil {
				res.Status = response.Error
				res.Message = err.Error()
				return nil
			}
		}

		if isActive == plan.Inactive && pln.Status == plan.Active {
			if err := service.publishPlan(ctx, pln, false); err != nil {
				res.Status = response.Error
				res.Message = err.Error()
				return nil
			}
		}

		res.Status = response.Success
		res.Data = &protoPlan.PlanData{
			Plan: pln,
		}
	}
	return nil
}
