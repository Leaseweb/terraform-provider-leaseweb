package entity

type Memory struct {
	Value float64
	Unit  string
}

func NewMemory(value float64, unit string) Memory {
	return Memory{Value: value, Unit: unit}
}
