// Package public_cloud implements the public_cloud facade.
package public_cloud

import (
	"context"
	"errors"
	"log"
	"slices"

	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain/public_cloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/ports"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/enum"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/value_object"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/facades/shared"
	dataSourceModel "github.com/leaseweb/terraform-provider-leaseweb/internal/provider/data_sources/public_cloud/model"
	resourceModel "github.com/leaseweb/terraform-provider-leaseweb/internal/provider/resources/public_cloud/model"
)

var ErrContractTermCannotBeZero = public_cloud.ErrContractTermCannotBeZero
var ErrContractTermMustBeZero = public_cloud.ErrContractTermMustBeZero

type CannotBeTerminatedReason string

// PublicCloudFacade handles all communication between provider & the core.
type PublicCloudFacade struct {
	publicCloudService ports.PublicCloudService
}

// GetAllInstances retrieve all instances.
func (p PublicCloudFacade) GetAllInstances(ctx context.Context) (
	*dataSourceModel.Instances,
	*shared.FacadeError,
) {
	instances, err := p.publicCloudService.GetAllInstances(ctx)
	if err != nil {
		return nil, shared.NewFromServicesError("GetAllInstances", err)
	}

	return &instances, nil
}

// LaunchInstance creates an instance.
func (p PublicCloudFacade) LaunchInstance(
	plan resourceModel.Instance,
	ctx context.Context,
) (*resourceModel.Instance, *shared.FacadeError) {
	instance, err := p.publicCloudService.LaunchInstance(plan, ctx)
	if err != nil {
		return nil, shared.NewFromServicesError("LaunchInstance", err)
	}

	return instance, nil
}

// DeleteInstance deletes an instance.
func (p PublicCloudFacade) DeleteInstance(
	id string,
	ctx context.Context,
) *shared.FacadeError {
	err := p.publicCloudService.DeleteInstance(id, ctx)
	if err != nil {
		return shared.NewFromServicesError("DeleteInstance", err)
	}

	return nil
}

// GetInstance returns instance details.
func (p PublicCloudFacade) GetInstance(
	id string,
	ctx context.Context,
) (*resourceModel.Instance, *shared.FacadeError) {
	instance, err := p.publicCloudService.GetInstance(id, ctx)
	if err != nil {
		return nil, shared.NewFromServicesError("GetInstance", err)
	}

	return instance, nil
}

// UpdateInstance updates an instance.
func (p PublicCloudFacade) UpdateInstance(
	plan resourceModel.Instance,
	ctx context.Context,
) (*resourceModel.Instance, *shared.FacadeError) {
	instance, err := p.publicCloudService.UpdateInstance(plan, ctx)
	if err != nil {
		return nil, shared.NewFromServicesError(
			"UpdateInstance",
			err,
		)
	}

	return instance, nil
}

// GetSshKeyRegularExpression returns regular expression used to validate ssh keys.
func (p PublicCloudFacade) GetSshKeyRegularExpression() string {
	return value_object.SshRegexp
}

// GetMinimumRootDiskSize returns the minimal valid rootDiskSize.
func (p PublicCloudFacade) GetMinimumRootDiskSize() int64 {
	return int64(value_object.MinRootDiskSize)
}

// GetMaximumRootDiskSize returns the maximum valid rootDiskSize.
func (p PublicCloudFacade) GetMaximumRootDiskSize() int64 {
	return int64(value_object.MaxRootDiskSize)
}

// GetRootDiskStorageTypes returns a list of valid rootDiskStorageTypes.
func (p PublicCloudFacade) GetRootDiskStorageTypes() []string {
	return enum.StorageTypeCentral.Values()
}

// GetBillingFrequencies returns a list of valid billing frequencies.
func (p PublicCloudFacade) GetBillingFrequencies() shared.IntMarkdownList {
	return shared.NewIntMarkdownList(enum.ContractBillingFrequencyThree.Values())
}

// GetContractTerms returns a list of valid contract terms.
func (p PublicCloudFacade) GetContractTerms() shared.IntMarkdownList {
	return shared.NewIntMarkdownList(enum.ContractTermThree.Values())
}

