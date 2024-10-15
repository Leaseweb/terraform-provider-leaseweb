package ports

import (
	"context"

	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain/public_cloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/services/errors"
	dataSourceModel "github.com/leaseweb/terraform-provider-leaseweb/internal/provider/data_sources/public_cloud/model"
	resourceModel "github.com/leaseweb/terraform-provider-leaseweb/internal/provider/resources/public_cloud/model"
)

// PublicCloudService gets data associated with public_cloud.
type PublicCloudService interface {
	// GetAllInstances gets all instances.
	GetAllInstances(ctx context.Context) (dataSourceModel.Instances, *errors.ServiceError)

	// GetInstance gets a single instance.
	GetInstance(
		id string,
		ctx context.Context,
	) (*resourceModel.Instance, *errors.ServiceError)

	// LaunchInstance creates an instance.
	LaunchInstance(
		plan resourceModel.Instance,
		ctx context.Context,
	) (*resourceModel.Instance, *errors.ServiceError)

	// UpdateInstance updates an instance.
	UpdateInstance(
		plan resourceModel.Instance,
		ctx context.Context,
	) (*resourceModel.Instance, *errors.ServiceError)

	// DeleteInstance deletes an instance.
	DeleteInstance(id string, ctx context.Context) *errors.ServiceError

	// GetAvailableInstanceTypesForUpdate gets all available instances types an instance can upgrade to.
	GetAvailableInstanceTypesForUpdate(
		id string,
		ctx context.Context,
	) ([]string, *errors.ServiceError)

	// GetRegions gets a list of all regions.
	GetRegions(ctx context.Context) ([]string, *errors.ServiceError)

	// GetAvailableInstanceTypesForRegion gets all available instances types for a specific region.
	GetAvailableInstanceTypesForRegion(
		region string,
		ctx context.Context,
	) ([]string, *errors.ServiceError)

	// CanInstanceBeTerminated determines if an instance can be terminated.
	CanInstanceBeTerminated(id string, ctx context.Context) (
		bool,
		*public_cloud.ReasonInstanceCannotBeTerminated,
		*errors.ServiceError,
	)
}
