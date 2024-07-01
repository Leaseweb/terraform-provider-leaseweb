package model

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func Test_newOperatingSystem(t *testing.T) {
	sdkOperatingSystem := publicCloud.NewOperatingSystemDetails(
		"id",
		"name",
		"version",
		"family",
		"flavour",
		"architecture",
		[]string{"one"},
		[]string{"storageType"},
	)

	operatingSystem, diags := newOperatingSystem(context.TODO(), *sdkOperatingSystem)

	assert.Nil(t, diags)

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

func TestOperatingSystem_attributeTypes(t *testing.T) {
	_, diags := types.ObjectValueFrom(
		context.TODO(),
		OperatingSystem{}.AttributeTypes(),
		OperatingSystem{},
	)

	assert.Nil(t, diags, "attributes should be correct")
}
