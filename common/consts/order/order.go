package order

const (
	// order types here
	BuyOrder         string = "buy"
	SellOrder        string = "sell"
	VirtualBuyOrder  string = "vbuy"
	VirutalSellOrder string = "vsell"
	UnknownOrderType string = "???"
)

func ValidateOrderType(ot string) bool {
	ots := [...]string{
		BuyOrder,
		SellOrder,
		VirtualBuyOrder,
		VirutalSellOrder,
	}

	for _, ty := range ots {
		if ty == ot {
			return true
		}
	}
	return false
}
