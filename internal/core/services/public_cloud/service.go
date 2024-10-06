// Package public_cloud implements services related to public_cloud instances
package public_cloud

import (
	"context"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/ports"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/services/errors"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/synced_map"
)

// Service fulfills the contract for ports.PublicCloudService.
type Service struct {
	publicCloudRepository ports.PublicCloudRepository

	cachedInstanceTypes     synced_map.SyncedMap[string, []publicCloud.InstanceType]
	cachedImages            synced_map.SyncedMap[string, publicCloud.ImageDetails]
	cachedRegions           synced_map.SyncedMap[string, []publicCloud.Region]
	cachedAutoScalingGroups synced_map.SyncedMap[string, publicCloud.AutoScalingGroupDetails]
	cachedLoadBalancers     synced_map.SyncedMap[string, publicCloud.LoadBalancerDetails]
}

func (srv *Service) GetAllInstances(ctx context.Context) (
	[]publicCloud.InstanceDetails,
	*errors.ServiceError,
) {
	var detailedInstances []publicCloud.InstanceDetails
	resultChan := make(chan publicCloud.InstanceDetails)
	errorChan := make(chan *errors.ServiceError)

	instanceListResult, err := srv.publicCloudRepository.GetAllInstances(ctx)
	if err != nil {
		return detailedInstances, errors.NewFromRepositoryError(
			"GetAllInstances",
			*err,
		)
	}

	for _, instancePage := range instanceListResult {
		for _, instance := range instancePage.Instances {
			go func(id string) {
				detailedInstance, err := srv.GetInstance(id, ctx)
				if err != nil {
					errorChan <- err
					return
				}
				resultChan <- *detailedInstance
			}(instance.Id)
		}
	}

	for i := 0; i < len(instanceListResult); i++ {
		select {
		case err := <-errorChan:
			return detailedInstances, err
		case res := <-resultChan:
			detailedInstances = append(detailedInstances, res)
		}
	}

	// Order the results as the channel result order is unpredictable.
	// detailedInstances.OrderById()
	return detailedInstances, nil
}

func (srv *Service) GetInstance(
	id string,
	ctx context.Context,
) (*publicCloud.InstanceDetails, *errors.ServiceError) {
	instance, err := srv.publicCloudRepository.GetInstance(id, ctx)
	if err != nil {
		return nil, errors.NewFromRepositoryError("GetInstance", *err)
	}

	return instance, nil
}

func (srv *Service) CreateInstance(
	opts publicCloud.LaunchInstanceOpts,
	ctx context.Context,
) (*publicCloud.InstanceDetails, *errors.ServiceError) {
	createdInstance, err := srv.publicCloudRepository.CreateInstance(opts, ctx)
	if err != nil {
		return nil, errors.NewFromRepositoryError("CreateInstance", *err)
	}

	// call GetInstance as createdInstance is created from instance and not instanceDetails
	return srv.GetInstance(createdInstance.Id, ctx)
}

func (srv *Service) UpdateInstance(
	id string,
	opts publicCloud.UpdateInstanceOpts,
	ctx context.Context,
) (*publicCloud.InstanceDetails, *errors.ServiceError) {
	updatedInstance, err := srv.publicCloudRepository.UpdateInstance(
		opts,
		id,
		ctx,
	)
	if err != nil {
		return nil, errors.NewFromRepositoryError("UpdateInstance", *err)
	}

	return updatedInstance, nil
}

func (srv *Service) DeleteInstance(
	id string,
	ctx context.Context,
) *errors.ServiceError {
	err := srv.publicCloudRepository.DeleteInstance(id, ctx)
	if err != nil {
		return errors.NewError("DeleteInstance", err)
	}

	return nil
}

func (srv *Service) GetAvailableInstanceTypesForUpdate(
	id string,
	ctx context.Context,
) ([]publicCloud.InstanceTypes, *errors.ServiceError) {
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
	[]publicCloud.Region,
	*errors.ServiceError,
) {
	var regions []publicCloud.Region

	cachedRegions, ok := srv.cachedRegions.Get("all")
	if ok {
		return cachedRegions, nil
	}

	regionsListResult, err := srv.publicCloudRepository.GetRegions(ctx)
	if err != nil {
		return nil, errors.NewFromRepositoryError("GetRegions", *err)
	}

	for _, regionList := range regionsListResult {
		for _, region := range regionList.Regions {
			regions = append(regions, region)
		}
	}

	srv.cachedRegions.Set("all", regions)

	return regions, nil
}

