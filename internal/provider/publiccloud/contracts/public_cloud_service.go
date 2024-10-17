package contracts

import (
	"context"

	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/shared/service/errors"
)

// PublicCloudService gets data associated with public_cloud.
type PublicCloudService interface {
	// ValidateContractTerm checks if the passed combination of contractTerm & contractType is valid.
	ValidateContractTerm(contractTerm int64, contractType string) error

	// DoesRegionExist checks if the region exists.
	DoesRegionExist(
		region string,
		ctx context.Context,
	) (bool, []string, *errors.ServiceError)
}
