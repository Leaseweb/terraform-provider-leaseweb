// Package to_domain_entity implements adapters to convert dedicated_server sdk models to domain entities.
package to_domain_entity

import (
	"github.com/leaseweb/leaseweb-go-sdk/dedicatedServer"
	domain "github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain/dedicated_server"
)

// AdaptDedicatedServer adapts dedicatedServer.Server to dedicated_server.DedicatedServer.
func AdaptDedicatedServer(sdkDedicatedServer dedicatedServer.Server) domain.DedicatedServer {
	return domain.NewDedicatedServer(
		sdkDedicatedServer.GetId(),
		sdkDedicatedServer.GetAssetId(),
		sdkDedicatedServer.GetSerialNumber(),
		adaptRack(sdkDedicatedServer.GetRack()),
		adaptLocation(sdkDedicatedServer.GetLocation()),
		adaptFeatureAvailability(sdkDedicatedServer.GetFeatureAvailability()),
		adaptContract(sdkDedicatedServer.GetContract()),
		adaptPorts(sdkDedicatedServer.GetPowerPorts()),
		adaptPrivateNetworks(sdkDedicatedServer.GetPrivateNetworks()),
		adaptNetworkInterfaces(sdkDedicatedServer.GetNetworkInterfaces()),
		adaptSpecs(sdkDedicatedServer.GetSpecs()),
	)
}

func adaptRack(sdkRack dedicatedServer.Rack) domain.Rack {
	return domain.NewRack(
		sdkRack.GetId(),
		sdkRack.GetCapacity(),
		sdkRack.GetType(),
	)
}

func adaptLocation(sdkLocation dedicatedServer.Location) domain.Location {
	return domain.NewLocation(
		sdkLocation.GetRack(),
		sdkLocation.GetSite(),
		sdkLocation.GetSuite(),
		sdkLocation.GetUnit(),
	)
}

func adaptFeatureAvailability(sdkFeatureAvailability dedicatedServer.FeatureAvailability) domain.FeatureAvailability {
	return domain.NewFeatureAvailability(
		sdkFeatureAvailability.GetAutomation(),
		sdkFeatureAvailability.GetIpmiReboot(),
		sdkFeatureAvailability.GetPowerCycle(),
		sdkFeatureAvailability.GetPrivateNetwork(),
		sdkFeatureAvailability.GetRemoteManagement(),
	)
}

func adaptContract(sdkContract dedicatedServer.Contract) domain.Contract {
	return domain.NewContract(
		sdkContract.GetId(),
		sdkContract.GetCustomerId(),
		sdkContract.GetDeliveryStatus(),
		sdkContract.GetReference(),
		sdkContract.GetSalesOrgId(),
	)
}

func adaptPorts(sdkPorts []dedicatedServer.Port) domain.Ports {
	ports := domain.Ports{}
	for _, port := range sdkPorts {
		ports = append(ports, adaptPort(port))
	}
	return ports
}

func adaptPort(sdkPort dedicatedServer.Port) domain.Port {
	return domain.NewPort(
		sdkPort.GetName(),
		sdkPort.GetPort(),
	)
}

func adaptPrivateNetworks(sdkPrivateNetwork []dedicatedServer.PrivateNetwork) domain.PrivateNetworks {
	privateNetworks := domain.PrivateNetworks{}
	for _, privateNetwork := range sdkPrivateNetwork {
		privateNetworks = append(privateNetworks, adaptPrivateNetwork(privateNetwork))
	}
	return privateNetworks
}

func adaptPrivateNetwork(sdkPrivateNetwork dedicatedServer.PrivateNetwork) domain.PrivateNetwork {
	return domain.NewPrivateNetwork(
		sdkPrivateNetwork.GetId(),
		sdkPrivateNetwork.GetStatus(),
		sdkPrivateNetwork.GetSubnet(),
		sdkPrivateNetwork.GetVlanId(),
		int32(sdkPrivateNetwork.GetLinkSpeed()),
	)
}

