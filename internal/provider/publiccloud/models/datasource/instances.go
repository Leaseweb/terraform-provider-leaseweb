package datasource

import (
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
)

type Instances struct {
	Instances []Instance `tfsdk:"instances"`
}

func NewInstances(sdkInstances []publicCloud.Instance) Instances {
	var instances Instances

	for _, sdkInstance := range sdkInstances {
		instance := NewInstance(sdkInstance)
		instances.Instances = append(instances.Instances, instance)
	}

	return instances
}
