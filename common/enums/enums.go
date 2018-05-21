package enums

type Exchange int
type Side int

const (
	Binance         Exchange = 0
	Bittrex         Exchange = 1
	Poloniex        Exchange = 2
	Kucoin          Exchange = 3
	UnknownExchange Exchange = 4
)

const (
	Buy         Side = 0
	Sell        Side = 1
	UnknownSide Side = 2
)

func NewExchange(ex string) Exchange {
	switch ex {
	case "binance":
		return Binance
	case "bittrex":
		return Bittrex
	case "kucoin":
		return Kucoin
	case "poloniex":
		return Poloniex
	default:
		return UnknownExchange
	}
}

func (exchange Exchange) String() string {
	// ... operator counts how many items in the array
	names := [...]string{
		"binance",
		"bittrex",
		"kucoin",
		"poloniex"}
	// → `day`: It's one of the
	// values of Weekday constants.
	// If the constant is Sunday,
	// then day is 0.
	// prevent panicking in case of
	// `day` is out of range of Weekday
	if exchange < Binance || exchange > Poloniex {
		return "Unknown"
	}
	// return the name of a Weekday
	// constant from the names array
	// above.
	return names[exchange]
}

func NewSideFromString(side string) Side {
	switch side {
	case "buy":
		return Buy
	case "sell":
		return Sell
	default:
		return UnknownSide
	}
}
