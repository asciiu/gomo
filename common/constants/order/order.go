package order

const (
	// order types here
	LimitOrder       string = "limit"
	MarketOrder      string = "market"
	VirtualOrder     string = "paper"
	UnknownOrderType string = "???"
)

func ValidateOrderType(ot string) bool {
	ots := [...]string{
		LimitOrder,
		MarketOrder,
		VirtualOrder,
	}

	for _, ty := range ots {
		if ty == ot {
			return true
		}
	}
	return false
}
