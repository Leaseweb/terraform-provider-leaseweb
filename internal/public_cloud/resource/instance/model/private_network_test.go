package model

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/domain"
)

func Test_newPrivateNetwork(t *testing.T) {
	privateNetwork := domain.NewPrivateNetwork(
		"id",
		"status",
		"subnet",
	)

	got, diags := newPrivateNetwork(context.TODO(), privateNetwork)

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
