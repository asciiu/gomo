package main

import (
	"context"
	"fmt"

	evt "github.com/asciiu/gomo/common/proto/events"
	pb "github.com/asciiu/gomo/order-service/proto/order"
	micro "github.com/micro/go-micro"
)

// PublishOrder will publish an order event
func PublishOrder(publisher micro.Publisher, order *pb.Order) error {
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
		return fmt.Errorf("could not publish order %s", err.Error())
	}
	return nil
}

// CheckBalance ...
func CheckBalance() error {
	return nil
}
