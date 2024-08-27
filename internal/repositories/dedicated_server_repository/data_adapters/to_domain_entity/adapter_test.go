package to_domain_entity

import (
	"testing"

	"github.com/leaseweb/leaseweb-go-sdk/dedicatedServer"
	domain "github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain/dedicated_server"
	"github.com/stretchr/testify/assert"
)

func Test_adaptRack(t *testing.T) {
	rack := dedicatedServer.NewRack()
	rack.SetId("id")
	rack.SetCapacity("cap")
	rack.SetType("type")

	got := adaptRack(*rack)
	want := domain.NewRack("id", "cap", "type")
	assert.Equal(t, want, got)
}

func Test_adaptLocation(t *testing.T) {
	location := dedicatedServer.NewLocation()
	location.SetRack("rack")
	location.SetSite("site")
	location.SetSuite("suite")
	location.SetUnit("unit")

	got := adaptLocation(*location)
	want := domain.NewLocation("rack", "site", "suite", "unit")
	assert.Equal(t, want, got)
}

func Test_adaptFeatureAvailability(t *testing.T) {
	featureAvailability := dedicatedServer.NewFeatureAvailability()
	featureAvailability.SetAutomation(true)
	featureAvailability.SetIpmiReboot(false)
	featureAvailability.SetPowerCycle(true)
	featureAvailability.SetPrivateNetwork(false)
	featureAvailability.SetRemoteManagement(true)

	got := adaptFeatureAvailability(*featureAvailability)
	want := domain.NewFeatureAvailability(true, false, true, false, true)
	assert.Equal(t, want, got)
}

func Test_adaptContract(t *testing.T) {
	contract := dedicatedServer.NewContract()

	contract.SetId("id")
	contract.SetCustomerId("customer_id")
	contract.SetDeliveryStatus("status")
	contract.SetReference("ref")
	contract.SetSalesOrgId("sales")

	got := adaptContract(*contract)
	want := domain.NewContract("id", "customer_id", "status", "ref", "sales")
	assert.Equal(t, want, got)
}

func Test_adaptPorts(t *testing.T) {
	port1 := dedicatedServer.NewPort()
	port1.SetName("name1")
	port1.SetPort("1111")

	port2 := dedicatedServer.NewPort()
	port2.SetName("name2")
	port2.SetPort("2222")

	got := adaptPorts([]dedicatedServer.Port{*port1, *port2})
	want := domain.Ports{
		domain.NewPort("name1", "1111"),
		domain.NewPort("name2", "2222"),
	}
	assert.Equal(t, want, got)
}

func Test_adaptPort(t *testing.T) {
	port := dedicatedServer.NewPort()
	port.SetName("name")
	port.SetPort("1212")

	got := adaptPort(*port)
	want := domain.NewPort("name", "1212")
	assert.Equal(t, want, got)
}

func Test_adaptPrivateNetwork(t *testing.T) {
	privateNetwork := dedicatedServer.NewPrivateNetwork()
	privateNetwork.SetId("id")
	privateNetwork.SetStatus("status")
	privateNetwork.SetSubnet("subnet")
	privateNetwork.SetVlanId("vlanid")
	privateNetwork.SetLinkSpeed(12)

	got := adaptPrivateNetwork(*privateNetwork)
	want := domain.NewPrivateNetwork("id", "status", "subnet", "vlanid", 12)
	assert.Equal(t, want, got)
}

func Test_adaptPrivateNetworks(t *testing.T) {
	privateNetwork1 := dedicatedServer.NewPrivateNetwork()
	privateNetwork1.SetId("id1")
	privateNetwork1.SetStatus("status1")
	privateNetwork1.SetSubnet("subnet1")
	privateNetwork1.SetVlanId("vlanid1")
	privateNetwork1.SetLinkSpeed(1)

	privateNetwork2 := dedicatedServer.NewPrivateNetwork()
	privateNetwork2.SetId("id2")
	privateNetwork2.SetStatus("status2")
	privateNetwork2.SetSubnet("subnet2")
	privateNetwork2.SetVlanId("vlanid2")
	privateNetwork2.SetLinkSpeed(2)

	got := adaptPrivateNetworks([]dedicatedServer.PrivateNetwork{*privateNetwork1, *privateNetwork2})
	want := domain.PrivateNetworks{
		domain.NewPrivateNetwork("id1", "status1", "subnet1", "vlanid1", 1),
		domain.NewPrivateNetwork("id2", "status2", "subnet2", "vlanid2", 2),
	}
	assert.Equal(t, want, got)
}

