// Package public_cloud implements services related to public_cloud instances
package public_cloud

import (
	"context"
	"fmt"
	"slices"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/ports"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/services/errors"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/services/public_cloud/data_adapters/to_data_source_model"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/services/public_cloud/data_adapters/to_opts"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/services/public_cloud/data_adapters/to_resource_model"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/synced_map"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/value_object"
	dataSourceModel "github.com/leaseweb/terraform-provider-leaseweb/internal/provider/data_sources/public_cloud/model"
	resourceModel "github.com/leaseweb/terraform-provider-leaseweb/internal/provider/resources/public_cloud/model"
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
type Service struct {
	publicCloudRepository ports.PublicCloudRepository
	adaptInstances        func(sdkInstances []publicCloud.Instance) dataSourceModel.Instances
	adaptInstanceDetails  func(
		sdkInstance publicCloud.InstanceDetails,
		ctx context.Context,
	) (*resourceModel.Instance, error)
	adaptInstance func(
		sdkInstance publicCloud.Instance,
		ctx context.Context,
	) (*resourceModel.Instance, error)
	adaptToLaunchInstanceOpts func(
		instance resourceModel.Instance,
		ctx context.Context,
	) (*publicCloud.LaunchInstanceOpts, error)
	adaptToUpdateInstanceOpts func(
		instance resourceModel.Instance,
		ctx context.Context,
	) (*publicCloud.UpdateInstanceOpts, error)
	cachedInstanceTypes synced_map.SyncedMap[string, []string]
	cachedRegions       synced_map.SyncedMap[string, []string]
}

func (srv *Service) GetAllInstances(ctx context.Context) (
	dataSourceModel.Instances,
	*errors.ServiceError,
) {
	instances, err := srv.publicCloudRepository.GetAllInstances(ctx)
	if err != nil {
		return dataSourceModel.Instances{}, errors.NewFromRepositoryError(
			"GetAllInstances",
			*err,
		)
	}

	return srv.adaptInstances(instances), nil
}

func (srv *Service) GetInstance(
	id string,
	ctx context.Context,
) (*resourceModel.Instance, *errors.ServiceError) {
	instance, err := srv.getSdkInstance(id, ctx)
	if err != nil {
		return nil, errors.NewError("GetInstance", *err)
	}

	instanceResourceModel, adaptErr := srv.adaptInstanceDetails(*instance, ctx)
	if adaptErr != nil {
		return nil, errors.NewError("GetInstance", adaptErr)
	}

	return instanceResourceModel, nil
}

func (srv *Service) getSdkInstance(
	id string,
	ctx context.Context,
) (*publicCloud.InstanceDetails, *errors.ServiceError) {
	instance, err := srv.publicCloudRepository.GetInstance(id, ctx)
	if err != nil {
		return nil, errors.NewFromRepositoryError("GetInstance", *err)
	}

	return instance, nil
}

func (srv *Service) LaunchInstance(
	plan resourceModel.Instance,
	ctx context.Context,
) (*resourceModel.Instance, *errors.ServiceError) {

	opts, err := srv.adaptToLaunchInstanceOpts(plan, ctx)
	if err != nil {
		return nil, errors.NewError("LaunchInstance", err)
	}

	instance, repositoryErr := srv.publicCloudRepository.LaunchInstance(*opts, ctx)
	if repositoryErr != nil {
		return nil, errors.NewFromRepositoryError("LaunchInstance", *repositoryErr)
	}

	instanceResourceModel, adaptErr := srv.adaptInstance(*instance, ctx)
	if adaptErr != nil {
		return nil, errors.NewError("LaunchInstance", adaptErr)
	}

	return instanceResourceModel, nil
}

func (srv *Service) UpdateInstance(
	plan resourceModel.Instance,
	ctx context.Context,
) (*resourceModel.Instance, *errors.ServiceError) {
	opts, err := srv.adaptToUpdateInstanceOpts(plan, ctx)
	if err != nil {
		return nil, errors.NewError("UpdateInstance", err)
	}

	instance, repositoryErr := srv.publicCloudRepository.UpdateInstance(
		plan.Id.ValueString(),
		*opts,
		ctx,
	)
	if repositoryErr != nil {
		return nil, errors.NewFromRepositoryError("UpdateInstance", *repositoryErr)
	}

	instanceResourceModel, adaptErr := srv.adaptInstanceDetails(*instance, ctx)
	if adaptErr != nil {
		return nil, errors.NewError("UpdateInstance", adaptErr)
	}
	return instanceResourceModel, nil
}

func (srv *Service) DeleteInstance(
	id string,
	ctx context.Context,
) *errors.ServiceError {
	err := srv.publicCloudRepository.DeleteInstance(id, ctx)
	if err != nil {
		return errors.NewFromRepositoryError("DeleteInstance", *err)
	}

	return nil
}

func (srv *Service) GetAvailableInstanceTypesForUpdate(
	id string,
	ctx context.Context,
) ([]string, *errors.ServiceError) {
	instanceTypes, err := srv.publicCloudRepository.GetAvailableInstanceTypesForUpdate(
		id,
		ctx,
	)
	if err != nil {
		return nil, errors.NewFromRepositoryError(
			"GetAvailableInstanceTypesForUpdate",
			*err,
		)
	}

	return instanceTypes, nil
}

func (srv *Service) GetRegions(ctx context.Context) (
	[]string,
	*errors.ServiceError,
) {
	regions, ok := srv.cachedRegions.Get("all")
	if ok {
		return regions, nil
	}

	regions, err := srv.publicCloudRepository.GetRegions(ctx)
	if err != nil {
		return nil, errors.NewFromRepositoryError("GetRegions", *err)
	}

	srv.cachedRegions.Set("all", regions)

	return regions, nil
}

func (srv *Service) GetAvailableInstanceTypesForRegion(
	region string,
	ctx context.Context,
) ([]string, *errors.ServiceError) {
	cachedInstanceTypes, ok := srv.cachedInstanceTypes.Get(region)
	if ok {
		return cachedInstanceTypes, nil
	}

	instanceTypes, err := srv.publicCloudRepository.GetInstanceTypesForRegion(
		region,
		ctx,
	)
	if err != nil {
		return nil, errors.NewFromRepositoryError(
			"populateMissingInstanceAttributes",
			*err,
		)
	}

	srv.cachedInstanceTypes.Set(region, instanceTypes)

	return instanceTypes, nil
}

func (srv *Service) CanInstanceBeTerminated(id string, ctx context.Context) (
	bool,
	*string,
	*errors.ServiceError,
) {
	instance, err := srv.getSdkInstance(id, ctx)
	if err != nil {
		return false, nil, errors.NewError("CanInstanceBeTerminated", err)
	}

	if instance.State == publicCloud.STATE_CREATING || instance.State == publicCloud.STATE_DESTROYING || instance.State == publicCloud.STATE_DESTROYED {
		reason := fmt.Sprintf("state is %q", instance.State)

		return false, &reason, nil
	}

	if instance.Contract.EndsAt.Get() != nil {
		reason := fmt.Sprintf("contract.endsAt is %q", instance.Contract.EndsAt.Get())

		return false, &reason, nil
	}

	return true, nil, nil
}

func (srv *Service) GetBillingFrequencies() shared.IntMarkdownList {
	var convertedBillingFrequencies []int
	// Have to add 0 manually here, no way round it.
	convertedBillingFrequencies = append(convertedBillingFrequencies, 0)

	for _, billingFrequency := range publicCloud.AllowedBillingFrequencyEnumValues {
		convertedBillingFrequencies = append(convertedBillingFrequencies, int(billingFrequency))
	}

	return shared.NewIntMarkdownList(convertedBillingFrequencies)
}

func (srv *Service) GetContractTerms() shared.IntMarkdownList {
	var convertedContractTerms []int

	for _, contractTerm := range publicCloud.AllowedContractTermEnumValues {
		convertedContractTerms = append(convertedContractTerms, int(contractTerm))
	}

	return shared.NewIntMarkdownList(convertedContractTerms)
}

func (srv *Service) GetContractTypes() []string {
	var convertedContractTypes []string

	for _, contractType := range publicCloud.AllowedContractTypeEnumValues {
		convertedContractTypes = append(convertedContractTypes, string(contractType))
	}

	return convertedContractTypes
}

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

func (srv *Service) GetMinimumRootDiskSize() int64 {
	return int64(value_object.MinRootDiskSize)
}

func (srv *Service) GetMaximumRootDiskSize() int64 {
	return int64(value_object.MaxRootDiskSize)
}

func (srv *Service) GetRootDiskStorageTypes() []string {
	var convertedStates []string

	for _, state := range publicCloud.AllowedStorageTypeEnumValues {
		convertedStates = append(convertedStates, string(state))
	}

	return convertedStates
}

func (srv *Service) DoesRegionExist(
	region string,
	ctx context.Context,
) (bool, []string, *errors.ServiceError) {
	regions, err := srv.GetRegions(ctx)
	if err != nil {
		return false, nil, errors.NewError(
			"DoesRegionExist",
			err,
		)
	}

	if slices.Contains(regions, region) {
		return true, regions, nil
	}

	return false, regions, nil
}

func (srv *Service) IsInstanceTypeAvailableForRegion(
	instanceType string,
	region string,
	ctx context.Context,
) (bool, []string, *errors.ServiceError) {
	instanceTypes, err := srv.GetAvailableInstanceTypesForRegion(
		region,
		ctx,
	)
	if err != nil {
		return false, nil, errors.NewError(
			"IsInstanceTypeAvailableForRegion",
			err,
		)
	}

	return slices.Contains(instanceTypes, instanceType), instanceTypes, nil
}

func (srv *Service) CanInstanceTypeBeUsedWithInstance(
	instanceId string,
	currentInstanceType string,
	instanceType string,
	ctx context.Context,
) (bool, []string, *errors.ServiceError) {
	instanceTypes, err := srv.GetAvailableInstanceTypesForUpdate(
		instanceId,
		ctx,
	)
	if err != nil {
		return false, nil, errors.NewError(
			"CanInstanceTypeBeUsedWithInstance",
			err,
		)
	}

	instanceTypes = append(instanceTypes, currentInstanceType)

	return slices.Contains(instanceTypes, instanceType), instanceTypes, nil
}

func New(publicCloudRepository ports.PublicCloudRepository) Service {
	return Service{
		publicCloudRepository:     publicCloudRepository,
		adaptInstances:            to_data_source_model.AdaptInstances,
		adaptInstanceDetails:      to_resource_model.AdaptInstanceDetails,
		adaptInstance:             to_resource_model.AdaptInstance,
		adaptToLaunchInstanceOpts: to_opts.AdaptToLaunchInstanceOpts,
		adaptToUpdateInstanceOpts: to_opts.AdaptToUpdateInstanceOpts,
		cachedInstanceTypes:       synced_map.NewSyncedMap[string, []string](),
		cachedRegions:             synced_map.NewSyncedMap[string, []string](),
	}
}
