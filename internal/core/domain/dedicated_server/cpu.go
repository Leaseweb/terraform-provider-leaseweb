package dedicated_server

type Cpu struct {
	Quantity int32
	Type     string
}

func NewCpu(quantity int32, Type string) Cpu {
	return Cpu{
		Quantity: quantity,
		Type:     Type,
	}
}
