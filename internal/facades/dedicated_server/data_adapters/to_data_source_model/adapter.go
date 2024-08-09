package to_data_source_model

import (
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	domain "github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain/dedicated_server"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/data_sources/dedicated_server/model"
)

func AdaptDedicatedServers(domainDedicatedServers domain.DedicatedServers) model.DedicatedServers {
	var dedicatedServers model.DedicatedServers

	for _, domainDedicatedServer := range domainDedicatedServers {
		dedicated_server := adaptDedicatedServer(domainDedicatedServer)
		dedicatedServers.DedicatedServers = append(dedicatedServers.DedicatedServers, dedicated_server)
	}

	return dedicatedServers
}

func adaptDedicatedServer(domainDedicatedServer domain.DedicatedServer) model.DedicatedServer {
	dedicatedServer := model.DedicatedServer{Id: basetypes.NewStringValue(domainDedicatedServer.Id)}
	return dedicatedServer
}
