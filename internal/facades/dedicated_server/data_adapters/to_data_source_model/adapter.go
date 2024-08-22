package to_data_source_model

import (
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	domain "github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain/dedicated_server"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/data_sources/dedicated_server/model"
)

func AdaptDedicatedServers(domainDedicatedServers domain.DedicatedServers) model.DedicatedServers {
	var dedicatedServers model.DedicatedServers

	for _, domainDedicatedServer := range domainDedicatedServers {
		dedicatedServer := adaptDedicatedServer(domainDedicatedServer)
		dedicatedServers.DedicatedServers = append(dedicatedServers.DedicatedServers, dedicatedServer)
	}

	return dedicatedServers
}

func adaptDedicatedServer(dedicatedServer domain.DedicatedServer) model.DedicatedServer {
	return model.DedicatedServer{
		Id:                  basetypes.NewStringValue(dedicatedServer.Id),
		AssetId:             basetypes.NewStringValue(dedicatedServer.AssetId),
		SerialNumber:        basetypes.NewStringValue(dedicatedServer.SerialNumber),
		Rack:                adaptRack(dedicatedServer.Rack),
		Location:            adaptLocation(dedicatedServer.Location),
		FeatureAvailability: adaptFeatureAvailability(dedicatedServer.FeatureAvailability),
		Contract:            adaptContract(dedicatedServer.Contract),
		PowerPorts:          adaptPorts(dedicatedServer.PowerPorts),
		PrivateNetworks:     adaptPrivateNetworks(dedicatedServer.PrivateNetworks),
		NetworkInterfaces:   adaptNetworkInterfaces(dedicatedServer.NetworkInterfaces),
		Specs:               adaptSpecs(dedicatedServer.Specs),
	}
}

func adaptRack(rack domain.Rack) model.Rack {
	return model.Rack{
		Id:       basetypes.NewStringValue(rack.Id),
		Capacity: basetypes.NewStringValue(rack.Capacity),
		Type:     basetypes.NewStringValue(rack.Type),
	}
}

func adaptLocation(location domain.Location) model.Location {
	return model.Location{
		Rack:  basetypes.NewStringValue(location.Rack),
		Site:  basetypes.NewStringValue(location.Site),
		Suite: basetypes.NewStringValue(location.Suite),
		Unit:  basetypes.NewStringValue(location.Unit),
	}
}

func adaptFeatureAvailability(featureAvailability domain.FeatureAvailability) model.FeatureAvailability {
	return model.FeatureAvailability{
		Automation:       basetypes.NewBoolValue(featureAvailability.Automation),
		IpmiReboot:       basetypes.NewBoolValue(featureAvailability.IpmiReboot),
		PowerCycle:       basetypes.NewBoolValue(featureAvailability.PowerCycle),
		PrivateNetwork:   basetypes.NewBoolValue(featureAvailability.PrivateNetwork),
		RemoteManagement: basetypes.NewBoolValue(featureAvailability.RemoteManagement),
	}
}

func adaptContract(contract domain.Contract) model.Contract {
	return model.Contract{
		Id:             basetypes.NewStringValue(contract.Id),
		CustomerId:     basetypes.NewStringValue(contract.CustomerId),
		DeliveryStatus: basetypes.NewStringValue(contract.DeliveryStatus),
		Reference:      basetypes.NewStringValue(contract.Reference),
		SalesOrgId:     basetypes.NewStringValue(contract.SalesOrgId),
	}
}

func adaptPrivateNetworks(domainPrivateNetworks domain.PrivateNetworks) model.PrivateNetworks {
	privateNetworks := model.PrivateNetworks{}
	for _, privateNetwork := range domainPrivateNetworks {
		privateNetworks = append(privateNetworks, adaptPrivateNetwork(privateNetwork))
	}
	return privateNetworks
}

