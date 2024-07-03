package model

import (
	"testing"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func TestInstances_Populate(t *testing.T) {
	t.Run("instance is set properly", func(t *testing.T) {
		instanceDetails := publicCloud.InstanceDetails{
			Id:               "instanceId",
			AutoScalingGroup: *publicCloud.NewNullableAutoScalingGroup(nil),
		}

		instances := Instances{}
		err := instances.Populate(
			[]publicCloud.InstanceDetails{instanceDetails},
			nil,
			nil,
		)

		assert.Nil(t, err)
		assert.Equal(
			t,
			"instanceId",
			instances.Instances[0].Id.ValueString(),
			"instance should be set",
		)
	})

	t.Run("related autoScalingGroupDetails cannot be found", func(t *testing.T) {
		instanceDetails := publicCloud.InstanceDetails{
			Id:               "instanceId",
			AutoScalingGroup: *publicCloud.NewNullableAutoScalingGroup(&publicCloud.AutoScalingGroup{}),
		}

		instances := Instances{}
		err := instances.Populate(
			[]publicCloud.InstanceDetails{instanceDetails},
			nil,
			nil,
		)

		assert.NotNil(t, err)
	})

	t.Run("related loadBalancerDetails cannot be found", func(t *testing.T) {
		instanceDetails := publicCloud.InstanceDetails{
			Id:               "instanceId",
			AutoScalingGroup: *publicCloud.NewNullableAutoScalingGroup(&publicCloud.AutoScalingGroup{Id: "autoScalingGroupId"}),
		}

		instances := Instances{}
		err := instances.Populate(
			[]publicCloud.InstanceDetails{instanceDetails},
			[]publicCloud.AutoScalingGroupDetails{{
				Id:           "autoScalingGroupId",
				LoadBalancer: *publicCloud.NewNullableLoadBalancer(&publicCloud.LoadBalancer{Id: "loadBalancerId"}),
			}},
			nil,
		)

		assert.NotNil(t, err)
	})
}

func Test_getAutoScalingGroupDetailsForInstance(t *testing.T) {
	t.Run(
		"autoScalingGroup not set in instanceDetails",
		func(t *testing.T) {
			got, err := getAutoScalingGroupDetailsForInstance(
				publicCloud.InstanceDetails{},
				nil,
			)
			assert.Nil(t, err)
			assert.Nil(t, got)
		},
	)

	t.Run(
		"autoScalingGroup is set in instanceDetails and is found in autoScalingGroups",
		func(t *testing.T) {
			got, err := getAutoScalingGroupDetailsForInstance(
				publicCloud.InstanceDetails{AutoScalingGroup: *publicCloud.NewNullableAutoScalingGroup(&publicCloud.AutoScalingGroup{Id: "autoScalingGroupId"})},
				[]publicCloud.AutoScalingGroupDetails{{Id: "autoScalingGroupId"}},
			)
			assert.Nil(t, err)
			assert.Equal(
				t,
				&publicCloud.AutoScalingGroupDetails{Id: "autoScalingGroupId"},
				got,
			)
		},
	)

	t.Run(
		"autoScalingGroup is set in instanceDetails and is not found in autoScalingGroups",
		func(t *testing.T) {
			got, err := getAutoScalingGroupDetailsForInstance(
				publicCloud.InstanceDetails{AutoScalingGroup: *publicCloud.NewNullableAutoScalingGroup(&publicCloud.AutoScalingGroup{Id: "autoScalingGroupId"})},
				nil,
			)

			assert.Nil(t, got)
			assert.NotNil(t, err)
		},
	)
}

func Test_getLoadBalancerDetailsForAutoScalingGroup(t *testing.T) {
	t.Run(
		"loadBalancer not set in autoScalingGroupDetails",
		func(t *testing.T) {
			got, err := getLoadBalancerDetailsForAutoScalingGroup(
				publicCloud.AutoScalingGroupDetails{},
				nil,
			)
			assert.Nil(t, err)
			assert.Nil(t, got)
		},
	)

	t.Run(
		"loadBalancer is set in autoScalingGroup and is found in loadBalancers",
		func(t *testing.T) {
			got, err := getLoadBalancerDetailsForAutoScalingGroup(
				publicCloud.AutoScalingGroupDetails{LoadBalancer: *publicCloud.NewNullableLoadBalancer(&publicCloud.LoadBalancer{Id: "loadBalancerId"})},
				[]publicCloud.LoadBalancerDetails{{Id: "loadBalancerId"}},
			)
			assert.Nil(t, err)
			assert.Equal(
				t,
				&publicCloud.LoadBalancerDetails{Id: "loadBalancerId"},
				got,
			)
		},
	)

	t.Run(
		"loadBalancer is set in autoScalingGroup and is not found in loadBalancers",
		func(t *testing.T) {
			got, err := getLoadBalancerDetailsForAutoScalingGroup(
				publicCloud.AutoScalingGroupDetails{LoadBalancer: *publicCloud.NewNullableLoadBalancer(&publicCloud.LoadBalancer{Id: "loadBalancerId"})},
				nil,
			)

			assert.Nil(t, got)
			assert.NotNil(t, err)
		},
	)
}