func Test_adaptNetworkInterface(t *testing.T) {
	port := dedicatedServer.NewPort()
	port.SetName("name")
	port.SetPort("1111")

	networkInterface := dedicatedServer.NewNetworkInterface()
	networkInterface.SetMac("mac")
	networkInterface.SetIp("ip")
	networkInterface.SetGateway("gateway")
	networkInterface.SetLocationId("id")
	networkInterface.SetNullRouted(true)
	networkInterface.SetPorts([]dedicatedServer.Port{*port})

	got := adaptNetworkInterface(*networkInterface)
	want := domain.NewNetworkInterface(
		"mac",
		"ip",
		"gateway",
		"id",
		true,
		domain.Ports{domain.NewPort("name", "1111")},
	)
	assert.Equal(t, want, got)
}

func Test_adaptNetworkInterfaces(t *testing.T) {

	public := dedicatedServer.NewNetworkInterface()
	public.SetMac("m1")
	public.SetIp("i1")
	public.SetGateway("g1")
	public.SetLocationId("l1")
	public.SetNullRouted(true)
	port1 := dedicatedServer.NewPort()
	port1.SetName("n1")
	port1.SetPort("p1")
	public.SetPorts([]dedicatedServer.Port{*port1})

	internal := dedicatedServer.NewNetworkInterface()
	internal.SetMac("m2")
	internal.SetIp("i2")
	internal.SetGateway("g2")
	internal.SetLocationId("l2")
	internal.SetNullRouted(false)
	port2 := dedicatedServer.NewPort()
	port2.SetName("n2")
	port2.SetPort("p2")
	internal.SetPorts([]dedicatedServer.Port{*port2})

	remote := dedicatedServer.NewNetworkInterface()
	remote.SetMac("m3")
	remote.SetIp("i3")
	remote.SetGateway("g3")
	remote.SetLocationId("l3")
	remote.SetNullRouted(true)
	port3 := dedicatedServer.NewPort()
	port3.SetName("n3")
	port3.SetPort("p3")
	remote.SetPorts([]dedicatedServer.Port{*port3})

	networkInterfaces := dedicatedServer.NewNetworkInterfaces()
	networkInterfaces.SetPublic(*public)
	networkInterfaces.SetInternal(*internal)
	networkInterfaces.SetRemoteManagement(*remote)

	got := adaptNetworkInterfaces(*networkInterfaces)
	want := domain.NewNetworkInterfaces(
		domain.NewNetworkInterface("m1", "i1", "g1", "l1", true, domain.Ports{domain.NewPort("n1", "p1")}),
		domain.NewNetworkInterface("m2", "i2", "g2", "l2", false, domain.Ports{domain.NewPort("n2", "p2")}),
		domain.NewNetworkInterface("m3", "i3", "g3", "l3", true, domain.Ports{domain.NewPort("n3", "p3")}),
	)
	assert.Equal(t, want, got)
}

func Test_adaptRam(t *testing.T) {
	ram := dedicatedServer.NewRam()
	ram.SetSize(12)
	ram.SetUnit("unit")

	got := adaptRam(*ram)
	want := domain.NewRam(12, "unit")
	assert.Equal(t, want, got)
}

func Test_adaptCpu(t *testing.T) {
	cpu := dedicatedServer.NewCpu()
	cpu.SetQuantity(1)
	cpu.SetType("type")

	got := adaptCpu(*cpu)
	want := domain.NewCpu(1, "type")
	assert.Equal(t, want, got)
}

func Test_adaptPciCard(t *testing.T) {
	pciCard := dedicatedServer.NewPciCard()
	pciCard.SetDescription("description")

	got := adaptPciCard(*pciCard)
	want := domain.NewPciCard("description")
	assert.Equal(t, want, got)
}

