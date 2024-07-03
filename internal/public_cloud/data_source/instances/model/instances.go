package model

import (
	"errors"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
)

var ErrAutoScalingGroupNotFound = errors.New("auto scaling group not found")
var ErrLoadBalancerNotFound = errors.New("load balancer not found")

type Instances struct {
	Instances []instance `tfsdk:"instances"`
}

func (m *Instances) Populate(
	sdkInstances []publicCloud.InstanceDetails,
	sdkAutoScalingGroups []publicCloud.AutoScalingGroupDetails,
	sdkLoadBalancers []publicCloud.LoadBalancerDetails,
) error {
	var loadBalancerDetails *publicCloud.LoadBalancerDetails
	var autoScalingGroupDetails *publicCloud.AutoScalingGroupDetails
	var err error

	for _, sdkInstance := range sdkInstances {
		autoScalingGroupDetails, err = getAutoScalingGroupDetailsForInstance(
			sdkInstance,
			sdkAutoScalingGroups,
		)
		if err != nil {
			return err
		}

		if autoScalingGroupDetails != nil {
			loadBalancerDetails, err = getLoadBalancerDetailsForAutoScalingGroup(
				*autoScalingGroupDetails,
				sdkLoadBalancers,
			)
			if err != nil {
				return err
			}
		}

		instance := newInstance(sdkInstance, autoScalingGroupDetails, loadBalancerDetails)
		m.Instances = append(m.Instances, instance)
	}

	return nil
}

// Get related autoScalingGroupDetails if it exists.
func getAutoScalingGroupDetailsForInstance(
	sdkInstance publicCloud.InstanceDetails,
	sdkAutoScalingGroups []publicCloud.AutoScalingGroupDetails,
) (*publicCloud.AutoScalingGroupDetails, error) {
	sdkInstanceAutoScalingGroup, _ := sdkInstance.GetAutoScalingGroupOk()
	if sdkInstanceAutoScalingGroup == nil {
		return nil, nil
	}

	for _, sdkAutoScalingGroup := range sdkAutoScalingGroups {
		if sdkAutoScalingGroup.Id == sdkInstanceAutoScalingGroup.Id {
			return &sdkAutoScalingGroup, nil
		}
	}

	return nil, ErrAutoScalingGroupNotFound
}

// Get related loadBalancerDetails if it exists.
func getLoadBalancerDetailsForAutoScalingGroup(
	sdkAutoScalingGroup publicCloud.AutoScalingGroupDetails,
	sdkLoadBalancers []publicCloud.LoadBalancerDetails,
) (*publicCloud.LoadBalancerDetails, error) {
	sdkAutoScalingGroupLoadBalancer, _ := sdkAutoScalingGroup.GetLoadBalancerOk()
	if sdkAutoScalingGroupLoadBalancer == nil {
		return nil, nil
	}

	for _, sdkLoadBalancer := range sdkLoadBalancers {
		if sdkLoadBalancer.Id == sdkAutoScalingGroupLoadBalancer.Id {
			return &sdkLoadBalancer, nil
		}
	}

	return nil, ErrLoadBalancerNotFound
}
