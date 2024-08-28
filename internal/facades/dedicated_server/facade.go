package dedicated_server

import (
	"context"

	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/ports"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/facades/dedicated_server/data_adapters/to_data_source_model"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/facades/shared"
	dataSourceModel "github.com/leaseweb/terraform-provider-leaseweb/internal/provider/data_sources/dedicated_server/model"
)

// DedicatedServerFacade handles all communication between provider & the core.
type DedicatedServerFacade struct {
	dedicatedServerService ports.DedicatedServerService
}

// GetAllDedicatedServers retrieves model.DedicatedServers.
func (d DedicatedServerFacade) GetAllDedicatedServers(ctx context.Context) (
	*dataSourceModel.DedicatedServers,
	*shared.FacadeError,
) {
	dedicatedServers, err := d.dedicatedServerService.GetAllDedicatedServers(ctx)
	if err != nil {
		return nil, shared.NewFromServicesError("GetAllDedicatedServers", err)
	}

	dataSourceDedicatedServers := to_data_source_model.AdaptDedicatedServers(dedicatedServers)

	return &dataSourceDedicatedServers, nil
}

// GetAllControlPanels retrieves model.ControlPanels.
func (d DedicatedServerFacade) GetAllControlPanels(ctx context.Context) (
	*dataSourceModel.ControlPanels,
	*shared.FacadeError,
) {
	controlPanels, err := d.dedicatedServerService.GetAllControlPanels(ctx)
	if err != nil {
		return nil, shared.NewFromServicesError("GetAllControlPanels", err)
	}

	dataSourceControlPanels := to_data_source_model.AdaptControlPanels(controlPanels)

	return &dataSourceControlPanels, nil
}

func New(dedicatedServerService ports.DedicatedServerService) DedicatedServerFacade {
	return DedicatedServerFacade{
		dedicatedServerService: dedicatedServerService,
	}
}
