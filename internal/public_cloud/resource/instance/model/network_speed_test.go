package model

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/domain"
)

func Test_newNetworkSpeed(t *testing.T) {
	networkSpeed := domain.NewNetworkSpeed(23, "unit")

	got := newNetworkSpeed(networkSpeed)

	assert.Equal(
		t,
		"unit",
		got.Unit.ValueString(),
		"unit should be set",
	)
	assert.Equal(
		t,
		int64(23),
		got.Value.ValueInt64(),
		"value should be set",
	)
}

func TestNetworkSpeed_attributeTypes(t *testing.T) {
	_, diags := types.ObjectValueFrom(
		context.TODO(),
		NetworkSpeed{}.AttributeTypes(),
		NetworkSpeed{},
	)

	assert.Nil(t, diags, "attributes should be correct")
}
