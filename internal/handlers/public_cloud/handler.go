package public_cloud

import (
	"context"
	"errors"
	"log"
	"time"

	"terraform-provider-leaseweb/internal/core/domain"
	"terraform-provider-leaseweb/internal/core/ports"
	"terraform-provider-leaseweb/internal/core/shared/enum"
	"terraform-provider-leaseweb/internal/core/shared/value_object"
	"terraform-provider-leaseweb/internal/handlers/shared"
	dataSourceModel "terraform-provider-leaseweb/internal/provider/data_sources/public_cloud/model"
	resourceModel "terraform-provider-leaseweb/internal/provider/resources/public_cloud/model"
)

var ErrContractTermCannotBeZero = domain.ErrContractTermCannotBeZero
var ErrContractTermMustBeZero = domain.ErrContractTermMustBeZero

type PublicCloudHandler struct {
	publicCloudService             ports.PublicCloudService
	convertInstanceToResourceModel func(
		instance domain.Instance,
		ctx context.Context,
	) (*resourceModel.Instance, error)
	convertInstancesToDataSourceModel func(
		instances domain.Instances,
	) dataSourceModel.Instances
	convertInstanceResourceModelToCreateInstanceOpts func(
		instance resourceModel.Instance,
		ctx context.Context,
	) (*domain.Instance, error)
	convertInstanceResourceModelToUpdateInstanceOpts func(
		instance resourceModel.Instance,
		ctx context.Context,
	) (*domain.Instance, error)
}

// GetAllInstances Retrieve all instances
func (h PublicCloudHandler) GetAllInstances(ctx context.Context) (
	*dataSourceModel.Instances,
	*shared.HandlerError,
) {
	instances, err := h.publicCloudService.GetAllInstances(ctx)
	if err != nil {
		return nil, shared.NewFromServicesError("GetAllInstances", err)
	}

	dataSourceInstances := convertInstancesToDataSourceModel(instances)

	return &dataSourceInstances, nil
}

