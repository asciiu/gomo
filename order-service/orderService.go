package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	bp "github.com/asciiu/gomo/balance-service/proto/balance"
	evt "github.com/asciiu/gomo/common/proto/events"
	orderRepo "github.com/asciiu/gomo/order-service/db/sql"
	pb "github.com/asciiu/gomo/order-service/proto/order"
	micro "github.com/micro/go-micro"
)

// MinBalance needed to submit order
const MinBalance = 0.00001000

// OrderService ...
type OrderService struct {
	DB      *sql.DB
	Client  bp.BalanceServiceClient
	NewBuy  micro.Publisher
	NewSell micro.Publisher
}

// These are all private functions
func (service *OrderService) publishOrder(publisher micro.Publisher, order *pb.Order) error {
	// process order here
	orderEvent := evt.OrderEvent{
		Exchange:     "Binance",
		OrderId:      order.OrderId,
		UserId:       order.UserId,
		ApiKeyId:     order.ApiKeyId,
		MarketName:   order.MarketName,
		BaseQuantity: order.BaseQuantity,
		Side:         order.Side,
		Conditions:   order.Conditions,
		Status:       "pending",
	}

	if err := publisher.Publish(context.Background(), &orderEvent); err != nil {
		return fmt.Errorf("publish error: %s %s", err, orderEvent)
	}
	log.Printf("publishOrder: %v", orderEvent)
	return nil
}

// private: validateBalance
func (service *OrderService) validateBalance(ctx context.Context, currency, userID, apikeyID string) error {
	balRequest := bp.GetUserBalanceRequest{
		UserId:   userID,
		ApiKeyId: apikeyID,
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

// LoadBuyOrder ...
func (service *OrderService) LoadBuyOrder(ctx context.Context, order *pb.Order) error {

	currencies := strings.Split(order.MarketName, "-")
	baseCurrency := currencies[1]

	if err := service.validateBalance(ctx, baseCurrency, order.UserId, order.ApiKeyId); err != nil {
		return err
	}

	if err := service.publishOrder(service.NewBuy, order); err != nil {
		return err
	}

	return nil
}

// LoadSellOrder ...
func (service *OrderService) LoadSellOrder(ctx context.Context, order *pb.Order) error {
	currencies := strings.Split(order.MarketName, "-")
	currency := currencies[0]

	if err := service.validateBalance(ctx, currency, order.UserId, order.ApiKeyId); err != nil {
		return err
	}

	if err := service.publishOrder(service.NewSell, order); err != nil {
		return err
	}

	return nil
}

// AddOrders ...
func (service *OrderService) AddOrders(ctx context.Context, req *pb.OrdersRequest, response *pb.OrderListResponse) error {
	orders := make([]*pb.Order, 0)
	requestOrders := req.Orders

	// begin with head order
	order := requestOrders[0]

	parentOrderID := order.ParentOrderId

	// if zero for parentOrderId we have must load this order immediately
	if parentOrderID == "00000000-0000-0000-0000-000000000000" {

		currencies := strings.Split(order.MarketName, "-")
		var currency string

		if order.Side == "buy" {
			// base currency will be second
			currency = currencies[1]
		} else {
			currency = currencies[0]
		}

		if err := service.validateBalance(ctx, currency, order.UserId, order.ApiKeyId); err != nil {
			response.Status = "fail"
			response.Message = err.Error()
			return nil
		}

		o, error := orderRepo.InsertOrder(service.DB, order)

		if error != nil {
			response.Status = "error"
			response.Message = "ecountered error on Insert: " + error.Error()
			return nil
		}

		orders = append(orders, o)
		var pub micro.Publisher

		if o.Side == "buy" {
			pub = service.NewBuy
		} else {
			pub = service.NewSell
		}

		if err := service.publishOrder(pub, o); err != nil {
			response.Status = "error"
			response.Message = "could not publish order: " + err.Error()
			return nil
		}

		requestOrders = requestOrders[1:]
		parentOrderID = o.OrderId
	}

	// loop through and insert the rest of the chain
	for i := 0; i < len(requestOrders); i++ {
		order = requestOrders[i]

		// assign the parent order id for following orders
		order.ParentOrderId = parentOrderID

		o, error := orderRepo.InsertOrder(service.DB, order)

		if error != nil {
			response.Status = "error"
			response.Message = "could not insert order: " + error.Error()
			return nil
		}

		orders = append(orders, o)
		parentOrderID = o.OrderId
	}

	response.Status = "success"
	response.Data = &pb.UserOrdersData{
		Orders: orders,
	}
	return nil
}

func (service *OrderService) GetUserOrder(ctx context.Context, req *pb.GetUserOrderRequest, res *pb.OrderResponse) error {
	order, error := orderRepo.FindOrderById(service.DB, req)

	switch {
	case error == nil:
		res.Status = "success"
		res.Data = &pb.UserOrderData{
			Order: order,
		}
		return nil
	default:
		res.Status = "error"
		res.Message = error.Error()
		return error
	}
}

func (service *OrderService) GetUserOrders(ctx context.Context, req *pb.GetUserOrdersRequest, res *pb.OrderListResponse) error {
	orders, error := orderRepo.FindOrdersByUserId(service.DB, req)

	switch {
	case error == nil:
		res.Status = "success"
		res.Data = &pb.UserOrdersData{
			Orders: orders,
		}
		return nil
	default:
		res.Status = "error"
		res.Message = error.Error()
		return error
	}
}

func (service *OrderService) RemoveOrder(ctx context.Context, req *pb.RemoveOrderRequest, res *pb.OrderResponse) error {
	error := orderRepo.DeleteOrder(service.DB, req.OrderId)
	switch {
	case error == nil:
		res.Status = "success"
		return nil
	default:
		res.Status = "error"
		res.Message = error.Error()
		return error
	}
}

func (service *OrderService) UpdateOrder(ctx context.Context, req *pb.OrderRequest, res *pb.OrderResponse) error {
	order, error := orderRepo.UpdateOrder(service.DB, req)
	switch {
	case error == nil:
		res.Status = "success"
		res.Data = &pb.UserOrderData{
			Order: order,
		}
		return nil
	default:
		res.Status = "error"
		res.Message = error.Error()
		return error
	}
}
