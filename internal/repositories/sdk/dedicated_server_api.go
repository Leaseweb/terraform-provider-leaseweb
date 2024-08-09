package sdk

import (
	"context"

	"github.com/leaseweb/leaseweb-go-sdk/dedicatedServer"
)

// DedicatedServerApi contains all methods that the sdk must support.
type DedicatedServerApi interface {
	// TODO: ApiGetServerListRequest or ApiGetDedicatedServerListRequest
	// TODO: GetServerList or GetDedicatedServerList
	GetServerList(ctx context.Context) dedicatedServer.ApiGetServerListRequest
}
