package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/leaseweb/leaseweb-go-sdk/dedicatedServer"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
)

var (
	_ resource.Resource              = &credentialResource{}
	_ resource.ResourceWithConfigure = &credentialResource{}
)

type credentialResource struct {
	// TODO: Refactor this part, apiKey shouldn't be here.
	apiKey string
	client dedicatedServer.DedicatedServerAPI
}

type credentialResourceData struct {
	DedicatedServerId types.String `tfsdk:"dedicated_server_id"`
	Username          types.String `tfsdk:"username"`
	Type              types.String `tfsdk:"type"`
	Password          types.String `tfsdk:"password"`
}

func NewCredentialResource() resource.Resource {
	return &credentialResource{}
}

func (d *credentialResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dedicated_server_credential"
}

func (d *credentialResource) authContext(ctx context.Context) context.Context {
	return context.WithValue(
		ctx,
		dedicatedServer.ContextAPIKeys,
		map[string]dedicatedServer.APIKey{
			"X-LSW-Auth": {Key: d.apiKey, Prefix: ""},
		},
	)
}

func (d *credentialResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (d *credentialResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"dedicated_server_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the dedicated server.",
			},
			"username": schema.StringAttribute{
				Required:    true,
				Description: `The username for the credentials`,
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "The type of the credential.",
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"OPERATING_SYSTEM", "CONTROL_PANEL", "REMOTE_MANAGEMENT", "RESCUE_MODE", "SWITCH", "PDU", "FIREWALL", "LOAD_BALANCER"}...),
				},
			},
			"password": schema.StringAttribute{
				Required:    true,
				Description: `The password for the credentials`,
			},
		},
	}
}

func (d *credentialResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data credentialResourceData
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	opts := dedicatedServer.NewCreateServerCredentialOpts(
		data.Password.ValueString(),
		data.Type.ValueString(),
		data.Username.ValueString(),
	)
	request := d.client.CreateServerCredential(d.authContext(ctx), data.DedicatedServerId.ValueString()).CreateServerCredentialOpts(*opts)
	result, _, err := request.Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf(
				"Error creating credential for dedicated_server_id: %q",
				data.DedicatedServerId.ValueString(),
			),
			err.Error(),
		)
		return
	}

	data = credentialResourceData{
		DedicatedServerId: data.DedicatedServerId,
		Type:              types.StringValue(result.GetType()),
		Password:          types.StringValue(result.GetPassword()),
		Username:          types.StringValue(result.GetUsername()),
	}
	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *credentialResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data credentialResourceData
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := d.client.GetServerCredential(d.authContext(ctx), data.DedicatedServerId.ValueString(), data.Type.ValueString(), data.Username.ValueString())
	result, _, err := request.Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf(
				"Error reading credential with username: %q and dedicated_server_id: %q",
				data.Username.ValueString(),
				data.DedicatedServerId.ValueString(),
			),
			err.Error(),
		)
		return
	}

	data = credentialResourceData{
		DedicatedServerId: data.DedicatedServerId,
		Type:              types.StringValue(result.GetType()),
		Password:          types.StringValue(result.GetPassword()),
		Username:          types.StringValue(result.GetUsername()),
	}
	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *credentialResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data credentialResourceData
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// NOTE: If type/username modified, The API will create a new credential
	opts := dedicatedServer.NewUpdateServerCredentialOpts(
		data.Password.ValueString(),
	)
	request := d.client.UpdateServerCredential(d.authContext(ctx), data.DedicatedServerId.ValueString(), data.Type.ValueString(), data.Username.ValueString()).UpdateServerCredentialOpts(*opts)
	result, _, err := request.Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf(
				"Error updating credential with username: %q and dedicated_server_id: %q",
				data.Username.ValueString(),
				data.DedicatedServerId.ValueString(),
			),
			err.Error(),
		)
		return
	}

	data = credentialResourceData{
		DedicatedServerId: data.DedicatedServerId,
		Type:              types.StringValue(result.GetType()),
		Password:          types.StringValue(result.GetPassword()),
		Username:          types.StringValue(result.GetUsername()),
	}
	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *credentialResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data credentialResourceData
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := d.client.DeleteServerCredential(d.authContext(ctx), data.DedicatedServerId.ValueString(), data.Type.ValueString(), data.Username.ValueString())
	_, err := request.Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf(
				"Error deleting credential with username: %q and dedicated_server_id: %q",
				data.Username.ValueString(),
				data.DedicatedServerId.ValueString(),
			),
			err.Error(),
		)
		return
	}
}
