package enums

import (
	"testing"
)

func TestExchangeModelNew(t *testing.T) {
	exchange := NewExchange("Binance")

	if exchange != Binance {
		t.Errorf("Binance != Binance enum")
	}
}

func TestExchangeModelUnknown(t *testing.T) {
	exchange := NewExchange("Kex")

	if exchange != Unknown {
		t.Errorf("unknown exchange failed")
	}
}

func TestExchangeModelStringConversion(t *testing.T) {
	exchange := NewExchange("Poloniex")

	if exchange.String() != "Poloniex" {
		t.Errorf("to string conversion failed")
	}
}
