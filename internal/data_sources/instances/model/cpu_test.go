package model

import (
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_newCpu(t *testing.T) {
	sdkCpu := publicCloud.NewCpu()
	sdkCpu.SetValue(1)
	sdkCpu.SetUnit("unit")

	cpu := newCpu(*sdkCpu)

	assert.Equal(t, int64(1), cpu.Value.ValueInt64(), "value should be set")
	assert.Equal(t, "unit", cpu.Unit.ValueString(), "unit should be set")
}
