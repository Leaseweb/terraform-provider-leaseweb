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
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/enum_utils"
)

var (
	_ datasource.DataSource              = &dataSource{}
	_ datasource.DataSourceWithConfigure = &dataSource{}
)

type dataSource struct {
	BaseDataSourceConfig
}

type model struct {
	DedicatedServerID types.String `tfsdk:"dedicated_server_id"`
	Username          types.String `tfsdk:"username"`
	Password          types.String `tfsdk:"password"`
	Type              types.String `tfsdk:"type"`
}

func (d *dataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dedicated_server_credential"
}

func (d *dataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"dedicated_server_id": schema.StringAttribute{
				Description: "The ID of a server",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "The type of the credential. Valid options are " + enum_utils.ConvertStringSliceToMarkdown(availableTypes()),
				Validators: []validator.String{
					stringvalidator.OneOf(availableTypes()...),
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

func (d *dataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data model

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serverID := data.DedicatedServerID.ValueString()
	credType := dedicatedServer.CredentialType(data.Type.ValueString())
	username := data.Username.ValueString()
	client := d.Client.DedicatedServer

	credential, _, err := client.API().GetServerCredential(
		client.AuthContext(ctx), serverID, credType, username,
	).Execute()

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
	return &dataSource{}
}

func availableTypes() []string {
	return enum_utils.ConvertStringSliceToValues(dedicatedServer.AllowedCredentialTypeEnumValues)
}
