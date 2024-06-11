package model

import (
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
)

type Instances struct {
	Instances []instance `tfsdk:"instances"`
}

func (m *Instances) Populate(sdkInstances []publicCloud.Instance) {
	for _, sdkInstance := range sdkInstances {
		instance := newInstance(sdkInstance)
		m.Instances = append(m.Instances, instance)
	}
}
