package ports

import (
	"context"

	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/services/errors"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/value_object"
)

// PublicCloudService gets data associated with public_cloud.
type PublicCloudService interface {
	// GetAllInstances gets all instances.
	GetAllInstances(ctx context.Context) (domain.Instances, *errors.ServiceError)

	// GetInstance gets a single instance.
	GetInstance(
		id value_object.Uuid,
		ctx context.Context,
	) (*domain.Instance, *errors.ServiceError)

	// CreateInstance creates an instance.
	CreateInstance(
		instance domain.Instance,
		ctx context.Context,
	) (*domain.Instance, *errors.ServiceError)

	// UpdateInstance updates an instance.
	UpdateInstance(
		instance domain.Instance,
		ctx context.Context,
	) (*domain.Instance, *errors.ServiceError)

	// DeleteInstance deletes an instance.
	DeleteInstance(id value_object.Uuid, ctx context.Context) *errors.ServiceError

	// GetAvailableInstanceTypesForUpdate gets all available instances types an instance can upgrade to.
	GetAvailableInstanceTypesForUpdate(
		id value_object.Uuid,
		ctx context.Context,
	) (domain.InstanceTypes, *errors.ServiceError)

	// GetRegions gets a list of all regions.
	GetRegions(
		ctx context.Context,
	) (domain.Regions, *errors.ServiceError)

	// GetAvailableInstanceTypesForRegion gets all available instances types for a specific region.
	GetAvailableInstanceTypesForRegion(
		region string,
		ctx context.Context,
	) (domain.InstanceTypes, *errors.ServiceError)
}
