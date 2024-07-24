package ports

import (
	"context"

	"terraform-provider-leaseweb/internal/core/domain"
	"terraform-provider-leaseweb/internal/core/services/shared"
	"terraform-provider-leaseweb/internal/core/shared/value_object"
)

type PublicCloudService interface {
	// GetAllInstances Get all instances.
	GetAllInstances(ctx context.Context) (domain.Instances, *shared.ServiceError)

	// GetInstance Get a single instance.
	GetInstance(
		id value_object.Uuid,
		ctx context.Context,
	) (*domain.Instance, *shared.ServiceError)

	// CreateInstance Create an instance.
	CreateInstance(
		instance domain.Instance,
		ctx context.Context,
	) (*domain.Instance, *shared.ServiceError)

	// UpdateInstance Update an instance.
	UpdateInstance(
		instance domain.Instance,
		ctx context.Context,
	) (*domain.Instance, *shared.ServiceError)

	// DeleteInstance Delete an instance.
	DeleteInstance(id value_object.Uuid, ctx context.Context) *shared.ServiceError

	// GetAvailableInstanceTypesForUpdate Get all available instances types an instance can upgrade to.
	GetAvailableInstanceTypesForUpdate(
		id value_object.Uuid,
		ctx context.Context,
	) (domain.InstanceTypes, *shared.ServiceError)

	// GetRegions Get a list of all regions.
	GetRegions(
		ctx context.Context,
	) (domain.Regions, *shared.ServiceError)
}
