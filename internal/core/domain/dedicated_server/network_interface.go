package dedicated_server

type NetworkInterface struct {
	Mac        string
	Ip         string
	Gateway    string
	NullRouted bool
	Ports      Ports
	LocationId string
}

func NewNetworkInterface(mac, ip, gateway, locationId string, nullRouted bool, ports Ports) NetworkInterface {
	return NetworkInterface{
		Mac:        mac,
		Ip:         ip,
		Gateway:    gateway,
		LocationId: locationId,
		NullRouted: nullRouted,
		Ports:      ports,
	}
}
