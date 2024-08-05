package enum

import (
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/enum_utils"
)

type AutoScalingGroupState string

func (a AutoScalingGroupState) String() string {
	return string(a)
}

const (
	AutoScalingGroupStateActive     AutoScalingGroupState = "ACTIVE"
	AutoScalingGroupStateCreating   AutoScalingGroupState = "CREATING"
	AutoScalingGroupStateCreated    AutoScalingGroupState = "CREATED"
	AutoScalingGroupStateDestroyed  AutoScalingGroupState = "DESTROYED"
	AutoScalingGroupStateDestroying AutoScalingGroupState = "DESTROYING"
	AutoScalingGroupStateScaling    AutoScalingGroupState = "SCALING"
	AutoScalingGroupStateUpdating   AutoScalingGroupState = "UPDATING"
)

var autoScalingGroupStates = []AutoScalingGroupState{
	AutoScalingGroupStateActive,
	AutoScalingGroupStateCreating,
	AutoScalingGroupStateCreated,
	AutoScalingGroupStateDestroyed,
	AutoScalingGroupStateDestroying,
	AutoScalingGroupStateScaling,
	AutoScalingGroupStateUpdating,
}

func NewAutoScalingGroupState(s string) (AutoScalingGroupState, error) {
	return enum_utils.FindEnumForString(
		s,
		autoScalingGroupStates,
		AutoScalingGroupStateActive,
	)
}
