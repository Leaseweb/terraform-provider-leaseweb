package model

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/domain/entity"
	"terraform-provider-leaseweb/internal/core/shared/value_object"
)

func Test_newAutoScalingGroup(t *testing.T) {
	desiredAmount := 1
	createdAt, _ := time.Parse(time.RFC3339, "2019-09-08T00:00:00Z")
	updatedAt, _ := time.Parse(time.RFC3339, "2020-09-08T00:00:00Z")
	startsAt, _ := time.Parse(time.RFC3339, "2010-09-08T00:00:00Z")
	endsAt, _ := time.Parse(time.RFC3339, "2011-09-08T00:00:00Z")
	minimumAmount := 2
	maximumAmount := 3
	cpuThreshold := 4
	warmupTime := 5
	cooldownTime := 6
	id := value_object.NewGeneratedUuid()
	reference, _ := value_object.NewAutoScalingGroupReference("reference")
	loadBalancerId := value_object.NewGeneratedUuid()

	autoScalingGroup := entity.NewAutoScalingGroup(
		id,
		"type",
		"state",
		"region",
		*reference,
		createdAt,
		updatedAt,
		entity.AutoScalingGroupOptions{
			DesiredAmount: &desiredAmount,
			StartsAt:      &startsAt,
			EndsAt:        &endsAt,
			MinimumAmount: &minimumAmount,
			MaximumAmount: &maximumAmount,
			CpuThreshold:  &cpuThreshold,
			WarmupTime:    &warmupTime,
			CoolDownTime:  &cooldownTime,
			LoadBalancer:  &entity.LoadBalancer{Id: loadBalancerId, StartedAt: &time.Time{}},
		},
	)

	got, diags := newAutoScalingGroup(context.TODO(), autoScalingGroup)

	assert.Nil(t, diags)
	assert.Equal(t, id.String(), got.Id.ValueString())
	assert.Equal(t, "type", got.Type.ValueString())
	assert.Equal(t, "state", got.State.ValueString())
	assert.Equal(t, int64(1), got.DesiredAmount.ValueInt64())
	assert.Equal(t, "region", got.Region.ValueString())
	assert.Equal(t, "reference", got.Reference.ValueString())
	assert.Equal(
		t,
		"2019-09-08 00:00:00 +0000 UTC",
		got.CreatedAt.ValueString(),
	)
	assert.Equal(
		t,
		"2020-09-08 00:00:00 +0000 UTC",
		got.UpdatedAt.ValueString(),
	)
	assert.Equal(
		t,
		"2010-09-08 00:00:00 +0000 UTC",
		got.StartsAt.ValueString(),
	)
	assert.Equal(
		t,
		"2011-09-08 00:00:00 +0000 UTC",
		got.EndsAt.ValueString(),
	)
	assert.Equal(t, int64(2), got.MinimumAmount.ValueInt64())
	assert.Equal(t, int64(3), got.MaximumAmount.ValueInt64())
	assert.Equal(t, int64(4), got.CpuThreshold.ValueInt64())
	assert.Equal(t, int64(5), got.WarmupTime.ValueInt64())
	assert.Equal(t, int64(6), got.CooldownTime.ValueInt64())

	loadBalancer := LoadBalancer{}
	got.LoadBalancer.As(
		context.TODO(),
		&loadBalancer,
		basetypes.ObjectAsOptions{},
	)
	assert.Equal(t, loadBalancerId.String(), loadBalancer.Id.ValueString())
}
