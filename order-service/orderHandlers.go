package main

import (
	"context"
	"database/sql"
	"fmt"

	evt "github.com/asciiu/gomo/common/proto/events"
)

// OrderFilledReceiver handles order filled events
type OrderFilledReceiver struct {
	DB *sql.DB
}

// ProcessEvent handles OrderEvents. These events are published by when an order was filled.
func (receiver *OrderFilledReceiver) ProcessEvent(ctx context.Context, orderEvent *evt.OrderEvent) error {
	fmt.Println(orderEvent)

	return nil
}
