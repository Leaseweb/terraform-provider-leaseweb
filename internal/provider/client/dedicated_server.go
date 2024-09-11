package client

import (
	"context"

	sdkDedicatedServer "github.com/leaseweb/leaseweb-go-sdk/dedicatedServer"
)

type dedicatedServer struct {
	optional Optional
	apiKey   string
}

func (c *dedicatedServer) API() sdkDedicatedServer.DedicatedServerAPI {
	configuration := sdkDedicatedServer.NewConfiguration()

	if c.optional.Host != nil {
		configuration.Host = *c.optional.Host
	}
	if c.optional.Scheme != nil {
		configuration.Scheme = *c.optional.Scheme
	}

	return sdkDedicatedServer.NewAPIClient(configuration).DedicatedServerAPI
}

func (c *dedicatedServer) AuthContext(ctx context.Context) context.Context {
	return context.WithValue(
		ctx,
		sdkDedicatedServer.ContextAPIKeys,
		map[string]sdkDedicatedServer.APIKey{
			"X-LSW-Auth": {Key: c.apiKey, Prefix: ""},
		},
	)
}

func GetDedicatedServerClient(optional Optional, apiKey string) dedicatedServer {
	return dedicatedServer{optional: optional, apiKey: apiKey}
}
