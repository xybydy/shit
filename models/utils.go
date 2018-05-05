package models

// Helper function of models package.
// Checks if given workstation is in given workstation list.
// reverse parameter make the function acts reversed.
func IsIn(a *Workstation, list Workstations, reverse bool) bool {
	for _, b := range list {
		if b == a {
			if reverse {
				return false
			}
			return true
		}
	}
	if reverse {
		return true
	}
	return false
}
