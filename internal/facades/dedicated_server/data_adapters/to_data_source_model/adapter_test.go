package to_data_source_model

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	domain "github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain/dedicated_server"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/data_sources/dedicated_server/model"
	"github.com/stretchr/testify/assert"
)

func Test_adaptRack(t *testing.T) {
	rack := domain.NewRack("id", "cap", "type")

	got := adaptRack(rack)
	want := model.Rack{
		Id:       basetypes.NewStringValue("id"),
		Capacity: basetypes.NewStringValue("cap"),
		Type:     basetypes.NewStringValue("type"),
	}
	assert.Equal(t, want, got)
}

func Test_adaptLocation(t *testing.T) {
	location := domain.NewLocation("rack", "site", "suite", "unit")

	got := adaptLocation(location)
	want := model.Location{
		Rack:  basetypes.NewStringValue("rack"),
		Site:  basetypes.NewStringValue("site"),
		Suite: basetypes.NewStringValue("suite"),
		Unit:  basetypes.NewStringValue("unit"),
	}
	assert.Equal(t, want, got)
}

func Test_adaptFeatureAvailability(t *testing.T) {
	featureAvailability := domain.NewFeatureAvailability(true, false, true, false, true)
	got := adaptFeatureAvailability(featureAvailability)
	want := model.FeatureAvailability{
		Automation:       basetypes.NewBoolValue(true),
		IpmiReboot:       basetypes.NewBoolValue(false),
		PowerCycle:       basetypes.NewBoolValue(true),
		PrivateNetwork:   basetypes.NewBoolValue(false),
		RemoteManagement: basetypes.NewBoolValue(true),
	}
	assert.Equal(t, want, got)

	featureAvailability = domain.NewFeatureAvailability(false, true, false, true, false)
	got = adaptFeatureAvailability(featureAvailability)
	want = model.FeatureAvailability{
		Automation:       basetypes.NewBoolValue(false),
		IpmiReboot:       basetypes.NewBoolValue(true),
		PowerCycle:       basetypes.NewBoolValue(false),
		PrivateNetwork:   basetypes.NewBoolValue(true),
		RemoteManagement: basetypes.NewBoolValue(false),
	}
	assert.Equal(t, want, got)
}

func Test_adaptContract(t *testing.T) {
	contract := domain.NewContract("id", "customer_id", "status", "ref", "sales")
	got := adaptContract(contract)
	want := model.Contract{
		Id:             basetypes.NewStringValue("id"),
		CustomerId:     basetypes.NewStringValue("customer_id"),
		DeliveryStatus: basetypes.NewStringValue("status"),
		Reference:      basetypes.NewStringValue("ref"),
		SalesOrgId:     basetypes.NewStringValue("sales"),
	}
	assert.Equal(t, want, got)
}

func Test_adaptPorts(t *testing.T) {
	port1 := domain.NewPort("name1", "1111")
	port2 := domain.NewPort("name2", "2222")
	ports := domain.Ports{port1, port2}
	got := adaptPorts(ports)
	want := model.Ports{
		model.Port{Name: basetypes.NewStringValue("name1"), Port: basetypes.NewStringValue("1111")},
		model.Port{Name: basetypes.NewStringValue("name2"), Port: basetypes.NewStringValue("2222")},
	}
	assert.Len(t, got, 2)
	assert.Equal(t, want, got)
}

func Test_adaptPort(t *testing.T) {
	port := domain.NewPort("name", "1111")
	got := adaptPort(port)
	want := model.Port{
		Name: basetypes.NewStringValue("name"),
		Port: basetypes.NewStringValue("1111"),
	}
	assert.Equal(t, want, got)
}

func Test_adaptPrivateNetwork(t *testing.T) {
	privateNetwork := domain.NewPrivateNetwork("id", "status", "subnet", "vlanid", 12)
	got := adaptPrivateNetwork(privateNetwork)
	want := model.PrivateNetwork{
		Id:        basetypes.NewStringValue("id"),
		Status:    basetypes.NewStringValue("status"),
		Subnet:    basetypes.NewStringValue("subnet"),
		VlanId:    basetypes.NewStringValue("vlanid"),
		LinkSpeed: basetypes.NewInt32Value(12),
	}
	assert.Equal(t, want, got)
}

