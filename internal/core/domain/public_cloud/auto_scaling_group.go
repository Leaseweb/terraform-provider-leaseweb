package public_cloud

import (
	"time"

	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/enum"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/value_object"
)

type AutoScalingGroup struct {
	Id            string
	Type          enum.AutoScalingGroupType
	State         enum.AutoScalingGroupState
	Region        string
	Reference     value_object.AutoScalingGroupReference
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DesiredAmount *int
	StartsAt      *time.Time
	EndsAt        *time.Time
	MinimumAmount *int
	MaximumAmount *int
	CpuThreshold  *int
	WarmupTime    *int
	CooldownTime  *int
}

type AutoScalingGroupOptions struct {
	DesiredAmount *int
	StartsAt      *time.Time
	EndsAt        *time.Time
	MinimumAmount *int
	MaximumAmount *int
	CpuThreshold  *int
	WarmupTime    *int
	CoolDownTime  *int
}

func NewAutoScalingGroup(
	id string,
	autoScalingGroupType enum.AutoScalingGroupType,
	state enum.AutoScalingGroupState,
	region string,
	reference value_object.AutoScalingGroupReference,
	createdAt time.Time,
	updatedAt time.Time,
	options AutoScalingGroupOptions,
) AutoScalingGroup {
	autoScalingGroup := AutoScalingGroup{
		Id:            id,
		Type:          autoScalingGroupType,
		State:         state,
		DesiredAmount: nil,
		Region:        region,
		Reference:     reference,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}

	autoScalingGroup.DesiredAmount = options.DesiredAmount
	autoScalingGroup.StartsAt = options.StartsAt
	autoScalingGroup.EndsAt = options.EndsAt
	autoScalingGroup.MinimumAmount = options.MinimumAmount
	autoScalingGroup.MaximumAmount = options.MaximumAmount
	autoScalingGroup.CpuThreshold = options.CpuThreshold
	autoScalingGroup.WarmupTime = options.WarmupTime
	autoScalingGroup.CooldownTime = options.CoolDownTime

	return autoScalingGroup
}
