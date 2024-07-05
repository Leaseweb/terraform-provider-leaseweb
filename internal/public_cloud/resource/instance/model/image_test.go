package model

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func Test_newImage(t *testing.T) {
	sdkImage := publicCloud.NewImageDetails(
		"id",
		"name",
		"version",
		"family",
		"flavour",
		"architecture",
		[]string{"one"},
		[]string{"storageType"},
	)

	image, diags := newImage(context.TODO(), *sdkImage)

	assert.Nil(t, diags)

	assert.Equal(
		t,
		"id",
		image.Id.ValueString(),
		"id should be set",
	)
	assert.Equal(
		t,
		"name",
		image.Name.ValueString(),
		"name should be set",
	)
	assert.Equal(
		t,
		"version",
		image.Version.ValueString(),
		"version should be set",
	)
	assert.Equal(
		t,
		"family",
		image.Family.ValueString(),
		"family should be set",
	)
	assert.Equal(
		t,
		"flavour",
		image.Flavour.ValueString(),
		"flavour should be set",
	)
	assert.Equal(
		t,
		"architecture",
		image.Architecture.ValueString(),
		"architecture should be set",
	)
	assert.Equal(
		t,
		[]types.String{basetypes.NewStringValue("one")},
		image.MarketApps,
		"marketApps should be set",
	)
	assert.Equal(
		t,
		[]types.String{basetypes.NewStringValue("storageType")},
		image.StorageTypes,
		"storageTypes should be set",
	)
}

func TestImage_attributeTypes(t *testing.T) {
	_, diags := types.ObjectValueFrom(
		context.TODO(),
		Image{}.AttributeTypes(),
		Image{},
	)

	assert.Nil(t, diags, "attributes should be correct")
}
