package dedicatedserver

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/dedicatedServer"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ datasource.DataSource              = &credentialDataSource{}
	_ datasource.DataSourceWithConfigure = &credentialDataSource{}
)

type credentialDataSource struct {
	name   string
	client dedicatedServer.DedicatedServerAPI
}

type credentialDataSourceModel struct {
	DedicatedServerID types.String `tfsdk:"dedicated_server_id"`
	Username          types.String `tfsdk:"username"`
	Password          types.String `tfsdk:"password"`
	Type              types.String `tfsdk:"type"`
}

func (c *credentialDataSource) Configure(
	_ context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	coreClient, ok := utils.GetDataSourceClient(req, resp)
	if !ok {
		return
	}

	c.client = coreClient.DedicatedServerAPI
}

func (c *credentialDataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, c.name)
}

func (c *credentialDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
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

func (c *credentialDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var data credentialDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serverID := data.DedicatedServerID.ValueString()
	credType := dedicatedServer.CredentialType(data.Type.ValueString())
	username := data.Username.ValueString()

	credential, response, err := c.client.GetServerCredential(
		ctx,
		serverID,
		credType,
		username,
	).Execute()

	if err != nil {
		summary := fmt.Sprintf(
			"Reading data %s for dedicated_server_id %q",
			c.name,
			serverID,
		)
		utils.Error(ctx, &resp.Diagnostics, summary, err, response)
		return
	}

	data.Password = types.StringValue(credential.GetPassword())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func NewCredentialDataSource() datasource.DataSource {
	return &credentialDataSource{
		name: "dedicated_server_credential",
	}
}
