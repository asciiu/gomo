package enums

type OrderType int

const (
	BuyOrder         OrderType = 0
	SellOrder        OrderType = 1
	PaperBuyOrder    OrderType = 2
	PaperSellOrder   OrderType = 3
	UnknownOrderType OrderType = 4
)

func OrderTypeFromString(ex string) OrderType {
	switch ex {
	case "BuyOrder":
		return BuyOrder
	case "SellOrder":
		return SellOrder
	case "PaperBuyOrder":
		return PaperBuyOrder
	case "PaperSellOrder":
		return PaperSellOrder
	default:
		return UnknownOrderType
	}
}

func (ot OrderType) String() string {
	// ... operator counts how many items in the array
	names := [...]string{
		"BuyOrder",
		"SellOrder",
		"PaperBuyOrder",
		"PaperSellOrder",
		"UnknownOrderType",
	}
	return names[ot]
}
