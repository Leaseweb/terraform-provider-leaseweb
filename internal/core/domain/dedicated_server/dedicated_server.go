package domain

type DedicatedServer struct {
	Id string
}

func NewDedicatedServer(id string) DedicatedServer {
	dedicatedServer := DedicatedServer{Id: id}
	return dedicatedServer
}
