package to_resource_model

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	domain "github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain/dedicated_server"
	resourceModel "github.com/leaseweb/terraform-provider-leaseweb/internal/provider/resources/dedicated_server/model"
	"github.com/stretchr/testify/assert"
)

func Test_adaptRack(t *testing.T) {
	rack := domain.NewRack("id", "cap", "type")

	got := adaptRack(rack)
	want := resourceModel.Rack{
		Id:       basetypes.NewStringValue("id"),
		Capacity: basetypes.NewStringValue("cap"),
		Type:     basetypes.NewStringValue("type"),
	}
	assert.Equal(t, want, got)
}

func Test_adaptLocation(t *testing.T) {
	location := domain.NewLocation("rack", "site", "suite", "unit")

	got := adaptLocation(location)
	want := resourceModel.Location{
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
	want := resourceModel.FeatureAvailability{
		Automation:       basetypes.NewBoolValue(true),
		IpmiReboot:       basetypes.NewBoolValue(false),
		PowerCycle:       basetypes.NewBoolValue(true),
		PrivateNetwork:   basetypes.NewBoolValue(false),
		RemoteManagement: basetypes.NewBoolValue(true),
	}
	assert.Equal(t, want, got)

	featureAvailability = domain.NewFeatureAvailability(false, true, false, true, false)
	got = adaptFeatureAvailability(featureAvailability)
	want = resourceModel.FeatureAvailability{
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
	want := resourceModel.Contract{
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
	want := resourceModel.Ports{
		resourceModel.Port{Name: basetypes.NewStringValue("name1"), Port: basetypes.NewStringValue("1111")},
		resourceModel.Port{Name: basetypes.NewStringValue("name2"), Port: basetypes.NewStringValue("2222")},
	}
	assert.Len(t, got, 2)
	assert.Equal(t, want, got)
}

func Test_adaptPort(t *testing.T) {
	port := domain.NewPort("name", "1111")
	got := adaptPort(port)
	want := resourceModel.Port{
		Name: basetypes.NewStringValue("name"),
		Port: basetypes.NewStringValue("1111"),
	}
	assert.Equal(t, want, got)
}

func Test_adaptPrivateNetwork(t *testing.T) {
	privateNetwork := domain.NewPrivateNetwork("id", "status", "subnet", "vlanid", 12)
	got := adaptPrivateNetwork(privateNetwork)
	want := resourceModel.PrivateNetwork{
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
	want := resourceModel.NetworkInterface{
		Mac:        basetypes.NewStringValue("mac"),
		Ip:         basetypes.NewStringValue("ip"),
		Gateway:    basetypes.NewStringValue("gateway"),
		LocationId: basetypes.NewStringValue("loc_id"),
		NullRouted: basetypes.NewBoolValue(true),
		Ports: resourceModel.Ports{
			resourceModel.Port{
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
	want := resourceModel.Ram{
		Size: basetypes.NewInt32Value(12),
		Unit: basetypes.NewStringValue("gb"),
	}
	assert.Equal(t, want, got)
}

func Test_adaptCpu(t *testing.T) {
	cpu := domain.NewCpu(12, "type")
	got := adaptCpu(cpu)
	want := resourceModel.Cpu{
		Quantity: basetypes.NewInt32Value(12),
		Type:     basetypes.NewStringValue("type"),
	}
	assert.Equal(t, want, got)
}

func Test_adaptPciCard(t *testing.T) {
	pciCard := domain.NewPciCard("description")
	got := adaptPciCard(pciCard)
	want := resourceModel.PciCard{
		Description: basetypes.NewStringValue("description"),
	}
	assert.Equal(t, want, got)
}

func Test_adaptPciCards(t *testing.T) {
	pciCards := domain.PciCards{{Description: "description1"}, {Description: "description2"}}
	got := adaptPciCards(pciCards)
	want := resourceModel.PciCards{
		resourceModel.PciCard{Description: basetypes.NewStringValue("description1")},
		resourceModel.PciCard{Description: basetypes.NewStringValue("description2")},
	}
	assert.Len(t, got, 2)
	assert.Equal(t, want, got)
}

func Test_adaptHdds(t *testing.T) {
	hdds := domain.Hdds{domain.Hdd{Id: "id1"}, domain.Hdd{Id: "id2"}}
	got := adaptHdds(hdds)
	assert.Len(t, got, 2)
	assert.Equal(t, "id1", got[0].Id.ValueString())
	assert.Equal(t, "id2", got[1].Id.ValueString())
}

func Test_adaptHdd(t *testing.T) {
	hdd := domain.NewHdd("id1", "type1", "unit1", "per1", 12, 13)
	got := adaptHdd(hdd)
	want := resourceModel.Hdd{
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

func TestAdaptDedicatedServer(t *testing.T) {

	specs := domain.NewSpecs(
		"chassis",
		true,
		domain.Cpu{Quantity: 1},
		domain.Ram{Size: 2},
		domain.Hdds{domain.Hdd{Id: "id1"}},
		domain.PciCards{domain.NewPciCard("d")},
	)

	server := domain.NewDedicatedServer(
		"id",
		"assetId",
		"sn",
		domain.Rack{Id: "rid"},
		domain.Location{Rack: "rack"},
		domain.FeatureAvailability{IpmiReboot: false},
		domain.Contract{Id: "cid"},
		domain.Ports{domain.Port{Name: "name1"}, domain.Port{Name: "name2"}},
		domain.PrivateNetworks{domain.PrivateNetwork{Id: "pid"}},
		domain.NetworkInterfaces{Public: domain.NetworkInterface{Mac: "public"}},
		specs,
	)

	got := AdaptDedicatedServer(server)
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
