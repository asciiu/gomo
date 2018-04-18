package main

import (
	"context"
	"database/sql"
	"log"
	"strings"

	bp "github.com/asciiu/gomo/balance-service/proto/balance"
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

func (service *OrderService) AddOrder(ctx context.Context, req *pb.OrderRequest, response *pb.OrderResponse) error {
	// we will always assume the market trading pairs will be
	// the currency-base currency: e.g. ADA-BTC
	baseCurrency := strings.Split(req.MarketName, "-")[1]

	// is there enough balance
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
		response.Message = "insufficient balance " + baseCurrency
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
