package main

import (
	protoOrder "github.com/asciiu/gomo/plan-service/proto/order"
)

// MinBalance needed to submit order
const MinBalance = 0.00001000

const (
	// order types here
	Delete string = "delete"
	New    string = "new"
	Update string = "update"
)

func ValidateOrderAction(ot string) bool {
	ots := [...]string{
		Delete,
		New,
		Update,
	}

	for _, ty := range ots {
		if ty == ot {
			return true
		}
	}
	return false
}

// New plans can only have a single order with parent order number == 0.
func ValidateSingleRootNode(orderRequests []*protoOrder.NewOrderRequest) bool {
	count := 0
	for _, o := range orderRequests {
		if o.ParentOrderID == "00000000-0000-0000-0000-000000000000" {
			count += 1
		}
	}
	return count == 1
}

// Validate connected tree when tested from root node
func ValidateConnectedRoutes(orderRequests []*protoOrder.NewOrderRequest) bool {
	orderIDs := make([]string, 0, len(orderRequests))
	orderIDs = append(orderIDs, "00000000-0000-0000-0000-000000000000")

	for _, o := range orderRequests {
		orderIDs = append(orderIDs, o.OrderID)
	}

	found := false
	for _, o := range orderRequests {
		for _, n := range orderIDs {
			if o.ParentOrderID == n {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// limit node count for new requests to 10
func ValidateNodeCount(orderRequests []*protoOrder.NewOrderRequest) bool {
	return len(orderRequests) <= 10
}
