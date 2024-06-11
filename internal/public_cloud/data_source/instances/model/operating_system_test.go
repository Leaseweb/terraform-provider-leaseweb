package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_newOperatingSystem(t *testing.T) {
	sdkOperatingSystem := publicCloud.NewOperatingSystem()
	sdkOperatingSystem.SetId("id")
	sdkOperatingSystem.SetName("name")
	sdkOperatingSystem.SetVersion("version")
	sdkOperatingSystem.SetFamily("family")
	sdkOperatingSystem.SetFlavour("flavour")
	sdkOperatingSystem.SetArchitecture("architecture")
	sdkOperatingSystem.SetMarketApps([]string{"one"})
	sdkOperatingSystem.SetStorageTypes([]string{"storageType"})

	operatingSystem := newOperatingSystem(*sdkOperatingSystem)

	assert.Equal(
		t,
		"id",
		operatingSystem.Id.ValueString(),
		"id should be set",
	)
	assert.Equal(
		t,
		"name",
		operatingSystem.Name.ValueString(),
		"name should be set",
	)
	assert.Equal(
		t,
		"version",
		operatingSystem.Version.ValueString(),
		"version should be set",
	)
	assert.Equal(
		t,
		"family",
		operatingSystem.Family.ValueString(),
		"family should be set",
	)
	assert.Equal(
		t,
		"flavour",
		operatingSystem.Flavour.ValueString(),
		"flavour should be set",
	)
	assert.Equal(
		t,
		"architecture",
		operatingSystem.Architecture.ValueString(),
		"architecture should be set",
	)
	assert.Equal(
		t,
		[]types.String{basetypes.NewStringValue("one")},
		operatingSystem.MarketApps,
		"marketApps should be set",
	)
	assert.Equal(
		t,
		[]types.String{basetypes.NewStringValue("storageType")},
		operatingSystem.StorageTypes,
		"storageTypes should be set",
	)
}
