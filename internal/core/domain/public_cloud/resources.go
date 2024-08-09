package public_cloud

type Resources struct {
	Cpu                 Cpu
	Memory              Memory
	PublicNetworkSpeed  NetworkSpeed
	PrivateNetworkSpeed NetworkSpeed
}

func NewResources(
	cpu Cpu,
	memory Memory,
	publicNetworkSpeed NetworkSpeed,
	privateNetworkSpeed NetworkSpeed,
) Resources {
	return Resources{
		Cpu:                 cpu,
		Memory:              memory,
		PublicNetworkSpeed:  publicNetworkSpeed,
		PrivateNetworkSpeed: privateNetworkSpeed,
	}
}
