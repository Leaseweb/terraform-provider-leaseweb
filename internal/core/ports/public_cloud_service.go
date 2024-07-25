package ports

import (
	"context"

	"terraform-provider-leaseweb/internal/core/domain"
	"terraform-provider-leaseweb/internal/core/services/shared"
	"terraform-provider-leaseweb/internal/core/shared/value_object"
)

// PublicCloudService gets data associated with public_cloud.
type PublicCloudService interface {
	// GetAllInstances gets all instances.
	GetAllInstances(ctx context.Context) (domain.Instances, *shared.ServiceError)

	// GetInstance gets a single instance.
	GetInstance(
		id value_object.Uuid,
		ctx context.Context,
	) (*domain.Instance, *shared.ServiceError)

	// CreateInstance creates an instance.
	CreateInstance(
		instance domain.Instance,
		ctx context.Context,
	) (*domain.Instance, *shared.ServiceError)

	// UpdateInstance updates an instance.
	UpdateInstance(
		instance domain.Instance,
		ctx context.Context,
	) (*domain.Instance, *shared.ServiceError)

	// DeleteInstance deletes an instance.
	DeleteInstance(id value_object.Uuid, ctx context.Context) *shared.ServiceError

	// GetAvailableInstanceTypesForUpdate gets all available instances types an instance can upgrade to.
	GetAvailableInstanceTypesForUpdate(
		id value_object.Uuid,
		ctx context.Context,
	) (domain.InstanceTypes, *shared.ServiceError)

	// GetRegions gets a list of all regions.
	GetRegions(
		ctx context.Context,
	) (domain.Regions, *shared.ServiceError)
}
