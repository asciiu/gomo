package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	balances "github.com/asciiu/gomo/balance-service/proto/balance"
	"github.com/asciiu/gomo/common/constants/response"
	"github.com/asciiu/gomo/common/constants/side"
	"github.com/asciiu/gomo/common/constants/status"
	evt "github.com/asciiu/gomo/common/proto/events"
	keys "github.com/asciiu/gomo/key-service/proto/key"
	orderRepo "github.com/asciiu/gomo/order-service/db/sql"
	plans "github.com/asciiu/gomo/plan-service/proto/plan"
	micro "github.com/micro/go-micro"
)

// MinBalance needed to submit order
const MinBalance = 0.00001000

// PlanService ...
type PlanService struct {
	DB        *sql.DB
	Client    balances.BalanceServiceClient
	KeyClient keys.KeyServiceClient
	NewPlan   micro.Publisher
	//NewBuy    micro.Publisher
	//NewSell   micro.Publisher
}

// private: This is where the order events are published to the rest of the system
// this function should only be callable from within the PlanService
func (service *PlanService) publishPlan(ctx context.Context, order *plans.Plan) error {

	currencies := strings.Split(order.MarketName, "-")
	currency := currencies[0]

	// convert order to order event
	orderEvent := evt.PlanEvent{
		Exchange:   order.Exchange,
		PlanID:     order.PlanID,
		UserID:     order.UserID,
		Key:        order.Key,
		Secret:     order.Secret,
		KeyID:      order.KeyID,
		MarketName: order.MarketName,
		Currency:   currency,
		Quantity:   order.BaseQuantity,
		Price:      order.Price,
		Side:       order.Side,
		PlanType:   order.PlanType,
		Conditions: order.Conditions,
		Status:     order.Status,
	}

	//if err := publisher.Publish(context.Background(), &orderEvent); err != nil {
	if err := service.NewPlan.Publish(context.Background(), &orderEvent); err != nil {
		return fmt.Errorf("publish error: %s -- orderEvent: %+v", err, &orderEvent)
	}
	log.Printf("publish order event -- %+v\n", &orderEvent)
	return nil
}

// private: validateBalance
func (service *PlanService) validateBalance(ctx context.Context, currency, userID, apikeyID string) error {
	balRequest := balances.GetUserBalanceRequest{
		UserID:   userID,
		KeyID:    apikeyID,
		Currency: currency,
	}

	balResponse, err := service.Client.GetUserBalance(ctx, &balRequest)
	if err != nil {
		return fmt.Errorf("ecountered error from GetUserBalance: %s", err.Error())
	}

	if balResponse.Data.Balance.Available < MinBalance {
		return fmt.Errorf("insufficient %s balance %.8f required", currency, MinBalance)
	}
	return nil
}

// LoadBuyPlan should not be invoked by the client. This function was designed to load an
// order after an order was filled.
func (service *PlanService) LoadPlan(ctx context.Context, order *plans.Plan) error {

	currencies := strings.Split(order.MarketName, "-")
	// default market currency
	currency := currencies[0]
	if order.Side == side.Buy {
		// buy uses base currency
		currency = currencies[1]
	}

	if err := service.validateBalance(ctx, currency, order.UserID, order.KeyID); err != nil {
		return err
	}

	if err := service.publishPlan(ctx, order); err != nil {
		return err
	}

	return nil
}

// AddPlans returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object.
func (service *PlanService) AddPlans(ctx context.Context, req *plans.PlansRequest, res *plans.PlanListResponse) error {
	ordrs := make([]*plans.Plan, 0)
	requestPlans := req.Plans

	// begin with head order
	order := requestPlans[0]

	parentPlanID := order.ParentPlanID

	// if zero for parentPlanID we must load this order immediately
	if parentPlanID == "00000000-0000-0000-0000-000000000000" {

		currencies := strings.Split(order.MarketName, "-")
		var currency string

		if order.Side == side.Buy {
			// base currency will be second
			currency = currencies[1]
		} else {
			currency = currencies[0]
		}

		if err := service.validateBalance(ctx, currency, order.UserID, order.KeyID); err != nil {
			res.Status = response.Fail
			res.Message = err.Error()
			return nil
		}

		// this order will be published to the executor so mark its status as Open
		o, error := orderRepo.InsertPlan(service.DB, order, status.Active)

		if error != nil {
			res.Status = response.Error
			res.Message = "ecountered error on Insert: " + error.Error()
			return nil
		}

		ordrs = append(ordrs, o)

		// we need to get the key for this order so we can publish it
		// to the engines
		keyReq := keys.GetUserKeyRequest{
			UserID: o.UserID,
			KeyID:  o.KeyID,
		}
		keyResponse, _ := service.KeyClient.GetUserKey(ctx, &keyReq)
		if keyResponse.Status != response.Success {
			return fmt.Errorf("key is invalid for order -- %s, %#v", keyResponse.Message, order)
		}

		o.Key = keyResponse.Data.Key.Key
		o.Secret = keyResponse.Data.Key.Secret

		if err := service.publishPlan(ctx, o); err != nil {
			res.Status = response.Error
			res.Message = "could not publish order: " + err.Error()
			return nil
		}

		requestPlans = requestPlans[1:]
		parentPlanID = o.PlanID
	}

	// loop through and insert the rest of the chain
	for i := 0; i < len(requestPlans); i++ {
		order = requestPlans[i]

		// assign the parent order id for following plans
		order.ParentPlanID = parentPlanID

		o, error := orderRepo.InsertPlan(service.DB, order, status.Pending)

		if error != nil {
			res.Status = response.Error
			res.Message = "could not insert order: " + error.Error()
			return nil
		}

		ordrs = append(ordrs, o)
		parentPlanID = o.PlanID
	}

	res.Status = response.Success
	res.Data = &plans.UserPlansData{
		Plans: ordrs,
	}
	return nil
}

// GetUserPlan returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object.
func (service *PlanService) GetUserPlan(ctx context.Context, req *plans.GetUserPlanRequest, res *plans.PlanResponse) error {
	order, error := orderRepo.FindPlanByID(service.DB, req)

	switch {
	case error == nil:
		res.Status = response.Success
		res.Data = &plans.UserPlanData{
			Plan: order,
		}
	default:
		res.Status = response.Error
		res.Message = error.Error()
	}
	return nil
}

// GetUserPlans returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object.
func (service *PlanService) GetUserPlans(ctx context.Context, req *plans.GetUserPlansRequest, res *plans.PlanListResponse) error {
	ordrs, error := orderRepo.FindPlansByUserID(service.DB, req)

	switch {
	case error == nil:
		res.Status = response.Success
		res.Data = &plans.UserPlansData{
			Plans: ordrs,
		}
	default:
		res.Status = response.Error
		res.Message = error.Error()
	}
	return nil
}

// RemovePlan returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object.
func (service *PlanService) DeletePlan(ctx context.Context, req *plans.RemovePlanRequest, res *plans.PlanResponse) error {
	error := orderRepo.DeletePlan(service.DB, req.PlanID)
	switch {
	case error == nil:
		res.Status = response.Success
	default:
		res.Status = response.Error
		res.Message = error.Error()
	}
	return nil
}

// UpdatePlan returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object.
func (service *PlanService) UpdatePlan(ctx context.Context, req *plans.PlanRequest, res *plans.PlanResponse) error {
	order, error := orderRepo.UpdatePlan(service.DB, req)
	switch {
	case error == nil:
		res.Status = response.Success
		res.Data = &plans.UserPlanData{
			Plan: order,
		}
	default:
		res.Status = response.Error
		res.Message = error.Error()
	}
	return nil
}
