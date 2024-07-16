package model

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/domain"
)

func Test_newCpu(t *testing.T) {
	entityCpu := domain.NewCpu(1, "unit")
	got := newCpu(entityCpu)

	assert.Equal(t, int64(1), got.Value.ValueInt64(), "value should be set")
	assert.Equal(t, "unit", got.Unit.ValueString(), "unit should be set")
}

func TestCpu_attributeTypes(t *testing.T) {
	_, diags := types.ObjectValueFrom(
		context.TODO(),
		Cpu{}.AttributeTypes(),
		Cpu{},
	)

	assert.Nil(t, diags, "attributes should be correct")
}
