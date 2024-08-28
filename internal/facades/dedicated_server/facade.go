// Package dedicated_server implements the dedicated_server facade.
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

// GetAllDedicatedServers retrieve all dedicated servers.
func (h DedicatedServerFacade) GetAllDedicatedServers(ctx context.Context) (
	*dataSourceModel.DedicatedServers,
	*shared.FacadeError,
) {
	dedicatedServers, err := h.dedicatedServerService.GetAllDedicatedServers(ctx)
	if err != nil {
		return nil, shared.NewFromServicesError("GetAllDedicatedServers", err)
	}

	dataSourceDedicatedServers := to_data_source_model.AdaptDedicatedServers(*dedicatedServers)

	return &dataSourceDedicatedServers, nil
}

func NewDedicatedServerFacade(dedicatedServerService ports.DedicatedServerService) DedicatedServerFacade {
	return DedicatedServerFacade{
		dedicatedServerService: dedicatedServerService,
	}
}
