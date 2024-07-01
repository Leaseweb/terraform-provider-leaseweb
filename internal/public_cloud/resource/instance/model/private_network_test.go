package model

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func Test_newPrivateNetwork(t *testing.T) {
	sdkPrivateNetwork := publicCloud.NewPrivateNetwork(
		"id",
		"status",
		"subnet",
	)

	got, diags := newPrivateNetwork(context.TODO(), *sdkPrivateNetwork)

	assert.Nil(t, diags)

	assert.Equal(
		t,
		"id",
		got.Id.ValueString(),
		"id should be set",
	)
	assert.Equal(
		t,
		"status",
		got.Status.ValueString(),
		"status should be set",
	)
	assert.Equal(
		t,
		"subnet",
		got.Subnet.ValueString(),
		"subnet should be set",
	)
}

func TestPrivateNetwork_attributeTypes(t *testing.T) {
	_, diags := types.ObjectValueFrom(
		context.TODO(),
		PrivateNetwork{}.AttributeTypes(),
		PrivateNetwork{},
	)

	assert.Nil(t, diags, "attributes should be correct")
}
