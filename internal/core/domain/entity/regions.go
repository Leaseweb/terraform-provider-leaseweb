package entity

type Regions []Region

func (r Regions) Contains(region string) bool {
	for _, r := range r {
		if r.Name == region {
			return true
		}
	}

	return false
}
