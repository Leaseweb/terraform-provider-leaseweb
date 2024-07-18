package ports

import (
	"context"

	"terraform-provider-leaseweb/internal/core/domain"
	"terraform-provider-leaseweb/internal/core/shared/value_object"
	"terraform-provider-leaseweb/internal/repositories/shared"
)

type PublicCloudRepository interface {
	GetAllInstances(ctx context.Context) (domain.Instances, *shared.RepositoryError)

	GetInstance(id value_object.Uuid, ctx context.Context) (*domain.Instance, *shared.RepositoryError)

	CreateInstance(
		instance domain.Instance,
		ctx context.Context,
	) (*domain.Instance, *shared.RepositoryError)

	UpdateInstance(
		instance domain.Instance,
		ctx context.Context,
	) (*domain.Instance, *shared.RepositoryError)

	DeleteInstance(id value_object.Uuid, ctx context.Context) *shared.RepositoryError

	GetAutoScalingGroup(
		id value_object.Uuid,
		ctx context.Context,
	) (*domain.AutoScalingGroup, *shared.RepositoryError)

	GetLoadBalancer(
		id value_object.Uuid,
		ctx context.Context,
	) (*domain.LoadBalancer, *shared.RepositoryError)

	GetAvailableInstanceTypesForUpdate(
		id value_object.Uuid,
		ctx context.Context,
	) (domain.InstanceTypes, *shared.RepositoryError)

	GetRegions(
		ctx context.Context,
	) (domain.Regions, *shared.RepositoryError)
}
