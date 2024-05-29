package client

import (
	"context"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
)

type Client struct {
	SdkClient *publicCloud.APIClient
	Token     string
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

	return &Client{SdkClient: publicCloud.NewAPIClient(configuration), Token: token}
}

func (c *Client) AuthContext() context.Context {
	return context.WithValue(
		context.Background(),
		publicCloud.ContextAPIKeys,
		map[string]publicCloud.APIKey{
			"X-LSW-Auth": {Key: c.Token, Prefix: ""},
		},
	)
}
