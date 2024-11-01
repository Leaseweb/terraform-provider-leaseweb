// Package client implements access to facades.
package client

import (
	"github.com/leaseweb/leaseweb-go-sdk/dedicatedServer"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
)

const userAgent = "leaseweb-terraform"

// The Client handles instantiation of the SDK.
type Client struct {
	PublicCloudAPI     publicCloud.PublicCloudAPI
	DedicatedServerAPI dedicatedServer.DedicatedServerAPI
}

type Optional struct {
	Host   *string
	Scheme *string
}

func NewClient(token string, optional Optional) Client {
	publicCloudCfg := publicCloud.NewConfiguration()
	dedicatedServerCfg := dedicatedServer.NewConfiguration()

	if optional.Host != nil {
		publicCloudCfg.Host = *optional.Host
		dedicatedServerCfg.Host = *optional.Host
	}
	if optional.Scheme != nil {
		publicCloudCfg.Scheme = *optional.Scheme
		dedicatedServerCfg.Scheme = *optional.Scheme
	}

	publicCloudCfg.AddDefaultHeader("X-LSW-Auth", token)
	publicCloudCfg.UserAgent = userAgent

	dedicatedServerCfg.AddDefaultHeader("X-LSW-Auth", token)
	dedicatedServerCfg.UserAgent = userAgent

	publicCloudApi := publicCloud.NewAPIClient(publicCloudCfg)
	dedicatedServerApi := dedicatedServer.NewAPIClient(dedicatedServerCfg)

	return Client{
		PublicCloudAPI:     publicCloudApi.PublicCloudAPI,
		DedicatedServerAPI: dedicatedServerApi.DedicatedServerAPI,
	}
}
