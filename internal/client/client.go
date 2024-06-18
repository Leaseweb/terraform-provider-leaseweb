package client

import (
	"context"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
)

type Client struct {
	PublicCloudClient *publicCloud.APIClient
	Token             string
}

type Options struct {
	Host   string
	Scheme string
}

func NewClient(token string, options *Options) *Client {
	configuration := publicCloud.NewConfiguration()

	if options.Host != "" {
		configuration.Host = options.Host
	}
	if options.Scheme != "" {
		configuration.Scheme = options.Scheme
	}

	return &Client{PublicCloudClient: publicCloud.NewAPIClient(configuration), Token: token}
}

func (c *Client) AuthContext(ctx context.Context) context.Context {
	return context.WithValue(
		ctx,
		publicCloud.ContextAPIKeys,
		map[string]publicCloud.APIKey{
			"X-LSW-Auth": {Key: c.Token, Prefix: ""},
		},
	)
}
