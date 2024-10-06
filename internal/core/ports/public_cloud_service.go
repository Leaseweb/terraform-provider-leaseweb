package ports

import (
	"context"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/services/errors"
)

// PublicCloudService gets data associated with public_cloud.
type PublicCloudService interface {
	// GetAllInstances gets all instances.
	GetAllInstances(ctx context.Context) (
		[]publicCloud.InstanceDetails,
		*errors.ServiceError,
	)

	// GetInstance gets a single instance.
	GetInstance(
		id string,
		ctx context.Context,
	) (*publicCloud.InstanceDetails, *errors.ServiceError)

	// CreateInstance creates an instance.
	CreateInstance(
		opts publicCloud.LaunchInstanceOpts,
		ctx context.Context,
	) (*publicCloud.InstanceDetails, *errors.ServiceError)

	// UpdateInstance updates an instance.
	UpdateInstance(
		id string,
		opts publicCloud.UpdateInstanceOpts,
		ctx context.Context,
	) (*publicCloud.InstanceDetails, *errors.ServiceError)

	// DeleteInstance deletes an instance.
	DeleteInstance(id string, ctx context.Context) *errors.ServiceError

	// GetAvailableInstanceTypesForUpdate gets all available instances types an instance can upgrade to.
	GetAvailableInstanceTypesForUpdate(
		id string,
		ctx context.Context,
	) ([]publicCloud.InstanceTypes, *errors.ServiceError)

	// GetRegions gets a list of all regions.
	GetRegions(
		ctx context.Context,
	) ([]publicCloud.Region, *errors.ServiceError)

	// GetAvailableInstanceTypesForRegion gets all available instances types for a specific region.
	GetAvailableInstanceTypesForRegion(
		region string,
		ctx context.Context,
	) ([]publicCloud.InstanceType, *errors.ServiceError)
}
