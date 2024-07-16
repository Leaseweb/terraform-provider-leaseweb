package model

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/domain"
)

func Test_newMemory(t *testing.T) {
	memory := domain.NewMemory(1, "unit")

	got := newMemory(memory)

	assert.Equal(t, float64(1), got.Value.ValueFloat64(), "value should be set")
	assert.Equal(t, "unit", got.Unit.ValueString(), "unit should be set")
}

func TestMemory_attributeTypes(t *testing.T) {
	_, diags := types.ObjectValueFrom(
		context.TODO(),
		Memory{}.AttributeTypes(),
		Memory{},
	)

	assert.Nil(t, diags, "attributes should be correct")
}
