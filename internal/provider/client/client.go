// Package client implements access to facades.
package client

import (
	"github.com/leaseweb/leaseweb-go-sdk/v2/dedicatedserver"
	"github.com/leaseweb/leaseweb-go-sdk/v2/publiccloud"
)

const userAgentBase = "leaseweb-terraform"

// The Client handles instantiation of the SDK.
type Client struct {
	PubliccloudAPI     publiccloud.PubliccloudAPI
	DedicatedserverAPI dedicatedserver.DedicatedserverAPI
}

type Optional struct {
	Host   *string
	Scheme *string
}

func NewClient(token string, optional Optional, version string) Client {
	publiccloudCFG := publiccloud.NewConfiguration()
	dedicatedserverCFG := dedicatedserver.NewConfiguration()

	if optional.Host != nil {
		publiccloudCFG.Host = *optional.Host
		dedicatedserverCFG.Host = *optional.Host
	}
	if optional.Scheme != nil {
		publiccloudCFG.Scheme = *optional.Scheme
		dedicatedserverCFG.Scheme = *optional.Scheme
	}

	userAgent := userAgentBase + "-" + version

	publiccloudCFG.AddDefaultHeader("X-LSW-Auth", token)
	publiccloudCFG.UserAgent = userAgent

	dedicatedserverCFG.AddDefaultHeader("X-LSW-Auth", token)
	dedicatedserverCFG.UserAgent = userAgent

	publiccloudAPI := publiccloud.NewAPIClient(publiccloudCFG)
	dedicatedserverAPI := dedicatedserver.NewAPIClient(dedicatedserverCFG)

	return Client{
		PubliccloudAPI:     publiccloudAPI.PubliccloudAPI,
		DedicatedserverAPI: dedicatedserverAPI.DedicatedserverAPI,
	}
}