func Test_adaptPrivateNetworks(t *testing.T) {
	privateNetwork := domain.PrivateNetworks{
		domain.PrivateNetwork{Id: "id1"},
		domain.PrivateNetwork{Id: "id2"},
	}

	got := adaptPrivateNetworks(privateNetwork)
	assert.Len(t, got, 2)
	assert.Equal(t, "id1", got[0].Id.ValueString())
	assert.Equal(t, "id2", got[1].Id.ValueString())
}

func Test_adaptNetworkInterface(t *testing.T) {
	port := domain.NewPort("name", "1111")
	ports := domain.Ports{port}
	networkInterface := domain.NewNetworkInterface("mac", "ip", "gateway", "loc_id", true, ports)

	got := adaptNetworkInterface(networkInterface)
	want := model.NetworkInterface{
		Mac:        basetypes.NewStringValue("mac"),
		Ip:         basetypes.NewStringValue("ip"),
		Gateway:    basetypes.NewStringValue("gateway"),
		LocationId: basetypes.NewStringValue("loc_id"),
		NullRouted: basetypes.NewBoolValue(true),
		Ports: model.Ports{
			model.Port{
				Name: basetypes.NewStringValue("name"),
				Port: basetypes.NewStringValue("1111"),
			},
		},
	}
	assert.Equal(t, want, got)
}

func Test_adaptNetworkInterfaces(t *testing.T) {
	networkInterfaces := domain.NetworkInterfaces{
		Public: domain.NetworkInterface{Mac: "public"},
	}

	got := adaptNetworkInterfaces(networkInterfaces)
	assert.Equal(t, "public", got.Public.Mac.ValueString())
}

func Test_adaptRam(t *testing.T) {
	ram := domain.NewRam(12, "gb")
	got := adaptRam(ram)
	want := model.Ram{
		Size: basetypes.NewInt32Value(12),
		Unit: basetypes.NewStringValue("gb"),
	}
	assert.Equal(t, want, got)
}

func Test_adaptCpu(t *testing.T) {
	cpu := domain.NewCpu(12, "type")
	got := adaptCpu(cpu)
	want := model.Cpu{
		Quantity: basetypes.NewInt32Value(12),
		Type:     basetypes.NewStringValue("type"),
	}
	assert.Equal(t, want, got)
}

func Test_adaptPciCard(t *testing.T) {
	pciCard := domain.NewPciCard("description")
	got := adaptPciCard(pciCard)
	want := model.PciCard{
		Description: basetypes.NewStringValue("description"),
	}
	assert.Equal(t, want, got)
}

func Test_adaptPciCards(t *testing.T) {
	pciCard1 := domain.NewPciCard("description1")
	pciCard2 := domain.NewPciCard("description2")
	pciCards := domain.PciCards{pciCard1, pciCard2}
	got := adaptPciCards(pciCards)
	want := model.PciCards{
		model.PciCard{Description: basetypes.NewStringValue("description1")},
		model.PciCard{Description: basetypes.NewStringValue("description2")},
	}
	assert.Len(t, got, 2)
	assert.Equal(t, want, got)
}

func Test_adaptHdds(t *testing.T) {
	hdds := domain.Hdds{
		domain.Hdd{Id: "id1"},
		domain.Hdd{Id: "id2"},
	}
	got := adaptHdds(hdds)
	assert.Len(t, got, 2)
	assert.Equal(t, "id1", got[0].Id.ValueString())
	assert.Equal(t, "id2", got[1].Id.ValueString())
}

func Test_adaptHdd(t *testing.T) {
	hdd := domain.NewHdd("id1", "type1", "unit1", "per1", 12, 13)
	got := adaptHdd(hdd)
	want := model.Hdd{
		Id:              basetypes.NewStringValue("id1"),
		Type:            basetypes.NewStringValue("type1"),
		Unit:            basetypes.NewStringValue("unit1"),
		PerformanceType: basetypes.NewStringValue("per1"),
		Amount:          basetypes.NewInt32Value(12),
		Size:            basetypes.NewFloat32Value(13),
	}
	assert.Equal(t, want, got)
}

func Test_adaptSpecs(t *testing.T) {
	specs := domain.NewSpecs(
		"chassis",
		true,
		domain.Cpu{Quantity: 1},
		domain.Ram{Size: 2},
		domain.Hdds{domain.Hdd{Id: "id"}},
		domain.PciCards{domain.NewPciCard("d")},
	)
	got := adaptSpecs(specs)
	assert.Equal(t, int32(2), got.Ram.Size.ValueInt32())
	assert.Equal(t, int32(1), got.Cpu.Quantity.ValueInt32())
	assert.Equal(t, "id", got.Hdds[0].Id.ValueString())
	assert.Equal(t, "d", got.PciCards[0].Description.ValueString())
	assert.Equal(t, "chassis", got.Chassis.ValueString())
	assert.Equal(t, true, got.HardwareRaidCapable.ValueBool())
}

