package model

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/domain/entity"
)

func Test_newIso(t *testing.T) {
	entityIso := entity.NewIso("id", "name")
	got, diags := newIso(context.TODO(), entityIso)

	assert.Nil(t, diags)

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
}

func TestIso_attributeTypes(t *testing.T) {
	_, diags := types.ObjectValueFrom(
		context.TODO(),
		Iso{}.AttributeTypes(),
		Iso{},
	)

	assert.Nil(t, diags, "attributes should be correct")
}
