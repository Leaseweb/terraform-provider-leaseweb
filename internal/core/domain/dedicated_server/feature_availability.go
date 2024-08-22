package dedicated_server

type FeatureAvailability struct {
	Automation       bool
	IpmiReboot       bool
	PowerCycle       bool
	PrivateNetwork   bool
	RemoteManagement bool
}

func NewFeatureAvailability(automation, ipmiReboot, powerCycle, privateNetwork, remoteManagement bool) FeatureAvailability {
	return FeatureAvailability{
		Automation:       automation,
		IpmiReboot:       ipmiReboot,
		PowerCycle:       powerCycle,
		PrivateNetwork:   privateNetwork,
		RemoteManagement: remoteManagement,
	}
}
