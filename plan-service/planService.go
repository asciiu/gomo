package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	balances "github.com/asciiu/gomo/balance-service/proto/balance"
	"github.com/asciiu/gomo/common/constants/key"
	orderConstants "github.com/asciiu/gomo/common/constants/order"
	"github.com/asciiu/gomo/common/constants/plan"
	"github.com/asciiu/gomo/common/constants/response"
	"github.com/asciiu/gomo/common/constants/side"
	"github.com/asciiu/gomo/common/constants/status"
	evt "github.com/asciiu/gomo/common/proto/events"
	keys "github.com/asciiu/gomo/key-service/proto/key"
	planRepo "github.com/asciiu/gomo/plan-service/db/sql"
	protoOrder "github.com/asciiu/gomo/plan-service/proto/order"
	protoPlan "github.com/asciiu/gomo/plan-service/proto/plan"
	"github.com/google/uuid"
	"github.com/lib/pq"
	micro "github.com/micro/go-micro"
)

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
	// activeOrder := evt.ActiveOrderEvent{
	// 	//Exchange:        plan.Exchange,
	// 	OrderID: planOrder.OrderID,
	// 	PlanID:  plan.PlanID,
	// 	UserID:  plan.UserID,
	// 	//BaseBalance:     plan.BaseBalance,
	// 	//CurrencyBalance: plan.CurrencyBalance,
	// 	BalancePercent: planOrder.PercentBalance,
	// 	KeyID:          plan.KeyID,
	// 	Key:            plan.Key,
	// 	Secret:         plan.KeySecret,
	// 	//MarketName:      plan.MarketName,
	// 	Side:      planOrder.Side,
	// 	OrderType: planOrder.OrderType,
	// 	Price:     planOrder.LimitPrice,
	// 	//NextOrderID:     planOrder.NextOrderID,
	// 	Revision:    isRevision,
	// 	OrderStatus: planOrder.Status,
	// 	Triggers:    triggers,
	// }

	// if err := service.OrderPub.Publish(context.Background(), &activeOrder); err != nil {
	// 	return fmt.Errorf("publish error: %s -- ActiveOrderEvent %+v", err, &activeOrder)
	// }
	//log.Printf("publish active order -- %+v\n", &activeOrder)
	return nil
}

// private: validateBalance
func (service *PlanService) validateBalance(ctx context.Context, currency string, balance float64, userID string, apikeyID string) (bool, error) {
	balRequest := balances.GetUserBalanceRequest{
		UserID:   userID,
		KeyID:    apikeyID,
		Currency: currency,
	}

	balResponse, err := service.Client.GetUserBalance(ctx, &balRequest)
	if err != nil {
		return false, fmt.Errorf("ecountered error from GetUserBalance: %s", err.Error())
	}

	if balResponse.Data.Balance.Available < balance {
		return false, nil
	}
	return true, nil
}

// LoadPlanOrder will activate an order (i.e. send a plan order) to the execution engine to process.
func (service *PlanService) LoadPlanOrder(ctx context.Context, plan *protoPlan.Plan, isRevision bool) error {

	// planOrder := plan.Orders[0]
	// currencies := strings.Split(plan.MarketName, "-")
	// // default market currency
	// currency := currencies[0]
	// balance := plan.CurrencyBalance
	// if planOrder.Side == side.Buy {
	// 	// buy uses base currency
	// 	currency = currencies[1]
	// 	balance = plan.BaseBalance
	// }

	// if err := service.validateBalance(ctx, currency, balance, plan.UserID, plan.KeyID); err != nil {
	// 	return err
	// }

	// if err := service.publishPlan(ctx, plan, isRevision); err != nil {
	// 	return err
	// }

	return nil
}

func (service *PlanService) fetchKeys(keyIDs []string) ([]*keys.Key, error) {
	request := keys.GetKeysRequest{
		KeyIDs: keyIDs}

	r, _ := service.KeyClient.GetKeys(context.Background(), &request)
	if r.Status != response.Success {
		if r.Status == response.Fail {
			return nil, fmt.Errorf(r.Message)
		}
		if r.Status == response.Error {
			return nil, fmt.Errorf(r.Message)
		}
		if r.Status == response.Nonentity {
			return nil, fmt.Errorf("invalid keys")
		}
	}

	return r.Data.Keys, nil
}

