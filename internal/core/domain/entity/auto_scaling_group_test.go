package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/shared/value_object"
	"terraform-provider-leaseweb/internal/core/shared/value_object/enum"
)

func TestNewAutoScalingGroup(t *testing.T) {
	t.Run("required values are set", func(t *testing.T) {
		id := value_object.NewGeneratedUuid()
		createdAt := time.Now()
		updatedAt := time.Now()
		reference, _ := value_object.NewAutoScalingGroupReference("reference")

		got := NewAutoScalingGroup(
			id,
			enum.AutoScalingGroupTypeCpuBased,
			enum.AutoScalingGroupStateScaling,
			"region",
			*reference,
			createdAt,
			updatedAt,
			AutoScalingGroupOptions{})

		assert.Equal(t, id, got.Id)
		assert.Equal(t, enum.AutoScalingGroupTypeCpuBased, got.Type)
		assert.Equal(t, enum.AutoScalingGroupStateScaling, got.State)
		assert.Equal(t, "region", got.Region)
		assert.Equal(t, "reference", got.Reference.String())
		assert.Equal(t, createdAt, got.CreatedAt)
		assert.Equal(t, updatedAt, got.UpdatedAt)

		assert.Nil(t, got.DesiredAmount)
		assert.Nil(t, got.StartsAt)
		assert.Nil(t, got.EndsAt)
		assert.Nil(t, got.MinimumAmount)
		assert.Nil(t, got.MaximumAmount)
		assert.Nil(t, got.CpuThreshold)
		assert.Nil(t, got.WarmupTime)
		assert.Nil(t, got.CooldownTime)
		assert.Nil(t, got.LoadBalancer)
	})

	t.Run("optional values are set", func(t *testing.T) {
		reference, _ := value_object.NewAutoScalingGroupReference("")

		desiredAmount := 5
		startsAt := time.Now()
		endsAt := time.Now()
		minimumAmount := 6
		maximumAmount := 7
		cpuThreshold := 8
		WarmupTime := 9
		CoolDownTime := 10

		got := NewAutoScalingGroup(
			value_object.NewGeneratedUuid(),
			enum.AutoScalingGroupTypeCpuBased,
			enum.AutoScalingGroupStateScaling,
			"region",
			*reference,
			time.Now(),
			time.Now(),
			AutoScalingGroupOptions{
				DesiredAmount: &desiredAmount,
				StartsAt:      &startsAt,
				EndsAt:        &endsAt,
				MinimumAmount: &minimumAmount,
				MaximumAmount: &maximumAmount,
				CpuThreshold:  &cpuThreshold,
				WarmupTime:    &WarmupTime,
				CoolDownTime:  &CoolDownTime,
				LoadBalancer:  &LoadBalancer{Region: "loadBalancerRegion"},
			})

		assert.Equal(t, desiredAmount, *got.DesiredAmount)
		assert.Equal(t, startsAt, *got.StartsAt)
		assert.Equal(t, endsAt, *got.EndsAt)
		assert.Equal(t, minimumAmount, *got.MinimumAmount)
		assert.Equal(t, maximumAmount, *got.MaximumAmount)
		assert.Equal(t, cpuThreshold, *got.CpuThreshold)
		assert.Equal(t, WarmupTime, *got.WarmupTime)
		assert.Equal(t, CoolDownTime, *got.CooldownTime)
		assert.Equal(t, "loadBalancerRegion", got.LoadBalancer.Region)
	})

}
