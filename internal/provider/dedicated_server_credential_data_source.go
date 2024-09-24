package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/dedicatedServer"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
)

var (
	_ datasource.DataSource              = &dedicatedServerCredentialDataSource{}
	_ datasource.DataSourceWithConfigure = &dedicatedServerCredentialDataSource{}
)

type dedicatedServerCredentialDataSource struct {
	// TODO: Refactor this part, apiKey shouldn't be here.
	apiKey string
	client dedicatedServer.DedicatedServerAPI
}

type dedicatedServerCredentialDataSourceModel struct {
	DedicatedServerID types.String `tfsdk:"dedicated_server_id"`
	Username          types.String `tfsdk:"username"`
	Password          types.String `tfsdk:"password"`
	Type              types.String `tfsdk:"type"`
}

func (d *dedicatedServerCredentialDataSource) authContext(ctx context.Context) context.Context {
	return context.WithValue(
		ctx,
		dedicatedServer.ContextAPIKeys,
		map[string]dedicatedServer.APIKey{
			"X-LSW-Auth": {Key: d.apiKey, Prefix: ""},
		},
	)
}

func (d *dedicatedServerCredentialDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	configuration := dedicatedServer.NewConfiguration()

	// TODO: Refactor this part, ProviderData can be managed directly, not within client.
	coreClient, ok := req.ProviderData.(client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf(
				"Expected client.Client, got: %T. Please report this issue to the provider developers.",
				req.ProviderData,
			),
		)
		return
	}

	d.apiKey = coreClient.ProviderData.ApiKey
	if coreClient.ProviderData.Host != nil {
		configuration.Host = *coreClient.ProviderData.Host
	}
	if coreClient.ProviderData.Scheme != nil {
		configuration.Scheme = *coreClient.ProviderData.Scheme
	}

	apiClient := dedicatedServer.NewAPIClient(configuration)
	d.client = apiClient.DedicatedServerAPI
}

func (d *dedicatedServerCredentialDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dedicated_server_credential"
}

func (d *dedicatedServerCredentialDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"dedicated_server_id": schema.StringAttribute{
				Description: "The ID of a server",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "The type of the credential. Valid options are \n  - *OPERATING_SYSTEM*\n  - *RESCUE_MODE*\n  - *REMOTE_MANAGEMENT*\n  - *CONTROL_PANEL*\n  - *SWITCH*\n  - *PDU*\n  - *FIREWALL*\n  - *LOAD_BALANCER*\n  - *VNC*\n  - *TEMPORARY_OPERATING_SYSTEM*\n  - *VPN_USER*\n  - *COMBINATION_LOCK*\n  - *DATABASE*\n",
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"OPERATING_SYSTEM", "RESCUE_MODE", "REMOTE_MANAGEMENT", "CONTROL_PANEL", "SWITCH", "PDU", "FIREWALL", "LOAD_BALANCER", "VNC", "TEMPORARY_OPERATING_SYSTEM", "VPN_USER", "COMBINATION_LOCK", "DATABASE"}...),
				},
			},
			"username": schema.StringAttribute{
				Description: "The username for the credentials",
				Required:    true,
			},
			"password": schema.StringAttribute{
				Description: "The password for the credentials",
				Computed:    true,
				Sensitive:   true,
			},
		},
	}
}

func (d *dedicatedServerCredentialDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data dedicatedServerCredentialDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serverID := data.DedicatedServerID.ValueString()
	credType := dedicatedServer.CredentialType(data.Type.ValueString())
	username := data.Username.ValueString()

	credential, _, err := d.client.GetServerCredential(d.authContext(ctx), serverID, credType, username).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading data dedicated_server_credential for server %q", serverID),
			err.Error(),
		)
		return
	}

	data.Password = types.StringValue(credential.GetPassword())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func NewDedicatedServerCredentialDataSource() datasource.DataSource {
	return &dedicatedServerCredentialDataSource{}
}
