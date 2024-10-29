// Package client implements access to facades.
package client

import (
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
)

const userAgent = "leaseweb-terraform"

// ProviderData TODO: Refactor this part, data can be managed directly, not within client.
type ProviderData struct {
	ApiKey string
	Host   *string
	Scheme *string
}

// The Client handles instantiation of the SDK.
type Client struct {
	ProviderData   ProviderData
	PublicCloudAPI publicCloud.PublicCloudAPI
}

type Optional struct {
	Host   *string
	Scheme *string
}

func NewClient(token string, optional Optional) Client {
	cfg := publicCloud.NewConfiguration()
	if optional.Host != nil {
		cfg.Host = *optional.Host
	}
	if optional.Scheme != nil {
		cfg.Scheme = *optional.Scheme
	}
	cfg.AddDefaultHeader("X-LSW-Auth", token)
	cfg.UserAgent = userAgent
	publicCloudApi := publicCloud.NewAPIClient(cfg)

	return Client{
		ProviderData: ProviderData{
			ApiKey: token,
			Host:   optional.Host,
			Scheme: optional.Scheme,
		},
		PublicCloudAPI: publicCloudApi.PublicCloudAPI,
	}
}
