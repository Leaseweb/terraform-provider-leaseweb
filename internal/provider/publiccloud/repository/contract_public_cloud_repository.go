package repository

import (
	"context"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/shared/repository"
)

// PublicCloudRepositoryContract is used to connect to public_cloud api.
type PublicCloudRepositoryContract interface {
	// GetAllInstances Retrieve all instances from the public cloud api.
	GetAllInstances(ctx context.Context) (
		[]publicCloud.Instance,
		*repository.RepositoryError,
	)

	// GetInstance retrieves instance details from the public cloud api.
	GetInstance(
		id string,
		ctx context.Context,
	) (*publicCloud.InstanceDetails, *repository.RepositoryError)

	// LaunchInstance launches a new instance in the public cloud api.
	LaunchInstance(
		opts publicCloud.LaunchInstanceOpts,
		ctx context.Context,
	) (*publicCloud.Instance, *repository.RepositoryError)

	// UpdateInstance updates an instance in the public cloud api.
	UpdateInstance(
		id string,
		opts publicCloud.UpdateInstanceOpts,
		ctx context.Context,
	) (*publicCloud.InstanceDetails, *repository.RepositoryError)

	// DeleteInstance deletes an instance in the public cloud api.
	DeleteInstance(id string, ctx context.Context) *repository.RepositoryError

	// GetAvailableInstanceTypesForUpdate return all possible instances types that an instance is allowed to upgrade to from the public cloud api.
	GetAvailableInstanceTypesForUpdate(
		id string,
		ctx context.Context,
	) ([]string, *repository.RepositoryError)

	// GetRegions gets a list of all regions from the public cloud api.
	GetRegions(
		ctx context.Context,
	) ([]string, *repository.RepositoryError)

	// GetInstanceTypesForRegion gets all instance types for a specific region.
	GetInstanceTypesForRegion(
		region string,
		ctx context.Context,
	) ([]string, *repository.RepositoryError)
}
