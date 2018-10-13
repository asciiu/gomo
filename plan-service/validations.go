package main

import (
	constPlan "github.com/asciiu/gomo/plan-service/constants"
	protoOrder "github.com/asciiu/gomo/plan-service/proto/order"
	protoPlan "github.com/asciiu/gomo/plan-service/proto/plan"
	"github.com/google/uuid"
)

// plans must have a title
func ValidateTitle(title string) bool {
	return title != ""
}

// validates user specified plan status
func ValidatePlanInputStatus(pstatus string) bool {
	pistats := [...]string{
		constPlan.Active,
		constPlan.Inactive,
		constPlan.Historic,
	}

	for _, stat := range pistats {
		if stat == pstatus {
			return true
		}
	}
	return false
}

// defines valid input for plan status when updating an executed plan (a.k.a. plan with a filled order)
func ValidateUpdatePlanStatus(pstatus string) bool {
	pistats := [...]string{
		constPlan.Active,
		constPlan.Inactive,
		constPlan.Closed,
	}

	for _, stat := range pistats {
		if stat == pstatus {
			return true
		}
	}
	return false
}

// New plans can only have a single order with parent order number == 0.
func ValidateSingleRootNode(orderRequests []*protoOrder.NewOrderRequest) bool {
	count := 0
	for _, o := range orderRequests {
		if o.ParentOrderID == uuid.Nil.String() || o.ParentOrderID == "" {
			count += 1
		}
	}
	return count == 1
}

// plan must have at least one order
func ValidateMinOrder(orderRequests []*protoOrder.NewOrderRequest) bool {
	return len(orderRequests) > 0
}

func ValidateConnectedRoutesFromParent(parentOrderID string, orderRequests []*protoOrder.NewOrderRequest) bool {
	orderIDs := make([]string, 0, len(orderRequests)+1)
	orderIDs = append(orderIDs, parentOrderID)

	for _, o := range orderRequests {
		orderIDs = append(orderIDs, o.OrderID)
	}

	for _, o := range orderRequests {
		found := false
		// check connected graph
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

// All orders must contain triggers
func ValidateOrderTrigger(orderRequests []*protoOrder.NewOrderRequest) bool {
	for _, o := range orderRequests {
		if len(o.Triggers) == 0 {
			return false
		}
	}

	return true
}

// Child orders must have a valid parent order ID
func ValidateChildNodes(orderRequests []*protoOrder.NewOrderRequest) bool {
	for _, o := range orderRequests {
		if o.ParentOrderID == uuid.Nil.String() {
			return false
		}
	}

	return true
}

// limit node count for new requests to 10
func ValidateNodeCount(orderRequests []*protoOrder.NewOrderRequest) bool {
	return len(orderRequests) <= 10
}

// validate non zero currency balance
func ValidateNoneZeroBalance(planRequest *protoPlan.NewPlanRequest) bool {
	if planRequest.CommittedCurrencyAmount > 0 || planRequest.Orders[0].InitialCurrencyBalance > 0 {
		return true
	}
	return false
}

func ValidatePaperOrders(orderRequests []*protoOrder.NewOrderRequest) bool {
	for _, o := range orderRequests {
		if o.OrderType != constPlan.PaperOrder {
			return false
		}
	}
	return true
}

func ValidateNotPaperOrders(orderRequests []*protoOrder.NewOrderRequest) bool {
	for _, o := range orderRequests {
		if o.OrderType == constPlan.PaperOrder {
			return false
		}
	}
	return true
}

func ValidateOrderType(ot string) bool {
	ots := [...]string{
		constPlan.LimitOrder,
		constPlan.MarketOrder,
	}

	for _, ty := range ots {
		if ty == ot {
			return true
		}
	}
	return false
}

func ValidateOrderSide(os string) bool {
	ots := [...]string{
		constPlan.Buy,
		constPlan.Sell,
	}

	for _, ty := range ots {
		if ty == os {
			return true
		}
	}
	return false
}

// New order request cannot overwrite a filled order
func ValidateNonExecutedOrder(porders []*protoOrder.Order, rorders []*protoOrder.NewOrderRequest) bool {
	for _, nor := range rorders {
		for _, po := range porders {
			if nor.OrderID == po.OrderID && po.Status == constPlan.Filled {
				return false
			}
		}
	}
	return true
}
