package dedicated_server

type DedicatedServer struct {
	Id                  string
	AssetId             string
	SerialNumber        string
	Rack                Rack
	Location            Location
	FeatureAvailability FeatureAvailability
	Contract            Contract
	PowerPorts          Ports
	PrivateNetworks     PrivateNetworks
	NetworkInterfaces   NetworkInterfaces
	Specs               Specs
}

func NewDedicatedServer(
	id,
	assetId,
	serialNumber string,
	rack Rack,
	location Location,
	featureAvailability FeatureAvailability,
	contract Contract,
	powerPorts Ports,
	privateNetworks PrivateNetworks,
	networkInterfaces NetworkInterfaces,
	specs Specs,
) DedicatedServer {
	dedicatedServer := DedicatedServer{
		Id:                  id,
		AssetId:             assetId,
		SerialNumber:        serialNumber,
		Rack:                rack,
		Location:            location,
		FeatureAvailability: featureAvailability,
		Contract:            contract,
		PowerPorts:          powerPorts,
		PrivateNetworks:     privateNetworks,
		NetworkInterfaces:   networkInterfaces,
		Specs:               specs,
	}
	return dedicatedServer
}
