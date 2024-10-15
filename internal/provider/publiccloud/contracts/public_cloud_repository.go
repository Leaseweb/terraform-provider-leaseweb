package contracts

import (
	"context"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/repositories/shared"
)

// PublicCloudRepository is used to connect to public_cloud api.
type PublicCloudRepository interface {
	// GetAllInstances Retrieve all instances from the public cloud api.
	GetAllInstances(ctx context.Context) (
		[]publicCloud.Instance,
		*shared.RepositoryError,
	)

	// GetInstance retrieves instance details from the public cloud api.
	GetInstance(
		id string,
		ctx context.Context,
	) (*publicCloud.InstanceDetails, *shared.RepositoryError)

	// LaunchInstance launches a new instance in the public cloud api.
	LaunchInstance(
		opts publicCloud.LaunchInstanceOpts,
		ctx context.Context,
	) (*publicCloud.Instance, *shared.RepositoryError)

	// UpdateInstance updates an instance in the public cloud api.
	UpdateInstance(
		id string,
		opts publicCloud.UpdateInstanceOpts,
		ctx context.Context,
	) (*publicCloud.InstanceDetails, *shared.RepositoryError)

	// DeleteInstance deletes an instance in the public cloud api.
	DeleteInstance(id string, ctx context.Context) *shared.RepositoryError

	// GetAvailableInstanceTypesForUpdate return all possible instances types that an instance is allowed to upgrade to from the public cloud api.
	GetAvailableInstanceTypesForUpdate(
		id string,
		ctx context.Context,
	) ([]string, *shared.RepositoryError)

	// GetRegions gets a list of all regions from the public cloud api.
	GetRegions(
		ctx context.Context,
	) ([]string, *shared.RepositoryError)

	// GetInstanceTypesForRegion gets all instance types for a specific region.
	GetInstanceTypesForRegion(
		region string,
		ctx context.Context,
	) ([]string, *shared.RepositoryError)
}