func adaptPrivateNetwork(privateNetwork domain.PrivateNetwork) model.PrivateNetwork {
	return model.PrivateNetwork{
		Id:        basetypes.NewStringValue(privateNetwork.Id),
		LinkSpeed: basetypes.NewInt32Value(privateNetwork.LinkSpeed),
		Status:    basetypes.NewStringValue(privateNetwork.Status),
		Subnet:    basetypes.NewStringValue(privateNetwork.Subnet),
		VlanId:    basetypes.NewStringValue(privateNetwork.VlanId),
	}
}

func adaptNetworkInterfaces(networInterfaces domain.NetworkInterfaces) model.NetworkInterfaces {
	return model.NetworkInterfaces{
		Public:           adaptNetworkInterface(networInterfaces.Public),
		Internal:         adaptNetworkInterface(networInterfaces.Internal),
		RemoteManagement: adaptNetworkInterface(networInterfaces.RemoteManagement),
	}
}

func adaptNetworkInterface(networInterface domain.NetworkInterface) model.NetworkInterface {
	return model.NetworkInterface{
		Mac:        basetypes.NewStringValue(networInterface.Mac),
		Ip:         basetypes.NewStringValue(networInterface.Ip),
		Gateway:    basetypes.NewStringValue(networInterface.Gateway),
		NullRouted: basetypes.NewBoolValue(networInterface.NullRouted),
		Ports:      adaptPorts(networInterface.Ports),
		LocationId: basetypes.NewStringValue(networInterface.LocationId),
	}
}

func adaptPorts(domainPorts domain.Ports) model.Ports {
	ports := model.Ports{}
	for _, port := range domainPorts {
		ports = append(ports, adaptPort(port))
	}
	return ports
}

func adaptPort(port domain.Port) model.Port {
	return model.Port{
		Name: basetypes.NewStringValue(port.Name),
		Port: basetypes.NewStringValue(port.Port),
	}
}

func adaptSpecs(specs domain.Specs) model.Specs {
	return model.Specs{
		Chassis:             basetypes.NewStringValue(specs.Chassis),
		HardwareRaidCapable: basetypes.NewBoolValue(specs.HardwareRaidCapable),
		Cpu:                 adaptCpu(specs.Cpu),
		Ram:                 adaptRam(specs.Ram),
		Hdds:                adaptHdds(specs.Hdds),
		PciCards:            adaptPciCards(specs.PciCards),
	}
}

func adaptCpu(cpu domain.Cpu) model.Cpu {
	return model.Cpu{
		Quantity: basetypes.NewInt32Value(cpu.Quantity),
		Type:     basetypes.NewStringValue(cpu.Type),
	}
}

func adaptRam(ram domain.Ram) model.Ram {
	return model.Ram{
		Size: basetypes.NewInt32Value(ram.Size),
		Unit: basetypes.NewStringValue(ram.Unit),
	}
}

func adaptHdds(domainHdds domain.Hdds) model.Hdds {
	hdds := model.Hdds{}
	for _, hdd := range domainHdds {
		hdds = append(hdds, adaptHdd(hdd))
	}
	return hdds
}

func adaptHdd(hdd domain.Hdd) model.Hdd {
	return model.Hdd{
		Id:              basetypes.NewStringValue(hdd.Id),
		Amount:          basetypes.NewInt32Value(hdd.Amount),
		Size:            basetypes.NewFloat32Value(hdd.Size),
		Type:            basetypes.NewStringValue(hdd.Type),
		Unit:            basetypes.NewStringValue(hdd.Unit),
		PerformanceType: basetypes.NewStringValue(hdd.PerformanceType),
	}
}

func adaptPciCards(domainPciCards domain.PciCards) model.PciCards {
	pciCards := model.PciCards{}
	for _, domainPciCard := range domainPciCards {
		pciCards = append(pciCards, adaptPciCard(domainPciCard))
	}
	return pciCards
}

func adaptPciCard(pciCard domain.PciCard) model.PciCard {
	return model.PciCard{
		Description: basetypes.NewStringValue(pciCard.String()),
	}
}
