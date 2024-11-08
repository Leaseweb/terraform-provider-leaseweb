package publiccloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ datasource.DataSourceWithConfigure = &credentialDataSource{}
)

type credentialDataSource struct {
	name   string
	client publicCloud.PublicCloudAPI
}

func NewCredentialDataSource() datasource.DataSource {
	return &credentialDataSource{
		name: "public_cloud_credential",
	}
}

type credentialDataSourceModel struct {
	InstanceID types.String `tfsdk:"instance_id"`
	Username   types.String `tfsdk:"username"`
	Password   types.String `tfsdk:"password"`
	Type       types.String `tfsdk:"type"`
}

func (d *credentialDataSource) Configure(
	_ context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	coreClient, ok := utils.GetDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = coreClient.PublicCloudAPI
}

func (d *credentialDataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, d.name)
}

func (d *credentialDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"instance_id": schema.StringAttribute{
				Description: "The ID of the instance.",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "The type of the credential. Valid options are " + utils.StringTypeArrayToMarkdown(publicCloud.AllowedCredentialTypeEnumValues),
				Validators: []validator.String{
					stringvalidator.OneOf(utils.AdaptStringTypeArrayToStringArray(publicCloud.AllowedCredentialTypeEnumValues)...),
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

func (d *credentialDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {

	var state credentialDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceID := state.InstanceID.ValueString()
	type_ := state.Type.ValueString()
	username := state.Username.ValueString()

	credential, response, err := d.client.GetCredential(
		ctx,
		instanceID,
		type_,
		username,
	).Execute()

	if err != nil {
		summary := fmt.Sprintf(
			"Reading data %s for instance_id %q",
			d.name,
			instanceID,
		)
		// TODO: Need change after a proper error logging implementation.
		resp.Diagnostics.AddError(summary, utils.NewError(response, err).Error())
		tflog.Error(ctx, fmt.Sprintf("%s %s", summary, utils.NewError(response, err).Error()))
		return
	}

	state.Password = types.StringValue(credential.GetPassword())
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
