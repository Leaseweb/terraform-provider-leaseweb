package domain

type NetworkSpeed struct {
	Value int
	Unit  string
}

func NewNetworkSpeed(value int, unit string) NetworkSpeed {
	return NetworkSpeed{Value: value, Unit: unit}
}
