package entity

type Region struct {
	Name     string
	Location string
}

func NewRegion(name string, location string) Region {
	return Region{Name: name, Location: location}
}
