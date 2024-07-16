package model

import (
	"terraform-provider-leaseweb/internal/core/domain/entity"
)

type Instances struct {
	Instances []instance `tfsdk:"instances"`
}

func (m *Instances) Populate(entityInstances entity.Instances) {
	for _, entityInstance := range entityInstances {
		instance := newInstance(entityInstance)
		m.Instances = append(m.Instances, instance)
	}
}
