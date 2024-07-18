package public_cloud

import (
	"context"

	"terraform-provider-leaseweb/internal/core/domain"
	"terraform-provider-leaseweb/internal/core/ports"
	"terraform-provider-leaseweb/internal/core/shared/enum"
	"terraform-provider-leaseweb/internal/core/shared/value_object"
	dataSourceModel "terraform-provider-leaseweb/internal/provider/data_sources/public_cloud/model"
	resourceModel "terraform-provider-leaseweb/internal/provider/resources/public_cloud/model"
)

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

func (h PublicCloudHandler) GetAllInstances(ctx context.Context) (
	*dataSourceModel.Instances,
	error,
) {
	instances, err := h.publicCloudService.GetAllInstances(ctx)
	if err != nil {
		return nil, err
	}

	dataSourceInstances := convertInstancesToDataSourceModel(instances)

	return &dataSourceInstances, nil
}

func (h PublicCloudHandler) CreateInstance(
	plan resourceModel.Instance,
	ctx context.Context,
) (*resourceModel.Instance, error) {

	createInstanceOpts, err := h.convertInstanceResourceModelToCreateInstanceOpts(
		plan,
		ctx,
	)
	if err != nil {
		return nil, err
	}

	createdInstance, err := h.publicCloudService.CreateInstance(
		*createInstanceOpts,
		ctx,
	)
	if err != nil {
		return nil, err
	}

	return h.convertInstanceToResourceModel(*createdInstance, ctx)
}

func (h PublicCloudHandler) DeleteInstance(id string, ctx context.Context) error {
	instanceId, err := value_object.NewUuid(id)
	if err != nil {
		return err
	}

	err = h.publicCloudService.DeleteInstance(*instanceId, ctx)
	if err != nil {
		return err
	}

	return nil
}

func (h PublicCloudHandler) GetAvailableInstanceTypesForUpdate(
	id string,
	ctx context.Context,
) (
	*domain.InstanceTypes, error) {
	instanceId, err := value_object.NewUuid(id)
	if err != nil {
		return nil, err
	}

	instanceTypes, err := h.publicCloudService.GetAvailableInstanceTypesForUpdate(
		*instanceId,
		ctx,
	)
	if err != nil {
		return nil, err
	}

	return &instanceTypes, nil
}

func (h PublicCloudHandler) GetRegions(ctx context.Context) (
	*domain.Regions,
	error,
) {
	regions, err := h.publicCloudService.GetRegions(ctx)
	if err != nil {
		return nil, err
	}

	return &regions, nil
}

func (h PublicCloudHandler) GetInstance(
	id string,
	ctx context.Context,
) (*resourceModel.Instance, error) {
	instanceId, err := value_object.NewUuid(id)
	if err != nil {
		return nil, err
	}

	instance, err := h.publicCloudService.GetInstance(*instanceId, ctx)
	if err != nil {
		return nil, err
	}

	return h.convertInstanceToResourceModel(*instance, ctx)
}

func (h PublicCloudHandler) UpdateInstance(
	plan resourceModel.Instance,
	ctx context.Context,
) (*resourceModel.Instance, error) {

	updateInstanceOpts, err := h.convertInstanceResourceModelToUpdateInstanceOpts(
		plan,
		ctx,
	)
	if err != nil {
		return nil, err
	}

	updatedInstance, err := h.publicCloudService.UpdateInstance(
		*updateInstanceOpts,
		ctx,
	)
	if err != nil {
		return nil, err
	}

	return h.convertInstanceToResourceModel(*updatedInstance, ctx)
}

func (h PublicCloudHandler) GetImageIds() []string {
	return enum.Debian1064Bit.Values()
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
