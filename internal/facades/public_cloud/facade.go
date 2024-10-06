// Package public_cloud implements the public_cloud facade.
package public_cloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/ports"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/facades/public_cloud/data_adapters"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/facades/shared"
	dataSourceModel "github.com/leaseweb/terraform-provider-leaseweb/internal/provider/data_sources/public_cloud/model"
	resourceModel "github.com/leaseweb/terraform-provider-leaseweb/internal/provider/resources/public_cloud/model"
)

type CannotBeTerminatedReason string

// PublicCloudFacade handles all communication between provider & the core.
type PublicCloudFacade struct {
	publicCloudService           ports.PublicCloudService
	adaptInstanceToResourceModel func(
		instance publicCloud.InstanceDetails,
		ctx context.Context,
	) (*resourceModel.Instance, error)
	adaptInstancesToDataSourceModel func(
		instances []publicCloud.InstanceDetails,
	) dataSourceModel.Instances
	adaptToCreateInstanceOpts func(
		instance resourceModel.Instance,
		ctx context.Context,
	) (*publicCloud.LaunchInstanceOpts, error)
	adaptToUpdateInstanceOpts func(
		instance resourceModel.Instance,
		ctx context.Context,
	) (*publicCloud.UpdateInstanceOpts, error)
}

// GetAllInstances retrieve all instances.
func (p PublicCloudFacade) GetAllInstances(ctx context.Context) (
	*dataSourceModel.Instances,
	*shared.FacadeError,
) {
	instances, err := p.publicCloudService.GetAllInstances(ctx)
	if err != nil {
		return nil, shared.NewFromServicesError("GetAllInstances", err)
	}

	dataSourceInstances := data_adapters.AdaptInstances(instances)

	return &dataSourceInstances, nil
}

// CreateInstance creates an instance.
func (p PublicCloudFacade) CreateInstance(
	plan resourceModel.Instance,
	ctx context.Context,
) (*resourceModel.Instance, *shared.FacadeError) {
	var region resourceModel.Region

	plan.Region.As(ctx, &region, basetypes.ObjectAsOptions{})

	createInstanceOpts, err := p.adaptToCreateInstanceOpts(
		plan,
		ctx,
	)
	if err != nil {
		return nil, shared.NewError("CreateInstance", err)
	}

	createdInstance, serviceErr := p.publicCloudService.CreateInstance(
		*createInstanceOpts,
		ctx,
	)
	if serviceErr != nil {
		return nil, shared.NewFromServicesError("CreateInstance", serviceErr)
	}

	instance, err := p.adaptInstanceToResourceModel(*createdInstance, ctx)
	if err != nil {
		return nil, shared.NewError("CreateInstance", err)
	}

	return instance, nil
}

// DeleteInstance deletes an instance.
func (p PublicCloudFacade) DeleteInstance(
	id string,
	ctx context.Context,
) *shared.FacadeError {
	serviceErr := p.publicCloudService.DeleteInstance(id, ctx)
	if serviceErr != nil {
		return shared.NewFromServicesError("DeleteInstance", serviceErr)
	}

	return nil
}

// GetInstance returns instance details.
func (p PublicCloudFacade) GetInstance(
	id string,
	ctx context.Context,
) (*resourceModel.Instance, *shared.FacadeError) {
	instance, serviceErr := p.publicCloudService.GetInstance(id, ctx)
	if serviceErr != nil {
		return nil, shared.NewFromServicesError("GetInstance", serviceErr)
	}

	convertedInstance, err := p.adaptInstanceToResourceModel(*instance, ctx)
	if err != nil {
		return nil, shared.NewError("GetInstance", err)
	}

	return convertedInstance, nil
}

// UpdateInstance updates an instance.
func (p PublicCloudFacade) UpdateInstance(
	plan resourceModel.Instance,
	ctx context.Context,
) (*resourceModel.Instance, *shared.FacadeError) {
	updateInstanceOpts, conversionError := p.adaptToUpdateInstanceOpts(
		plan,
		ctx,
	)
	if conversionError != nil {
		return nil, shared.NewError("UpdateInstance", conversionError)
	}

	updatedInstance, updateInstanceErr := p.publicCloudService.UpdateInstance(
		plan.Id.ValueString(),
		*updateInstanceOpts,
		ctx,
	)
	if updateInstanceErr != nil {
		return nil, shared.NewFromServicesError(
			"UpdateInstance",
			updateInstanceErr,
		)
	}

	convertedInstance, conversionError := p.adaptInstanceToResourceModel(
		*updatedInstance,
		ctx,
	)
	if conversionError != nil {
		return nil, shared.NewError("UpdateInstance", conversionError)
	}

	return convertedInstance, nil
}

func NewPublicCloudFacade(publicCloudService ports.PublicCloudService) PublicCloudFacade {
	return PublicCloudFacade{
		publicCloudService:              publicCloudService,
		adaptInstanceToResourceModel:    data_adapters.AdaptInstance,
		adaptInstancesToDataSourceModel: data_adapters.AdaptInstances,
		adaptToCreateInstanceOpts:       to_opts.AdaptToCreateInstanceOpts,
		adaptToUpdateInstanceOpts:       to_opts.AdaptToUpdateInstanceOpts,
	}
}