// Get autoScalingGroupDetails.
func (srv *Service) getAutoScalingGroup(
	id string,
	ctx context.Context,
) (*publicCloud.AutoScalingGroupDetails, *errors.ServiceError) {
	cachedAutoScalingGroup, ok := srv.cachedAutoScalingGroups.Get(id)
	if ok {
		return &cachedAutoScalingGroup, nil
	}

	autoScalingGroup, err := srv.publicCloudRepository.GetAutoScalingGroup(
		id,
		ctx,
	)
	if err != nil {
		return nil, errors.NewFromRepositoryError(
			"getAutoScalingGroup",
			*err,
		)
	}

	srv.cachedAutoScalingGroups.Set(id, *autoScalingGroup)

	return autoScalingGroup, nil
}

func (srv *Service) getLoadBalancer(
	id string,
	ctx context.Context,
) (*publicCloud.LoadBalancerDetails, *errors.ServiceError) {
	cachedLoadBalancer, ok := srv.cachedLoadBalancers.Get(id)
	if ok {
		return &cachedLoadBalancer, nil
	}

	loadBalancer, err := srv.publicCloudRepository.GetLoadBalancer(id, ctx)
	if err != nil {
		return nil, errors.NewFromRepositoryError(
			"getLoadBalancer",
			*err,
		)
	}

	srv.cachedLoadBalancers.Set(id, *loadBalancer)

	return loadBalancer, nil
}

// Get imageDetails.
func (srv *Service) getImage(
	id string,
	ctx context.Context,
) (*publicCloud.ImageDetails, *errors.ServiceError) {
	cachedImage, ok := srv.cachedImages.Get(id)
	if ok {
		return &cachedImage, nil
	}

	imageListResults, err := srv.publicCloudRepository.GetAllImages(ctx)
	if err != nil {
		return nil, errors.NewFromRepositoryError(
			"getImage",
			*err,
		)
	}

	for _, imageList := range imageListResults {
		for _, image := range imageList.Images {
			srv.cachedImages.Set(image.Id, image)
		}
	}

	image, ok := srv.cachedImages.Get(id)
	if ok {
		return &image, nil
	}

	return nil, nil
}

func (srv *Service) GetAvailableInstanceTypesForRegion(
	region string,
	ctx context.Context,
) ([]publicCloud.InstanceType, *errors.ServiceError) {
	var instanceTypes []publicCloud.InstanceType

	cachedInstanceTypes, ok := srv.cachedInstanceTypes.Get(region)
	if ok {
		return cachedInstanceTypes, nil
	}

	instanceTypesPages, err := srv.publicCloudRepository.GetInstanceTypesForRegion(
		region,
		ctx,
	)
	if err != nil {
		return nil, errors.NewFromRepositoryError(
			"GetAvailableInstanceTypesForRegion",
			*err,
		)
	}

	for _, instanceTypePage := range instanceTypesPages {
		for _, instanceType := range instanceTypePage.GetInstanceTypes() {
			instanceTypes = append(instanceTypes, instanceType)
		}
	}

	srv.cachedInstanceTypes.Set(region, instanceTypes)

	return instanceTypes, nil
}

func (srv *Service) getInstanceType(
	name string,
	region string,
	ctx context.Context,
) (*publicCloud.InstanceType, *errors.ServiceError) {

	instanceTypes, serviceErr := srv.GetAvailableInstanceTypesForRegion(
		region,
		ctx,
	)
	if serviceErr != nil {
		return nil, errors.NewError("GetInstanceType", *serviceErr)
	}

	for _, instanceType := range instanceTypes {
		if string(instanceType.Name) == name {
			return &instanceType, nil
		}
	}

	return nil, nil
}

func (srv *Service) getRegion(
	name string,
	ctx context.Context,
) (*publicCloud.Region, *errors.ServiceError) {
	regions, err := srv.GetRegions(ctx)
	if err != nil {
		return nil, errors.NewError("GetRegion", err)
	}

	for _, region := range regions {
		if string(region.Name) == name {
			return &region, nil
		}
	}

	return nil, nil
}

func New(publicCloudRepository ports.PublicCloudRepository) Service {
	return Service{
		publicCloudRepository:   publicCloudRepository,
		cachedInstanceTypes:     synced_map.NewSyncedMap[string, []publicCloud.InstanceType](),
		cachedImages:            synced_map.NewSyncedMap[string, publicCloud.ImageDetails](),
		cachedRegions:           synced_map.NewSyncedMap[string, []publicCloud.Region](),
		cachedAutoScalingGroups: synced_map.NewSyncedMap[string, publicCloud.AutoScalingGroupDetails](),
		cachedLoadBalancers:     synced_map.NewSyncedMap[string, publicCloud.LoadBalancerDetails](),
	}
}
