package publiccloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	publicCloudSDK "github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ resource.Resource              = &credentialResource{}
	_ resource.ResourceWithConfigure = &credentialResource{}
)

type credentialResource struct {
	// TODO: Refactor this part, apiKey shouldn't be here.
	apiKey string
	client publicCloudSDK.PublicCloudAPI
}

type credentialResourceData struct {
	InstanceID types.String `tfsdk:"instance_id"`
	Username   types.String `tfsdk:"username"`
	Type       types.String `tfsdk:"type"`
	Password   types.String `tfsdk:"password"`
}

func NewCredentialResource() resource.Resource {
	return &credentialResource{}
}

func (d *credentialResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_public_cloud_credential"
}

func (d *credentialResource) authContext(ctx context.Context) context.Context {
	return context.WithValue(
		ctx,
		publicCloudSDK.ContextAPIKeys,
		map[string]publicCloudSDK.APIKey{
			"X-LSW-Auth": {Key: d.apiKey, Prefix: ""},
		},
	)
}

func (d *credentialResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	configuration := publicCloudSDK.NewConfiguration()

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

	apiClient := publicCloudSDK.NewAPIClient(configuration)
	d.client = apiClient.PublicCloudAPI
}

func (d *credentialResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the instance.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"username": schema.StringAttribute{
				Required:    true,
				Description: `The username for the credentials`,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: `The type of the credential. Valid options are: "OPERATING_SYSTEM", "CONTROL_PANEL"`,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"OPERATING_SYSTEM", "CONTROL_PANEL"}...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
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

	opts := publicCloudSDK.NewStoreCredentialOpts(
		publicCloudSDK.CredentialType(data.Type.ValueString()),
		data.Username.ValueString(),
		data.Password.ValueString(),
	)
	request := d.client.StoreCredential(d.authContext(ctx), data.InstanceID.ValueString()).StoreCredentialOpts(*opts)
	result, response, err := request.Execute()
	if err != nil {
		summary := fmt.Sprintf("Error creating credential with username: %q and instance_id: %q", data.Username.ValueString(), data.InstanceID.ValueString())
		resp.Diagnostics.AddError(summary, utils.NewError(response, err).Error())
		tflog.Error(ctx, fmt.Sprintf("%s %s", summary, utils.NewError(response, err).Error()))
		return
	}

	data = credentialResourceData{
		InstanceID: data.InstanceID,
		Type:       types.StringValue(string(result.GetType())),
		Password:   types.StringValue(result.GetPassword()),
		Username:   types.StringValue(result.GetUsername()),
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

	request := d.client.GetCredential(d.authContext(ctx), data.InstanceID.ValueString(), data.Type.ValueString(), data.Username.ValueString())
	result, response, err := request.Execute()
	if err != nil {
		summary := fmt.Sprintf("Error reading credential with username: %q and instance_id: %q", data.Username.ValueString(), data.InstanceID.ValueString())
		resp.Diagnostics.AddError(summary, utils.NewError(response, err).Error())
		tflog.Error(ctx, fmt.Sprintf("%s %s", summary, utils.NewError(response, err).Error()))
		return
	}

	data = credentialResourceData{
		InstanceID: data.InstanceID,
		Type:       types.StringValue(string(result.GetType())),
		Password:   types.StringValue(result.GetPassword()),
		Username:   types.StringValue(result.GetUsername()),
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

	opts := publicCloudSDK.NewUpdateCredentialOpts(
		data.Password.ValueString(),
	)
	request := d.client.UpdateCredential(d.authContext(ctx), data.InstanceID.ValueString(), data.Type.ValueString(), data.Username.ValueString()).UpdateCredentialOpts(*opts)
	result, response, err := request.Execute()
	if err != nil {
		summary := fmt.Sprintf("Error updating credential with username: %q and instance_id: %q", data.Username.ValueString(), data.InstanceID.ValueString())
		resp.Diagnostics.AddError(summary, utils.NewError(response, err).Error())
		tflog.Error(ctx, fmt.Sprintf("%s %s", summary, utils.NewError(response, err).Error()))
		return
	}

	data = credentialResourceData{
		InstanceID: data.InstanceID,
		Type:       types.StringValue(string(result.GetType())),
		Password:   types.StringValue(result.GetPassword()),
		Username:   types.StringValue(result.GetUsername()),
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

	request := d.client.DeleteCredential(d.authContext(ctx), data.InstanceID.ValueString(), data.Type.ValueString(), data.Username.ValueString())
	response, err := request.Execute()
	if err != nil {
		summary := fmt.Sprintf("Error deleting credential with username: %q and instance_id: %q", data.Username.ValueString(), data.InstanceID.ValueString())
		resp.Diagnostics.AddError(summary, utils.NewError(response, err).Error())
		tflog.Error(ctx, fmt.Sprintf("%s %s", summary, utils.NewError(response, err).Error()))
		return
	}
}
