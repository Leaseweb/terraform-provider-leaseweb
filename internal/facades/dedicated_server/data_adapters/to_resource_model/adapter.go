package to_resource_model

import (
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	domain "github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain/dedicated_server"
	resourcesModel "github.com/leaseweb/terraform-provider-leaseweb/internal/provider/resources/dedicated_server/model"
)

// AdaptDedicatedServer adapts domain.DedicatedServer to resourcesModel.DedicatedServer.
func AdaptDedicatedServer(dedicatedServer domain.DedicatedServer) resourcesModel.DedicatedServer {
	return resourcesModel.DedicatedServer{
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

func adaptRack(rack domain.Rack) resourcesModel.Rack {
	return resourcesModel.Rack{
		Id:       basetypes.NewStringValue(rack.Id),
		Capacity: basetypes.NewStringValue(rack.Capacity),
		Type:     basetypes.NewStringValue(rack.Type),
	}
}

func adaptLocation(location domain.Location) resourcesModel.Location {
	return resourcesModel.Location{
		Rack:  basetypes.NewStringValue(location.Rack),
		Site:  basetypes.NewStringValue(location.Site),
		Suite: basetypes.NewStringValue(location.Suite),
		Unit:  basetypes.NewStringValue(location.Unit),
	}
}

func adaptFeatureAvailability(featureAvailability domain.FeatureAvailability) resourcesModel.FeatureAvailability {
	return resourcesModel.FeatureAvailability{
		Automation:       basetypes.NewBoolValue(featureAvailability.Automation),
		IpmiReboot:       basetypes.NewBoolValue(featureAvailability.IpmiReboot),
		PowerCycle:       basetypes.NewBoolValue(featureAvailability.PowerCycle),
		PrivateNetwork:   basetypes.NewBoolValue(featureAvailability.PrivateNetwork),
		RemoteManagement: basetypes.NewBoolValue(featureAvailability.RemoteManagement),
	}
}

func adaptContract(contract domain.Contract) resourcesModel.Contract {
	return resourcesModel.Contract{
		Id:             basetypes.NewStringValue(contract.Id),
		CustomerId:     basetypes.NewStringValue(contract.CustomerId),
		DeliveryStatus: basetypes.NewStringValue(contract.DeliveryStatus),
		Reference:      basetypes.NewStringValue(contract.Reference),
		SalesOrgId:     basetypes.NewStringValue(contract.SalesOrgId),
	}
}

func adaptPrivateNetworks(domainPrivateNetworks domain.PrivateNetworks) resourcesModel.PrivateNetworks {
	privateNetworks := resourcesModel.PrivateNetworks{}
	for _, privateNetwork := range domainPrivateNetworks {
		privateNetworks = append(privateNetworks, adaptPrivateNetwork(privateNetwork))
	}
	return privateNetworks
}

func adaptPrivateNetwork(privateNetwork domain.PrivateNetwork) resourcesModel.PrivateNetwork {
	return resourcesModel.PrivateNetwork{
		Id:        basetypes.NewStringValue(privateNetwork.Id),
		LinkSpeed: basetypes.NewInt32Value(privateNetwork.LinkSpeed),
		Status:    basetypes.NewStringValue(privateNetwork.Status),
		Subnet:    basetypes.NewStringValue(privateNetwork.Subnet),
		VlanId:    basetypes.NewStringValue(privateNetwork.VlanId),
	}
}

func adaptNetworkInterfaces(networInterfaces domain.NetworkInterfaces) resourcesModel.NetworkInterfaces {
	return resourcesModel.NetworkInterfaces{
		Public:           adaptNetworkInterface(networInterfaces.Public),
		Internal:         adaptNetworkInterface(networInterfaces.Internal),
		RemoteManagement: adaptNetworkInterface(networInterfaces.RemoteManagement),
	}
}

func adaptNetworkInterface(networInterface domain.NetworkInterface) resourcesModel.NetworkInterface {
	return resourcesModel.NetworkInterface{
		Mac:        basetypes.NewStringValue(networInterface.Mac),
		Ip:         basetypes.NewStringValue(networInterface.Ip),
		Gateway:    basetypes.NewStringValue(networInterface.Gateway),
		NullRouted: basetypes.NewBoolValue(networInterface.NullRouted),
		Ports:      adaptPorts(networInterface.Ports),
		LocationId: basetypes.NewStringValue(networInterface.LocationId),
	}
}

func adaptPorts(domainPorts domain.Ports) resourcesModel.Ports {
	ports := resourcesModel.Ports{}
	for _, port := range domainPorts {
		ports = append(ports, adaptPort(port))
	}
	return ports
}

func adaptPort(port domain.Port) resourcesModel.Port {
	return resourcesModel.Port{
		Name: basetypes.NewStringValue(port.Name),
		Port: basetypes.NewStringValue(port.Port),
	}
}

func adaptSpecs(specs domain.Specs) resourcesModel.Specs {
	return resourcesModel.Specs{
		Chassis:             basetypes.NewStringValue(specs.Chassis),
		HardwareRaidCapable: basetypes.NewBoolValue(specs.HardwareRaidCapable),
		Cpu:                 adaptCpu(specs.Cpu),
		Ram:                 adaptRam(specs.Ram),
		Hdds:                adaptHdds(specs.Hdds),
		PciCards:            adaptPciCards(specs.PciCards),
	}
}

func adaptCpu(cpu domain.Cpu) resourcesModel.Cpu {
	return resourcesModel.Cpu{
		Quantity: basetypes.NewInt32Value(cpu.Quantity),
		Type:     basetypes.NewStringValue(cpu.Type),
	}
}

func adaptRam(ram domain.Ram) resourcesModel.Ram {
	return resourcesModel.Ram{
		Size: basetypes.NewInt32Value(ram.Size),
		Unit: basetypes.NewStringValue(ram.Unit),
	}
}

func adaptHdds(domainHdds domain.Hdds) resourcesModel.Hdds {
	hdds := resourcesModel.Hdds{}
	for _, hdd := range domainHdds {
		hdds = append(hdds, adaptHdd(hdd))
	}
	return hdds
}

func adaptHdd(hdd domain.Hdd) resourcesModel.Hdd {
	return resourcesModel.Hdd{
		Id:              basetypes.NewStringValue(hdd.Id),
		Amount:          basetypes.NewInt32Value(hdd.Amount),
		Size:            basetypes.NewFloat32Value(hdd.Size),
		Type:            basetypes.NewStringValue(hdd.Type),
		Unit:            basetypes.NewStringValue(hdd.Unit),
		PerformanceType: basetypes.NewStringValue(hdd.PerformanceType),
	}
}

func adaptPciCards(domainPciCards domain.PciCards) resourcesModel.PciCards {
	pciCards := resourcesModel.PciCards{}
	for _, domainPciCard := range domainPciCards {
		pciCards = append(pciCards, adaptPciCard(domainPciCard))
	}
	return pciCards
}

func adaptPciCard(pciCard domain.PciCard) resourcesModel.PciCard {
	return resourcesModel.PciCard{
		Description: basetypes.NewStringValue(pciCard.String()),
	}
}

// AdaptOperatingSystems adapts domain.OperatingSystems to resourcesModel.OperatingSystems.
func AdaptOperatingSystems(domainOs domain.OperatingSystems) (modelOs resourcesModel.OperatingSystems) {
	for _, os := range domainOs {
		modelOs.OperatingSystems = append(modelOs.OperatingSystems, adaptOperatingSystem(os))
	}

	return
}

func adaptOperatingSystem(domainOs domain.OperatingSystem) resourcesModel.OperatingSystem {
	return resourcesModel.OperatingSystem{
		Id:   basetypes.NewStringValue(domainOs.Id),
		Name: basetypes.NewStringValue(domainOs.Name),
	}
}

// AdaptControlPanels adapts domain.ControlPanels to resourcesModel.ControlPanels.
func AdaptControlPanels(domainControlPanels domain.ControlPanels) resourcesModel.ControlPanels {
	var controlPanels resourcesModel.ControlPanels

	for _, domainControlPanel := range domainControlPanels {
		controlPanels.ControlPanels = append(controlPanels.ControlPanels, adaptControlPanel(domainControlPanel))
	}

	return controlPanels
}

func adaptControlPanel(domainControlPanel domain.ControlPanel) resourcesModel.ControlPanel {
	return resourcesModel.ControlPanel{
		Id:   basetypes.NewStringValue(domainControlPanel.Id),
		Name: basetypes.NewStringValue(domainControlPanel.Name),
	}
}
