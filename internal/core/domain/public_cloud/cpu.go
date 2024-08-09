package public_cloud

type Cpu struct {
	Value int
	Unit  string
}

func NewCpu(value int, unit string) Cpu {
	return Cpu{Value: value, Unit: unit}
}
