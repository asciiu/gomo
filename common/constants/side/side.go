package side

const (
	// order sides
	Buy         string = "buy"
	Sell        string = "sell"
	UnknownSide string = "???"
)

func ValidateSide(ot string) bool {
	ots := [...]string{
		Buy,
		Sell,
	}

	for _, ty := range ots {
		if ty == ot {
			return true
		}
	}
	return false
}
