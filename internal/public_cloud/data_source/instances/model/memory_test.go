package model

import (
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_newMemory(t *testing.T) {
	sdkMemory := publicCloud.NewMemory()
	sdkMemory.SetValue(1)
	sdkMemory.SetUnit("unit")

	memory := newMemory(*sdkMemory)

	assert.Equal(t, float64(1), memory.Value.ValueFloat64(), "value should be set")
	assert.Equal(t, "unit", memory.Unit.ValueString(), "unit should be set")
}
