package utils

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
)

// GetResourceClient returns client.Client with true if it can be found. Otherwise, return nil, false.
func GetResourceClient(req resource.ConfigureRequest, resp *resource.ConfigureResponse) (*client.Client, bool) {
	if req.ProviderData == nil {
		return nil, false
	}

	coreClient, ok := req.ProviderData.(client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf(
				"Expected client.Client, got: %T. Please report this issue to the provider developers.",
				req.ProviderData,
			),
		)

		return nil, false
	}

	return &coreClient, true
}

// GetDataSourceClient returns client.Client with true if it can be found. Otherwise, return nil, false.
func GetDataSourceClient(req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) (*client.Client, bool) {
	if req.ProviderData == nil {
		return nil, false
	}

	coreClient, ok := req.ProviderData.(client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf(
				"Expected client.Client, got: %T. Please report this issue to the provider developers.",
				req.ProviderData,
			),
		)

		return nil, false
	}

	return &coreClient, true
}
