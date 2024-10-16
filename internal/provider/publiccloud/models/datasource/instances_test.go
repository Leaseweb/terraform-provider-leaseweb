package datasource

import (
	"testing"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func Test_newInstances(t *testing.T) {
	sdkInstances := []publicCloud.Instance{
		{Id: "id"},
	}

	got := NewInstances(sdkInstances)

	assert.Len(t, got.Instances, 1)
	assert.Equal(t, "id", got.Instances[0].Id.ValueString())
}
