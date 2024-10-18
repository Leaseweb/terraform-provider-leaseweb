package publiccloud

import (
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
)

type DataSourceModelInstances struct {
	Instances []DataSourceModelInstance `tfsdk:"instances"`
}

func newDataSourceModelInstances(sdkInstances []publicCloud.Instance) DataSourceModelInstances {
	var instances DataSourceModelInstances

	for _, sdkInstance := range sdkInstances {
		instance := newDataSourceModelInstance(sdkInstance)
		instances.Instances = append(instances.Instances, instance)
	}

	return instances
}
