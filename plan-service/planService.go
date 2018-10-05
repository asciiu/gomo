package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	constAccount "github.com/asciiu/gomo/account-service/constants"
	protoActivity "github.com/asciiu/gomo/activity-bulletin/proto"
	protoAnalytics "github.com/asciiu/gomo/analytics-service/proto/analytics"
	constRes "github.com/asciiu/gomo/common/constants/response"
	protoEvt "github.com/asciiu/gomo/common/proto/events"
	commonUtil "github.com/asciiu/gomo/common/util"
	constPlan "github.com/asciiu/gomo/plan-service/constants"
	micro "github.com/micro/go-micro"

	protoAccount "github.com/asciiu/gomo/account-service/proto/account"
	protoBalance "github.com/asciiu/gomo/account-service/proto/balance"
	protoEngine "github.com/asciiu/gomo/execution-engine/proto/engine"
	repoPlan "github.com/asciiu/gomo/plan-service/db/sql"
	protoOrder "github.com/asciiu/gomo/plan-service/proto/order"
	protoPlan "github.com/asciiu/gomo/plan-service/proto/plan"
	protoTrade "github.com/asciiu/gomo/plan-service/proto/trade"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

const precision = 8

// PlanService ...
type PlanService struct {
	DB              *sql.DB
	AccountClient   protoAccount.AccountServiceClient
	EngineClient    protoEngine.ExecutionEngineClient
	AnalyticsClient protoAnalytics.AnalyticsServiceClient
	NotifyPub       micro.Publisher
}

