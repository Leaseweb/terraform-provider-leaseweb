package entity

type NetworkSpeed struct {
	Value int64
	Unit  string
}

func NewNetworkSpeed(value int64, unit string) NetworkSpeed {
	return NetworkSpeed{Value: value, Unit: unit}
}
