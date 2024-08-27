package dedicated_server

type Hdd struct {
	Id              string
	Amount          int32
	Size            float32
	Type            string
	Unit            string
	PerformanceType string
}

func NewHdd(id, Type, unit, performanceType string, amount int32, size float32) Hdd {
	return Hdd{
		Id:              id,
		Amount:          amount,
		Size:            size,
		Type:            Type,
		Unit:            unit,
		PerformanceType: performanceType,
	}
}
