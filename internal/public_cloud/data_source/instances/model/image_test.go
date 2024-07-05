package model

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func Test_newImage(t *testing.T) {
	sdkImageDetails := publicCloud.NewImageDetails(
		"id",
		"name",
		"version",
		"family",
		"flavour",
		"architecture",
		[]string{"one"},
		[]string{"storageType"},
	)

	got := newImage(*sdkImageDetails)

	assert.Equal(
		t,
		"id",
		got.Id.ValueString(),
		"id should be set",
	)
	assert.Equal(
		t,
		"name",
		got.Name.ValueString(),
		"name should be set",
	)
	assert.Equal(
		t,
		"version",
		got.Version.ValueString(),
		"version should be set",
	)
	assert.Equal(
		t,
		"family",
		got.Family.ValueString(),
		"family should be set",
	)
	assert.Equal(
		t,
		"flavour",
		got.Flavour.ValueString(),
		"flavour should be set",
	)
	assert.Equal(
		t,
		"architecture",
		got.Architecture.ValueString(),
		"architecture should be set",
	)
	assert.Equal(
		t,
		[]types.String{basetypes.NewStringValue("one")},
		got.MarketApps,
		"marketApps should be set",
	)
	assert.Equal(
		t,
		[]types.String{basetypes.NewStringValue("storageType")},
		got.StorageTypes,
		"storageTypes should be set",
	)
}