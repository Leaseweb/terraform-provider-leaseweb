package model

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func Test_newCpu(t *testing.T) {
	sdkCpu := publicCloud.NewCpu(1, "unit")
	got := newCpu(sdkCpu)

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
