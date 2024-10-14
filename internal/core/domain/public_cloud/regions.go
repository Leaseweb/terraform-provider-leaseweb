package public_cloud

type Regions []string

func (r Regions) Contains(region string) bool {
	for _, r := range r {
		if r == region {
			return true
		}
	}

	return false
}
