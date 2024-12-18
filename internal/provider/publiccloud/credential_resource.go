package publiccloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/v3/publiccloud"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ resource.Resource              = &credentialResource{}
	_ resource.ResourceWithConfigure = &credentialResource{}
)

type credentialResource struct {
	name   string
	client publiccloud.PubliccloudAPI
}

type credentialResourceModel struct {
	InstanceID types.String `tfsdk:"instance_id"`
	Username   types.String `tfsdk:"username"`
	Type       types.String `tfsdk:"type"`
	Password   types.String `tfsdk:"password"`
}

func NewCredentialResource() resource.Resource {
	return &credentialResource{
		name: "public_cloud_credential",
	}
}

func (c *credentialResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, c.name)
}

func (c *credentialResource) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}

	coreClient, ok := req.ProviderData.(client.Client)

	if !ok {
		utils.ConfigError(&resp.Diagnostics, req.ProviderData)
		return
	}

	c.client = coreClient.PubliccloudAPI
}

func (c *credentialResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: utils.BetaDescription,
		Attributes: map[string]schema.Attribute{
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: `The ID of the instance.`,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"username": schema.StringAttribute{
				Required:    true,
				Description: `The username for the credentials`,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: `The type of the credential. Valid options are ` + utils.StringTypeArrayToMarkdown(publiccloud.AllowedCredentialTypeEnumValues),
				Validators: []validator.String{
					stringvalidator.OneOf(utils.AdaptStringTypeArrayToStringArray(publiccloud.AllowedCredentialTypeEnumValues)...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"password": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: `The password for the credentials`,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
	}
}

func (c *credentialResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan credentialResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	opts := publiccloud.NewStoreCredentialOpts(
		publiccloud.CredentialType(plan.Type.ValueString()),
		plan.Username.ValueString(),
		plan.Password.ValueString(),
	)
	request := c.client.StoreInstanceCredential(
		ctx,
		plan.InstanceID.ValueString(),
	).StoreCredentialOpts(*opts)
	result, response, err := request.Execute()
	if err != nil {
		utils.SdkError(ctx, &resp.Diagnostics, err, response)
		return
	}

	resp.Diagnostics.Append(
		resp.State.Set(
			ctx,
			credentialResourceModel{
				InstanceID: plan.InstanceID,
				Type:       types.StringValue(string(result.GetType())),
				Password:   types.StringValue(result.GetPassword()),
				Username:   types.StringValue(result.GetUsername()),
			},
		)...,
	)
}

func (c *credentialResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state credentialResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := c.client.GetInstanceCredential(
		ctx,
		state.InstanceID.ValueString(),
		publiccloud.CredentialType(state.Type.ValueString()),
		state.Username.ValueString(),
	)
	result, response, err := request.Execute()
	if err != nil {
		utils.SdkError(ctx, &resp.Diagnostics, err, response)
		return
	}

	resp.Diagnostics.Append(
		resp.State.Set(
			ctx,
			credentialResourceModel{
				InstanceID: state.InstanceID,
				Type:       types.StringValue(string(result.GetType())),
				Password:   types.StringValue(result.GetPassword()),
				Username:   types.StringValue(result.GetUsername()),
			},
		)...,
	)
}

func (c *credentialResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan credentialResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	opts := publiccloud.NewUpdateCredentialOpts(
		plan.Password.ValueString(),
	)
	request := c.client.UpdateInstanceCredential(
		ctx,
		plan.InstanceID.ValueString(),
		publiccloud.CredentialType(plan.Type.ValueString()),
		plan.Username.ValueString(),
	).UpdateCredentialOpts(*opts)
	result, response, err := request.Execute()
	if err != nil {
		utils.SdkError(ctx, &resp.Diagnostics, err, response)
		return
	}

	resp.Diagnostics.Append(
		resp.State.Set(
			ctx,
			credentialResourceModel{
				InstanceID: plan.InstanceID,
				Type:       types.StringValue(string(result.GetType())),
				Password:   types.StringValue(result.GetPassword()),
				Username:   types.StringValue(result.GetUsername()),
			},
		)...,
	)
}

func (c *credentialResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state credentialResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := c.client.DeleteInstanceCredential(
		ctx,
		state.InstanceID.ValueString(),
		publiccloud.CredentialType(state.Type.ValueString()),
		state.Username.ValueString(),
	)
	response, err := request.Execute()
	if err != nil {
		utils.SdkError(ctx, &resp.Diagnostics, err, response)
	}
}
