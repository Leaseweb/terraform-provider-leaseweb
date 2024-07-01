package model

import (
	"testing"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func Test_newCpu(t *testing.T) {
	sdkCpu := publicCloud.Cpu{Value: 1, Unit: "unit"}
	got := newCpu(sdkCpu)

	assert.Equal(t, int64(1), got.Value.ValueInt64(), "value should be set")
	assert.Equal(t, "unit", got.Unit.ValueString(), "unit should be set")
}
