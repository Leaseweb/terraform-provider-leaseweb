package model

import (
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInstances_Populate(t *testing.T) {
	instance := publicCloud.NewInstance()
	instance.SetId("instanceId")
	instance.SetResources(*publicCloud.NewInstanceResources())
	instance.SetOperatingSystem(*publicCloud.NewOperatingSystem())

	instances := Instances{}
	instances.Populate([]publicCloud.Instance{*instance})

	assert.Equal(
		t,
		"instanceId",
		instances.Instances[0].Id.ValueString(),
		"instance should be set",
	)
}
