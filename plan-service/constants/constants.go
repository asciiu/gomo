package constants

const (
	// plan or order may be active
	Active string = "active"
	// plan is closed means cannot be updated
	Closed string = "closed"
	// soft delete plans
	Deleted string = "deleted"
	// plan or order is not running
	Inactive string = "inactive"
	// historic plan was never mean't to run
	Historic string = "historic"

	// order can be aborted status
	Aborted string = "aborted"
	// order failed to execute due to exchange error
	Failed string = "failed"
	// order was successfully completed
	Filled string = "filled"

	// order sides
	Buy  string = "buy"
	Sell string = "sell"

	// a limit order requires a price
	LimitOrder string = "limit"
	// fill quantity immediately
	MarketOrder string = "market"
	// virtual never executed
	PaperOrder string = "paper"
)
