// Package dedicated_server implements the dedicated_server facade.
package dedicated_server

import (
	"context"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/enum"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/facades/dedicated_server/data_adapters/to_domain_entity"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/facades/dedicated_server/data_adapters/to_resource_model"
	resourceModel "github.com/leaseweb/terraform-provider-leaseweb/internal/provider/resources/dedicated_server/model"

	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/ports"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/facades/dedicated_server/data_adapters/to_data_source_model"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/facades/shared"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/data_sources/dedicated_server/model"
)

// DedicatedServerFacade handles all communication between provider & the core.
type DedicatedServerFacade struct {
	dedicatedServerService ports.DedicatedServerService
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

	dataSourceDedicatedServers := to_data_source_model.AdaptDedicatedServers(dedicatedServers)

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

// CreateNotificationSettingBandwidth creates notification setting bandwidth.
func (f DedicatedServerFacade) CreateNotificationSettingBandwidth(
	plan resourceModel.NotificationSettingBandwidth,
	ctx context.Context,
) (*resourceModel.NotificationSettingBandwidth, *shared.FacadeError) {

	createNotificationSettingBandwidthOpts := to_domain_entity.AdaptToCreateNotificationSettingBandwidthOpts(
		plan,
	)

	createNotificationSettingBandwidth, serviceErr := f.dedicatedServerService.CreateNotificationSettingBandwidth(
		*createNotificationSettingBandwidthOpts,
		ctx,
	)
	if serviceErr != nil {
		return nil, shared.NewFromServicesError("CreateNotificationSettingBandwidth", serviceErr)
	}

	notificationSettingBandwidth, err := to_resource_model.AdaptNotificationSettingBandwidth(ctx, *createNotificationSettingBandwidth)
	if err != nil {
		return nil, shared.NewError("CreateNotificationSettingBandwidth", err)
	}

	return notificationSettingBandwidth, nil
}

// GetFrequencies returns a list of valid frequencies.
func (f DedicatedServerFacade) GetFrequencies() []string {
	return enum.NotificationSettingBandwidthFrequencyDaily.Values()
}

// GetUnits returns a list of valid units.
func (f DedicatedServerFacade) GetUnits() []string {
	return enum.NotificationSettingBandwidthUnitMbps.Values()
}

func New(dedicatedServerService ports.DedicatedServerService) DedicatedServerFacade {
	return DedicatedServerFacade{
		dedicatedServerService: dedicatedServerService,
	}
}
