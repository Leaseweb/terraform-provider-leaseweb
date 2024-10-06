package provider

import (
	"context"

	dedicatedServerSdk "github.com/leaseweb/leaseweb-go-sdk/dedicatedServer"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/dedicatedserver"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/dedicatedserver/services"
)

// ProviderData TODO: Refactor this part, ProviderData can be managed directly, not within client.
type ProviderData struct {
	ApiKey string
	Host   *string
	Scheme *string
}

// The Client handles instantiation of the services.
type Client struct {
	ProviderData    ProviderData
	DedicatedServer dedicatedserver.DedicatedServer
}

// AuthContext should be refactored
func (c Client) AuthContext(ctx context.Context) context.Context {
	return context.WithValue(
		ctx,
		publicCloud.ContextAPIKeys,
		map[string]publicCloud.APIKey{
			"X-LSW-Auth": {Key: c.ProviderData.ApiKey, Prefix: ""},
		},
	)
}

type Optional struct {
	Host   *string
	Scheme *string
}

func NewClient(token string, optional Optional) Client {
	configuration := dedicatedServerSdk.NewConfiguration()

	if optional.Host != nil {
		configuration.Host = *optional.Host
	}
	if optional.Scheme != nil {
		configuration.Scheme = *optional.Scheme
	}

	client := *dedicatedServerSdk.NewAPIClient(configuration)

	return Client{
		ProviderData: ProviderData{
			ApiKey: token,
			Host:   optional.Host,
			Scheme: optional.Scheme,
		},
		DedicatedServer: services.NewDedicatedServer(
			client.DedicatedServerAPI.UpdateServerReference,
			client.DedicatedServerAPI.PowerServerOn,
			client.DedicatedServerAPI.PowerServerOff,
		),
	}
}
