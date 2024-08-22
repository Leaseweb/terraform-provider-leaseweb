package dedicated_server

type Location struct {
	Rack  string
	Site  string
	Suite string
	Unit  string
}

func NewLocation(rack, site, suite, unit string) Location {
	return Location{
		Rack:  rack,
		Site:  site,
		Suite: suite,
		Unit:  unit,
	}
}
