package dedicated_server

type OperatingSystem struct {
	Id   string
	Name string
}

func NewOperatingSystem(id string, name string) OperatingSystem {
	return OperatingSystem{Id: id, Name: name}
}
