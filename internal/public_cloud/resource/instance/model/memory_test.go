package model

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func Test_newMemory(t *testing.T) {
	sdkMemory := publicCloud.NewMemory(1, "unit")

	got := newMemory(sdkMemory)

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
