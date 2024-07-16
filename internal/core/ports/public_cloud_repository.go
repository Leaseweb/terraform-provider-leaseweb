package ports

import (
	"context"

	"terraform-provider-leaseweb/internal/core/domain"
	"terraform-provider-leaseweb/internal/core/shared/value_object"
)

type PublicCloudRepository interface {
	GetAllInstances(ctx context.Context) (domain.Instances, error)

	GetInstance(id value_object.Uuid, ctx context.Context) (*domain.Instance, error)

	CreateInstance(
		instance domain.Instance,
		ctx context.Context,
	) (*domain.Instance, error)

	UpdateInstance(
		instance domain.Instance,
		ctx context.Context,
	) (*domain.Instance, error)

	DeleteInstance(id value_object.Uuid, ctx context.Context) error

	GetAutoScalingGroup(
		id value_object.Uuid,
		ctx context.Context,
	) (*domain.AutoScalingGroup, error)

	GetLoadBalancer(
		id value_object.Uuid,
		ctx context.Context,
	) (*domain.LoadBalancer, error)

	GetAvailableInstanceTypesForUpdate(
		id value_object.Uuid,
		ctx context.Context,
	) (domain.InstanceTypes, error)

	GetRegions(
		ctx context.Context,
	) (domain.Regions, error)
}
