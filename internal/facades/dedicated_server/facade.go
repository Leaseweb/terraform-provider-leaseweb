package dedicated_server

import (
	"context"

	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/ports"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/facades/dedicated_server/data_adapters/to_data_source_model"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/facades/shared"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/data_sources/dedicated_server/model"
)

// DedicatedServerFacade handles all communication between provider & the core.
type DedicatedServerFacade struct {
	dedicatedServerService ports.DedicatedServerService
}

// GetAllDedicatedServers retrieve all dedicated servers.
func (f DedicatedServerFacade) GetAllDedicatedServers(ctx context.Context) (
	*model.DedicatedServers,
	*shared.FacadeError,
) {
	dedicatedServers, err := f.dedicatedServerService.GetAllDedicatedServers(ctx)
	if err != nil {
		return nil, shared.NewFromServicesError("GetAllDedicatedServers", err)
	}

	dataSourceDedicatedServers := to_data_source_model.AdaptDedicatedServers(*dedicatedServers)

	return &dataSourceDedicatedServers, nil
}

// GetAllOperatingSystems retrieve all operating systems.
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

func New(dedicatedServerService ports.DedicatedServerService) DedicatedServerFacade {
	return DedicatedServerFacade{
		dedicatedServerService: dedicatedServerService,
	}
}
