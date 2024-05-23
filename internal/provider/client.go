package provider

import (
	"context"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
)

type leasewebProviderClient struct {
	Host  string
	Token string
}

func (c *leasewebProviderClient) Client() *publicCloud.APIClient {
	configuration := publicCloud.NewConfiguration()
	configuration.Host = c.Host

	return publicCloud.NewAPIClient(configuration)
}

func (c *leasewebProviderClient) AuthContext() context.Context {
	return context.WithValue(
		context.Background(),
		publicCloud.ContextAPIKeys,
		map[string]publicCloud.APIKey{
			"X-LSW-Auth": {Key: c.Token, Prefix: ""},
		},
	)
}
