package ports

import (
	"context"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/repositories/shared"
)

// PublicCloudRepository is used to connect to public_cloud api.
type PublicCloudRepository interface {
	// GetAllInstances Retrieve all instances from the public cloud api.
	GetAllInstances(ctx context.Context) (
		[]publicCloud.GetInstanceListResult,
		*shared.RepositoryError,
	)

	// GetInstance retrieves instance details from the public cloud api.
	GetInstance(
		id string,
		ctx context.Context,
	) (*publicCloud.InstanceDetails, *shared.RepositoryError)

	// CreateInstance creates a new instance in the public cloud api.
	CreateInstance(
		opts publicCloud.LaunchInstanceOpts,
		ctx context.Context,
	) (*publicCloud.Instance, *shared.RepositoryError)

	// UpdateInstance updates an instance in the public cloud api.
	UpdateInstance(
		opts publicCloud.UpdateInstanceOpts,
		id string,
		ctx context.Context,
	) (*publicCloud.InstanceDetails, *shared.RepositoryError)

	// DeleteInstance deletes an instance in the public cloud api.
	DeleteInstance(id string, ctx context.Context) *shared.RepositoryError

	// GetAutoScalingGroup Get autoScalingGroup details from the public cloud api.
	GetAutoScalingGroup(
		id string,
		ctx context.Context,
	) (*publicCloud.AutoScalingGroupDetails, *shared.RepositoryError)

	// GetLoadBalancer gets load balancer details from the public cloud api.
	GetLoadBalancer(
		id string,
		ctx context.Context,
	) (*publicCloud.LoadBalancerDetails, *shared.RepositoryError)

	// GetAvailableInstanceTypesForUpdate gets all possible instances types an instance is allowed to upgrade to from the public cloud api.
	GetAvailableInstanceTypesForUpdate(
		id string,
		ctx context.Context,
	) ([]publicCloud.InstanceTypes, *shared.RepositoryError)

	// GetRegions gets a list of all regions from the public cloud api.
	GetRegions(
		ctx context.Context,
	) ([]publicCloud.GetRegionListResult, *shared.RepositoryError)

	// GetInstanceTypesForRegion gets all instance types for a specific region.
	GetInstanceTypesForRegion(
		region string,
		ctx context.Context,
	) ([]publicCloud.InstanceTypes, *shared.RepositoryError)

	// GetAllImages gets all available images.
	GetAllImages(ctx context.Context) ([]publicCloud.GetImageListResult, *shared.RepositoryError)
}
