package public_cloud_repository

import (
	"context"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
)

type publicCloudApi interface {
	GetInstanceList(ctx context.Context) publicCloud.ApiGetInstanceListRequest

	GetInstance(
		ctx context.Context,
		instanceId string,
	) publicCloud.ApiGetInstanceRequest

	GetAutoScalingGroup(
		ctx context.Context,
		autoScalingGroupId string,
	) publicCloud.ApiGetAutoScalingGroupRequest

	GetLoadBalancer(
		ctx context.Context,
		loadBalancerId string,
	) publicCloud.ApiGetLoadBalancerRequest

	LaunchInstance(ctx context.Context) publicCloud.ApiLaunchInstanceRequest

	UpdateInstance(
		ctx context.Context,
		instanceId string,
	) publicCloud.ApiUpdateInstanceRequest

	TerminateInstance(
		ctx context.Context,
		instanceId string,
	) publicCloud.ApiTerminateInstanceRequest

	GetUpdateInstanceTypeList(
		ctx context.Context,
		instanceId string,
	) publicCloud.ApiGetUpdateInstanceTypeListRequest

	GetRegionList(ctx context.Context) publicCloud.ApiGetRegionListRequest
}
