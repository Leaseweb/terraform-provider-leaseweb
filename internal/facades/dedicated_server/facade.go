// Package dedicated_server implements the dedicated_server facade.
package dedicated_server

import (
	"context"

	domain "github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain/dedicated_server"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/ports"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/facades/dedicated_server/data_adapters/to_data_source_model"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/facades/dedicated_server/data_adapters/to_resource_model"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/facades/shared"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/data_sources/dedicated_server/model"
	resourceModel "github.com/leaseweb/terraform-provider-leaseweb/internal/provider/resources/dedicated_server/model"
)

// DedicatedServerFacade handles all communication between provider & the core.
type DedicatedServerFacade struct {
	dedicatedServerService                 ports.DedicatedServerService
	adaptDedicatedServersToDatasourceModel func(dedicatedServers domain.DedicatedServers) model.DedicatedServers
	adaptDedicatedServerToResourceModel    func(dedicatedServer domain.DedicatedServer) resourceModel.DedicatedServer
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

func New(dedicatedServerService ports.DedicatedServerService) DedicatedServerFacade {
	return DedicatedServerFacade{
		dedicatedServerService:                 dedicatedServerService,
		adaptDedicatedServersToDatasourceModel: to_data_source_model.AdaptDedicatedServers,
		adaptDedicatedServerToResourceModel:    to_resource_model.AdaptDedicatedServer,
	}
}
