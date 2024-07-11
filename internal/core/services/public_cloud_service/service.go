package public_cloud_service

import (
	"context"
	"fmt"

	"terraform-provider-leaseweb/internal/core/domain/entity"
	"terraform-provider-leaseweb/internal/core/ports"
	"terraform-provider-leaseweb/internal/core/shared/value_object"
)

type Service struct {
	publicCloudRepository ports.PublicCloudRepository
}

func (srv Service) GetAllInstances(ctx context.Context) (
	entity.Instances,
	error,
) {
	instances, err := srv.publicCloudRepository.GetAllInstances(ctx)
	if err != nil {
		return entity.Instances{}, fmt.Errorf(
			"failed to retrieve instances from repository: %w",
			err,
		)
	}

	return instances, nil
}

func (srv Service) GetInstance(
	id value_object.Uuid,
	ctx context.Context,
) (*entity.Instance, error) {
	instance, err := srv.publicCloudRepository.GetInstance(id, ctx)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to retrieve instance %q from repository of: %w",
			id,
			err,
		)
	}

	return instance, nil
}

func (srv Service) CreateInstance(
	instance entity.Instance,
	ctx context.Context,
) (*entity.Instance, error) {
	createdInstance, err := srv.publicCloudRepository.CreateInstance(instance, ctx)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create instance %q in repository: %w",
			instance.Id,
			err,
		)
	}

	return createdInstance, nil
}

func (srv Service) UpdateInstance(
	instance entity.Instance,
	ctx context.Context,
) (*entity.Instance, error) {
	updatedInstance, err := srv.publicCloudRepository.UpdateInstance(instance, ctx)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to update instance %q in repository: %w",
			instance.Id,
			err,
		)
	}

	return updatedInstance, nil
}

func (srv Service) DeleteInstance(
	id value_object.Uuid,
	ctx context.Context,
) error {
	err := srv.publicCloudRepository.DeleteInstance(id, ctx)
	if err != nil {
		return fmt.Errorf(
			"failed to delete instance %q in repository: %w",
			id,
			err,
		)
	}

	return nil
}

func New(publicCloudRepository ports.PublicCloudRepository) Service {
	return Service{publicCloudRepository: publicCloudRepository}
}
