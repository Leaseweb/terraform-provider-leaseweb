package model

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func TestInstance_Populate(t *testing.T) {
	startedAt, _ := time.Parse(time.RFC3339, "2019-09-08T00:00:00Z")

	sdkInstance := publicCloud.NewInstance()
	sdkInstance.SetId("id")
	sdkInstance.SetEquipmentId("equipmentId")
	sdkInstance.SetSalesOrgId("salesOrgId")
	sdkInstance.SetCustomerId("customerId")
	sdkInstance.SetRegion("region")
	sdkInstance.SetOperatingSystem(*publicCloud.NewOperatingSystem())
	sdkInstance.SetState("state")
	sdkInstance.SetProductType("productType")
	sdkInstance.SetHasPublicIpV4(true)
	sdkInstance.SetincludesPrivateNetwork(false)
	sdkInstance.SetRootDiskSize(32)
	sdkInstance.SetRootDiskStorageType("rootDiskStorageType")
	sdkInstance.SetStartedAt(startedAt)
	sdkInstance.SetContract(*publicCloud.NewContract())
	sdkInstance.SetIso(*publicCloud.NewIso())
	sdkInstance.SetMarketAppId("marketAppId")
	sdkInstance.SetPrivateNetwork(*publicCloud.NewPrivateNetwork())
	sdkInstance.SetResources(*publicCloud.NewInstanceResources())

	sdkInstanceTypeName, _ := publicCloud.NewInstanceTypeNameFromValue("lsw.m5a.4xlarge")
	sdkInstance.SetType(*sdkInstanceTypeName)

	sdkContract := publicCloud.NewContract()
	sdkContract.SetType("contract")
	sdkInstance.SetContract(*sdkContract)

	sdkOperatingSystem := publicCloud.NewOperatingSystem()
	sdkOperatingSystem.SetId("operatingSystemId")
	sdkInstance.SetOperatingSystem(*sdkOperatingSystem)

	sdkIso := publicCloud.NewIso()
	sdkIso.SetId("isoId")
	sdkInstance.SetIso(*sdkIso)

	sdkPrivateNetwork := publicCloud.NewPrivateNetwork()
	sdkPrivateNetwork.SetPrivateNetworkId("privateNetworkId")
	sdkInstance.SetPrivateNetwork(*sdkPrivateNetwork)

	sdkIp := publicCloud.NewIp()
	sdkIp.SetIp("1.2.3.4")

	sdkInstance.SetIps([]publicCloud.Ip{*sdkIp})

	sdkCpu := publicCloud.NewCpu()
	sdkCpu.SetUnit("cpu")

	sdkResources := publicCloud.NewInstanceResources()
	sdkResources.SetCpu(*sdkCpu)
	sdkInstance.SetResources(*sdkResources)

	instance := Instance{}
	instance.Populate(sdkInstance, context.TODO())

	assert.Equal(
		t,
		"id",
		instance.Id.ValueString(),
		"id should be set",
	)
	assert.Equal(
		t,
		"equipmentId",
		instance.EquipmentId.ValueString(),
		"equipmentId should be set",
	)
	assert.Equal(
		t,
		"salesOrgId",
		instance.SalesOrgId.ValueString(),
		"salesOrgId should be set",
	)
	assert.Equal(
		t,
		"customerId",
		instance.CustomerId.ValueString(),
		"customerId should be set",
	)
	assert.Equal(
		t,
		"region",
		instance.Region.ValueString(),
		"region should be set",
	)
	assert.Equal(
		t,
		"state",
		instance.State.ValueString(),
		"state should be set",
	)
	assert.Equal(
		t,
		"productType",
		instance.ProductType.ValueString(),
		"productType should be set",
	)
	assert.Equal(
		t,
		true,
		instance.HasPublicIpv4.ValueBool(),
		"hasPublicIpv should be set",
	)
	assert.Equal(
		t,
		false,
		instance.HasPrivateNetwork.ValueBool(),
		"hasPrivateNetwork should be set",
	)
	assert.Equal(
		t,
		"lsw.m5a.4xlarge",
		instance.Type.ValueString(),
		"type should be set",
	)
	assert.Equal(
		t,
		int64(32),
		instance.RootDiskSize.ValueInt64(),
		"rootDiskSize should be set",
	)
	assert.Equal(
		t,
		"rootDiskStorageType",
		instance.RootDiskStorageType.ValueString(),
		"rootDiskStorageType should be set",
	)
	assert.Equal(
		t,
		"2019-09-08 00:00:00 +0000 UTC",
		instance.StartedAt.ValueString(),
		"startedAt should be set",
	)
	assert.Equal(
		t,
		"marketAppId",
		instance.MarketAppId.ValueString(),
		"marketAppId should be set",
	)

	operatingSystem := OperatingSystem{}
	instance.OperatingSystem.As(
		context.TODO(),
		&operatingSystem,
		basetypes.ObjectAsOptions{},
	)
	assert.Equal(
		t,
		"operatingSystemId",
		operatingSystem.Id.ValueString(),
		"operating_system should be set",
	)

	contract := Contract{}
	instance.Contract.As(context.TODO(), &contract, basetypes.ObjectAsOptions{})
	assert.Equal(
		t,
		"contract",
		contract.Type.ValueString(),
		"contract should be set",
	)

	iso := Iso{}
	instance.Iso.As(context.TODO(), &iso, basetypes.ObjectAsOptions{})
	assert.Equal(
		t,
		"isoId",
		iso.Id.ValueString(),
		"iso should be set",
	)

	privateNetwork := PrivateNetwork{}
	instance.PrivateNetwork.As(
		context.TODO(),
		&privateNetwork,
		basetypes.ObjectAsOptions{},
	)
	assert.Equal(
		t,
		"privateNetworkId",
		privateNetwork.Id.ValueString(),
		"privateNetwork should be set",
	)

	var ips []Ip
	instance.Ips.ElementsAs(context.TODO(), &ips, false)
	assert.Equal(
		t,
		"1.2.3.4",
		ips[0].Ip.ValueString(),
		"ip should be set",
	)

	resources := Resources{}
	cpu := Cpu{}
	instance.Resources.As(context.TODO(), &resources, basetypes.ObjectAsOptions{})
	resources.Cpu.As(context.TODO(), &cpu, basetypes.ObjectAsOptions{})
	assert.Equal(
		t,
		"cpu",
		cpu.Unit.ValueString(),
		"privateNetwork should be set",
	)
}
