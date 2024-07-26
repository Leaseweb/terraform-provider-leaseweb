package domain

type Regions []Region

func (r Regions) Contains(region string) bool {
	for _, r := range r {
		if r.Name == region {
			return true
		}
	}

	return false
}

func (r Regions) ToArray() []string {
	var values []string
	for _, region := range r {
		values = append(values, region.String())
	}

	return values
}