// private: This is where the order events are published to the rest of the system
// this function should only be callable from within the PlanService. When a plan is
// published the first order of the plan will be emmitted as an ActiveOrderEvent to the
// system.
//
// VERY IMPORTANT: In theory, this should never be used to republish filled orders.
func (service *PlanService) publishPlan(ctx context.Context, plan *protoPlan.Plan) error {
	newOrders := make([]*protoEngine.Order, 0)
	activeDepth := plan.LastExecutedPlanDepth + 1
	accountIDs := make([]string, 0)

	for _, order := range plan.Orders {

		// only add orders that are at the active plan depth
		if order.PlanDepth != activeDepth {
			continue
		}

		triggers := make([]*protoEngine.Trigger, 0)
		for _, t := range order.Triggers {
			trig := protoEngine.Trigger{
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
		newOrderEvt := protoEngine.Order{
			OrderID:     order.OrderID,
			Exchange:    order.Exchange,
			MarketName:  order.MarketName,
			Side:        order.Side,
			LimitPrice:  order.LimitPrice,
			OrderType:   order.OrderType,
			OrderStatus: order.Status,
			AccountID:   order.AccountID,
			AccountType: order.AccountType,
			KeyPublic:   order.KeyPublic,
			KeySecret:   order.KeySecret,
			Triggers:    triggers,
		}

		// gather unique account ids here
		found := false
		for _, a := range accountIDs {
			if a == order.AccountID {
				found = true
			}
		}
		if !found {
			accountIDs = append(accountIDs, order.AccountID)
		}

		newOrders = append(newOrders, &newOrderEvt)
	}

	// no orders to publish
	if len(newOrders) == 0 {
		return nil
	}

	newPlan := protoEngine.NewPlanRequest{
		PlanID:                plan.PlanID,
		UserID:                plan.UserID,
		ActiveCurrencyBalance: plan.ActiveCurrencyBalance,
		ActiveCurrencySymbol:  plan.ActiveCurrencySymbol,
		CloseOnComplete:       plan.CloseOnComplete,
		Orders:                newOrders,
	}

	response, err := service.EngineClient.AddPlan(ctx, &newPlan)
	if err != nil {
		return fmt.Errorf("could not publish the plan %s %s", newPlan.PlanID, err.Error())
	}

	log.Printf("published plan -- %s\n", response.Data.Plans[0].PlanID)

	// update the order status
	for _, o := range newOrders {
		if _, _, err := repoPlan.UpdateOrderStatusAndBalance(service.DB, o.OrderID, constPlan.Active, plan.ActiveCurrencyBalance); err != nil {
			log.Println("could not update order status to active -- ", err.Error())
		}
	}

	return nil
}

// private: validateBalance
// returns true, nil when the balance can be validated
func (service *PlanService) validateAvailableBalance(ctx context.Context, currency string, amount float64, userID string, accountID string) (bool, error) {
	validateRequest := protoBalance.ValidateBalanceRequest{
		UserID:          userID,
		AccountID:       accountID,
		CurrencySymbol:  currency,
		RequestedAmount: amount,
	}

	valResponse, _ := service.AccountClient.ValidateAvailableBalance(ctx, &validateRequest)
	if valResponse.Status != constRes.Success {
		return false, fmt.Errorf("ecountered error from ValidateAccountBalance: %s", valResponse.Message)
	}

	return valResponse.Data, nil
}

// is the available on exchange less than the active currency balance
func (service *PlanService) validateLockedBalance(ctx context.Context, currency string, amount float64, userID string, accountID string) (bool, error) {
	validateRequest := protoBalance.ValidateBalanceRequest{
		UserID:          userID,
		AccountID:       accountID,
		CurrencySymbol:  currency,
		RequestedAmount: amount,
	}

	valResponse, _ := service.AccountClient.ValidateLockedBalance(ctx, &validateRequest)
	if valResponse.Status != constRes.Success {
		return false, fmt.Errorf("ecountered error from ValidateAccountBalance: %s", valResponse.Message)
	}

	return valResponse.Data, nil
}

// ContinuePlan will activate an order (i.e. send a plan order) to the execution engine to process.
func (service *PlanService) ContinuePlan(ctx context.Context, plan *protoPlan.Plan) error {

	currency := plan.ActiveCurrencySymbol
	balance := plan.ActiveCurrencyBalance
	// assume for now that all orders in the plan use the same account ID. This may not be true in the future
	accountID := plan.Orders[0].AccountID

	valid, err := service.validateLockedBalance(ctx, currency, balance, plan.UserID, accountID)
	if err != nil {
		return err
	}
	if !valid {
		return errors.New(fmt.Sprintf("insufficient %s balance for requested amount %.8f", currency, balance))
	}

	if err := service.publishPlan(ctx, plan); err != nil {
		return err
	}

	return nil
}

func (service *PlanService) fetchKeys(accountIDs []string) ([]*protoAccount.AccountKey, error) {
	request := protoAccount.GetAccountKeysRequest{
		AccountIDs: accountIDs}

	r, _ := service.AccountClient.GetAccountKeys(context.Background(), &request)
	if r.Status != constRes.Success {
		if r.Status == constRes.Fail {
			return nil, fmt.Errorf(r.Message)
		}
		if r.Status == constRes.Error {
			return nil, fmt.Errorf(r.Message)
		}
		if r.Status == constRes.Nonentity {
			return nil, fmt.Errorf("invalid accounts")
		}
	}

	return r.Data.Keys, nil
}

func (service *PlanService) HandleAccountDeleted(ctx context.Context, evt *protoEvt.DeletedAccountEvent) error {
	// Flowy's Gospel: 20180912
	// Our first version will be using the same account per plan.
	// However, the intended design allows for multiple accounts.
	// The account is associated at the order.
	// since you currently cannot use different accounts within a plan
	// close the plan if the account was deleted. Close only if the plan status is active or inactive

	plans, err := repoPlan.FindAccountPlans(service.DB, evt.AccountID)
	if err != nil {
		log.Println("HandleAccountDeleted error on FindAccountOrders: ", err.Error())
	}

	// close these plans
	for _, plan := range plans {
		if plan.Status == constPlan.Active || plan.Status == constPlan.Inactive {
			// close plans that are active or inactive only
			if repoPlan.UpdatePlanStatus(service.DB, plan.PlanID, constPlan.Closed) != nil {
				log.Printf("could not close plan -- %s\n", plan.PlanID)
			}
		}
	}

	return nil
}

// Handle a completed order event
func (service *PlanService) HandleCompletedOrder(ctx context.Context, completedOrderEvent *protoEvt.CompletedOrderEvent) error {
	log.Printf("completed event -- %+v\n", completedOrderEvent)

	now := string(pq.FormatTimestamp(time.Now().UTC()))

	notification := protoActivity.Activity{
		UserID:      completedOrderEvent.UserID,
		ObjectID:    completedOrderEvent.PlanID,
		Type:        "plan",
		Timestamp:   now,
		Title:       completedOrderEvent.MarketName,
		Subtitle:    completedOrderEvent.Side,
		Description: completedOrderEvent.Details,
		Details:     fmt.Sprintf("{orderID: %s}", completedOrderEvent.OrderID),
	}

	log.Printf("%+v\n", notification)

	// notify the user of completed order
	if service.NotifyPub != nil {
		if err := service.NotifyPub.Publish(context.Background(), &notification); err != nil {
			log.Println("could not publish notification: ", err)
		}
	}

	planID, depth, err := repoPlan.UpdateOrderStatus(service.DB, completedOrderEvent.OrderID, completedOrderEvent.Status)
	if err != nil {
		log.Println("could not update order status -- ", err.Error())
		return nil
	}

	if completedOrderEvent.Status == constPlan.Filled {
		if err := repoPlan.InsertTradeResult(service.DB, &protoTrade.Trade{
			TradeID:                  uuid.New().String(),
			OrderID:                  completedOrderEvent.OrderID,
			InitialCurrencySymbol:    completedOrderEvent.InitialCurrencySymbol,
			InitialCurrencyBalance:   completedOrderEvent.InitialCurrencyBalance,
			InitialCurrencyTraded:    completedOrderEvent.InitialCurrencyTraded,
			InitialCurrencyRemainder: completedOrderEvent.InitialCurrencyRemainder,
			InitialCurrencyPrice:     completedOrderEvent.ExchangePrice,
			FinalCurrencySymbol:      completedOrderEvent.FinalCurrencySymbol,
			FinalCurrencyBalance:     completedOrderEvent.FinalCurrencyBalance,
			FeeCurrencySymbol:        completedOrderEvent.FeeCurrencySymbol,
			FeeCurrencyAmount:        completedOrderEvent.FeeCurrencyAmount,
			ExchangeTime:             completedOrderEvent.ExchangeTime,
			Side:                     completedOrderEvent.Side,
			CreatedOn:                now,
			UpdatedOn:                now,
		}); err != nil {
			log.Println("could not log trade results -- ", err.Error())
		}

		// TODO error check these in case the account service is unreachable
		// remove lock on initial balance
		changeReq := protoBalance.ChangeBalanceRequest{
			UserID:         completedOrderEvent.UserID,
			AccountID:      completedOrderEvent.AccountID,
			CurrencySymbol: completedOrderEvent.InitialCurrencySymbol,
			Amount:         -completedOrderEvent.InitialCurrencyBalance,
		}
		service.AccountClient.ChangeLockedBalance(ctx, &changeReq)

		// add remainder to initial balance
		if completedOrderEvent.InitialCurrencyRemainder > 0 {
			changeReq.Amount = completedOrderEvent.InitialCurrencyRemainder
			service.AccountClient.ChangeAvailableBalance(ctx, &changeReq)
		}

		// add the final currency balance to the locked balance
		changeReq.CurrencySymbol = completedOrderEvent.FinalCurrencySymbol
		changeReq.Amount = completedOrderEvent.FinalCurrencyBalance
		service.AccountClient.ChangeLockedBalance(ctx, &changeReq)

		if err := repoPlan.UpdateTriggerResults(service.DB,
			completedOrderEvent.TriggerID,
			completedOrderEvent.TriggeredPrice,
			completedOrderEvent.TriggeredCondition,
			completedOrderEvent.ExchangeTime); err != nil {
			log.Println("completed order error trying to update the trigger -- ", err.Error())
			return nil
		}

		if err := repoPlan.UpdateOrderResults(service.DB,
			completedOrderEvent.OrderID,
			completedOrderEvent.ExchangeOrderID,
			completedOrderEvent.ExchangeTime,
			completedOrderEvent.InitialCurrencyTraded,
			completedOrderEvent.InitialCurrencyRemainder,
			completedOrderEvent.FinalCurrencyBalance,
			completedOrderEvent.FeeCurrencyAmount,
			completedOrderEvent.ExchangePrice,
			completedOrderEvent.FinalCurrencySymbol,
			completedOrderEvent.FeeCurrencySymbol); err != nil {
			log.Println("completed order error trying to update the order -- ", err.Error())
			return nil
		}

		initTime, err := repoPlan.UpdatePlanContext(service.DB,
			planID,
			completedOrderEvent.OrderID,
			completedOrderEvent.Exchange,
			completedOrderEvent.FinalCurrencySymbol,
			completedOrderEvent.FinalCurrencyBalance,
			depth)
		if err != nil {
			log.Println("completed order error trying to update the plan context -- ", err.Error())
			return nil
		}

		// first order and init time was not set
		if depth == 1 && initTime == "" {
			userCurrencySymbol, err := repoPlan.FindPlanUserCurrencySymbol(service.DB, completedOrderEvent.PlanID)
			var initialAssetValue float64
			if err == nil && userCurrencySymbol != "" {
				convertReq := protoAnalytics.ConversionRequest{
					Exchange:    completedOrderEvent.Exchange,
					From:        completedOrderEvent.InitialCurrencySymbol,
					FromAmount:  completedOrderEvent.InitialCurrencyBalance,
					To:          userCurrencySymbol,
					AtTimestamp: now,
				}
				convertRes, _ := service.AnalyticsClient.ConvertCurrency(context.Background(), &convertReq)
				initialAssetValue = convertRes.Data.ConvertedAmount
			}

			if err := repoPlan.UpdatePlanInitTimestamp(service.DB, planID, now, initialAssetValue); err != nil {
				log.Println("could not update plan init time -- ", err.Error())
				return nil
			}
		}

		// load the child orders of this completed order
		nextPlanOrders, err := repoPlan.FindChildOrders(service.DB, planID, completedOrderEvent.OrderID)

		switch {
		case err == sql.ErrNoRows:
			// close status of plan when no more orders and CloseOnComplete is true
			if completedOrderEvent.CloseOnComplete {
				if repoPlan.UpdatePlanStatus(service.DB, planID, constPlan.Closed) != nil {
					log.Printf("could not close plan -- %s\n", planID)
				}

				// plan was closed, release the locked funds
				changeReq := protoBalance.ChangeBalanceRequest{
					UserID:         completedOrderEvent.UserID,
					AccountID:      completedOrderEvent.AccountID,
					CurrencySymbol: completedOrderEvent.FinalCurrencySymbol,
					Amount:         completedOrderEvent.FinalCurrencyBalance,
				}
				service.AccountClient.UnlockBalance(ctx, &changeReq)
			}

		case err != nil:
			log.Println("completed order error on FindChildOrders -- ", err.Error())

		default:
			// load new plan order with false - it is not a revision of an active order
			if err := service.ContinuePlan(ctx, nextPlanOrders); err != nil {
				log.Println("could not load the plan orders -- ", err.Error())
			}
		}
	} else if completedOrderEvent.Status == constPlan.Failed {
		if err := repoPlan.UpdateOrderErrors(service.DB, completedOrderEvent.OrderID, completedOrderEvent.Details); err != nil {
			log.Println("could not update order error -- ", err.Error())
		}

		// in theory this should not be required because the plan should close
		// to free the locked funds

		// unlock the initial balance
		//changeReq := protoBalance.ChangeBalanceRequest{
		//	UserID:         completedOrderEvent.UserID,
		//	AccountID:      completedOrderEvent.AccountID,
		//	CurrencySymbol: completedOrderEvent.InitialCurrencySymbol,
		//	Amount:         -completedOrderEvent.InitialCurrencyBalance,
		//}
		//service.AccountClient.ChangeLockedBalance(ctx, &changeReq)

		//// add remainder to initial balance
		//if completedOrderEvent.InitialCurrencyRemainder > 0 {
		//	changeReq.Amount = completedOrderEvent.InitialCurrencyRemainder
		//	service.AccountClient.ChangeAvailableBalance(ctx, &changeReq)
		//}
	}

	return nil
}

// Used to populate the engine on engine start. The engine will broadcast an EngineStartEvent.
func (service *PlanService) HandleStartEngine(ctx context.Context, engine *protoEvt.EngineStartEvent) error {

	plans, err := repoPlan.FindActivePlans(service.DB)
	if err != nil {
		log.Println("could not find active plans -- ", err)
	}

	// must sleep before sending off to execution engine
	// because engine might not have fully started yet
	time.Sleep(5 * time.Second)

	// TODO we need to explore a different approach here that is more efficient.
	for _, plan := range plans {
		// load the active orders - these are not revisions of active orders since it is assumed
		// the engine is asking to reload them from the DB
		if err = service.ContinuePlan(ctx, plan); err != nil {
			log.Println("load plan error -- ", err)
		}
	}

	return nil
}

// AddPlans returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object. MarketName example: ADA-BTC where BTC is base.
func (service *PlanService) NewPlan(ctx context.Context, req *protoPlan.NewPlanRequest, res *protoPlan.PlanResponse) error {

	switch {
	case !ValidatePlanInputStatus(req.Status):
		res.Status = constRes.Fail
		res.Message = "plan status must be active, inactive, or historic"
		return nil
	case !ValidateMinOrder(req.Orders):
		res.Status = constRes.Fail
		res.Message = "at least one order required for a new plan."
		return nil
	case !ValidateOrderTrigger(req.Orders):
		res.Status = constRes.Fail
		res.Message = "orders must have triggers"
		return nil
	case !ValidateSingleRootNode(req.Orders):
		res.Status = constRes.Fail
		res.Message = "multiple root nodes found, only one is allowed"
		return nil
	case !ValidateConnectedRoutesFromParent(uuid.Nil.String(), req.Orders):
		res.Status = constRes.Fail
		res.Message = "an order does not have a valid parent_order_id in your request"
		return nil
	case !ValidateNodeCount(req.Orders):
		res.Status = constRes.Fail
		res.Message = "you can only post 10 inactive nodes at a time!"
		return nil
	case !ValidateNoneZeroBalance(req.Orders):
		res.Status = constRes.Fail
		res.Message = "non zero initialCurrencyBalance required for root order!"
		return nil
	case req.Orders[0].OrderType == constPlan.PaperOrder && !ValidatePaperOrders(req.Orders):
		res.Status = constRes.Fail
		res.Message = "you cannot add a market/limit order to a plan that will begin with a paper order"
		return nil
	}

	// fetch all order accounts
	accountIDs := make([]string, 0, len(req.Orders))
	for _, or := range req.Orders {
		accountIDs = append(accountIDs, or.AccountID)
	}

	akeys, err := service.fetchKeys(accountIDs)
	if err != nil {
		if strings.Contains(err.Error(), "invalid input") {
			res.Status = constRes.Fail
			res.Message = fmt.Sprintf("valid accountID required for each order")
			return nil
		}

		msg := fmt.Sprintf("ecountered error when fetching keys: %s\n", err.Error())
		log.Println(msg)

		res.Status = constRes.Error
		res.Message = msg
		return nil
	}

	none := uuid.Nil.String()
	planID := uuid.New()
	now := string(pq.FormatTimestamp(time.Now().UTC()))
	newOrders := make([]*protoOrder.Order, 0, len(req.Orders))
	exchange := ""
	depth := uint32(0)
	totalDepth := uint32(0)

	for _, or := range req.Orders {
		orderStatus := constPlan.Inactive

		if or.MarketName == "" || or.AccountID == "" {
			res.Status = constRes.Fail
			res.Message = "missing marketName/accountID for order"
			return nil
		}
		if !strings.Contains(or.MarketName, "-") {
			res.Status = constRes.Fail
			res.Message = "marketName must be currency-base: e.g. ADA-BTC"
			return nil
		}
		if !ValidateOrderType(or.OrderType) {
			res.Status = constRes.Fail
			res.Message = "market, or limit required for order type"
			return nil
		}
		if !ValidateOrderSide(or.Side) {
			res.Status = constRes.Fail
			res.Message = "buy or sell required for order side"
			return nil
		}

		// root order starts at depth 1
		if or.ParentOrderID == none {
			depth = 1
		} else {

			for _, o := range newOrders {
				if o.OrderID == or.ParentOrderID {
					depth = o.PlanDepth + 1
					break
				}
			}

			// this will happen if the parent order for this order cannot be found
			// it essentially means the orders are not in the correct order - sorted from parent to child
			if depth == 0 {
				res.Status = constRes.Fail
				res.Message = "the orders must be sorted by plan depth, i.e. parents must be before children"
				return nil
			}
		}

		if depth > totalDepth {
			totalDepth = depth
		}

		// root order status should be active
		if or.ParentOrderID == none && req.Status == constPlan.Active {
			orderStatus = constPlan.Active
		}

		// assign exchange name from account
		foundKey := false
		accType := constAccount.AccountPaper
		for _, ky := range akeys {
			if ky.AccountID == or.AccountID {
				foundKey = true
				exchange = ky.Exchange
				accType = ky.AccountType

				if ky.Status != constAccount.AccountValid {
					res.Status = constRes.Fail
					res.Message = "using an invalid account. They key for the account is no longer valid. Possibly because it changed or was removed from the exchange."
					return nil
				}
			}
		}
		// if no match on account ID, then the request sent in an account ID that does not exist
		if !foundKey {
			res.Status = constRes.Fail
			res.Message = fmt.Sprintf("the requested account %s does not exist", or.AccountID)
			return nil
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
		if or.Side == constPlan.Sell {
			symbol = symbolPair[0]
		}

		order := protoOrder.Order{
			AccountID:              or.AccountID,
			AccountType:            accType,
			OrderID:                or.OrderID,
			OrderPriority:          or.OrderPriority,
			OrderType:              or.OrderType,
			OrderTemplateID:        or.OrderTemplateID,
			ParentOrderID:          or.ParentOrderID,
			PlanID:                 planID.String(),
			PlanDepth:              depth,
			Side:                   or.Side,
			LimitPrice:             commonUtil.ToFixedFloor(or.LimitPrice, precision),
			Exchange:               exchange,
			MarketName:             or.MarketName,
			InitialCurrencySymbol:  symbol,
			InitialCurrencyBalance: commonUtil.ToFixedFloor(or.InitialCurrencyBalance, precision),
			Status:                 orderStatus,
			Grupo:                  or.Grupo,
			Triggers:               triggers,
			CreatedOn:              now,
			UpdatedOn:              now,
		}
		newOrders = append(newOrders, &order)
	}

	currencySymbol := newOrders[0].InitialCurrencySymbol
	currencyBalance := newOrders[0].InitialCurrencyBalance
	accountID := newOrders[0].AccountID

	// is the initial amount in the account balance?
	validBalance, err := service.validateAvailableBalance(ctx, currencySymbol, currencyBalance, req.UserID, accountID)
	if err != nil {
		msg := fmt.Sprintf("validateAvailableBalance error %s", err.Error())
		log.Println(msg)

		res.Status = constRes.Error
		res.Message = msg
		return nil
	}
	if !validBalance {
		res.Status = constRes.Fail
		res.Message = fmt.Sprintf("insufficient %s balance, requested amount %.8f exceeds available balance", currencySymbol, currencyBalance)
		return nil
	}

	lockRequest := protoBalance.ChangeBalanceRequest{
		AccountID:      accountID,
		UserID:         req.UserID,
		CurrencySymbol: currencySymbol,
		Amount:         currencyBalance,
	}

	count := repoPlan.FindUserPlanCount(service.DB, req.UserID)
	title := req.Title
	if title == "" {
		title = fmt.Sprintf("Trade %d", count+1)
	}

	// base currency defaults to USDT
	baseCurrencySymbol := "USDT"
	if req.UserCurrencySymbol != "" {
		baseCurrencySymbol = req.UserCurrencySymbol
	}

	// TODO if req.IntialTimestamp compute the initial currency balance using the baseCurrencySymbol

	pln := protoPlan.Plan{
		PlanID:                 planID.String(),
		PlanTemplateID:         req.PlanTemplateID,
		UserID:                 req.UserID,
		Title:                  title,
		TotalDepth:             totalDepth,
		UserCurrencySymbol:     baseCurrencySymbol,
		ActiveCurrencySymbol:   newOrders[0].InitialCurrencySymbol,
		ActiveCurrencyBalance:  newOrders[0].InitialCurrencyBalance,
		InitialCurrencySymbol:  newOrders[0].InitialCurrencySymbol,
		InitialCurrencyBalance: newOrders[0].InitialCurrencyBalance,
		InitialTimestamp:       req.InitialTimestamp,
		Exchange:               newOrders[0].Exchange,
		LastExecutedPlanDepth:  0,
		LastExecutedOrderID:    none,
		Orders:                 newOrders,
		UserPlanNumber:         count + 1,
		Status:                 req.Status,
		CloseOnComplete:        req.CloseOnComplete,
		CreatedOn:              now,
		UpdatedOn:              now,
	}

	err = repoPlan.InsertPlan(service.DB, &pln)
	if err != nil {
		msg := fmt.Sprintf("insert plan failed %s", err.Error())
		log.Println(msg)

		res.Status = constRes.Error
		res.Message = msg
		return nil
	}

	// activate first plan order if plan is active
	if pln.Status == constPlan.Active {
		p := pln
		p.Orders = []*protoOrder.Order{pln.Orders[0]}

		// assign keys
		for _, ky := range akeys {
			if ky.AccountID == p.Orders[0].AccountID {
				p.Orders[0].KeyPublic = ky.KeyPublic
				p.Orders[0].KeySecret = ky.KeySecret
			}
		}

		// publish this plan to the engine
		if err := service.publishPlan(ctx, &p); err != nil {
			// TODO return a warning here
			res.Status = constRes.Error
			res.Message = "could not publish first order: " + err.Error()
			return nil
		}
	}

	r, _ := service.AccountClient.LockBalance(ctx, &lockRequest)
	if r.Status != constRes.Success {
		log.Println("could not lock the initial balance within update -- ", err.Error())
	}

	res.Status = constRes.Success
	res.Data = &protoPlan.PlanData{Plan: &pln}

	return nil
}

// GetUserPlan returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object.
func (service *PlanService) GetUserPlan(ctx context.Context, req *protoPlan.GetUserPlanRequest, res *protoPlan.PlanResponse) error {
	plan, error := repoPlan.FindPlanOrders(service.DB, req)

	switch {
	case error == sql.ErrNoRows:
		res.Status = constRes.Nonentity
		res.Message = fmt.Sprintf("planID not found %s", req.PlanID)
	case error != nil:
		res.Status = constRes.Error
		res.Message = error.Error()
	// case plan.totalDepth < req.PlanDepth:
	// 	res.Status = constRes.Nonentity
	// 	res.Message = "plan depth out of bounds, max depth is %s"
	case error == nil:
		res.Status = constRes.Success
		res.Data = &protoPlan.PlanData{Plan: plan}
	}

	return nil
}

// GetUserPlans returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object.
func (service *PlanService) GetUserPlans(ctx context.Context, req *protoPlan.GetUserPlansRequest, res *protoPlan.PlansPageResponse) error {
	// search by userID and status
	page, err := repoPlan.FindUserPlansWithStatus(service.DB, req.UserID, req.Status, req.Page, req.PageSize)

	if err == nil {
		res.Status = constRes.Success
		res.Data = page
	} else {
		res.Status = constRes.Error
		res.Message = err.Error()
	}

	return nil
}

// We can delete plans that have no filled orders and that are inactive. This becomes an abort plan
// if the plan status is active.
func (service *PlanService) DeletePlan(ctx context.Context, req *protoPlan.DeletePlanRequest, res *protoPlan.PlanResponse) error {
	pln, err := repoPlan.FindPlanWithUnexecutedOrders(service.DB, req.PlanID)
	switch {
	case err == sql.ErrNoRows:
		res.Status = constRes.Nonentity
		res.Message = fmt.Sprintf("planID not found %s", req.PlanID)
		return nil

	case err != nil:
		res.Status = constRes.Error
		res.Message = fmt.Sprintf("unexpected error in DeletePlan: %s", err.Error())
		return nil

	case pln.Status == constPlan.Active:

		req := protoEngine.KillRequest{PlanID: pln.PlanID}
		service.EngineClient.KillPlan(ctx, &req)
	}

	pln.Status = constPlan.Deleted
	err = repoPlan.UpdatePlanStatus(service.DB, req.PlanID, pln.Status)

	if err != nil {
		res.Status = constRes.Error
		res.Message = err.Error()
	}

	res.Status = constRes.Success
	res.Data = &protoPlan.PlanData{
		Plan: pln,
	}

	return nil
}

// UpdatePlan returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object.
func (service *PlanService) UpdatePlan(ctx context.Context, req *protoPlan.UpdatePlanRequest, res *protoPlan.PlanResponse) error {
	// load current state of plan
	// the plan should be paused long before UpdatePlan is called
	// this function assumes that the plan is inactive
	pln, err := repoPlan.FindPlanWithUnexecutedOrders(service.DB, req.PlanID)

	switch {
	case err == sql.ErrNoRows:
		res.Status = constRes.Nonentity
		res.Message = fmt.Sprintf("planID not found %s", req.PlanID)
		return nil
	case err != nil:
		msg := fmt.Sprintf("FindPlanWithUnexecutedOrders error: %s", err.Error())
		log.Println(msg)
		res.Status = constRes.Error
		res.Message = fmt.Sprintf(err.Error())
		return nil
	case len(req.Orders) == 0 && pln.LastExecutedPlanDepth == 0:
		res.Status = constRes.Fail
		res.Message = "to plan or not to plan. That is the question. A plan must have at least 1 order."
		return nil
	case !ValidateNonExecutedOrder(pln.Orders, req.Orders):
		res.Status = constRes.Fail
		res.Message = "an order has executed."
		return nil
	case !ValidateOrderTrigger(req.Orders):
		res.Status = constRes.Fail
		res.Message = "orders must have triggers"
		return nil
	case !ValidatePlanInputStatus(req.Status):
		res.Status = constRes.Fail
		res.Message = "plan status must be active, inactive"
		return nil
	case !ValidateConnectedRoutesFromParent(pln.LastExecutedOrderID, req.Orders):
		// all orders must be connected using parentOrderID
		res.Status = constRes.Fail
		res.Message = "this ain't no tree! All orders must be connected using the parentOrderID relationship."
		return nil
	case pln.LastExecutedPlanDepth == 0 && !ValidateNoneZeroBalance(req.Orders):
		// you must commit a balance for the plan in the first order
		res.Status = constRes.Fail
		res.Message = "the initialCurrencyBalance must be set for the root order"
		return nil
	case pln.LastExecutedPlanDepth == 0 && !ValidateSingleRootNode(req.Orders):
		// you can't start a plan without a root order
		res.Status = constRes.Fail
		res.Message = "multiple root nodes found, only one is allowed"
		return nil
	case pln.LastExecutedPlanDepth > 0 && !ValidateChildNodes(req.Orders):
		// update on an executed tree can only append child orders
		res.Status = constRes.Fail
		res.Message = fmt.Sprintf("an order's parentOrderID is %s. This plan already has an executed root order.", uuid.Nil.String())
		return nil
	case !ValidateNodeCount(req.Orders):
		res.Status = constRes.Fail
		res.Message = "you can only apply 10 inactive orders at a time!"
		return nil
	case pln.LastExecutedPlanDepth == 0 && req.Orders[0].OrderType == constPlan.PaperOrder && !ValidatePaperOrders(req.Orders):
		// can't mix real orders to a paper plan
		res.Status = constRes.Fail
		res.Message = "you cannot append market or limit orders to a plan that will begin with a paper order"
		return nil
	}

	// fetch the keys
	accountIDs := make([]string, 0, len(req.Orders))
	for _, or := range req.Orders {
		accountIDs = append(accountIDs, or.AccountID)
	}
	akeys := make([]*protoAccount.AccountKey, 0)

	if len(req.Orders) > 0 {
		akeys, err = service.fetchKeys(accountIDs)
		if err != nil {
			res.Status = constRes.Error
			res.Message = fmt.Sprintf("ecountered error when fetching account keys: %s", err.Error())
			return nil
		}
	}

	none := uuid.Nil.String()
	now := string(pq.FormatTimestamp(time.Now().UTC()))
	newOrders := make([]*protoOrder.Order, 0, len(req.Orders))
	exchange := ""
	keyPublic := ""
	keySecret := ""
	depth := pln.LastExecutedPlanDepth
	totalDepth := depth
	lockRequest := new(protoBalance.ChangeBalanceRequest)

	for _, or := range req.Orders {
		orderStatus := constPlan.Inactive

		if or.MarketName == "" || or.AccountID == "" {
			res.Status = constRes.Fail
			res.Message = "missing marketName/accountID for order"
			return nil
		}
		if !strings.Contains(or.MarketName, "-") {
			res.Status = constRes.Fail
			res.Message = "marketName must be currency-base: e.g. ADA-BTC"
			return nil
		}
		if !ValidateOrderType(or.OrderType) {
			res.Status = constRes.Fail
			res.Message = "market or limit required for order type"
			return nil
		}
		if !ValidateOrderSide(or.Side) {
			res.Status = constRes.Fail
			res.Message = "buy or sell required for order side"
			return nil
		}

		// compute the depth for the order
		if or.ParentOrderID == pln.LastExecutedOrderID {
			depth++
		} else {
			for _, o := range newOrders {
				if o.OrderID == or.ParentOrderID {
					depth = o.PlanDepth + 1
					break
				}
			}

			// this will happen if the parent order for this order cannot be found
			// the depth would not change
			if depth == pln.LastExecutedPlanDepth {
				res.Status = constRes.Fail
				res.Message = "the orders must be sorted by plan depth, i.e. parents must be before children"
				return nil
			}
		}

		if depth > totalDepth {
			totalDepth = depth
		}

		if or.ParentOrderID == none && req.Status == constPlan.Active {
			orderStatus = constPlan.Active
		}

		// assign exchange name from key
		foundKey := false
		accType := constAccount.AccountPaper
		for _, ky := range akeys {
			if ky.AccountID == or.AccountID {
				exchange = ky.Exchange
				keyPublic = ky.KeyPublic
				keySecret = ky.KeySecret
				foundKey = true
				accType = ky.AccountType

				if ky.Status != constAccount.AccountValid {
					res.Status = constRes.Fail
					res.Message = fmt.Sprintf("the key for account %s is invalid.", ky.AccountID)
					return nil
				}
			}
		}
		// if no match on account ID, then the request sent in an account ID that does not exist
		if !foundKey {
			res.Status = constRes.Fail
			res.Message = fmt.Sprintf("the requested account %s does not exist", or.AccountID)
			return nil
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
		if or.Side == constPlan.Sell {
			currencySymbol = symbolPair[0]
		}

		// validate the balance for non paper orders that are set to use a predefined balance
		if or.InitialCurrencyBalance > 0 {
			validBalance, err := service.validateAvailableBalance(ctx, currencySymbol, or.InitialCurrencyBalance, req.UserID, or.AccountID)
			if err != nil {
				res.Status = constRes.Error
				res.Message = fmt.Sprintf("failed to validate the currency balance for %s: %s", currencySymbol, err.Error())
				return nil
			}
			if !validBalance {
				res.Status = constRes.Fail
				res.Message = fmt.Sprintf("insufficient %s balance, %.8f requested in orderID: %s", currencySymbol, or.InitialCurrencyBalance, or.OrderID)
				return nil
			}

			lockRequest.AccountID = or.AccountID
			lockRequest.UserID = req.UserID
			lockRequest.CurrencySymbol = currencySymbol
			lockRequest.Amount = or.InitialCurrencyBalance
		}

		order := protoOrder.Order{
			AccountID:              or.AccountID,
			AccountType:            accType,
			KeyPublic:              keyPublic,
			KeySecret:              keySecret,
			OrderID:                or.OrderID,
			OrderPriority:          or.OrderPriority,
			OrderType:              or.OrderType,
			OrderTemplateID:        or.OrderTemplateID,
			ParentOrderID:          or.ParentOrderID,
			PlanID:                 req.PlanID,
			PlanDepth:              depth,
			Side:                   or.Side,
			LimitPrice:             commonUtil.ToFixedFloor(or.LimitPrice, precision),
			Exchange:               exchange,
			MarketName:             or.MarketName,
			InitialCurrencySymbol:  currencySymbol,
			InitialCurrencyBalance: commonUtil.ToFixedFloor(or.InitialCurrencyBalance, precision),
			Grupo:     or.Grupo,
			Status:    orderStatus,
			Triggers:  triggers,
			CreatedOn: now,
			UpdatedOn: now,
		}
		newOrders = append(newOrders, &order)
	}

	// if the plan is currently active we need to kill the plan by its ID in the engine
	// the length of orders in the existing plan should include the last executed if any
	// more orders than that means there are active orders currency running

	activeOrder := false
	for _, previous := range pln.Orders {
		if previous.Status == constPlan.Active {
			activeOrder = true
		}
	}

	// only key plans that are active and that have an active order
	if pln.Status == constPlan.Active && activeOrder {
		req := protoEngine.KillRequest{PlanID: pln.PlanID}
		_, err := service.EngineClient.KillPlan(ctx, &req)
		if err != nil {
			res.Status = constRes.Fail
			res.Message = "you cannot update this plan because an order for this plan has updated. To avoid seeing this message again try pausing your plan before you update it."
			return nil
		}
	}

	txn, err := service.DB.Begin()
	if err != nil {
		return err
	}

	// Overwrite the entire unexecuted portion of the plan tree with the new orders above.
	// Gather all previous orderIDs for this plan so we can drop them from the DB.
	orderIDs := make([]string, 0)
	for _, o := range pln.Orders {
		if o.Status != constPlan.Filled {
			orderIDs = append(orderIDs, o.OrderID)
		}
	}

	// if root has not executed yet and we have new orders
	if pln.LastExecutedPlanDepth == 0 && len(newOrders) > 0 {
		pln.ActiveCurrencySymbol = newOrders[0].InitialCurrencySymbol
		pln.ActiveCurrencyBalance = newOrders[0].InitialCurrencyBalance
		pln.InitialCurrencySymbol = newOrders[0].InitialCurrencySymbol
		pln.InitialCurrencyBalance = newOrders[0].InitialCurrencyBalance
		pln.Exchange = newOrders[0].Exchange
		if err := repoPlan.UpdatePlanContextTxn(txn, ctx, pln.PlanID, pln.ActiveCurrencySymbol, pln.Exchange, pln.ActiveCurrencyBalance); err != nil {
			txn.Rollback()
			res.Status = constRes.Error
			res.Message = "error encountered while updating the plan context: " + err.Error()
			return nil
		}
	}
	if pln.UserCurrencySymbol != req.UserCurrencySymbol {
		pln.UserCurrencySymbol = req.UserCurrencySymbol
		if err := repoPlan.UpdatePlanBaseCurrencyTxn(txn, ctx, pln.PlanID, pln.UserCurrencySymbol); err != nil {
			txn.Rollback()
			res.Status = constRes.Error
			res.Message = "error encountered while updating the plan base currency: " + err.Error()
			return nil
		}
	}
	if pln.Status != req.Status {
		pln.Status = req.Status
		if err := repoPlan.UpdatePlanStatusTxn(txn, ctx, pln.PlanID, pln.Status); err != nil {
			txn.Rollback()
			res.Status = constRes.Error
			res.Message = "error encountered while updating the plan status: " + err.Error()
			return nil
		}
		// TODO unlock the active currency? when close
	}
	if pln.InitialTimestamp != req.InitialTimestamp && req.InitialTimestamp != "" {

		// TODO if req.IntialTimestamp compute the initial currency balance using the baseCurrencySymbol
		pln.InitialTimestamp = req.InitialTimestamp
		if err := repoPlan.UpdatePlanInitTimeTxn(txn, ctx, pln.PlanID, pln.InitialTimestamp); err != nil {
			txn.Rollback()
			res.Status = constRes.Error
			res.Message = "error encountered while updating the plan init time: " + err.Error()
			return nil
		}
	}
	if pln.TotalDepth != totalDepth {
		pln.TotalDepth = totalDepth
		if err := repoPlan.UpdatePlanTotalDepthTxn(txn, ctx, pln.PlanID, pln.TotalDepth); err != nil {
			txn.Rollback()
			res.Status = constRes.Error
			res.Message = "error encountered while updating the plan total depth: " + err.Error()
			return nil
		}
	}
	if req.Title != "" {
		pln.Title = req.Title
		if err := repoPlan.UpdatePlanTitleTxn(txn, ctx, pln.PlanID, pln.Title); err != nil {
			txn.Rollback()
			res.Status = constRes.Error
			res.Message = "error encountered while updating the plan title: " + err.Error()
			return nil
		}
	}
	if pln.CloseOnComplete != req.CloseOnComplete {
		pln.CloseOnComplete = req.CloseOnComplete
		if err := repoPlan.UpdatePlanCloseOnCompleteTxn(txn, ctx, pln.PlanID, pln.CloseOnComplete); err != nil {
			txn.Rollback()
			res.Status = constRes.Error
			res.Message = "error encountered while updating the plan close on complete option: " + err.Error()
			return nil
		}

		// close the plan
		if pln.CloseOnComplete && len(req.Orders) == 0 {
			if err := repoPlan.UpdatePlanStatusTxn(txn, ctx, pln.PlanID, constPlan.Closed); err != nil {
				txn.Rollback()
				res.Status = constRes.Error
				res.Message = "error encountered while closing the plan: " + err.Error()
				return nil
			}

			for _, order := range pln.Orders {
				if order.OrderID == pln.LastExecutedOrderID {
					// release the locked funds from the last executed order
					changeReq := protoBalance.ChangeBalanceRequest{
						UserID:         pln.UserID,
						AccountID:      order.AccountID,
						CurrencySymbol: pln.ActiveCurrencySymbol,
						Amount:         pln.ActiveCurrencyBalance,
					}
					service.AccountClient.UnlockBalance(ctx, &changeReq)
					break
				}
			}
		}
	}

	if pln.PlanTemplateID != req.PlanTemplateID {
		pln.PlanTemplateID = req.PlanTemplateID
		if err := repoPlan.UpdatePlanTemplateTxn(txn, ctx, pln.PlanID, pln.PlanTemplateID); err != nil {
			txn.Rollback()
			res.Status = constRes.Error
			res.Message = "error encountered while updating the plan template: " + err.Error()
			return nil
		}
	}

	// keep the update timestamp of the plan in sync with the orders
	// no particular reason, but it could be useful in debugging
	if err := repoPlan.UpdatePlanTimestampTxn(txn, ctx, pln.PlanID, now); err != nil {
		txn.Rollback()
		res.Status = constRes.Error
		res.Message = "error encountered while updating the plan timestamp: " + err.Error()
		return nil
	}

	// drop current orders from the plan
	if err := repoPlan.DeleteOrders(txn, ctx, orderIDs); err != nil {
		txn.Rollback()
		res.Status = constRes.Error
		res.Message = "error while deleting the previous orders: " + err.Error()
		return nil
	}

	if len(newOrders) > 0 {
		// insert new orders for this plan
		if err := repoPlan.InsertOrders(txn, newOrders); err != nil {
			txn.Rollback()
			res.Status = constRes.Error
			res.Message = "insert orders error: " + err.Error()
			return nil
		}

		newTriggers := make([]*protoOrder.Trigger, 0, len(newOrders))
		for _, o := range newOrders {
			for _, t := range o.Triggers {
				newTriggers = append(newTriggers, t)
			}
		}

		if err := repoPlan.InsertTriggers(txn, newTriggers); err != nil {
			txn.Rollback()
			return errors.New("bulk triggers failed: " + err.Error())
		}
	}

	txn.Commit()

	// we overwrite the pln orders here so the plan order response will contain the
	// new order. Publish plan shall also publish the new orders
	pln.Orders = newOrders

	// activate first plan order if plan is active
	if pln.Status == constPlan.Active {
		// pub new plan for execution
		if err := service.publishPlan(ctx, pln); err != nil {
			res.Status = constRes.Error
			res.Message = "could not fully set this plan active, error was: " + err.Error()
			return nil
		}
	}

	if lockRequest.AccountID != "" {
		r, _ := service.AccountClient.LockBalance(ctx, lockRequest)
		if r.Status != constRes.Success {
			log.Println("could not lock the initial balance within update -- ", err.Error())
		}
	}

	res.Status = constRes.Success
	res.Data = &protoPlan.PlanData{Plan: pln}

	return nil
}
