package main

import (
	"context"
	"database/sql"
	"log"
	"strings"

	bp "github.com/asciiu/gomo/balance-service/proto/balance"
	"github.com/asciiu/gomo/common/enums"
	evt "github.com/asciiu/gomo/common/proto/events"
	orderRepo "github.com/asciiu/gomo/order-service/db/sql"
	pb "github.com/asciiu/gomo/order-service/proto/order"
	micro "github.com/micro/go-micro"
)

type OrderService struct {
	DB          *sql.DB
	Client      bp.BalanceServiceClient
	NewOrderPub micro.Publisher
}

// Add a buy order
func (service *OrderService) addBuyOrder(ctx context.Context, req *pb.OrderRequest, response *pb.OrderResponse) error {
	currencies := strings.Split(req.MarketName, "-")
	baseCurrency := currencies[1]

	balRequest := bp.GetUserBalanceRequest{
		UserId:   req.UserId,
		ApiKeyId: req.ApiKeyId,
		Currency: baseCurrency,
	}

	balResponse, err := service.Client.GetUserBalance(ctx, &balRequest)
	if err != nil {
		response.Status = "error"
		response.Message = "ecountered error from GetUserBalance: " + err.Error()
		// we need to return nil here in order to pass the appropriate status
		// and message to the client. If we return the err the response in
		// the client will be nil.
		return nil
	}

	if balResponse.Data.Balance.Available < req.BaseQuantity {
		response.Status = "fail"
		response.Message = "insufficient available balance: " + baseCurrency
		return nil
	}

	order, error := orderRepo.InsertOrder(service.DB, req)

	switch {
	case error == nil:
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

		if err := service.NewOrderPub.Publish(context.Background(), &orderEvent); err != nil {
			log.Println("publish warning: ", err, orderEvent)
		}

		response.Status = "success"
		response.Data = &pb.UserOrderData{
			Order: order,
		}
		return nil

	default:
		response.Status = "error"
		response.Message = "ecountered error on Insert: " + error.Error()
		return nil
	}
}

// add a sell order
func (service *OrderService) addSellOrder(ctx context.Context, req *pb.OrderRequest, response *pb.OrderResponse) error {
	currencies := strings.Split(req.MarketName, "-")
	currency := currencies[0]

	balRequest := bp.GetUserBalanceRequest{
		UserId:   req.UserId,
		ApiKeyId: req.ApiKeyId,
		Currency: currency,
	}

	balResponse, err := service.Client.GetUserBalance(ctx, &balRequest)
	if err != nil {
		response.Status = "error"
		response.Message = "ecountered error from GetUserBalance: " + err.Error()
		// we need to return nil here in order to pass the appropriate status
		// and message to the client. If we return the err the response in
		// the client will be nil.
		return nil
	}

	if balResponse.Data.Balance.Available < req.CurrencyQuantity {
		response.Status = "fail"
		response.Message = "sell quantity is greater than available balance for " + currency
		return nil
	}

	order, error := orderRepo.InsertOrder(service.DB, req)

	switch {
	case error == nil:
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

		if err := service.NewOrderPub.Publish(context.Background(), &orderEvent); err != nil {
			log.Println("publish warning: ", err, orderEvent)
		}

		response.Status = "success"
		response.Data = &pb.UserOrderData{
			Order: order,
		}
		return nil

	default:
		response.Status = "error"
		response.Message = "ecountered error on Insert: " + error.Error()
		return nil
	}
}

func (service *OrderService) AddOrder(ctx context.Context, req *pb.OrderRequest, response *pb.OrderResponse) error {
	side := enums.NewSideFromString(req.Side)

	switch side {
	case enums.Buy:
		return service.addBuyOrder(ctx, req, response)
	case enums.Sell:
		return service.addSellOrder(ctx, req, response)
	default:
		response.Status = "fail"
		response.Message = "side unknown: " + req.Side
		return nil
	}
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

func (service *OrderService) UpdateOrderStatus(ctx context.Context, req *pb.OrderStatusRequest, res *pb.OrderResponse) error {
	order, error := orderRepo.UpdateOrderStatus(service.DB, req)
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
