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
func (h PublicCloudFacade) GetAllInstances(ctx context.Context) (
	*dataSourceModel.Instances,
	*shared.FacadeError,
) {
	instances, err := h.publicCloudService.GetAllInstances(ctx)
	if err != nil {
		return nil, shared.NewFromServicesError("GetAllInstances", err)
	}

	dataSourceInstances := to_data_source_model.AdaptInstances(instances)

	return &dataSourceInstances, nil
}

// CreateInstance creates an instance.
func (h PublicCloudFacade) CreateInstance(
	plan resourceModel.Instance,
	ctx context.Context,
) (*resourceModel.Instance, *shared.FacadeError) {

	availableInstanceTypes, serviceError := h.publicCloudService.GetAvailableInstanceTypesForRegion(
		plan.Region.ValueString(),
		ctx,
	)
	if serviceError != nil {
		return nil, shared.NewError("CreateInstance", serviceError)
	}

	createInstanceOpts, err := h.adaptToCreateInstanceOpts(
		plan,
		availableInstanceTypes.ToArray(),
		ctx,
	)
	if err != nil {
		return nil, shared.NewError("CreateInstance", err)
	}

	createdInstance, serviceErr := h.publicCloudService.CreateInstance(
		*createInstanceOpts,
		ctx,
	)
	if serviceErr != nil {
		return nil, shared.NewFromServicesError("CreateInstance", serviceErr)
	}

	instance, err := h.adaptInstanceToResourceModel(*createdInstance, ctx)
	if err != nil {
		return nil, shared.NewError("CreateInstance", err)
	}

	return instance, nil
}

// DeleteInstance deletes an instance.
func (h PublicCloudFacade) DeleteInstance(
	id string,
	ctx context.Context,
) *shared.FacadeError {
	serviceErr := h.publicCloudService.DeleteInstance(id, ctx)
	if serviceErr != nil {
		return shared.NewFromServicesError("DeleteInstance", serviceErr)
	}

	return nil
}

// GetInstance returns instance details.
func (h PublicCloudFacade) GetInstance(
	id string,
	ctx context.Context,
) (*resourceModel.Instance, *shared.FacadeError) {
	instance, serviceErr := h.publicCloudService.GetInstance(id, ctx)
	if serviceErr != nil {
		return nil, shared.NewFromServicesError("GetInstance", serviceErr)
	}

	convertedInstance, err := h.adaptInstanceToResourceModel(*instance, ctx)
	if err != nil {
		return nil, shared.NewError("GetInstance", err)
	}

	return convertedInstance, nil
}

// UpdateInstance updates an instance.
func (h PublicCloudFacade) UpdateInstance(
	plan resourceModel.Instance,
	ctx context.Context,
) (*resourceModel.Instance, *shared.FacadeError) {
	availableInstanceTypes, repositoryErr := h.publicCloudService.GetAvailableInstanceTypesForUpdate(
		plan.Id.ValueString(),
		ctx,
	)
	if repositoryErr != nil {
		return nil, shared.NewError("UpdateInstance", repositoryErr)
	}

	instance, repositoryErr := h.publicCloudService.GetInstance(
		plan.Id.ValueString(),
		ctx,
	)
	if repositoryErr != nil {
		return nil, shared.NewError("UpdateInstance", repositoryErr)
	}

	updateInstanceOpts, conversionError := h.adaptToUpdateInstanceOpts(
		plan,
		availableInstanceTypes.ToArray(),
		instance.Type.Name,
		ctx,
	)
	if conversionError != nil {
		return nil, shared.NewError("UpdateInstance", conversionError)
	}

	updatedInstance, updateInstanceErr := h.publicCloudService.UpdateInstance(
		*updateInstanceOpts,
		ctx,
	)
	if updateInstanceErr != nil {
		return nil, shared.NewFromServicesError(
			"UpdateInstance",
			updateInstanceErr,
		)
	}

	convertedInstance, conversionError := h.adaptInstanceToResourceModel(
		*updatedInstance,
		ctx,
	)
	if conversionError != nil {
		return nil, shared.NewError("UpdateInstance", conversionError)
	}

	return convertedInstance, nil
}

// GetSshKeyRegularExpression returns regular expression used to validate ssh keys.
func (h PublicCloudFacade) GetSshKeyRegularExpression() string {
	return value_object.SshRegexp
}

// GetMinimumRootDiskSize returns the minimal valid rootDiskSize.
func (h PublicCloudFacade) GetMinimumRootDiskSize() int64 {
	return int64(value_object.MinRootDiskSize)
}

// GetMaximumRootDiskSize returns the maximum valid rootDiskSize.
func (h PublicCloudFacade) GetMaximumRootDiskSize() int64 {
	return int64(value_object.MaxRootDiskSize)
}

// GetRootDiskStorageTypes returns a list of valid rootDiskStorageTypes.
func (h PublicCloudFacade) GetRootDiskStorageTypes() []string {
	return enum.RootDiskStorageTypeCentral.Values()
}

// GetBillingFrequencies returns a list of valid billing frequencies.
func (h PublicCloudFacade) GetBillingFrequencies() []int64 {
	return shared.AdaptIntArrayToInt64Array(
		enum.ContractBillingFrequencyThree.Values(),
	)
}

// GetContractTerms returns a list of valid contract terms.
func (h PublicCloudFacade) GetContractTerms() []int64 {
	return shared.AdaptIntArrayToInt64Array(
		enum.ContractTermThree.Values(),
	)
}

// GetContractTypes returns a list of valid contract types.
func (h PublicCloudFacade) GetContractTypes() []string {
	return enum.ContractTypeHourly.Values()
}

// ValidateContractTerm checks if the passed combination of contractTerm & contractType is valid.
func (h PublicCloudFacade) ValidateContractTerm(
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
func (h PublicCloudFacade) DoesRegionExist(
	region string,
	ctx context.Context,
) (bool, []string, *shared.FacadeError) {
	regions, err := h.publicCloudService.GetRegions(ctx)
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
func (h PublicCloudFacade) IsInstanceTypeAvailableForRegion(
	instanceType string,
	region string,
	ctx context.Context,
) (bool, []string, error) {
	instanceTypes, err := h.publicCloudService.GetAvailableInstanceTypesForRegion(
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
// if the passed instanceType can be used with the passed instance.
func (h PublicCloudFacade) CanInstanceTypeBeUsedWithInstance(
	instanceId string,
	instanceType string,
	ctx context.Context,
) (bool, []string, error) {
	instanceTypes, serviceErr := h.publicCloudService.GetAvailableInstanceTypesForUpdate(
		instanceId,
		ctx,
	)
	if serviceErr != nil {
		return false, nil, shared.NewFromServicesError(
			"CanInstanceTypeBeUsedWithInstance",
			serviceErr,
		)
	}

	return instanceTypes.ContainsName(instanceType), instanceTypes.ToArray(), nil
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
