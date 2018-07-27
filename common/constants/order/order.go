package order

const (
	// order types here
	LimitOrder       string = "limit"
	MarketOrder      string = "market"
	PaperOrder       string = "paper"
	UnknownOrderType string = "???"

	NewOrder    string = "new"
	DeleteOrder string = "delete"
	UpdateOrder string = "update"
	Unchanged   string = "unchanged"
)
