package enum

type AutoScalingGroupType string

const (
	AutoScalingGroupTypeManual    AutoScalingGroupType = "ACTIVE"
	AutoScalingGroupTypeScheduled AutoScalingGroupType = "SCHEDULED"
	AutoScalingGroupTypeCpuBased  AutoScalingGroupType = "CPU_BASED"
)