// CreateInstance Creates an instance.
func (h PublicCloudHandler) CreateInstance(
	plan resourceModel.Instance,
	ctx context.Context,
) (*resourceModel.Instance, *shared.HandlerError) {

	createInstanceOpts, err := h.convertInstanceResourceModelToCreateInstanceOpts(
		plan,
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

	instance, err := h.convertInstanceToResourceModel(*createdInstance, ctx)
	if err != nil {
		return nil, shared.NewError("CreateInstance", err)
	}

	return instance, nil
}

// DeleteInstance Deletes an instance.
func (h PublicCloudHandler) DeleteInstance(
	id string,
	ctx context.Context,
) *shared.HandlerError {
	instanceId, err := value_object.NewUuid(id)
	if err != nil {
		return shared.NewError("DeleteInstance", err)
	}

	serviceErr := h.publicCloudService.DeleteInstance(*instanceId, ctx)
	if serviceErr != nil {
		return shared.NewFromServicesError("DeleteInstance", serviceErr)
	}

	return nil
}

// GetAvailableInstanceTypesForUpdate Gets all instance types an instance is allowed to upgrade to.
func (h PublicCloudHandler) GetAvailableInstanceTypesForUpdate(
	id string,
	ctx context.Context,
) (*domain.InstanceTypes, *shared.HandlerError) {
	instanceId, err := value_object.NewUuid(id)
	if err != nil {
		return nil, shared.NewError(
			"GetAvailableInstanceTypesForUpdate",
			err,
		)
	}

	instanceTypes, serviceErr := h.publicCloudService.GetAvailableInstanceTypesForUpdate(
		*instanceId,
		ctx,
	)
	if serviceErr != nil {
		return nil, shared.NewFromServicesError(
			"GetAvailableInstanceTypesForUpdate",
			serviceErr,
		)
	}

	return &instanceTypes, nil
}

// GetRegions Return a list of all regions.
func (h PublicCloudHandler) GetRegions(ctx context.Context) (
	*domain.Regions,
	*shared.HandlerError,
) {
	regions, err := h.publicCloudService.GetRegions(ctx)
	if err != nil {
		return nil, shared.NewFromServicesError("GetRegions", err)
	}

	return &regions, nil
}

// GetInstance Return instance details.
func (h PublicCloudHandler) GetInstance(
	id string,
	ctx context.Context,
) (*resourceModel.Instance, *shared.HandlerError) {
	instanceId, err := value_object.NewUuid(id)
	if err != nil {
		return nil, shared.NewError("GetInstance", err)
	}

	instance, serviceErr := h.publicCloudService.GetInstance(*instanceId, ctx)
	if serviceErr != nil {
		return nil, shared.NewFromServicesError("GetInstance", serviceErr)
	}

	convertedInstance, err := h.convertInstanceToResourceModel(*instance, ctx)
	if err != nil {
		return nil, shared.NewError("GetInstance", err)
	}

	return convertedInstance, nil
}

// UpdateInstance Update an instance.
func (h PublicCloudHandler) UpdateInstance(
	plan resourceModel.Instance,
	ctx context.Context,
) (*resourceModel.Instance, *shared.HandlerError) {

	updateInstanceOpts, err := h.convertInstanceResourceModelToUpdateInstanceOpts(
		plan,
		ctx,
	)
	if err != nil {
		return nil, shared.NewError("UpdateInstance", err)
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

	convertedInstance, err := h.convertInstanceToResourceModel(
		*updatedInstance,
		ctx,
	)
	if err != nil {
		return nil, shared.NewError("UpdateInstance", err)
	}

	return convertedInstance, nil
}

// GetImageIds Return a list of valid image ids.
func (h PublicCloudHandler) GetImageIds() []string {
	return enum.Debian1064Bit.Values()
}

// GetSshKeyRegularExpression Returns regular expression used to validate ssh keys.
func (h PublicCloudHandler) GetSshKeyRegularExpression() string {
	return value_object.SshRegexp
}

// GetMinimumRootDiskSize Return the minimal valid rootDiskSize.
func (h PublicCloudHandler) GetMinimumRootDiskSize() int64 {
	return int64(value_object.MinRootDiskSize)
}

// GetMaximumRootDiskSize Return the maximum valid rootDiskSize.
func (h PublicCloudHandler) GetMaximumRootDiskSize() int64 {
	return int64(value_object.MaxRootDiskSize)
}

// GetRootDiskStorageTypes Return a list of valid rootDiskStorageTypes.
func (h PublicCloudHandler) GetRootDiskStorageTypes() []string {
	return enum.RootDiskStorageTypeCentral.Values()
}

// GetBillingFrequencies Return a list of valid billing frequencies.
func (h PublicCloudHandler) GetBillingFrequencies() []int64 {
	return convertIntArrayToInt64(enum.ContractBillingFrequencyThree.Values())
}

// GetContractTerms Return a list of valid contract terms.
func (h PublicCloudHandler) GetContractTerms() []int64 {
	return convertIntArrayToInt64(enum.ContractTermThree.Values())
}

// GetContractTypes Return a list of valid contract types.
func (h PublicCloudHandler) GetContractTypes() []string {
	return enum.ContractTypeHourly.Values()
}

// ValidateContractTerm Checks if the passed combination of contractTerm & contractType is valid.
func (h PublicCloudHandler) ValidateContractTerm(
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

	_, err = domain.NewContract(
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
		case errors.Is(err, domain.ErrContractTermMustBeZero):
			return ErrContractTermMustBeZero
		case errors.Is(err, domain.ErrContractTermCannotBeZero):
			return ErrContractTermCannotBeZero
		default:
			log.Fatal(err)
		}
	}

	return nil
}

func NewPublicCloudHandler(publicCloudService ports.PublicCloudService) PublicCloudHandler {
	return PublicCloudHandler{
		publicCloudService:                               publicCloudService,
		convertInstanceToResourceModel:                   convertInstanceToResourceModel,
		convertInstancesToDataSourceModel:                convertInstancesToDataSourceModel,
		convertInstanceResourceModelToCreateInstanceOpts: convertInstanceResourceModelToCreateInstanceOpts,
		convertInstanceResourceModelToUpdateInstanceOpts: convertInstanceResourceModelToUpdateInstanceOpts,
	}
}
