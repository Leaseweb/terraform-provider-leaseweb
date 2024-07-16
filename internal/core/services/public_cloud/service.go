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
	var detailedInstances domain.Instances

	instances, err := srv.publicCloudRepository.GetAllInstances(ctx)
	if err != nil {
		return domain.Instances{}, fmt.Errorf("GetALlInstances: %w", err)
	}

	// Get instance details.
	for _, instance := range instances {
		detailedInstance, err := srv.GetInstance(instance.Id, ctx)
		if err != nil {
			return domain.Instances{}, fmt.Errorf("GetallAllInstances: %w", err)
		}

		detailedInstances = append(detailedInstances, *detailedInstance)
	}

	return detailedInstances, nil
}

func (srv Service) GetInstance(
	id value_object.Uuid,
	ctx context.Context,
) (*domain.Instance, error) {
	instance, err := srv.publicCloudRepository.GetInstance(id, ctx)
	if err != nil {
		return nil, fmt.Errorf("GetInstance: %w", err)
	}

	return srv.populateMissingInstanceAttributes(*instance, ctx)
}

func (srv Service) CreateInstance(
	instance domain.Instance,
	ctx context.Context,
) (*domain.Instance, error) {
	createdInstance, err := srv.publicCloudRepository.CreateInstance(instance, ctx)
	if err != nil {
		return nil, fmt.Errorf("CreateInstance: %w", err)
	}

	// call GetInstance as createdInstance is created from instance and not instanceDetails
	return srv.GetInstance(createdInstance.Id, ctx)
}

func (srv Service) UpdateInstance(
	instance domain.Instance,
	ctx context.Context,
) (*domain.Instance, error) {
	updatedInstance, err := srv.publicCloudRepository.UpdateInstance(instance, ctx)
	if err != nil {
		return nil, fmt.Errorf("UpdateInstance: %w", err)
	}

	return srv.populateMissingInstanceAttributes(*updatedInstance, ctx)
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

// Populate instance with autoScalingGroupDetails & loadBalancerDetails
func (srv Service) populateMissingInstanceAttributes(
	instance domain.Instance,
	ctx context.Context,
) (*domain.Instance, error) {
	// Get autoScalingGroupDetails.
	if instance.AutoScalingGroup != nil {
		autoScalingGroup, err := srv.publicCloudRepository.GetAutoScalingGroup(
			instance.AutoScalingGroup.Id,
			ctx,
		)
		if err != nil {
			return nil, fmt.Errorf("populateMissingInstanceAttributes: %w", err)
		}

		// Get loadBalancerDetails.
		if autoScalingGroup.LoadBalancer != nil {
			loadBalancer, err := srv.publicCloudRepository.GetLoadBalancer(
				instance.AutoScalingGroup.Id,
				ctx,
			)
			if err != nil {
				return nil, fmt.Errorf("populateMissingInstanceAttributes: %w", err)
			}
			autoScalingGroup.LoadBalancer = loadBalancer
		}

		instance.AutoScalingGroup = autoScalingGroup
	}

	return &instance, nil
}

func New(publicCloudRepository ports.PublicCloudRepository) Service {
	return Service{publicCloudRepository: publicCloudRepository}
}