func Test_adaptPciCards(t *testing.T) {
	pciCard1 := dedicatedServer.NewPciCard()
	pciCard1.SetDescription("d1")

	pciCard2 := dedicatedServer.NewPciCard()
	pciCard2.SetDescription("d2")

	got := adaptPciCards([]dedicatedServer.PciCard{*pciCard1, *pciCard2})
	want := domain.PciCards{
		domain.NewPciCard("d1"),
		domain.NewPciCard("d2"),
	}
	assert.Equal(t, want, got)
}

func Test_adaptHdds(t *testing.T) {
	hdd1 := dedicatedServer.NewHdd()
	hdd1.SetId("i1")
	hdd1.SetType("t1")
	hdd1.SetUnit("u1")
	hdd1.SetPerformanceType("p1")
	hdd1.SetAmount(1)
	hdd1.SetSize(2)

	hdd2 := dedicatedServer.NewHdd()
	hdd2.SetId("i2")
	hdd2.SetType("t2")
	hdd2.SetUnit("u2")
	hdd2.SetPerformanceType("p2")
	hdd2.SetAmount(3)
	hdd2.SetSize(4)

	got := adaptHdds([]dedicatedServer.Hdd{*hdd1, *hdd2})
	want := domain.Hdds{
		domain.NewHdd("i1", "t1", "u1", "p1", 1, 2),
		domain.NewHdd("i2", "t2", "u2", "p2", 3, 4),
	}
	assert.Equal(t, want, got)
}

func Test_adaptHdd(t *testing.T) {
	hdd := dedicatedServer.NewHdd()
	hdd.SetId("id")
	hdd.SetType("type")
	hdd.SetUnit("unit")
	hdd.SetAmount(12)
	hdd.SetSize(54)
	hdd.SetPerformanceType("per_type")

	got := adaptHdd(*hdd)
	want := domain.NewHdd("id", "type", "unit", "per_type", 12, 54)
	assert.Equal(t, want, got)
}

func Test_adaptSpecs(t *testing.T) {
	cpu := dedicatedServer.NewCpu()
	cpu.SetQuantity(1)
	cpu.SetType("type")

	ram := dedicatedServer.NewRam()
	ram.SetSize(2)
	ram.SetUnit("unit")

	hdd := dedicatedServer.NewHdd()
	hdd.SetId("i")
	hdd.SetType("t")
	hdd.SetUnit("u")
	hdd.SetPerformanceType("p")
	hdd.SetAmount(1)
	hdd.SetSize(2)

	pciCard := dedicatedServer.NewPciCard()
	pciCard.SetDescription("d")

	specs := dedicatedServer.NewServerSpecs()
	specs.SetChassis("chassis")
	specs.SetHardwareRaidCapable(true)
	specs.SetCpu(*cpu)
	specs.SetRam(*ram)
	specs.SetHdd([]dedicatedServer.Hdd{*hdd})
	specs.SetPciCards([]dedicatedServer.PciCard{*pciCard})

	got := adaptSpecs(*specs)
	want := domain.NewSpecs(
		"chassis",
		true,
		domain.NewCpu(1, "type"),
		domain.NewRam(2, "unit"),
		domain.Hdds{domain.NewHdd("i", "t", "u", "p", 1, 2)},
		domain.PciCards{domain.NewPciCard("d")},
	)
	assert.Equal(t, want, got)
}

