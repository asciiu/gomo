package plan

const (
	// the plan was aborted
	Aborted string = "aborted"
	// chain status
	// the chain of orders is running
	Active string = "active"
	// the chain of orders has completed
	Completed string = "completed"
	// this means the record was successfully delete from our DB
	Deleted string = "deleted"
	// the chain of orders is not running
	Inactive string = "inactive"
	// the chain failed
	Failed string = "failed"
	// past plan that user did outside of the platform
	Historic string = "historic"
	// wtf is this, a learning center for ants!
	UnknownChainStatus string = "???"
)

// validates user specified plan status
func ValidatePlanInputStatus(pstatus string) bool {
	pistats := [...]string{
		Active,
		Inactive,
		Historic,
	}

	for _, stat := range pistats {
		if stat == pstatus {
			return true
		}
	}
	return false
}

// defines valid input for plan status when updating an executed plan (a.k.a. plan with a filled order)
func ValidateUpdatePlanStatus(pstatus string) bool {
	pistats := [...]string{
		Active,
		Inactive,
	}

	for _, stat := range pistats {
		if stat == pstatus {
			return true
		}
	}
	return false
}
