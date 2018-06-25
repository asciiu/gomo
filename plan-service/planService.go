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
	NewPlan   micro.Publisher
	//NewBuy    micro.Publisher
	//NewSell   micro.Publisher
}

// private: This is where the order events are published to the rest of the system
// this function should only be callable from within the PlanService
func (service *PlanService) publishPlan(ctx context.Context, plan *protoPlan.Plan) error {

	currencies := strings.Split(plan.MarketName, "-")
	currency := currencies[0]

	// // convert order to order event
	orderEvent := evt.OrderEvent{
		Exchange: plan.Exchange,
		//PlanID:   plan.PlanID,
		UserID: plan.UserID,
		// 	Key:        order.Key,
		// 	Secret:     order.Secret,
		KeyID:      plan.KeyID,
		MarketName: plan.MarketName,
		Currency:   currency,
		// 	Quantity:   order.BaseQuantity,
		//Price: plan.Price,
		//Side:  order.Side,
		//  NextOrderID: ...
		// 	PlanType:   order.PlanType,
		// 	Conditions: order.Conditions,
		// 	Status:     order.Status,
	}

	// //if err := publisher.Publish(context.Background(), &orderEvent); err != nil {
	// if err := service.NewPlan.Publish(context.Background(), &orderEvent); err != nil {
	// 	return fmt.Errorf("publish error: %s -- orderEvent: %+v", err, &orderEvent)
	// }
	log.Printf("publish order event -- %+v\n", &orderEvent)
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

	// ask key service for key
	r, _ := service.KeyClient.GetUserKey(context.Background(), &getRequest)
	//	if keyResponse.Status != response.Success {
	//		return fmt.Errorf("key is invalid for order -- %s, %#v", keyResponse.Message, order)
	//	}

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

	// validate key
	// TODO we need to get the secret with this key
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

	if pln.Status == plan.Active {
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
