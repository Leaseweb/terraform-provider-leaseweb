package enum

import (
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/enum_utils"
)

type AutoScalingGroupType string

func (t AutoScalingGroupType) String() string {
	return string(t)
}

const (
	AutoScalingCpuTypeManual      AutoScalingGroupType = "MANUAL"
	AutoScalingGroupTypeScheduled AutoScalingGroupType = "SCHEDULED"
	AutoScalingGroupTypeCpuBased  AutoScalingGroupType = "CPU_BASED"
)

var autoScalingGroupTypes = []AutoScalingGroupType{
	AutoScalingCpuTypeManual,
	AutoScalingGroupTypeScheduled,
	AutoScalingGroupTypeCpuBased,
}

func NewAutoScalingGroupType(s string) (AutoScalingGroupType, error) {
	return enum_utils.FindEnumForString(s, autoScalingGroupTypes, AutoScalingCpuTypeManual)
}
