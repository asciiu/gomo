package main

import (
	"context"
	"database/sql"

	bp "github.com/asciiu/gomo/balance-service/proto/balance"
	orderRepo "github.com/asciiu/gomo/order-service/db/sql"
	pb "github.com/asciiu/gomo/order-service/proto/order"
)

type OrderService struct {
	DB     *sql.DB
	Client bp.BalanceServiceClient
}

func (service *OrderService) AddOrder(ctx context.Context, req *pb.OrderRequest, res *pb.OrderResponse) error {

	// is there enough balance
	//balRequest := bp.GetUserBalanceRequest{
	//	UserId:   req.UserId,
	//	ApiKeyId: req.ApiKeyId,
	//	Currency: req.MarketName,
	//}

	//balance, err := service.Client.GetUserBalance(ctx, &balRequest)
	//if err != nil {
	//	res.Status = "error"
	//	res.Message = err.Error()
	//	return err
	//}

	//// do we have enough of base currency to make this order?
	//if balance.Data.Balance.Available >= req.

	order, error := orderRepo.InsertOrder(service.DB, req)

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
