// Package public_cloud implements services related to public_cloud instances
package public_cloud

import (
	"context"
	"fmt"
	"slices"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/publiccloud/contracts"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/shared/service"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/shared/service/errors"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/shared/service/synced_map"
)

var ErrContractTermCannotBeZero = fmt.Errorf(
	"contract.term cannot be 0 when contract.type is %q",
	publicCloud.CONTRACTTYPE_MONTHLY,
)

var ErrContractTermMustBeZero = fmt.Errorf(
	"contract.term must be 0 when contract.type is %q",
	publicCloud.CONTRACTTYPE_HOURLY,
)

var minimumRootDiskSize = 5
var maximumRootDiskSize = 1000

// Service fulfills the contract for ports.PublicCloudService.
type Service struct {
	publicCloudRepository contracts.PublicCloudRepository
	cachedInstanceTypes   synced_map.SyncedMap[string, []string]
	cachedRegions         synced_map.SyncedMap[string, []string]
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

func (srv *Service) GetBillingFrequencies() service.IntMarkdownList {
	var convertedBillingFrequencies []int
	// Have to add 0 manually here, no way round it.
	convertedBillingFrequencies = append(convertedBillingFrequencies, 0)

	for _, billingFrequency := range publicCloud.AllowedBillingFrequencyEnumValues {
		convertedBillingFrequencies = append(convertedBillingFrequencies, int(billingFrequency))
	}

	return service.NewIntMarkdownList(convertedBillingFrequencies)
}

func (srv *Service) GetContractTerms() service.IntMarkdownList {
	var convertedContractTerms []int

	for _, contractTerm := range publicCloud.AllowedContractTermEnumValues {
		convertedContractTerms = append(convertedContractTerms, int(contractTerm))
	}

	return service.NewIntMarkdownList(convertedContractTerms)
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
	return int64(minimumRootDiskSize)
}

func (srv *Service) GetMaximumRootDiskSize() int64 {
	return int64(maximumRootDiskSize)
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

func New(publicCloudRepository contracts.PublicCloudRepository) Service {
	return Service{
		publicCloudRepository: publicCloudRepository,
		cachedInstanceTypes:   synced_map.NewSyncedMap[string, []string](),
		cachedRegions:         synced_map.NewSyncedMap[string, []string](),
	}
}
