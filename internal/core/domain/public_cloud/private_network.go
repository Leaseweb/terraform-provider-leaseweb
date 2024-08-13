package public_cloud

type PrivateNetwork struct {
	Id     string
	Status string
	Subnet string
}

func NewPrivateNetwork(
	id string,
	status string,
	subnet string,
) PrivateNetwork {
	return PrivateNetwork{Id: id, Status: status, Subnet: subnet}
}
