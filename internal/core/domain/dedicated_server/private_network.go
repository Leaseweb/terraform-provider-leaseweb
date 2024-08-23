package dedicated_server

type PrivateNetwork struct {
	Id        string
	LinkSpeed int32
	Status    string
	Subnet    string
	VlanId    string
}

func NewPrivateNetwork(id, status, subnet, vlanId string, linkSpeed int32) PrivateNetwork {
	return PrivateNetwork{
		Id:        id,
		LinkSpeed: linkSpeed,
		Status:    status,
		Subnet:    subnet,
		VlanId:    vlanId,
	}
}
