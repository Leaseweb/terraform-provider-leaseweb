// Package public_cloud implements the public_cloud facade.
package public_cloud

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain/public_cloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/ports"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/enum"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/value_object"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/facades/public_cloud/data_adapters/to_data_source_model"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/facades/public_cloud/data_adapters/to_domain_entity"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/facades/public_cloud/data_adapters/to_resource_model"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/facades/shared"
	dataSourceModel "github.com/leaseweb/terraform-provider-leaseweb/internal/provider/data_sources/public_cloud/model"
	resourceModel "github.com/leaseweb/terraform-provider-leaseweb/internal/provider/resources/public_cloud/model"
)

var ErrContractTermCannotBeZero = public_cloud.ErrContractTermCannotBeZero
var ErrContractTermMustBeZero = public_cloud.ErrContractTermMustBeZero

type CannotBeTerminatedReason string

// PublicCloudFacade handles all communication between provider & the core.
type PublicCloudFacade struct {
	publicCloudService           ports.PublicCloudService
	adaptInstanceToResourceModel func(
		instance public_cloud.Instance,
		ctx context.Context,
	) (*resourceModel.Instance, error)
	adaptInstancesToDataSourceModel func(
		instances public_cloud.Instances,
	) dataSourceModel.Instances
	adaptToCreateInstanceOpts func(
		instance resourceModel.Instance,
		allowedInstanceTypes []string,
		ctx context.Context,
	) (*public_cloud.Instance, error)
	adaptToUpdateInstanceOpts func(
		instance resourceModel.Instance,
		allowedInstanceTypes []string,
		currentInstanceType string,
		ctx context.Context,
	) (*public_cloud.Instance, error)
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

	dataSourceInstances := to_data_source_model.AdaptInstances(instances)

	return &dataSourceInstances, nil
}

// CreateInstance creates an instance.
func (p PublicCloudFacade) CreateInstance(
	plan resourceModel.Instance,
	ctx context.Context,
) (*resourceModel.Instance, *shared.FacadeError) {
	availableInstanceTypes, serviceError := p.publicCloudService.GetAvailableInstanceTypesForRegion(
		plan.Region.ValueString(),
		ctx,
	)
	if serviceError != nil {
		return nil, shared.NewError("CreateInstance", serviceError)
	}

	createInstanceOpts, err := p.adaptToCreateInstanceOpts(
		plan,
		availableInstanceTypes.ToArray(),
		ctx,
	)
	if err != nil {
		return nil, shared.NewError("CreateInstance", err)
	}

	createdInstance, serviceErr := p.publicCloudService.CreateInstance(
		*createInstanceOpts,
		ctx,
	)
	if serviceErr != nil {
		return nil, shared.NewFromServicesError("CreateInstance", serviceErr)
	}

	instance, err := p.adaptInstanceToResourceModel(*createdInstance, ctx)
	if err != nil {
		return nil, shared.NewError("CreateInstance", err)
	}

	return instance, nil
}

// DeleteInstance deletes an instance.
func (p PublicCloudFacade) DeleteInstance(
	id string,
	ctx context.Context,
) *shared.FacadeError {
	serviceErr := p.publicCloudService.DeleteInstance(id, ctx)
	if serviceErr != nil {
		return shared.NewFromServicesError("DeleteInstance", serviceErr)
	}

	return nil
}

// GetInstance returns instance details.
func (p PublicCloudFacade) GetInstance(
	id string,
	ctx context.Context,
) (*resourceModel.Instance, *shared.FacadeError) {
	instance, serviceErr := p.publicCloudService.GetInstance(id, ctx)
	if serviceErr != nil {
		return nil, shared.NewFromServicesError("GetInstance", serviceErr)
	}

	convertedInstance, err := p.adaptInstanceToResourceModel(*instance, ctx)
	if err != nil {
		return nil, shared.NewError("GetInstance", err)
	}

	return convertedInstance, nil
}

// UpdateInstance updates an instance.
func (p PublicCloudFacade) UpdateInstance(
	plan resourceModel.Instance,
	ctx context.Context,
) (*resourceModel.Instance, *shared.FacadeError) {
	availableInstanceTypes, repositoryErr := p.publicCloudService.GetAvailableInstanceTypesForUpdate(
		plan.Id.ValueString(),
		ctx,
	)
	if repositoryErr != nil {
		return nil, shared.NewError("UpdateInstance", repositoryErr)
	}

	instance, repositoryErr := p.publicCloudService.GetInstance(
		plan.Id.ValueString(),
		ctx,
	)
	if repositoryErr != nil {
		return nil, shared.NewError("UpdateInstance", repositoryErr)
	}

	updateInstanceOpts, conversionError := p.adaptToUpdateInstanceOpts(
		plan,
		availableInstanceTypes.ToArray(),
		instance.Type.Name,
		ctx,
	)
	if conversionError != nil {
		return nil, shared.NewError("UpdateInstance", conversionError)
	}

	updatedInstance, updateInstanceErr := p.publicCloudService.UpdateInstance(
		*updateInstanceOpts,
		ctx,
	)
	if updateInstanceErr != nil {
		return nil, shared.NewFromServicesError(
			"UpdateInstance",
			updateInstanceErr,
		)
	}

	convertedInstance, conversionError := p.adaptInstanceToResourceModel(
		*updatedInstance,
		ctx,
	)
	if conversionError != nil {
		return nil, shared.NewError("UpdateInstance", conversionError)
	}

	return convertedInstance, nil
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
		time.Now(),
		time.Now(),
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

	if regions.Contains(region) {
		return true, regions.ToArray(), nil
	}

	return false, regions.ToArray(), nil
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

	return instanceTypes.ContainsName(instanceType), instanceTypes.ToArray(), nil
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
	instanceTypes, serviceErr := p.publicCloudService.GetAvailableInstanceTypesForUpdate(
		instanceId,
		ctx,
	)
	instanceTypes = append(
		instanceTypes,
		public_cloud.InstanceType{Name: currentInstanceType},
	)

	if serviceErr != nil {
		return false, nil, shared.NewFromServicesError(
			"CanInstanceTypeBeUsedWithInstance",
			serviceErr,
		)
	}

	return instanceTypes.ContainsName(instanceType), instanceTypes.ToArray(), nil
}

// CanInstanceBeTerminated determines whether an instance can be terminated.
func (p PublicCloudFacade) CanInstanceBeTerminated(
	instanceId string,
	ctx context.Context,
) (bool, *CannotBeTerminatedReason, error) {

	instance, err := p.publicCloudService.GetInstance(instanceId, ctx)
	if err != nil {
		return false, nil, shared.NewError("CanInstanceBeTerminated", err)
	}

	canBeTerminated, instanceReason := instance.CanBeTerminated()
	if instanceReason != nil {
		reason := CannotBeTerminatedReason(*instanceReason)
		return canBeTerminated, &reason, nil
	}

	return canBeTerminated, nil, nil
}

func NewPublicCloudFacade(publicCloudService ports.PublicCloudService) PublicCloudFacade {
	return PublicCloudFacade{
		publicCloudService:              publicCloudService,
		adaptInstanceToResourceModel:    to_resource_model.AdaptInstance,
		adaptInstancesToDataSourceModel: to_data_source_model.AdaptInstances,
		adaptToCreateInstanceOpts:       to_domain_entity.AdaptToCreateInstanceOpts,
		adaptToUpdateInstanceOpts:       to_domain_entity.AdaptToUpdateInstanceOpts,
	}
}
