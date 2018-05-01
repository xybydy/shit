package models

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
