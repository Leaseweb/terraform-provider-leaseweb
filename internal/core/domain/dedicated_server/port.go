package dedicated_server

type Port struct {
	Name string
	Port string
}

func NewPort(name, port string) Port {
	return Port{
		Name: name,
		Port: port,
	}
}