func adaptNetworkInterfaces(sdkNetworkInterfaces dedicatedServer.NetworkInterfaces) domain.NetworkInterfaces {
	return domain.NewNetworkInterfaces(
		adaptNetworkInterface(sdkNetworkInterfaces.GetPublic()),
		adaptNetworkInterface(sdkNetworkInterfaces.GetInternal()),
		adaptNetworkInterface(sdkNetworkInterfaces.GetRemoteManagement()),
	)
}

func adaptNetworkInterface(sdkNetworkInterface dedicatedServer.NetworkInterface) domain.NetworkInterface {
	return domain.NewNetworkInterface(
		sdkNetworkInterface.GetMac(),
		sdkNetworkInterface.GetIp(),
		sdkNetworkInterface.GetGateway(),
		sdkNetworkInterface.GetLocationId(),
		sdkNetworkInterface.GetNullRouted(),
		adaptPorts(sdkNetworkInterface.GetPorts()),
	)
}

func adaptSpecs(sdkServerSpecs dedicatedServer.ServerSpecs) domain.Specs {
	return domain.NewSpecs(
		sdkServerSpecs.GetChassis(),
		sdkServerSpecs.GetHardwareRaidCapable(),
		adaptCpu(sdkServerSpecs.GetCpu()),
		adaptRam(sdkServerSpecs.GetRam()),
		adaptHdds(sdkServerSpecs.GetHdd()),
		adaptPciCards(sdkServerSpecs.GetPciCards()),
	)
}

func adaptRam(sdkRam dedicatedServer.Ram) domain.Ram {
	return domain.NewRam(
		sdkRam.GetSize(),
		sdkRam.GetUnit(),
	)
}

func adaptPciCards(sdkPciCards []dedicatedServer.PciCard) domain.PciCards {
	pciCards := domain.PciCards{}
	for _, pciCard := range sdkPciCards {
		pciCards = append(pciCards, adaptPciCard(pciCard))
	}
	return pciCards
}

func adaptPciCard(sdkPciCard dedicatedServer.PciCard) domain.PciCard {
	return domain.NewPciCard(
		sdkPciCard.GetDescription(),
	)
}

func adaptCpu(sdkCpu dedicatedServer.Cpu) domain.Cpu {
	return domain.NewCpu(
		sdkCpu.GetQuantity(),
		sdkCpu.GetType(),
	)
}

func adaptHdds(sdkHdds []dedicatedServer.Hdd) domain.Hdds {
	hdds := domain.Hdds{}
	for _, hdd := range sdkHdds {
		hdds = append(hdds, adaptHdd(hdd))
	}
	return hdds
}

func adaptHdd(sdkHdd dedicatedServer.Hdd) domain.Hdd {
	return domain.NewHdd(
		sdkHdd.GetId(),
		sdkHdd.GetType(),
		sdkHdd.GetUnit(),
		sdkHdd.GetPerformanceType(),
		sdkHdd.GetAmount(),
		sdkHdd.GetSize(),
	)
}

// AdaptOperatingSystem adapts []dedicatedServer.OperatingSystem to domain.OperatingSystems.
func AdaptOperatingSystems(sdkOs []dedicatedServer.OperatingSystem) (domainOs domain.OperatingSystems) {
	for _, os := range sdkOs {
		domainOs = append(domainOs, adaptOperatingSystem(os))
	}
	return
}

func adaptOperatingSystem(sdkOs dedicatedServer.OperatingSystem) domain.OperatingSystem {
	return domain.NewOperatingSystem(
		sdkOs.GetId(),
		sdkOs.GetName(),
	)
}

// AdaptControlPanels adapts dedicated_server.ControlPanels to dedicatedServer.ControlPanel[].
func AdaptControlPanels(sdkControlPanels []dedicatedServer.ControlPanel) domain.ControlPanels {
	controlPanels := domain.ControlPanels{}
	for _, sdkControlPanel := range sdkControlPanels {
		controlPanels = append(controlPanels, adaptControlPanel(sdkControlPanel))
	}
	return controlPanels
}

func adaptControlPanel(sdkControlPanel dedicatedServer.ControlPanel) domain.ControlPanel {
	return domain.NewControlPanel(
		sdkControlPanel.GetId(),
		sdkControlPanel.GetName(),
	)
}
