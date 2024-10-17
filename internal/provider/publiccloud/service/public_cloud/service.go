// Package public_cloud implements services related to public_cloud instances
package public_cloud

import (
	"fmt"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/shared/service/errors"
)

var ErrContractTermCannotBeZero = fmt.Errorf(
	"contract.term cannot be 0 when contract.type is %q",
	publicCloud.CONTRACTTYPE_MONTHLY,
)

var ErrContractTermMustBeZero = fmt.Errorf(
	"contract.term must be 0 when contract.type is %q",
	publicCloud.CONTRACTTYPE_HOURLY,
)

// Service fulfills the contract for ports.PublicCloudService.
type Service struct{}

func (srv *Service) ValidateContractTerm(
	contractTerm int64,
	contractType string,
) error {
	contractTermEnum, err := publicCloud.NewContractTermFromValue(int32(contractTerm))
	if err != nil {
		return errors.NewError("ValidateContractTerm", err)
	}
	contractTypeEnum, err := publicCloud.NewContractTypeFromValue(contractType)
	if err != nil {
		return errors.NewError("ValidateContractType", err)
	}

	if *contractTypeEnum == publicCloud.CONTRACTTYPE_MONTHLY && *contractTermEnum == publicCloud.CONTRACTTERM__0 {
		return ErrContractTermCannotBeZero
	}

	if *contractTypeEnum == publicCloud.CONTRACTTYPE_HOURLY && *contractTermEnum != publicCloud.CONTRACTTERM__0 {
		return ErrContractTermMustBeZero
	}

	return nil
}
