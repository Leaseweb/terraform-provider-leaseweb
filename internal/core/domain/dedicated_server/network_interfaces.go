package dedicated_server

type NetworkInterfaces struct {
	Public           NetworkInterface
	Internal         NetworkInterface
	RemoteManagement NetworkInterface
}

func NewNetworkInterfaces(public, internal, remoteManagement NetworkInterface) NetworkInterfaces {
	return NetworkInterfaces{
		Public:           public,
		Internal:         internal,
		RemoteManagement: remoteManagement,
	}
}