func Test_adaptDedicatedServer(t *testing.T) {

	specs := domain.NewSpecs(
		"chassis",
		true,
		domain.Cpu{Quantity: 1},
		domain.Ram{Size: 2},
		domain.Hdds{domain.Hdd{Id: "id1"}},
		domain.PciCards{domain.NewPciCard("d")},
	)
	featureAvailability := domain.FeatureAvailability{IpmiReboot: false}
	contract := domain.Contract{Id: "cid"}
	rack := domain.Rack{Id: "rid"}
	location := domain.Location{Rack: "rack"}
	ports := domain.Ports{domain.Port{Name: "name1"}, domain.Port{Name: "name2"}}
	privateNetwork := domain.PrivateNetworks{domain.PrivateNetwork{Id: "pid"}}
	networkInterfaces := domain.NetworkInterfaces{
		Public: domain.NetworkInterface{Mac: "public"},
	}

	server := domain.NewDedicatedServer(
		"id",
		"assetId",
		"sn",
		rack,
		location,
		featureAvailability,
		contract,
		ports,
		privateNetwork,
		networkInterfaces,
		specs,
	)

	got := adaptDedicatedServer(server)
	assert.Equal(t, "id", got.Id.ValueString())
	assert.Equal(t, "assetId", got.AssetId.ValueString())
	assert.Equal(t, "sn", got.SerialNumber.ValueString())
	assert.Equal(t, int32(2), got.Specs.Ram.Size.ValueInt32())
	assert.Equal(t, int32(1), got.Specs.Cpu.Quantity.ValueInt32())
	assert.Equal(t, "id1", got.Specs.Hdds[0].Id.ValueString())
	assert.Equal(t, "d", got.Specs.PciCards[0].Description.ValueString())
	assert.Equal(t, "chassis", got.Specs.Chassis.ValueString())
	assert.Equal(t, true, got.Specs.HardwareRaidCapable.ValueBool())
	assert.Equal(t, "cid", got.Contract.Id.ValueString())
	assert.Equal(t, "rid", got.Rack.Id.ValueString())
	assert.False(t, got.FeatureAvailability.IpmiReboot.ValueBool())
	assert.Equal(t, "rack", got.Location.Rack.ValueString())
	assert.Equal(t, "name1", got.PowerPorts[0].Name.ValueString())
	assert.Equal(t, "pid", got.PrivateNetworks[0].Id.ValueString())
	assert.Equal(t, "public", got.NetworkInterfaces.Public.Mac.ValueString())
}

func TestAdaptDedicatedServers(t *testing.T) {

	specs := domain.NewSpecs(
		"chassis",
		true,
		domain.Cpu{Quantity: 1},
		domain.Ram{Size: 2},
		domain.Hdds{domain.Hdd{Id: "id1"}},
		domain.PciCards{domain.NewPciCard("d")},
	)
	featureAvailability := domain.FeatureAvailability{IpmiReboot: false}
	contract := domain.Contract{Id: "cid"}
	rack := domain.Rack{Id: "rid"}
	location := domain.Location{Rack: "rack"}
	ports := domain.Ports{domain.Port{Name: "name1"}, domain.Port{Name: "name2"}}
	privateNetwork := domain.PrivateNetworks{domain.PrivateNetwork{Id: "pid"}}
	networkInterfaces := domain.NetworkInterfaces{
		Public: domain.NetworkInterface{Mac: "public"},
	}

	server1 := domain.NewDedicatedServer(
		"id1",
		"assetId1",
		"sn1",
		rack,
		location,
		featureAvailability,
		contract,
		ports,
		privateNetwork,
		networkInterfaces,
		specs,
	)

	server2 := domain.NewDedicatedServer(
		"id2",
		"assetId2",
		"sn2",
		rack,
		location,
		featureAvailability,
		contract,
		ports,
		privateNetwork,
		networkInterfaces,
		specs,
	)

	servers := domain.DedicatedServers{server1, server2}
	got := AdaptDedicatedServers(servers)
	assert.Len(t, got.DedicatedServers, 2)
	assert.Equal(t, "id1", got.DedicatedServers[0].Id.ValueString())
	assert.Equal(t, "id2", got.DedicatedServers[1].Id.ValueString())
}
