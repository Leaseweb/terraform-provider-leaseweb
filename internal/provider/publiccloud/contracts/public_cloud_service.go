package contracts

import (
	"context"

	dataSourceModel "github.com/leaseweb/terraform-provider-leaseweb/internal/provider/publiccloud/models/datasource"
	resourceModel "github.com/leaseweb/terraform-provider-leaseweb/internal/provider/publiccloud/models/resource"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/shared/service"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/shared/service/errors"
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
		*string,
		*errors.ServiceError,
	)

	// GetBillingFrequencies returns a list of valid billing frequencies.
	GetBillingFrequencies() service.IntMarkdownList

	// GetContractTerms returns a list of valid contract terms.
	GetContractTerms() service.IntMarkdownList

	// GetContractTypes returns a list of valid contract types.
	GetContractTypes() []string

	// ValidateContractTerm checks if the passed combination of contractTerm & contractType is valid.
	ValidateContractTerm(contractTerm int64, contractType string) error

	// GetMinimumRootDiskSize returns the minimal valid rootDiskSize.
	GetMinimumRootDiskSize() int64

	// GetMaximumRootDiskSize returns the maximum valid rootDiskSize.
	GetMaximumRootDiskSize() int64

	// GetRootDiskStorageTypes returns a list of valid rootDiskStorageTypes.
	GetRootDiskStorageTypes() []string

	// DoesRegionExist checks if the region exists.
	DoesRegionExist(
		region string,
		ctx context.Context,
	) (bool, []string, *errors.ServiceError)

	// IsInstanceTypeAvailableForRegion checks
	// if the instanceType is available for the region.
	IsInstanceTypeAvailableForRegion(
		instanceType string,
		region string,
		ctx context.Context,
	) (bool, []string, *errors.ServiceError)

	// CanInstanceTypeBeUsedWithInstance checks
	// if the passed instanceType can be used with the passed instance.
	// This is the case if:
	//   - instanceType is equal to currentInstanceType
	//   - instanceType is in available instanceTypes returned by service
	CanInstanceTypeBeUsedWithInstance(
		instanceId string,
		currentInstanceType string,
		instanceType string,
		ctx context.Context,
	) (bool, []string, *errors.ServiceError)
}
