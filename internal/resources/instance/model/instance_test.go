package model

import (
	"context"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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
	sdkInstance.SetType("type")
	sdkInstance.SetRootDiskSize(32)
	sdkInstance.SetRootDiskStorageType("rootDiskStorageType")
	sdkInstance.SetStartedAt(startedAt)
	sdkInstance.SetContract(*publicCloud.NewContract())
	sdkInstance.SetIso(*publicCloud.NewIso())
	sdkInstance.SetMarketAppId("marketAppId")
	sdkInstance.SetPrivateNetwork(*publicCloud.NewPrivateNetwork())
	sdkInstance.SetResources(*publicCloud.NewInstanceResources())

	instance := Instance{}
	instance.Populate(sdkInstance, context.Background())

	assert.Equal(t, "id", instance.Id.ValueString(), "id should be set")
	assert.Equal(t, "equipmentId", instance.EquipmentId.ValueString(), "equipmentId should be set")
	assert.Equal(t, "salesOrgId", instance.SalesOrgId.ValueString(), "salesOrgId should be set")
	assert.Equal(t, "customerId", instance.CustomerId.ValueString(), "customerId should be set")
	assert.Equal(t, "region", instance.Region.ValueString(), "region should be set")
	assert.Equal(t, "state", instance.State.ValueString(), "state should be set")
	assert.Equal(t, "productType", instance.ProductType.ValueString(), "productType should be set")
	assert.Equal(t, true, instance.HasPublicIpv4.ValueBool(), "hasPublicIpv should be set")
	assert.Equal(t, false, instance.HasPrivateNetwork.ValueBool(), "hasPrivateNetwork should be set")
	assert.Equal(t, "type", instance.Type.ValueString(), "type should be set")
	assert.Equal(t, int64(32), instance.RootDiskSize.ValueInt64(), "rootDiskSize should be set")
	assert.Equal(t, "\"rootDiskStorageType\"", instance.RootDiskStorageType.String(), "rootDiskStorageType should be set")
	assert.Equal(t, "2019-09-08 00:00:00 +0000 UTC", instance.StartedAt.ValueString(), "startedAt should be set")
	assert.Equal(t, "marketAppId", instance.MarketAppId.ValueString(), "marketAppId should be set")

}
