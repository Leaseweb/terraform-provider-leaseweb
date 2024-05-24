package provider

import (
	"context"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
)

type LeasewebProviderClient struct {
	Client *publicCloud.APIClient
	Token  string
}

type LeasewebProviderClientOptions struct {
	Host string
}

func NewLeasewebProviderClient(token string, options *LeasewebProviderClientOptions) *LeasewebProviderClient {
	configuration := publicCloud.NewConfiguration()

	if options.Host != "" {
		configuration.Host = options.Host
	}

	return &LeasewebProviderClient{Client: publicCloud.NewAPIClient(configuration), Token: token}
}

func (c *LeasewebProviderClient) AuthContext() context.Context {
	return context.WithValue(
		context.Background(),
		publicCloud.ContextAPIKeys,
		map[string]publicCloud.APIKey{
			"X-LSW-Auth": {Key: c.Token, Prefix: ""},
		},
	)
}
