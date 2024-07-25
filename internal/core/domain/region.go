package domain

type Region struct {
	Name     string
	Location string
}

func (r Region) String() string {
	return r.Name
}

func NewRegion(name string, location string) Region {
	return Region{Name: name, Location: location}
}