// AddPlans returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object. MarketName example: ADA-BTC where BTC is base.
func (service *PlanService) NewPlan(ctx context.Context, req *protoPlan.NewPlanRequest, res *protoPlan.PlanResponse) error {

	switch {
	case !ValidatePlanInputStatus(req.Status):
		res.Status = response.Fail
		res.Message = "plan status must be active, inactive, or historic"
		return nil
	case !ValidateMinOrder(req.Orders):
		res.Status = response.Fail
		res.Message = "at least one order required for a new plan."
		return nil
	case !ValidateSingleRootNode(req.Orders):
		res.Status = response.Fail
		res.Message = "multiple root nodes found, only one is allowed"
		return nil
	case !ValidateConnectedRoutesFromParent(uuid.Nil.String(), req.Orders):
		res.Status = response.Fail
		res.Message = "an order does not have a valid parent_order_id in your request"
		return nil
	case !ValidateNodeCount(req.Orders):
		res.Status = response.Fail
		res.Message = "you can only post 10 inactive nodes at a time!"
		return nil
	case !ValidateNoneZeroBalance(req.Orders):
		res.Status = response.Fail
		res.Message = "non zero activeCurrencyBalance required for root order!"
		return nil
	case req.Orders[0].OrderType == orderConstants.PaperOrder && !ValidatePaperOrders(req.Orders):
		res.Status = response.Fail
		res.Message = "you cannot add a market/limit order to a plan that will begin with a paper order"
		return nil
	case req.Orders[0].OrderType != orderConstants.PaperOrder && !ValidateNotPaperOrders(req.Orders):
		res.Status = response.Fail
		res.Message = "you cannot add paper orders to a plan that will begin with a market/limit order"
		return nil
	}

	// fetch all order keys
	keyIDs := make([]string, 0, len(req.Orders))
	for _, or := range req.Orders {
		keyIDs = append(keyIDs, or.KeyID)
	}
	kys, err := service.fetchKeys(keyIDs)
	if err != nil {
		if strings.Contains(err.Error(), "invalid input") {
			res.Status = response.Fail
			res.Message = fmt.Sprintf("valid keyID required for each order")
			return nil
		}

		msg := fmt.Sprintf("ecountered error when fetching keys: %s\n", err.Error())
		log.Println(msg)

		res.Status = response.Error
		res.Message = msg
		return nil
	}

	none := uuid.Nil.String()
	planID := uuid.New()
	now := string(pq.FormatTimestamp(time.Now().UTC()))
	newOrders := make([]*protoOrder.Order, 0, len(req.Orders))
	exchange := ""

	for _, or := range req.Orders {
		orderStatus := status.Inactive
		depth := uint32(1)

		if or.MarketName == "" || or.KeyID == "" {
			res.Status = response.Fail
			res.Message = "missing marketName/keyID for order"
			return nil
		}
		if !strings.Contains(or.MarketName, "-") {
			res.Status = response.Fail
			res.Message = "marketName must be currency-base: e.g. ADA-BTC"
			return nil
		}
		if !ValidateOrderType(or.OrderType) {
			res.Status = response.Fail
			res.Message = "market, limit, or paper required for order type"
			return nil
		}
		if !ValidateOrderSide(or.Side) {
			res.Status = response.Fail
			res.Message = "buy or sell required for order side"
			return nil
		}

		// compute the depth for the order
		if or.ParentOrderID != none {
			for _, o := range newOrders {
				if o.OrderID == or.ParentOrderID {
					depth = o.PlanDepth + 1
					break
				}
			}
		}

		if or.ParentOrderID == none && req.Status == plan.Active {
			orderStatus = status.Active
		}

		// assign exchange name from key
		for _, ky := range kys {
			if ky.KeyID == or.KeyID {
				exchange = ky.Exchange

				if ky.Status != key.Verified {
					res.Status = response.Fail
					res.Message = "using an unverified key!"
					return nil

				}
			}
		}

		// collect triggers for this order
		triggers := make([]*protoOrder.Trigger, 0, len(or.Triggers))
		for _, cond := range or.Triggers {
			triggerID := uuid.New()
			trigger := protoOrder.Trigger{
				TriggerID:         triggerID.String(),
				TriggerTemplateID: cond.TriggerTemplateID,
				OrderID:           or.OrderID,
				Index:             cond.Index,
				Title:             cond.Title,
				Name:              cond.Name,
				Code:              cond.Code,
				Actions:           cond.Actions,
				Triggered:         false,
				CreatedOn:         now,
				UpdatedOn:         now,
			}
			triggers = append(triggers, &trigger)
		}

		// market name will be Currency-Base: ADA-BTC
		symbolPair := strings.Split(or.MarketName, "-")
		symbol := symbolPair[1]
		if or.Side == side.Sell {
			symbol = symbolPair[0]
		}

		order := protoOrder.Order{
			KeyID:                 or.KeyID,
			OrderID:               or.OrderID,
			OrderPriority:         or.OrderPriority,
			OrderType:             or.OrderType,
			OrderTemplateID:       or.OrderTemplateID,
			ParentOrderID:         or.ParentOrderID,
			PlanID:                planID.String(),
			PlanDepth:             depth,
			Side:                  or.Side,
			LimitPrice:            or.LimitPrice,
			Exchange:              exchange,
			MarketName:            or.MarketName,
			ActiveCurrencySymbol:  symbol,
			ActiveCurrencyBalance: or.ActiveCurrencyBalance,
			Status:                orderStatus,
			Grupo:                 or.Grupo,
			Triggers:              triggers,
			CreatedOn:             now,
			UpdatedOn:             now,
		}
		newOrders = append(newOrders, &order)
	}

	currencySymbol := newOrders[0].ActiveCurrencySymbol
	currencyBalance := newOrders[0].ActiveCurrencyBalance
	keyID := newOrders[0].KeyID

	if newOrders[0].OrderType != orderConstants.PaperOrder {
		validBalance, err := service.validateBalance(ctx, currencySymbol, currencyBalance, req.UserID, keyID)
		if err != nil {
			msg := fmt.Sprintf("failed to validate the currency balance for %s: %s", currencySymbol, err.Error())
			log.Println(msg)

			res.Status = response.Error
			res.Message = msg
			return nil
		}
		if !validBalance {
			res.Status = response.Fail
			res.Message = fmt.Sprintf("insufficient %s balance, %.8f requested", currencySymbol, currencyBalance)
			return nil
		}
	}

	pln := protoPlan.Plan{
		PlanID:                planID.String(),
		PlanTemplateID:        req.PlanTemplateID,
		UserID:                req.UserID,
		ActiveCurrencySymbol:  newOrders[0].ActiveCurrencySymbol,
		ActiveCurrencyBalance: newOrders[0].ActiveCurrencyBalance,
		Exchange:              newOrders[0].Exchange,
		MarketName:            newOrders[0].MarketName,
		LastExecutedPlanDepth: 0,
		LastExecutedOrderID:   none,
		Orders:                newOrders,
		Status:                req.Status,
		CloseOnComplete:       req.CloseOnComplete,
		CreatedOn:             now,
		UpdatedOn:             now,
	}

	error := planRepo.InsertPlan(service.DB, &pln)
	if error != nil {
		msg := fmt.Sprintf("insert plan failed %s", err.Error())
		log.Println(msg)

		res.Status = response.Error
		res.Message = msg
		return nil
	}

	// activate first plan order if plan is active
	if pln.Status == plan.Active {
		// send key and secret with plan
		//pln.Key = ky.Key
		//pln.KeySecret = ky.Secret

		// this is a new plan
		if err := service.publishPlan(ctx, &pln, false); err != nil {
			// TODO return a warning here
			res.Status = response.Error
			res.Message = "could not publish first order: " + err.Error()
			return nil
		}
	}

	res.Status = response.Success
	res.Data = &protoPlan.PlanData{Plan: &pln}

	return nil
}

