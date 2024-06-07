package model

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_newCpu(t *testing.T) {
	sdkCpu := publicCloud.NewCpu()
	sdkCpu.SetValue(1)
	sdkCpu.SetUnit("unit")

	cpu := newCpu(sdkCpu)

	assert.Equal(t, int64(1), cpu.Value.ValueInt64(), "value should be set")
	assert.Equal(t, "unit", cpu.Unit.ValueString(), "unit should be set")
}

func TestCpu_attributeTypes(t *testing.T) {
	_, diags := types.ObjectValueFrom(context.TODO(), Cpu{}.attributeTypes(), Cpu{})

	assert.Nil(t, diags, "attributes should be correct")
}
