package ports

import (
	"context"

	"terraform-provider-leaseweb/internal/core/domain"
	"terraform-provider-leaseweb/internal/core/shared/value_object"
	"terraform-provider-leaseweb/internal/repositories/shared"
)

// PublicCloudRepository Used to connect to public_cloud api.
type PublicCloudRepository interface {
	// GetAllInstances Retrieve all instances from the public cloud api.
	GetAllInstances(ctx context.Context) (
		domain.Instances,
		*shared.RepositoryError,
	)

	// GetInstance Retrieve instance details from the public cloud api.
	GetInstance(
		id value_object.Uuid,
		ctx context.Context,
	) (*domain.Instance, *shared.RepositoryError)

	// CreateInstance Create a new instance in the public cloud api.
	CreateInstance(
		instance domain.Instance,
		ctx context.Context,
	) (*domain.Instance, *shared.RepositoryError)

	// UpdateInstance Update an instance in the public cloud api.
	UpdateInstance(
		instance domain.Instance,
		ctx context.Context,
	) (*domain.Instance, *shared.RepositoryError)

	// DeleteInstance Delete an instance in the public cloud api.
	DeleteInstance(id value_object.Uuid, ctx context.Context) *shared.RepositoryError

	// GetAutoScalingGroup Get autoScalingGroup details from the public cloud api.
	GetAutoScalingGroup(
		id value_object.Uuid,
		ctx context.Context,
	) (*domain.AutoScalingGroup, *shared.RepositoryError)

	// GetLoadBalancer Get load balancer details from the public cloud api.
	GetLoadBalancer(
		id value_object.Uuid,
		ctx context.Context,
	) (*domain.LoadBalancer, *shared.RepositoryError)

	// GetAvailableInstanceTypesForUpdate Get all possible instances types an instance is allowed to upgrade to from the public cloud api.
	GetAvailableInstanceTypesForUpdate(
		id value_object.Uuid,
		ctx context.Context,
	) (domain.InstanceTypes, *shared.RepositoryError)

	// GetRegions Get a list of all regions from the public cloud api.
	GetRegions(
		ctx context.Context,
	) (domain.Regions, *shared.RepositoryError)
}
