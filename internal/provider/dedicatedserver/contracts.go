package dedicatedserver

import (
	"context"

	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/dedicatedserver/services"
)

type DedicatedServer interface {
	Update(
		plan DedicatedServerModel,
		state *DedicatedServerModel,
		ctx context.Context,
	) *services.DedicatedServerError
}
