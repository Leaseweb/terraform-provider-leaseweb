// Package dedicated_server implements the dedicated_server facade.
package dedicated_server

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	domain "github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain/dedicated_server"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/ports"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/facades/dedicated_server/data_adapters/to_data_source_model"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/facades/dedicated_server/data_adapters/to_domain_entity"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/facades/dedicated_server/data_adapters/to_resource_model"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/facades/shared"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/data_sources/dedicated_server/model"
	resourceModel "github.com/leaseweb/terraform-provider-leaseweb/internal/provider/resources/dedicated_server/model"
)

// DedicatedServerFacade handles all communication between provider & the core.
type DedicatedServerFacade struct {
	dedicatedServerService                          ports.DedicatedServerService
	adaptDedicatedServersToDatasourceModel          func(dedicatedServers domain.DedicatedServers) model.DedicatedServers
	adaptDedicatedServerToResourceModel             func(dedicatedServer domain.DedicatedServer) resourceModel.DedicatedServer
	AdaptDataTrafficNotificationSetting             func(severId string, dataTrafficNotificationSetting domain.DataTrafficNotificationSetting) resourceModel.DataTrafficNotificationSetting
	adaptToCreateDataTrafficNotificationSettingOpts func(dataTrafficNotificationSetting resourceModel.DataTrafficNotificationSetting) domain.DataTrafficNotificationSetting
}

// GetAllDedicatedServers retrieves model.DedicatedServers.
func (f DedicatedServerFacade) GetAllDedicatedServers(ctx context.Context) (
	*model.DedicatedServers,
	*shared.FacadeError,
) {
	dedicatedServers, err := f.dedicatedServerService.GetAllDedicatedServers(ctx)
	if err != nil {
		return nil, shared.NewFromServicesError("GetAllDedicatedServers", err)
	}

	dataSourceDedicatedServers := f.adaptDedicatedServersToDatasourceModel(dedicatedServers)

	return &dataSourceDedicatedServers, nil
}

// GetAllOperatingSystems retrieve model.OperatingSystems.
func (f DedicatedServerFacade) GetAllOperatingSystems(ctx context.Context) (
	*model.OperatingSystems,
	*shared.FacadeError,
) {
	operatingSystems, err := f.dedicatedServerService.GetAllOperatingSystems(ctx)

	if err != nil {
		return nil, shared.NewFromServicesError("GetAllOperatingSystems", err)
	}

	dataSourceOperatingSystems := to_data_source_model.AdaptOperatingSystems(operatingSystems)

	return &dataSourceOperatingSystems, nil
}

// GetAllControlPanels retrieves model.ControlPanels.
func (f DedicatedServerFacade) GetAllControlPanels(ctx context.Context) (
	*model.ControlPanels,
	*shared.FacadeError,
) {
	controlPanels, err := f.dedicatedServerService.GetAllControlPanels(ctx)
	if err != nil {
		return nil, shared.NewFromServicesError("GetAllControlPanels", err)
	}

	dataSourceControlPanels := to_data_source_model.AdaptControlPanels(controlPanels)

	return &dataSourceControlPanels, nil
}

// GetDedicatedServer returns dedicated server details.
func (f DedicatedServerFacade) GetDedicatedServer(ctx context.Context, id string) (
	*resourceModel.DedicatedServer,
	*shared.FacadeError,
) {
	dedicatedServer, err := f.dedicatedServerService.GetDedicatedServer(ctx, id)
	if err != nil {
		return nil, shared.NewFromServicesError("GetDedicatedServer", err)
	}

	resourceDedicatedServer := f.adaptDedicatedServerToResourceModel(*dedicatedServer)

	return &resourceDedicatedServer, nil
}

// GetDataTrafficNotificationSetting returns data traffic notification setting details.
func (f DedicatedServerFacade) GetDataTrafficNotificationSetting(ctx context.Context, serverId string, dataTrafficNotificationSettingId string) (
	*resourceModel.DataTrafficNotificationSetting,
	*shared.FacadeError,
) {
	dataTrafficNotificationSetting, err := f.dedicatedServerService.GetDataTrafficNotificationSetting(ctx, serverId, dataTrafficNotificationSettingId)
	if err != nil {
		return nil, shared.NewFromServicesError("GetDataTrafficNotificationSetting", err)
	}

	resourceDataTrafficNotificationSetting := f.AdaptDataTrafficNotificationSetting(serverId, *dataTrafficNotificationSetting)

	return &resourceDataTrafficNotificationSetting, nil
}

