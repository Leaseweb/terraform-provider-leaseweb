package entity

import (
	"time"

	"github.com/google/uuid"
	"terraform-provider-leaseweb/internal/core/shared/value_object"
	"terraform-provider-leaseweb/internal/core/shared/value_object/enum"
)

type AutoScalingGroup struct {
	Id            uuid.UUID
	Type          enum.AutoScalingGroupType
	State         enum.State
	Region        string
	Reference     value_object.AutoScalingGroupReference
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DesiredAmount *int64
	StartsAt      *time.Time
	EndsAt        *time.Time
	MinimumAmount *int64
	MaximumAmount *int64
	CpuThreshold  *int64
	WarmupTime    *int64
	CooldownTime  *int64
	LoadBalancer  *LoadBalancer
}

type AutoScalingGroupOptions struct {
	DesiredAmount *int64
	StartsAt      *time.Time
	EndsAt        *time.Time
	MinimumAmount *int64
	MaximumAmount *int64
	CpuThreshold  *int64
	WarmupTime    *int64
	CoolDownTime  *int64
	LoadBalancer  *LoadBalancer
}

func NewAutoScalingGroup(
	id uuid.UUID,
	autoScalingGroupType enum.AutoScalingGroupType,
	state enum.State,
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
	autoScalingGroup.LoadBalancer = options.LoadBalancer

	return autoScalingGroup
}
