// Package public_cloud implements services related to public_cloud instances
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

	cachedInstanceTypes synced_map.SyncedMap[string, public_cloud.InstanceTypes]
	cachedRegions       synced_map.SyncedMap[string, public_cloud.Regions]
}

func (srv *Service) GetAllInstances(ctx context.Context) (
	public_cloud.Instances,
	*errors.ServiceError,
) {
	instances, err := srv.publicCloudRepository.GetAllInstances(ctx)
	if err != nil {
		return public_cloud.Instances{}, errors.NewFromRepositoryError(
			"GetAllInstances",
			*err,
		)
	}

	return instances, nil
}

func (srv *Service) GetInstance(
	id string,
	ctx context.Context,
) (*public_cloud.Instance, *errors.ServiceError) {
	instance, err := srv.publicCloudRepository.GetInstance(id, ctx)
	if err != nil {
		return nil, errors.NewFromRepositoryError("GetInstance", *err)
	}

	return instance, nil
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

func New(publicCloudRepository ports.PublicCloudRepository) Service {
	return Service{
		publicCloudRepository: publicCloudRepository,
		cachedInstanceTypes:   synced_map.NewSyncedMap[string, public_cloud.InstanceTypes](),
		cachedRegions:         synced_map.NewSyncedMap[string, public_cloud.Regions](),
	}
}