// CreateDataTrafficNotificationSetting creates a data traffic notification setting.
func (f DedicatedServerFacade) CreateDataTrafficNotificationSetting(
	ctx context.Context,
	plan resourceModel.DataTrafficNotificationSetting,
) (
	*resourceModel.DataTrafficNotificationSetting,
	*shared.FacadeError,
) {

	createDataTrafficNotificationSettingOpts := f.adaptToCreateDataTrafficNotificationSettingOpts(plan)

	createdDataTrafficNotificationSetting, serviceErr := f.dedicatedServerService.CreateDataTrafficNotificationSetting(
		ctx,
		plan.ServerId.ValueString(),
		createDataTrafficNotificationSettingOpts,
	)
	if serviceErr != nil {
		return nil, shared.NewFromServicesError("CreateDataTrafficNotificationSetting", serviceErr)
	}
	dataTrafficNotificationSetting := f.AdaptDataTrafficNotificationSetting(plan.ServerId.ValueString(), *createdDataTrafficNotificationSetting)

	return &dataTrafficNotificationSetting, nil
}

// UpdateDataTrafficNotificationSetting updates a data traffic notification setting.
func (f DedicatedServerFacade) UpdateDataTrafficNotificationSetting(
	ctx context.Context,
	plan resourceModel.DataTrafficNotificationSetting,
) (
	*resourceModel.DataTrafficNotificationSetting,
	*shared.FacadeError,
) {

	updateDataTrafficNotificationSettingOpts := domain.NewUpdateDataTrafficNotificationSetting(
		plan.Frequency.ValueString(),
		plan.Threshold.ValueString(),
		plan.Unit.ValueString(),
	)

	updatedDataTrafficNotificationSetting, serviceErr := f.dedicatedServerService.UpdateDataTrafficNotificationSetting(
		ctx,
		plan.ServerId.ValueString(),
		plan.Id.ValueString(),
		updateDataTrafficNotificationSettingOpts,
	)
	if serviceErr != nil {
		return nil, shared.NewFromServicesError("UpdateDataTrafficNotificationSetting", serviceErr)
	}

	dataTrafficNotificationSetting := resourceModel.DataTrafficNotificationSetting{
		Id:        plan.Id,
		ServerId:  plan.ServerId,
		Frequency: basetypes.NewStringValue(updatedDataTrafficNotificationSetting.Frequency),
		Threshold: basetypes.NewStringValue(updatedDataTrafficNotificationSetting.Threshold),
		Unit:      basetypes.NewStringValue(updatedDataTrafficNotificationSetting.Unit),
	}

	return &dataTrafficNotificationSetting, nil

}

// DeleteDataTrafficNotificationSetting deletes a data traffic notification setting.
func (f DedicatedServerFacade) DeleteDataTrafficNotificationSetting(ctx context.Context, serverId string, dataTrafficNotificationSettingId string) *shared.FacadeError {
	err := f.dedicatedServerService.DeleteDataTrafficNotificationSetting(ctx, serverId, dataTrafficNotificationSettingId)
	if err != nil {
		return shared.NewFromServicesError("DeleteDataTrafficNotificationSetting", err)
	}
	return nil
}

func New(dedicatedServerService ports.DedicatedServerService) DedicatedServerFacade {
	return DedicatedServerFacade{
		dedicatedServerService:                          dedicatedServerService,
		adaptDedicatedServersToDatasourceModel:          to_data_source_model.AdaptDedicatedServers,
		adaptToCreateDataTrafficNotificationSettingOpts: to_domain_entity.AdaptToCreateDataTrafficNotificationSettingOpts,
		adaptDedicatedServerToResourceModel:             to_resource_model.AdaptDedicatedServer,
		AdaptDataTrafficNotificationSetting:             to_resource_model.AdaptDataTrafficNotificationSetting,
	}
}
