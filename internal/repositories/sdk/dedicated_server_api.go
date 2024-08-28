package sdk

import (
	"context"

	"github.com/leaseweb/leaseweb-go-sdk/dedicatedServer"
)

// DedicatedServerApi contains all methods that the sdk must support.
type DedicatedServerApi interface {
	GetServerList(ctx context.Context) dedicatedServer.ApiGetServerListRequest

	GetControlPanelList(ctx context.Context) dedicatedServer.ApiGetControlPanelListRequest
}
