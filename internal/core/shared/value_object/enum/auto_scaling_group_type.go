package enum

type AutoScalingGroupType string

func (t AutoScalingGroupType) String() string {
	return string(t)
}

type AutoScalingGroupTypes []AutoScalingGroupType

const (
	AutoScalingCpuTypeManual      AutoScalingGroupType = "MANUAL"
	AutoScalingGroupTypeScheduled AutoScalingGroupType = "SCHEDULED"
	AutoScalingGroupTypeCpuBased  AutoScalingGroupType = "CPU_BASED"
)

var AutoScalingGroupTypeValues = AutoScalingGroupTypes{
	AutoScalingCpuTypeManual,
	AutoScalingGroupTypeScheduled,
	AutoScalingGroupTypeCpuBased,
}
