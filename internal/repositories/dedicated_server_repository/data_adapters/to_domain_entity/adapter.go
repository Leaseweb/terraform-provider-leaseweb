package to_domain_entity

import (
	"github.com/leaseweb/leaseweb-go-sdk/dedicatedServer"
	domain "github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain/dedicated_server"
)

// AdaptDedicatedServer adapts an dedicatedServer domain entity to an sdk dedicatedServer model.
func AdaptDedicatedServer(
	sdkDedicatedServer dedicatedServer.Server,
) (
	*domain.DedicatedServer,
	error,
) {

	dedicatedServer := domain.NewDedicatedServer(
		sdkDedicatedServer.GetId(),
	)

	return &dedicatedServer, nil
}