func Test_adaptDedicatedServer(t *testing.T) {
	cpu := dedicatedServer.NewCpu()
	cpu.SetQuantity(1)
	cpu.SetType("type")
	ram := dedicatedServer.NewRam()
	ram.SetSize(2)
	ram.SetUnit("unit")
	hdd := dedicatedServer.NewHdd()
	hdd.SetId("i")
	hdd.SetType("t")
	hdd.SetUnit("u")
	hdd.SetPerformanceType("p")
	hdd.SetAmount(1)
	hdd.SetSize(2)
	pciCard := dedicatedServer.NewPciCard()
	pciCard.SetDescription("d")
	specs := dedicatedServer.NewServerSpecs()
	specs.SetChassis("chassis")
	specs.SetHardwareRaidCapable(true)
	specs.SetCpu(*cpu)
	specs.SetRam(*ram)
	specs.SetHdd([]dedicatedServer.Hdd{*hdd})
	specs.SetPciCards([]dedicatedServer.PciCard{*pciCard})

	contract := dedicatedServer.NewContract()
	contract.SetId("id")
	contract.SetCustomerId("customer_id")
	contract.SetDeliveryStatus("status")
	contract.SetReference("ref")
	contract.SetSalesOrgId("sales")

	rack := dedicatedServer.NewRack()
	rack.SetId("id")
	rack.SetCapacity("cap")
	rack.SetType("type")

	featureAvailability := dedicatedServer.NewFeatureAvailability()
	featureAvailability.SetAutomation(true)
	featureAvailability.SetIpmiReboot(false)
	featureAvailability.SetPowerCycle(true)
	featureAvailability.SetPrivateNetwork(false)
	featureAvailability.SetRemoteManagement(true)

	location := dedicatedServer.NewLocation()
	location.SetRack("rack")
	location.SetSite("site")
	location.SetSuite("suite")
	location.SetUnit("unit")

	port := dedicatedServer.NewPort()
	port.SetName("name")
	port.SetPort("1212")

	privateNetwork := dedicatedServer.NewPrivateNetwork()
	privateNetwork.SetId("id")
	privateNetwork.SetStatus("status")
	privateNetwork.SetSubnet("subnet")
	privateNetwork.SetVlanId("vlanid")
	privateNetwork.SetLinkSpeed(12)

	public := dedicatedServer.NewNetworkInterface()
	public.SetMac("m1")
	public.SetIp("i1")
	public.SetGateway("g1")
	public.SetLocationId("l1")
	public.SetNullRouted(true)
	port1 := dedicatedServer.NewPort()
	port1.SetName("n1")
	port1.SetPort("p1")
	public.SetPorts([]dedicatedServer.Port{*port1})

	internal := dedicatedServer.NewNetworkInterface()
	internal.SetMac("m2")
	internal.SetIp("i2")
	internal.SetGateway("g2")
	internal.SetLocationId("l2")
	internal.SetNullRouted(false)
	port2 := dedicatedServer.NewPort()
	port2.SetName("n2")
	port2.SetPort("p2")
	internal.SetPorts([]dedicatedServer.Port{*port2})

	remote := dedicatedServer.NewNetworkInterface()
	remote.SetMac("m3")
	remote.SetIp("i3")
	remote.SetGateway("g3")
	remote.SetLocationId("l3")
	remote.SetNullRouted(true)
	port3 := dedicatedServer.NewPort()
	port3.SetName("n3")
	port3.SetPort("p3")
	remote.SetPorts([]dedicatedServer.Port{*port3})

	networkInterfaces := dedicatedServer.NewNetworkInterfaces()
	networkInterfaces.SetPublic(*public)
	networkInterfaces.SetInternal(*internal)
	networkInterfaces.SetRemoteManagement(*remote)

	server := dedicatedServer.NewServer()
	server.SetId("id")
	server.SetAssetId("aid")
	server.SetSerialNumber("sn")
	server.SetPrivateNetworks([]dedicatedServer.PrivateNetwork{*privateNetwork})
	server.SetFeatureAvailability(*featureAvailability)
	server.SetLocation(*location)
	server.SetContract(*contract)
	server.SetNetworkInterfaces(*networkInterfaces)
	server.SetRack(*rack)
	server.SetPowerPorts([]dedicatedServer.Port{*port})
	server.SetSpecs(*specs)

	got := AdaptDedicatedServer(*server)
	want := domain.NewDedicatedServer(
		"id",
		"aid",
		"sn",
		adaptRack(*rack),
		adaptLocation(*location),
		adaptFeatureAvailability(*featureAvailability),
		adaptContract(*contract),
		adaptPorts([]dedicatedServer.Port{*port}),
		adaptPrivateNetworks([]dedicatedServer.PrivateNetwork{*privateNetwork}),
		adaptNetworkInterfaces(*networkInterfaces),
		adaptSpecs(*specs),
	)
	assert.Equal(t, want, got)
}

func TestAdaptOperatingSystem(t *testing.T) {
	operatingSystem := dedicatedServer.NewOperatingSystem("id", "name")

	got := AdaptOperatingSystem(*operatingSystem)
	want := domain.NewOperatingSystem("id", "name")

	assert.Equal(t, want, got)
}
