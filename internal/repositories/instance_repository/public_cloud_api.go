package instance_repository

import (
	"context"
	"net/http"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
)

type publicCloudApi interface {
	GetInstanceList(ctx context.Context) publicCloud.ApiGetInstanceListRequest

	GetInstance(
		ctx context.Context,
		instanceId string,
	) publicCloud.ApiGetInstanceRequest

	GetInstanceListExecute(r publicCloud.ApiGetInstanceListRequest) (
		*publicCloud.GetInstanceListResult,
		*http.Response,
		error,
	)

	GetInstanceExecute(r publicCloud.ApiGetInstanceRequest) (
		*publicCloud.InstanceDetails,
		*http.Response,
		error,
	)

	GetAutoScalingGroup(
		ctx context.Context,
		autoScalingGroupId string,
	) publicCloud.ApiGetAutoScalingGroupRequest

	GetAutoScalingGroupExecute(r publicCloud.ApiGetAutoScalingGroupRequest) (
		*publicCloud.AutoScalingGroupDetails,
		*http.Response,
		error,
	)

	GetLoadBalancer(
		ctx context.Context,
		loadBalancerId string,
	) publicCloud.ApiGetLoadBalancerRequest

	GetLoadBalancerExecute(r publicCloud.ApiGetLoadBalancerRequest) (
		*publicCloud.LoadBalancerDetails,
		*http.Response,
		error,
	)
}
