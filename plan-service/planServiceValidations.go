package main

import (
	orderConstants "github.com/asciiu/gomo/common/constants/order"
	planConstants "github.com/asciiu/gomo/common/constants/plan"
	sideConstants "github.com/asciiu/gomo/common/constants/side"
	protoOrder "github.com/asciiu/gomo/plan-service/proto/order"
	"github.com/google/uuid"
)

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
func ValidateNoneZeroBalance(orderRequests []*protoOrder.NewOrderRequest) bool {
	for _, o := range orderRequests {
		if o.ParentOrderID == uuid.Nil.String() && o.ActiveCurrencyBalance > 0 {
			return true
		}
	}
	return false
}

// func ValidateNoneZeroBalance(orderRequests []*protoOrder.UpdateOrderRequest) bool {
// 	switch {
// 	case len(orderRequests) == 0:
// 		return false
// 	case orderRequests[0].ActiveCurrencyBalance <= 0.0:
// 		return false
// 	default:
// 		return true
// 	}
// }

func ValidatePaperOrders(orderRequests []*protoOrder.NewOrderRequest) bool {
	for _, o := range orderRequests {
		if o.OrderType != orderConstants.PaperOrder {
			return false
		}
	}
	return true
}

func ValidateNotPaperOrders(orderRequests []*protoOrder.NewOrderRequest) bool {
	for _, o := range orderRequests {
		if o.OrderType == orderConstants.PaperOrder {
			return false
		}
	}
	return true
}

func ValidateOrderType(ot string) bool {
	ots := [...]string{
		orderConstants.LimitOrder,
		orderConstants.MarketOrder,
		orderConstants.PaperOrder,
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
		sideConstants.Buy,
		sideConstants.Sell,
	}

	for _, ty := range ots {
		if ty == os {
			return true
		}
	}
	return false
}

// validates user specified plan status
func ValidatePlanInputStatus(pstatus string) bool {
	pistats := [...]string{
		planConstants.Active,
		planConstants.Inactive,
		planConstants.Historic,
	}

	for _, stat := range pistats {
		if stat == pstatus {
			return true
		}
	}
	return false
}

// defines valid input for plan status when updating an executed plan (a.k.a. plan with a filled order)
func ValidatePlanUpdateStatus(pstatus string) bool {
	pistats := [...]string{
		planConstants.Active,
		planConstants.Inactive,
	}

	for _, stat := range pistats {
		if stat == pstatus {
			return true
		}
	}
	return false
}
