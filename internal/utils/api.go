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

// PubliccloudResourceAPI contains reusable Configure & Metadata functions for resources that implement publiccloud.PubliccloudAPI.
type PubliccloudResourceAPI struct {
	Name   string
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

func (p *PubliccloudResourceAPI) Metadata(
	_ context.Context,
	request resource.MetadataRequest,
	response *resource.MetadataResponse,
) {
	response.TypeName = generateTypeName(request.ProviderTypeName, p.Name)
}

// PubliccloudDataSourceAPI contains reusable Configure & Metadata functions for data sources that implement publiccloud.PubliccloudAPI.
type PubliccloudDataSourceAPI struct {
	Name   string
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

func (p *PubliccloudDataSourceAPI) Metadata(
	_ context.Context,
	request datasource.MetadataRequest,
	response *datasource.MetadataResponse,
) {
	response.TypeName = generateTypeName(request.ProviderTypeName, p.Name)
}

// DedicatedserverResourceAPI contains reusable Configure & Metadata functions for resources that implement dedicatedserver.DedicatedserverAPI.
type DedicatedserverResourceAPI struct {
	Name   string
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

func (d *DedicatedserverResourceAPI) Metadata(
	_ context.Context,
	request resource.MetadataRequest,
	response *resource.MetadataResponse,
) {
	response.TypeName = generateTypeName(request.ProviderTypeName, d.Name)
}

// DedicatedserverDataSourceAPI contains reusable Configure & Metadata functions for data sources that implement dedicatedserver.DedicatedserverAPI.
type DedicatedserverDataSourceAPI struct {
	Name   string
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

func (d *DedicatedserverDataSourceAPI) Metadata(
	_ context.Context,
	request datasource.MetadataRequest,
	response *datasource.MetadataResponse,
) {
	response.TypeName = generateTypeName(request.ProviderTypeName, d.Name)
}
