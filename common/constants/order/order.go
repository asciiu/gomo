package order

const (
	// order types here
	LimitOrder       string = "limit"
	MarketOrder      string = "market"
	PaperOrder       string = "paper"
	UnknownOrderType string = "???"
)

func ValidateOrderType(ot string) bool {
	ots := [...]string{
		LimitOrder,
		MarketOrder,
		PaperOrder,
	}

	for _, ty := range ots {
		if ty == ot {
			return true
		}
	}
	return false
}
