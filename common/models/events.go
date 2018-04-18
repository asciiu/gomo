package models

type ExchangeEvent struct {
	Type       string
	EventTime  int
	MarketName string
	TradeID    int
	Price      string
	Quantity   string
	TradeTime  int
}
