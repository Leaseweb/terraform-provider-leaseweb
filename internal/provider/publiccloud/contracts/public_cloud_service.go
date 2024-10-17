package contracts

import (
	"context"

	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/shared/service/errors"
)

// PublicCloudService gets data associated with public_cloud.
type PublicCloudService interface {
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

	// ValidateContractTerm checks if the passed combination of contractTerm & contractType is valid.
	ValidateContractTerm(contractTerm int64, contractType string) error

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
}
