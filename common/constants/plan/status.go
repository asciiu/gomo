package plan

const (
	// chain status
	// the chain of orders is running
	Active string = "active"
	// the chain of orders has completed
	Completed string = "completed"
	// the chain of orders is not running
	Inactive string = "inactive"
	// the chain failed
	Failed string = "failed"
	// wtf is this, a learning center for ants!
	UnknownChainStatus string = "???"
)
