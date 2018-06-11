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
	orders "github.com/asciiu/gomo/order-service/proto/order"
	micro "github.com/micro/go-micro"
)

// MinBalance needed to submit order
const MinBalance = 0.00001000

// OrderService ...
type OrderService struct {
	DB        *sql.DB
	Client    balances.BalanceServiceClient
	KeyClient keys.KeyServiceClient
	NewOrder  micro.Publisher
	//NewBuy    micro.Publisher
	//NewSell   micro.Publisher
}

// private: This is where the order events are published to the rest of the system
// this function should only be callable from within the OrderService
func (service *OrderService) publishOrder(ctx context.Context, order *orders.Order) error {

	keyReq := keys.GetUserKeyRequest{
		UserID: order.UserID,
		KeyID:  order.KeyID,
	}

	keyResponse, _ := service.KeyClient.GetUserKey(ctx, &keyReq)
	if keyResponse.Status != response.Success {
		return fmt.Errorf("key is invalid for order -- %s, %#v", keyResponse.Message, order)
	}

	// convert order to order event
	orderEvent := evt.OrderEvent{
		Exchange:     order.Exchange,
		OrderID:      order.OrderID,
		UserID:       order.UserID,
		Key:          keyResponse.Data.Key.Key,
		Secret:       keyResponse.Data.Key.Secret,
		KeyID:        order.KeyID,
		MarketName:   order.MarketName,
		BaseQuantity: order.BaseQuantity,
		Side:         order.Side,
		OrderType:    order.OrderType,
		Conditions:   order.Conditions,
		Status:       order.Status,
	}

	//if err := publisher.Publish(context.Background(), &orderEvent); err != nil {
	if err := service.NewOrder.Publish(context.Background(), &orderEvent); err != nil {
		return fmt.Errorf("publish error -- %s -- %+v", err, &orderEvent)
	}
	log.Printf("publish order -- %+v\n", &orderEvent)
	return nil
}

// private: validateBalance
func (service *OrderService) validateBalance(ctx context.Context, currency, userID, apikeyID string) error {
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

// LoadBuyOrder should not be invoked by the client. This function was designed to load an
// order after an order was filled.
func (service *OrderService) LoadBuyOrder(ctx context.Context, order *orders.Order) error {

	currencies := strings.Split(order.MarketName, "-")
	baseCurrency := currencies[1]

	if err := service.validateBalance(ctx, baseCurrency, order.UserID, order.KeyID); err != nil {
		return err
	}

	if err := service.publishOrder(ctx, order); err != nil {
		return err
	}

	return nil
}

// LoadSellOrder should not be invoked by the client. This function was designed to load an
// order after an order was filled.
func (service *OrderService) LoadSellOrder(ctx context.Context, order *orders.Order) error {
	currencies := strings.Split(order.MarketName, "-")
	currency := currencies[0]

	if err := service.validateBalance(ctx, currency, order.UserID, order.KeyID); err != nil {
		return err
	}

	if err := service.publishOrder(ctx, order); err != nil {
		return err
	}

	return nil
}

// AddOrders returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object.
func (service *OrderService) AddOrders(ctx context.Context, req *orders.OrdersRequest, res *orders.OrderListResponse) error {
	ordrs := make([]*orders.Order, 0)
	requestOrders := req.Orders

	// begin with head order
	order := requestOrders[0]

	parentOrderID := order.ParentOrderID

	// if zero for parentOrderID we must load this order immediately
	if parentOrderID == "00000000-0000-0000-0000-000000000000" {

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
		o, error := orderRepo.InsertOrder(service.DB, order, status.Open)

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

		if err := service.publishOrder(ctx, o); err != nil {
			res.Status = response.Error
			res.Message = "could not publish order: " + err.Error()
			return nil
		}

		requestOrders = requestOrders[1:]
		parentOrderID = o.OrderID
	}

	// loop through and insert the rest of the chain
	for i := 0; i < len(requestOrders); i++ {
		order = requestOrders[i]

		// assign the parent order id for following orders
		order.ParentOrderID = parentOrderID

		o, error := orderRepo.InsertOrder(service.DB, order, status.Pending)

		if error != nil {
			res.Status = response.Error
			res.Message = "could not insert order: " + error.Error()
			return nil
		}

		ordrs = append(ordrs, o)
		parentOrderID = o.OrderID
	}

	res.Status = response.Success
	res.Data = &orders.UserOrdersData{
		Orders: ordrs,
	}
	return nil
}

// GetUserOrder returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object.
func (service *OrderService) GetUserOrder(ctx context.Context, req *orders.GetUserOrderRequest, res *orders.OrderResponse) error {
	order, error := orderRepo.FindOrderByID(service.DB, req)

	switch {
	case error == nil:
		res.Status = response.Success
		res.Data = &orders.UserOrderData{
			Order: order,
		}
	default:
		res.Status = response.Error
		res.Message = error.Error()
	}
	return nil
}

// GetUserOrders returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object.
func (service *OrderService) GetUserOrders(ctx context.Context, req *orders.GetUserOrdersRequest, res *orders.OrderListResponse) error {
	ordrs, error := orderRepo.FindOrdersByUserID(service.DB, req)

	switch {
	case error == nil:
		res.Status = response.Success
		res.Data = &orders.UserOrdersData{
			Orders: ordrs,
		}
	default:
		res.Status = response.Error
		res.Message = error.Error()
	}
	return nil
}

// RemoveOrder returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object.
func (service *OrderService) RemoveOrder(ctx context.Context, req *orders.RemoveOrderRequest, res *orders.OrderResponse) error {
	error := orderRepo.DeleteOrder(service.DB, req.OrderID)
	switch {
	case error == nil:
		res.Status = response.Success
	default:
		res.Status = response.Error
		res.Message = error.Error()
	}
	return nil
}

// UpdateOrder returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object.
func (service *OrderService) UpdateOrder(ctx context.Context, req *orders.OrderRequest, res *orders.OrderResponse) error {
	order, error := orderRepo.UpdateOrder(service.DB, req)
	switch {
	case error == nil:
		res.Status = response.Success
		res.Data = &orders.UserOrderData{
			Order: order,
		}
	default:
		res.Status = response.Error
		res.Message = error.Error()
	}
	return nil
}
