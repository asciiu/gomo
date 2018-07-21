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
func ValidateSingleRootNode(orderRequests []*protoOrder.OrderRequest) bool {
	count := 0
	for _, o := range orderRequests {
		if o.ParentOrderNumber == 0 {
			count += 1
		}
	}
	return count == 1
}

// Validate connected tree
func ValidateConnectedRoutes(orderRequests []*protoOrder.OrderRequest) bool {
	orderNumbers := make([]uint32, 0, len(orderRequests))
	orderNumbers = append(orderNumbers, 0)

	for _, o := range orderRequests {
		orderNumbers = append(orderNumbers, o.OrderNumber)
	}

	found := false
	for _, o := range orderRequests {
		for _, n := range orderNumbers {
			if o.ParentOrderNumber == n {
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
