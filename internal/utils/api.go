package utils

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/leaseweb/leaseweb-go-sdk/v2/dedicatedserver"
	"github.com/leaseweb/leaseweb-go-sdk/v2/publiccloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
)

func generateTypeName(providerTypeName string, name string) string {
	return fmt.Sprintf("%s_%s", providerTypeName, name)
}

// ResourceAPI contains a reusable Metadata function for resources.
type ResourceAPI struct {
	Name string
}

func (r ResourceAPI) Metadata(
	_ context.Context,
	request resource.MetadataRequest,
	response *resource.MetadataResponse,
) {
	response.TypeName = generateTypeName(request.ProviderTypeName, r.Name)
}

// DataSourceAPI contains a reusable Metadata function for data sources.
type DataSourceAPI struct {
	Name string
}

func (d DataSourceAPI) Metadata(
	_ context.Context,
	request datasource.MetadataRequest,
	response *datasource.MetadataResponse,
) {
	response.TypeName = generateTypeName(request.ProviderTypeName, d.Name)
}

func getCoreClient(
	providerData any,
	diagnostics *diag.Diagnostics,
) *client.Client {
	if providerData == nil {
		return nil
	}

	coreClient, ok := providerData.(client.Client)

	if !ok {
		diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf(
				"Expected client.Client, got: %T. Please report this issue to the provider developers.",
				providerData,
			),
		)
		return nil
	}

	return &coreClient
}

// PubliccloudResourceAPI contains ResourceAPI & the publiccloud.PubliccloudAPI client and a reusable Configure function.
type PubliccloudResourceAPI struct {
	ResourceAPI
	Client publiccloud.PubliccloudAPI
}

func (p *PubliccloudResourceAPI) Configure(
	_ context.Context,
	request resource.ConfigureRequest,
	response *resource.ConfigureResponse,
) {
	coreClient := getCoreClient(request.ProviderData, &response.Diagnostics)
	if coreClient == nil {
		return
	}

	p.Client = coreClient.PubliccloudAPI
}

// NewPubliccloudResourceAPI returns a new PubliccloudResourceAPI.
func NewPubliccloudResourceAPI(name string) PubliccloudResourceAPI {
	return PubliccloudResourceAPI{
		ResourceAPI: ResourceAPI{
			Name: name,
		},
	}
}

// PubliccloudDataSourceAPI contains DataSourceAPI & the publiccloud.PubliccloudAPI client and a reusable Configure function.
type PubliccloudDataSourceAPI struct {
	DataSourceAPI
	Client publiccloud.PubliccloudAPI
}

func (p *PubliccloudDataSourceAPI) Configure(
	_ context.Context,
	request datasource.ConfigureRequest,
	response *datasource.ConfigureResponse,
) {
	coreClient := getCoreClient(request.ProviderData, &response.Diagnostics)
	if coreClient == nil {
		return
	}

	p.Client = coreClient.PubliccloudAPI
}

// NewPubliccloudDataSourceAPI returns a new PubliccloudDataSourceAPI.
func NewPubliccloudDataSourceAPI(name string) PubliccloudDataSourceAPI {
	return PubliccloudDataSourceAPI{
		DataSourceAPI: DataSourceAPI{
			Name: name,
		},
	}
}

// DedicatedserverResourceAPI contains ResourceAPI & the dedicatedserver.DedicatedserverAPI client and a reusable Configure function.
type DedicatedserverResourceAPI struct {
	ResourceAPI
	Client dedicatedserver.DedicatedserverAPI
}

func (d *DedicatedserverResourceAPI) Configure(
	_ context.Context,
	request resource.ConfigureRequest,
	response *resource.ConfigureResponse,
) {
	coreClient := getCoreClient(request.ProviderData, &response.Diagnostics)
	if coreClient == nil {
		return
	}

	d.Client = coreClient.DedicatedserverAPI
}

// NewDedicatedserverResourceAPI returns a new DedicatedserverResourceAPI.
func NewDedicatedserverResourceAPI(name string) DedicatedserverResourceAPI {
	return DedicatedserverResourceAPI{
		ResourceAPI: ResourceAPI{
			Name: name,
		},
	}
}

// DedicatedserverDataSourceAPI contains DataSourceAPI & the dedicatedserver.DedicatedserverAPI client and a reusable Configure function.
type DedicatedserverDataSourceAPI struct {
	DataSourceAPI
	Client dedicatedserver.DedicatedserverAPI
}

func (d *DedicatedserverDataSourceAPI) Configure(
	_ context.Context,
	request datasource.ConfigureRequest,
	response *datasource.ConfigureResponse,
) {
	coreClient := getCoreClient(request.ProviderData, &response.Diagnostics)
	if coreClient == nil {
		return
	}

	d.Client = coreClient.DedicatedserverAPI
}

// NewDedicatedserverDataSourceAPI returns a new DedicatedserverDataSourceAPI.
func NewDedicatedserverDataSourceAPI(name string) DedicatedserverDataSourceAPI {
	return DedicatedserverDataSourceAPI{
		DataSourceAPI: DataSourceAPI{
			Name: name,
		},
	}
}
