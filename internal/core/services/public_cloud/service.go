package public_cloud

import (
	"context"
	"fmt"

	"terraform-provider-leaseweb/internal/core/domain"
	"terraform-provider-leaseweb/internal/core/ports"
	"terraform-provider-leaseweb/internal/core/shared/value_object"
)

type Service struct {
	publicCloudRepository ports.PublicCloudRepository
}

func (srv Service) GetAllInstances(ctx context.Context) (
	domain.Instances,
	error,
) {
	instances, err := srv.publicCloudRepository.GetAllInstances(ctx)
	if err != nil {
		return domain.Instances{}, fmt.Errorf("GetALlInstances: %w", err)
	}

	return instances, nil
}

func (srv Service) GetInstance(
	id value_object.Uuid,
	ctx context.Context,
) (*domain.Instance, error) {
	instance, err := srv.publicCloudRepository.GetInstance(id, ctx)
	if err != nil {
		return nil, fmt.Errorf("GetInstance: %w", err)
	}

	return instance, nil
}

func (srv Service) CreateInstance(
	instance domain.Instance,
	ctx context.Context,
) (*domain.Instance, error) {
	createdInstance, err := srv.publicCloudRepository.CreateInstance(instance, ctx)
	if err != nil {
		return nil, fmt.Errorf("CreateInstance: %w", err)
	}

	return createdInstance, nil
}

func (srv Service) UpdateInstance(
	instance domain.Instance,
	ctx context.Context,
) (*domain.Instance, error) {
	updatedInstance, err := srv.publicCloudRepository.UpdateInstance(instance, ctx)
	if err != nil {
		return nil, fmt.Errorf("UpdateInstance: %w", err)
	}

	return updatedInstance, nil
}

func (srv Service) DeleteInstance(
	id value_object.Uuid,
	ctx context.Context,
) error {
	err := srv.publicCloudRepository.DeleteInstance(id, ctx)
	if err != nil {
		return fmt.Errorf("DeleteInstance: %w", err)
	}

	return nil
}

func (srv Service) GetAvailableInstanceTypesForUpdate(
	id value_object.Uuid,
	ctx context.Context,
) (domain.InstanceTypes, error) {
	instanceTypes, err := srv.publicCloudRepository.GetAvailableInstanceTypesForUpdate(id, ctx)
	if err != nil {
		return nil, fmt.Errorf("GetAvailableInstanceTypesForUpdate: %w", err)
	}

	return instanceTypes, nil
}

func (srv Service) GetRegions(ctx context.Context) (domain.Regions, error) {
	regions, err := srv.publicCloudRepository.GetRegions(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetRegions: %w", err)
	}

	return regions, nil
}

func New(publicCloudRepository ports.PublicCloudRepository) Service {
	return Service{publicCloudRepository: publicCloudRepository}
}