// GetContractTypes returns a list of valid contract types.
func (p PublicCloudFacade) GetContractTypes() []string {
	return enum.ContractTypeHourly.Values()
}

// ValidateContractTerm checks if the passed combination of contractTerm & contractType is valid.
func (p PublicCloudFacade) ValidateContractTerm(
	contractTerm int64,
	contractType string,
) error {

	contractTermEnum, err := enum.NewContractTerm(int(contractTerm))
	if err != nil {
		return shared.NewError("ValidateContractTerm", err)
	}
	contractTypeEnum, err := enum.NewContractType(contractType)
	if err != nil {
		return shared.NewError("ValidateContractType", err)
	}

	_, err = public_cloud.NewContract(
		enum.ContractBillingFrequencySix,
		contractTermEnum,
		contractTypeEnum,
		enum.ContractStateActive,
		nil,
	)

	if err != nil {
		switch {
		case errors.Is(err, public_cloud.ErrContractTermMustBeZero):
			return ErrContractTermMustBeZero
		case errors.Is(err, public_cloud.ErrContractTermCannotBeZero):
			return ErrContractTermCannotBeZero
		default:
			log.Fatal(err)
		}
	}

	return nil
}

// DoesRegionExist checks if the region exists.
func (p PublicCloudFacade) DoesRegionExist(
	region string,
	ctx context.Context,
) (bool, []string, *shared.FacadeError) {
	regions, err := p.publicCloudService.GetRegions(ctx)
	if err != nil {
		return false, nil, shared.NewFromServicesError(
			"DoesRegionExist",
			err,
		)
	}

	if slices.Contains(regions, region) {
		return true, regions, nil
	}

	return false, regions, nil
}

// IsInstanceTypeAvailableForRegion checks if the instanceType is available for the region.
func (p PublicCloudFacade) IsInstanceTypeAvailableForRegion(
	instanceType string,
	region string,
	ctx context.Context,
) (bool, []string, error) {
	instanceTypes, err := p.publicCloudService.GetAvailableInstanceTypesForRegion(
		region,
		ctx,
	)
	if err != nil {
		return false, nil, shared.NewFromServicesError(
			"IsInstanceTypeAvailableForRegion",
			err,
		)
	}

	return slices.Contains(instanceTypes, instanceType), instanceTypes, nil
}

// CanInstanceTypeBeUsedWithInstance checks
// if the passed instanceType can be used with the passed instance. This is
// the case if:
//   - instanceType is equal to currentInstanceType
//   - instanceType is in available instanceTypes returned by service
func (p PublicCloudFacade) CanInstanceTypeBeUsedWithInstance(
	instanceId string,
	currentInstanceType string,
	instanceType string,
	ctx context.Context,
) (bool, []string, error) {
	instanceTypes, err := p.publicCloudService.GetAvailableInstanceTypesForUpdate(
		instanceId,
		ctx,
	)
	if err != nil {
		return false, nil, shared.NewFromServicesError(
			"CanInstanceTypeBeUsedWithInstance",
			err,
		)
	}

	instanceTypes = append(instanceTypes, currentInstanceType)

	return slices.Contains(instanceTypes, instanceType), instanceTypes, nil
}

// CanInstanceBeTerminated determines whether an instance can be terminated.
func (p PublicCloudFacade) CanInstanceBeTerminated(
	instanceId string,
	ctx context.Context,
) (bool, *CannotBeTerminatedReason, error) {
	canBeTerminated, reason, err := p.publicCloudService.CanInstanceBeTerminated(
		instanceId,
		ctx,
	)
	if err != nil {
		return false, nil, shared.NewFromServicesError("CanInstanceBeTerminated", err)
	}

	if reason != nil {
		reason := CannotBeTerminatedReason(*reason)
		return canBeTerminated, &reason, nil
	}

	return canBeTerminated, nil, nil
}

func NewPublicCloudFacade(publicCloudService ports.PublicCloudService) PublicCloudFacade {
	return PublicCloudFacade{
		publicCloudService: publicCloudService,
	}
}
