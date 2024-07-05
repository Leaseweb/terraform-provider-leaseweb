package entity

type Cpu struct {
	Value int64
	Unit  string
}

func NewCpu(value int64, unit string) Cpu {
	return Cpu{Value: value, Unit: unit}
}
