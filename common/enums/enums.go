package enums

type Exchange int

const (
	Binance  Exchange = 0
	Bittrex  Exchange = 1
	Poloniex Exchange = 2
	Unknown  Exchange = 3
)

func NewExchange(ex string) Exchange {
	switch ex {
	case "Binance":
		return Binance
	case "Bittrex":
		return Bittrex
	case "Poloniex":
		return Poloniex
	default:
		return Unknown
	}
}

func (exchange Exchange) String() string {
	// ... operator counts how many items in the array
	names := [...]string{
		"Binance",
		"Bittrex",
		"Poloniex",
		"Cryptopia"}
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
