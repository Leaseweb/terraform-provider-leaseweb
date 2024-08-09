package ports

import (
	"context"

	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain/public_cloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/value_object"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/repositories/shared"
)

// PublicCloudRepository is used to connect to public_cloud api.
type PublicCloudRepository interface {
	// GetAllInstances Retrieve all instances from the public cloud api.
	GetAllInstances(ctx context.Context) (
		public_cloud.Instances,
		*shared.RepositoryError,
	)

	// GetInstance retrieves instance details from the public cloud api.
	GetInstance(
		id value_object.Uuid,
		ctx context.Context,
	) (*public_cloud.Instance, *shared.RepositoryError)

	// CreateInstance creates a new instance in the public cloud api.
	CreateInstance(
		instance public_cloud.Instance,
		ctx context.Context,
	) (*public_cloud.Instance, *shared.RepositoryError)

	// UpdateInstance updates an instance in the public cloud api.
	UpdateInstance(
		instance public_cloud.Instance,
		ctx context.Context,
	) (*public_cloud.Instance, *shared.RepositoryError)

	// DeleteInstance deletes an instance in the public cloud api.
	DeleteInstance(id value_object.Uuid, ctx context.Context) *shared.RepositoryError

	// GetAutoScalingGroup Get autoScalingGroup details from the public cloud api.
	GetAutoScalingGroup(
		id value_object.Uuid,
		ctx context.Context,
	) (*public_cloud.AutoScalingGroup, *shared.RepositoryError)

	// GetLoadBalancer gets load balancer details from the public cloud api.
	GetLoadBalancer(
		id value_object.Uuid,
		ctx context.Context,
	) (*public_cloud.LoadBalancer, *shared.RepositoryError)

	// GetAvailableInstanceTypesForUpdate gets all possible instances types an instance is allowed to upgrade to from the public cloud api.
	GetAvailableInstanceTypesForUpdate(
		id value_object.Uuid,
		ctx context.Context,
	) (public_cloud.InstanceTypes, *shared.RepositoryError)

	// GetRegions gets a list of all regions from the public cloud api.
	GetRegions(
		ctx context.Context,
	) (public_cloud.Regions, *shared.RepositoryError)

	// GetInstanceTypesForRegion gets all instance types for a specific region.
	GetInstanceTypesForRegion(
		region string,
		ctx context.Context,
	) (public_cloud.InstanceTypes, *shared.RepositoryError)

	// GetAllImages gets all available images.
	GetAllImages(ctx context.Context) (public_cloud.Images, *shared.RepositoryError)
}