// GetUserPlan returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object.
func (service *PlanService) GetUserPlan(ctx context.Context, req *protoPlan.GetUserPlanRequest, res *protoPlan.PlanResponse) error {
	plan, error := planRepo.FindPlanOrders(service.DB, req)

	switch {
	case error == sql.ErrNoRows:
		res.Status = response.Nonentity
		res.Message = fmt.Sprintf("planID not found %s", req.PlanID)
	case error != nil:
		res.Status = response.Error
		res.Message = error.Error()
	// case plan.totalDepth < req.PlanDepth:
	// 	res.Status = response.Nonentity
	// 	res.Message = "plan depth out of bounds, max depth is %s"
	case error == nil:
		res.Status = response.Success
		res.Data = &protoPlan.PlanData{Plan: plan}
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
	pln, err := planRepo.FindPlanWithUnexecutedOrders(service.DB, req.PlanID)
	switch {
	case err == sql.ErrNoRows:
		res.Status = response.Nonentity
		res.Message = fmt.Sprintf("planID not found %s", req.PlanID)
		return nil

	case err != nil:
		res.Status = response.Error
		res.Message = fmt.Sprintf("unexpected error in DeletePlan: %s", err.Error())
		return nil

	case pln.Status == plan.Active:
		pln.Status = plan.PendingAbort
		err = planRepo.UpdatePlanStatus(service.DB, req.PlanID, pln.Status)

		if err != nil {
			res.Status = response.Error
			res.Message = err.Error()
			return nil
		}

		// 		// set the plan order status to aborted we are going to use
		// 		// this status in the execution engine to remove order from memory
		// 		pln.Orders[0].Status = status.Aborted
		// 		// publish this revision to the system so the plan order can be removed from execution
		// 		if err := service.publishPlan(ctx, pln, true); err != nil {
		// 			res.Status = response.Error
		// 			res.Message = fmt.Sprintf("failed to remove active plan order from execution: %s", err.Error())
		// 			return nil
		// 		}

		res.Status = response.Success
		res.Data = &protoPlan.PlanData{
			Plan: pln,
		}

	case pln.LastExecutedPlanDepth == 0 && pln.Status != plan.Active:
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
	}
	return nil
}

// UpdatePlan returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object.
func (service *PlanService) UpdatePlan(ctx context.Context, req *protoPlan.UpdatePlanRequest, res *protoPlan.PlanResponse) error {

	// TODO pause the plan here before going any further

	// load current state of plan
	// the plan should be paused long before UpdatePlan is called
	// this function assumes that the plan is inactive
	pln, err := planRepo.FindPlanWithUnexecutedOrders(service.DB, req.PlanID)
	fmt.Println(len(pln.Orders))
	fmt.Println(pln)

	switch {
	case err == sql.ErrNoRows:
		res.Status = response.Nonentity
		res.Message = fmt.Sprintf("planID not found %s", req.PlanID)
		return nil
	case err != nil:
		msg := fmt.Sprintf("FindPlanWithUnexecutedOrders error: %s", err.Error())
		log.Println(msg)
		res.Status = response.Error
		res.Message = fmt.Sprintf(err.Error())
		return nil
	case len(req.Orders) == 0 && pln.LastExecutedPlanDepth == 0:
		res.Status = response.Fail
		res.Message = "to plan or not to plan. That is the question. A plan must have at least 1 order."
		return nil
	case !ValidatePlanInputStatus(req.Status):
		res.Status = response.Fail
		res.Message = "plan status must be active, inactive"
		return nil
	case !ValidateConnectedRoutesFromParent(pln.LastExecutedOrderID, req.Orders):
		// all orders must be connected using parentOrderID
		res.Status = response.Fail
		res.Message = "this ain't no tree! All orders must be connected using the parentOrderID relationship."
		return nil
	case pln.LastExecutedPlanDepth == 0 && !ValidateNoneZeroBalance(req.Orders):
		// you must commit a balance for the plan in the first order
		res.Status = response.Fail
		res.Message = "the activeCurrencyBalance must be set for the root order"
		return nil
	case pln.LastExecutedPlanDepth == 0 && !ValidateSingleRootNode(req.Orders):
		// you can't start a plan without a root order
		res.Status = response.Fail
		res.Message = "multiple root nodes found, only one is allowed"
		return nil
	case pln.LastExecutedPlanDepth > 0 && !ValidateChildNodes(req.Orders):
		// update on an executed tree can only append child orders
		res.Status = response.Fail
		res.Message = fmt.Sprintf("an order's parentOrderID is %s. This plan already has an executed root order.", uuid.Nil.String())
		return nil
	case !ValidateNodeCount(req.Orders):
		res.Status = response.Fail
		res.Message = "you can only apply 10 inactive orders at a time!"
		return nil
	case pln.LastExecutedPlanDepth == 0 && req.Orders[0].OrderType == orderConstants.PaperOrder && !ValidatePaperOrders(req.Orders):
		// can't mix real orders to a paper plan
		res.Status = response.Fail
		res.Message = "you cannot append market or limit orders to a plan that will begin with a paper order"
		return nil
	case pln.LastExecutedPlanDepth > 0 && pln.Orders[0].OrderType != orderConstants.PaperOrder && !ValidateNotPaperOrders(req.Orders):
		// the executed plan was live - can't add paper orders to this plan
		res.Status = response.Fail
		res.Message = "you cannot add paper orders with a plan that has already executed a live order"
		return nil
	}

	// fetch all order keys
	keyIDs := make([]string, 0, len(req.Orders))
	for _, or := range req.Orders {
		keyIDs = append(keyIDs, or.KeyID)
	}
	kys := make([]*keys.Key, 0)

	if len(req.Orders) > 0 {
		kys, err = service.fetchKeys(keyIDs)
		if err != nil {
			res.Status = response.Error
			res.Message = fmt.Sprintf("ecountered error when fetching keys: %s", err.Error())
			return nil
		}
	}

	none := uuid.Nil.String()
	now := string(pq.FormatTimestamp(time.Now().UTC()))
	newOrders := make([]*protoOrder.Order, 0, len(req.Orders))
	exchange := ""

	for _, or := range req.Orders {
		orderStatus := status.Inactive
		depth := pln.LastExecutedPlanDepth + 1

		if or.MarketName == "" || or.KeyID == "" {
			res.Status = response.Fail
			res.Message = "missing marketName/keyID for order"
			return nil
		}
		if !strings.Contains(or.MarketName, "-") {
			res.Status = response.Fail
			res.Message = "marketName must be currency-base: e.g. ADA-BTC"
			return nil
		}
		if !ValidateOrderType(or.OrderType) {
			res.Status = response.Fail
			res.Message = "market, limit, or paper required for order type"
			return nil
		}
		if !ValidateOrderSide(or.Side) {
			res.Status = response.Fail
			res.Message = "buy or sell required for order side"
			return nil
		}

		// compute the depth for the order
		if or.ParentOrderID != none {
			for _, o := range newOrders {
				if o.OrderID == or.ParentOrderID {
					depth = o.PlanDepth + 1
					break
				}
			}
		}

		if or.ParentOrderID == none && req.Status == plan.Active {
			orderStatus = status.Active
		}

		// assign exchange name from key
		for _, ky := range kys {
			if ky.KeyID == or.KeyID {
				exchange = ky.Exchange

				if ky.Status != key.Verified {
					res.Status = response.Fail
					res.Message = "using an unverified key!"
					return nil

				}
			}
		}

		// collect triggers for this order
		triggers := make([]*protoOrder.Trigger, 0, len(or.Triggers))
		for _, cond := range or.Triggers {
			triggerID := uuid.New()
			trigger := protoOrder.Trigger{
				TriggerID:         triggerID.String(),
				TriggerTemplateID: cond.TriggerTemplateID,
				OrderID:           or.OrderID,
				Index:             cond.Index,
				Title:             cond.Title,
				Name:              cond.Name,
				Code:              cond.Code,
				Actions:           cond.Actions,
				Triggered:         false,
				CreatedOn:         now,
				UpdatedOn:         now,
			}
			triggers = append(triggers, &trigger)
		}

		// market name will be Currency-Base: ADA-BTC
		// the currency context of this order is dictated with the side of the
		// order. If you're buying, you're using the base (BTC). If
		// you're selling, you're using the currency (ADA).
		symbolPair := strings.Split(or.MarketName, "-")
		currencySymbol := symbolPair[1]
		if or.Side == side.Sell {
			currencySymbol = symbolPair[0]
		}

		// validate the balance for non paper orders that are set to use a predefined balance
		if or.OrderType != orderConstants.PaperOrder && or.ActiveCurrencyBalance > 0 {
			validBalance, err := service.validateBalance(ctx, currencySymbol, or.ActiveCurrencyBalance, req.UserID, or.KeyID)
			if err != nil {
				res.Status = response.Error
				res.Message = fmt.Sprintf("failed to validate the currency balance for %s: %s", currencySymbol, err.Error())
				return nil
			}
			if !validBalance {
				res.Status = response.Fail
				res.Message = fmt.Sprintf("insufficient %s balance, %.8f requested in orderID: %s", currencySymbol, or.ActiveCurrencyBalance, or.OrderID)
				return nil
			}
		}

		order := protoOrder.Order{
			KeyID:                 or.KeyID,
			OrderID:               or.OrderID,
			OrderPriority:         or.OrderPriority,
			OrderType:             or.OrderType,
			OrderTemplateID:       or.OrderTemplateID,
			ParentOrderID:         or.ParentOrderID,
			PlanID:                req.PlanID,
			PlanDepth:             depth,
			Side:                  or.Side,
			LimitPrice:            or.LimitPrice,
			Exchange:              exchange,
			MarketName:            or.MarketName,
			ActiveCurrencySymbol:  currencySymbol,
			ActiveCurrencyBalance: or.ActiveCurrencyBalance,
			Grupo:     or.Grupo,
			Status:    orderStatus,
			Triggers:  triggers,
			CreatedOn: now,
			UpdatedOn: now,
		}
		newOrders = append(newOrders, &order)
	}

	pln.Orders = newOrders

	txn, err := service.DB.Begin()
	if err != nil {
		return err
	}

	// Overwrite the entire unexecuted portion of the plan tree with the new orders above.
	// Gather all previous orderIDs for this plan so we can drop them from the DB.
	orderIDs := make([]string, 0)
	for _, o := range pln.Orders {
		if o.Status != status.Filled {
			orderIDs = append(orderIDs, o.OrderID)
		}
	}

	if pln.LastExecutedPlanDepth == 0 && len(newOrders) > 0 {
		pln.ActiveCurrencySymbol = newOrders[0].ActiveCurrencySymbol
		pln.ActiveCurrencyBalance = newOrders[0].ActiveCurrencyBalance
		pln.Exchange = newOrders[0].Exchange
		pln.MarketName = newOrders[0].MarketName
		if err := planRepo.UpdatePlanContextTxn(txn, ctx, pln.PlanID, pln.ActiveCurrencySymbol, pln.Exchange, pln.MarketName, pln.ActiveCurrencyBalance); err != nil {
			txn.Rollback()
			res.Status = response.Error
			res.Message = "error encountered while updating the plan context: " + err.Error()
			return nil

		}
	}
	if pln.Status != req.Status {
		pln.Status = req.Status
		if err := planRepo.UpdatePlanStatusTxn(txn, ctx, pln.PlanID, pln.Status); err != nil {
			txn.Rollback()
			res.Status = response.Error
			res.Message = "error encountered while updating the plan status: " + err.Error()
			return nil

		}
	}
	if pln.CloseOnComplete != req.CloseOnComplete {
		pln.CloseOnComplete = req.CloseOnComplete
		if err := planRepo.UpdatePlanCloseOnCompleteTxn(txn, ctx, pln.PlanID, pln.CloseOnComplete); err != nil {
			txn.Rollback()
			res.Status = response.Error
			res.Message = "error encountered while updating the plan close on complete option: " + err.Error()
			return nil

		}
	}
	if pln.PlanTemplateID != req.PlanTemplateID {
		pln.PlanTemplateID = req.PlanTemplateID
		if err := planRepo.UpdatePlanTemplateTxn(txn, ctx, pln.PlanID, pln.PlanTemplateID); err != nil {
			txn.Rollback()
			res.Status = response.Error
			res.Message = "error encountered while updating the plan template: " + err.Error()
			return nil

		}
	}

	// keep the update timestamp of the plan in sync with the orders
	// no particular reason, but it could be useful in debugging
	if err := planRepo.UpdatePlanTimestampTxn(txn, ctx, pln.PlanID, now); err != nil {
		txn.Rollback()
		res.Status = response.Error
		res.Message = "error encountered while updating the plan timestamp: " + err.Error()
		return nil
	}

	// drop current orders from the plan
	if err := planRepo.DeleteOrders(txn, ctx, orderIDs); err != nil {
		txn.Rollback()
		res.Status = response.Error
		res.Message = "error while deleting the previous orders: " + err.Error()
		return nil
	}

	if len(newOrders) > 0 {
		// insert new orders for this plan
		if err := planRepo.InsertOrders(txn, newOrders); err != nil {
			txn.Rollback()
			res.Status = response.Error
			res.Message = "insert orders error: " + err.Error()
			return nil
		}

		newTriggers := make([]*protoOrder.Trigger, 0, len(newOrders))
		for _, o := range newOrders {
			for _, t := range o.Triggers {
				newTriggers = append(newTriggers, t)
			}
		}

		if err := planRepo.InsertTriggers(txn, newTriggers); err != nil {
			txn.Rollback()
			return errors.New("bulk triggers failed: " + err.Error())
		}
	}

	txn.Commit()

	// activate first plan order if plan is active
	if pln.Status == plan.Active {
		// send key and secret with plan
		//pln.Key = ky.Key
		//pln.KeySecret = ky.Secret

		// this is a new plan
		if err := service.publishPlan(ctx, pln, false); err != nil {
			// TODO return a warning here
			res.Status = response.Error
			res.Message = "could not publish first order: " + err.Error()
			return nil
		}
	}

	res.Status = response.Success
	res.Data = &protoPlan.PlanData{Plan: pln}

	return nil
}
