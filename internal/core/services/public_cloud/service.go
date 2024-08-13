package public_cloud

import (
	"context"

	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain/public_cloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/ports"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/services/errors"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/synced_map"
)

// Service fulfills the contract for ports.PublicCloudService.
type Service struct {
	publicCloudRepository ports.PublicCloudRepository

	cachedInstanceTypes     synced_map.SyncedMap[string, public_cloud.InstanceTypes]
	cachedImages            synced_map.SyncedMap[string, public_cloud.Image]
	cachedRegions           synced_map.SyncedMap[string, public_cloud.Regions]
	cachedAutoScalingGroups synced_map.SyncedMap[string, public_cloud.AutoScalingGroup]
	cachedLoadBalancers     synced_map.SyncedMap[string, public_cloud.LoadBalancer]
}

func (srv *Service) GetAllInstances(ctx context.Context) (
	public_cloud.Instances,
	*errors.ServiceError,
) {
	var detailedInstances public_cloud.Instances
	resultChan := make(chan public_cloud.Instance)
	errorChan := make(chan *errors.ServiceError)

	instances, err := srv.publicCloudRepository.GetAllInstances(ctx)
	if err != nil {
		return public_cloud.Instances{}, errors.NewFromRepositoryError(
			"GetAllInstances",
			*err,
		)
	}

	for _, instance := range instances {
		go func(id string) {
			detailedInstance, err := srv.GetInstance(id, ctx)
			if err != nil {
				errorChan <- err
				return
			}
			resultChan <- *detailedInstance
		}(instance.Id)
	}

	for i := 0; i < len(instances); i++ {
		select {
		case err := <-errorChan:
			return public_cloud.Instances{}, err
		case res := <-resultChan:
			detailedInstances = append(detailedInstances, res)
		}
	}

	return detailedInstances, nil
}

func (srv *Service) GetInstance(
	id string,
	ctx context.Context,
) (*public_cloud.Instance, *errors.ServiceError) {
	instance, err := srv.publicCloudRepository.GetInstance(id, ctx)
	if err != nil {
		return nil, errors.NewFromRepositoryError("GetInstance", *err)
	}

	return srv.populateMissingInstanceAttributes(*instance, ctx)
}

func (srv *Service) CreateInstance(
	instance public_cloud.Instance,
	ctx context.Context,
) (*public_cloud.Instance, *errors.ServiceError) {
	createdInstance, err := srv.publicCloudRepository.CreateInstance(
		instance,
		ctx,
	)
	if err != nil {
		return nil, errors.NewFromRepositoryError("CreateInstance", *err)
	}

	// call GetInstance as createdInstance is created from instance and not instanceDetails
	return srv.GetInstance(createdInstance.Id, ctx)
}

func (srv *Service) UpdateInstance(
	instance public_cloud.Instance,
	ctx context.Context,
) (*public_cloud.Instance, *errors.ServiceError) {
	updatedInstance, err := srv.publicCloudRepository.UpdateInstance(
		instance,
		ctx,
	)
	if err != nil {
		return nil, errors.NewFromRepositoryError("UpdateInstance", *err)
	}

	return srv.populateMissingInstanceAttributes(*updatedInstance, ctx)
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
) (public_cloud.InstanceTypes, *errors.ServiceError) {
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
	public_cloud.Regions,
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

// Get autoScalingGroupDetails.
func (srv *Service) getAutoScalingGroup(
	id string,
	ctx context.Context,
) (*public_cloud.AutoScalingGroup, *errors.ServiceError) {
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

	// Get loadBalancerDetails.
	if autoScalingGroup.LoadBalancer != nil {
		loadBalancer, err := srv.getLoadBalancer(
			autoScalingGroup.LoadBalancer.Id,
			ctx,
		)
		if err != nil {
			return nil, errors.NewError("getAutoScalingGroup", *err)
		}
		autoScalingGroup.LoadBalancer = loadBalancer
	}

	srv.cachedAutoScalingGroups.Set(id, *autoScalingGroup)

	return autoScalingGroup, nil
}

func (srv *Service) getLoadBalancer(
	id string,
	ctx context.Context,
) (*public_cloud.LoadBalancer, *errors.ServiceError) {
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
) (*public_cloud.Image, *errors.ServiceError) {
	cachedImage, ok := srv.cachedImages.Get(id)
	if ok {
		return &cachedImage, nil
	}

	images, err := srv.publicCloudRepository.GetAllImages(ctx)
	if err != nil {
		return nil, errors.NewFromRepositoryError(
			"getImage",
			*err,
		)
	}

	for _, image := range images {
		srv.cachedImages.Set(id, image)
	}

	image, imageErr := images.FilterById(id)
	if imageErr != nil {
		return nil, errors.NewError(
			"getImage",
			imageErr,
		)
	}

	return image, nil
}

// Populate instance with missing details.
func (srv *Service) populateMissingInstanceAttributes(
	instance public_cloud.Instance,
	ctx context.Context,
) (*public_cloud.Instance, *errors.ServiceError) {
	if instance.AutoScalingGroup != nil {
		autoScalingGroup, err := srv.getAutoScalingGroup(
			instance.AutoScalingGroup.Id,
			ctx,
		)
		if err != nil {
			return nil, err
		}
		instance.AutoScalingGroup = autoScalingGroup
	}

	image, err := srv.getImage(instance.Image.Id, ctx)
	if err != nil {
		return nil, err
	}
	instance.Image = *image

	instanceType, err := srv.getInstanceType(
		instance.Type.Name,
		instance.Region,
		ctx,
	)
	if err != nil {
		return nil, err
	}
	instance.Type = *instanceType

	return &instance, nil
}

func (srv *Service) GetAvailableInstanceTypesForRegion(
	region string,
	ctx context.Context,
) (public_cloud.InstanceTypes, *errors.ServiceError) {
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

func (srv *Service) getInstanceType(
	name string,
	region string,
	ctx context.Context,
) (*public_cloud.InstanceType, *errors.ServiceError) {

	instanceTypes, serviceErr := srv.GetAvailableInstanceTypesForRegion(
		region,
		ctx,
	)
	if serviceErr != nil {
		return nil, errors.NewError("GetInstanceType", *serviceErr)
	}

	instanceType, err := instanceTypes.GetByName(name)
	if err != nil {
		return nil, errors.NewError("GetInstanceType", err)
	}

	return instanceType, nil
}

func New(publicCloudRepository ports.PublicCloudRepository) Service {
	return Service{
		publicCloudRepository:   publicCloudRepository,
		cachedInstanceTypes:     synced_map.NewSyncedMap[string, public_cloud.InstanceTypes](),
		cachedImages:            synced_map.NewSyncedMap[string, public_cloud.Image](),
		cachedRegions:           synced_map.NewSyncedMap[string, public_cloud.Regions](),
		cachedAutoScalingGroups: synced_map.NewSyncedMap[string, public_cloud.AutoScalingGroup](),
		cachedLoadBalancers:     synced_map.NewSyncedMap[string, public_cloud.LoadBalancer](),
	}
}
