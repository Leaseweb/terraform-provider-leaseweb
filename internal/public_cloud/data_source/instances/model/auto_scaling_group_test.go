package model

import (
	"testing"
	"time"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func Test_generateLoadBalancerDetails(t *testing.T) {
	t.Run("sdkLoadBalancerDetails is empty", func(t *testing.T) {
		got := generateLoadBalancerDetails(nil)
		assert.Nil(t, got)
	})

	t.Run("sdkLoadBalancerDetails is set", func(t *testing.T) {
		got := generateLoadBalancerDetails(&publicCloud.LoadBalancerDetails{})
		assert.NotNil(t, got)
	})
}

func Test_newAutoScalingGroup(t *testing.T) {
	desiredAmount := int32(1)
	createdAt, _ := time.Parse(time.RFC3339, "2019-09-08T00:00:00Z")
	updatedAt, _ := time.Parse(time.RFC3339, "2020-09-08T00:00:00Z")
	startsAt, _ := time.Parse(time.RFC3339, "2010-09-08T00:00:00Z")
	endsAt, _ := time.Parse(time.RFC3339, "2011-09-08T00:00:00Z")
	minimumAmount := int32(2)
	maximumAmount := int32(3)
	cpuThreshold := int32(4)
	warmupTime := int32(5)
	cooldownTime := int32(6)

	sdkAutoScalingGroupDetails := publicCloud.NewAutoScalingGroupDetails(
		"id",
		"type",
		"state",
		*publicCloud.NewNullableInt32(&desiredAmount),
		"region",
		"reference",
		createdAt,
		updatedAt,
		*publicCloud.NewNullableTime(&startsAt),
		*publicCloud.NewNullableTime(&endsAt),
		*publicCloud.NewNullableInt32(&minimumAmount),
		*publicCloud.NewNullableInt32(&maximumAmount),
		*publicCloud.NewNullableInt32(&cpuThreshold),
		*publicCloud.NewNullableInt32(&warmupTime),
		*publicCloud.NewNullableInt32(&cooldownTime),
		*publicCloud.NewNullableLoadBalancer(nil),
	)

	sdkLoadBalancerDetails := publicCloud.LoadBalancerDetails{Id: "loadBalancerId"}

	got := newAutoScalingGroup(
		*sdkAutoScalingGroupDetails,
		&sdkLoadBalancerDetails,
	)

	assert.Equal(t, "id", got.Id.ValueString())
	assert.Equal(t, "type", got.Type.ValueString())
	assert.Equal(t, "state", got.State.ValueString())
	assert.Equal(t, int64(1), got.DesiredAmount.ValueInt64())
	assert.Equal(t, "region", got.Region.ValueString())
	assert.Equal(t, "reference", got.Reference.ValueString())
	assert.Equal(t, "2019-09-08 00:00:00 +0000 UTC", got.CreatedAt.ValueString())
	assert.Equal(t, "2020-09-08 00:00:00 +0000 UTC", got.UpdatedAt.ValueString())
	assert.Equal(t, "2010-09-08 00:00:00 +0000 UTC", got.StartsAt.ValueString())
	assert.Equal(t, "2011-09-08 00:00:00 +0000 UTC", got.EndsAt.ValueString())
	assert.Equal(t, int64(2), got.MinimumAmount.ValueInt64())
	assert.Equal(t, int64(3), got.MaximumAmount.ValueInt64())
	assert.Equal(t, int64(4), got.CpuThreshold.ValueInt64())
	assert.Equal(t, int64(5), got.WarmupTime.ValueInt64())
	assert.Equal(t, int64(6), got.CooldownTime.ValueInt64())
	assert.Equal(t, "loadBalancerId", got.LoadBalancer.Id.ValueString())
}
