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
	// past plan that user did outside of the platform
	Past string = "past"
	// wtf is this, a learning center for ants!
	UnknownChainStatus string = "???"
)

// validates user specified plan status
func ValidatePlanInputStatus(pstatus string) bool {
	pistats := [...]string{
		Active,
		Inactive,
		Past,
	}

	for _, stat := range pistats {
		if stat == pstatus {
			return true
		}
	}
	return false
}
