package dedicated_server

type Rack struct {
	Id       string
	Capacity string
	Type     string
}

func NewRack(id, capacity, rackType string) Rack {
	return Rack{
		Id:       id,
		Capacity: capacity,
		Type:     rackType,
	}
}
