package enum

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
	return findEnumForString(s, autoScalingGroupTypes, AutoScalingCpuTypeManual)
}
