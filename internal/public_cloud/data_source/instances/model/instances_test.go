package model

import (
	"testing"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func TestInstances_Populate(t *testing.T) {
	instanceDetails := publicCloud.InstanceDetails{Id: "instanceId"}

	instances := Instances{}
	instances.Populate([]publicCloud.InstanceDetails{instanceDetails})

	assert.Equal(
		t,
		"instanceId",
		instances.Instances[0].Id.ValueString(),
		"instance should be set",
	)
}
