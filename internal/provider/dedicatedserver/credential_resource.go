package dedicatedserver

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/v2/dedicatedserver"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ resource.Resource              = &credentialResource{}
	_ resource.ResourceWithConfigure = &credentialResource{}
)

type credentialResource struct {
	name   string
	client dedicatedserver.DedicatedserverAPI
}

type credentialResourceModel struct {
	DedicatedServerId types.String `tfsdk:"dedicated_server_id"`
	Username          types.String `tfsdk:"username"`
	Type              types.String `tfsdk:"type"`
	Password          types.String `tfsdk:"password"`
}

func NewCredentialResource() resource.Resource {
	return &credentialResource{
		name: "dedicated_server_credential",
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
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf(
				"Expected client.Client, got: %T. Please report this issue to the provider developers.",
				req.ProviderData,
			),
		)

		return
	}

	c.client = coreClient.DedicatedserverAPI
}

func (c *credentialResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"dedicated_server_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the dedicated server.",
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
				Description: `The type of the credential. Valid options are: "OPERATING_SYSTEM", "CONTROL_PANEL", "REMOTE_MANAGEMENT", "RESCUE_MODE", "SWITCH", "PDU", "FIREWALL", "LOAD_BALANCER"`,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"OPERATING_SYSTEM", "CONTROL_PANEL", "REMOTE_MANAGEMENT", "RESCUE_MODE", "SWITCH", "PDU", "FIREWALL", "LOAD_BALANCER"}...),
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

	opts := dedicatedserver.NewCreateServerCredentialOpts(
		plan.Password.ValueString(),
		dedicatedserver.CredentialType(plan.Type.ValueString()),
		plan.Username.ValueString(),
	)
	request := c.client.CreateServerCredential(
		ctx,
		plan.DedicatedServerId.ValueString(),
	).CreateServerCredentialOpts(*opts)
	result, response, err := request.Execute()
	if err != nil {
		summary := fmt.Sprintf(
			"Creating resource %s for username %q and dedicated_server_id %q",
			c.name,
			plan.Username.ValueString(),
			plan.DedicatedServerId.ValueString(),
		)
		utils.Error(ctx, &resp.Diagnostics, summary, err, response)
		return
	}

	resp.Diagnostics.Append(
		resp.State.Set(
			ctx,
			credentialResourceModel{
				DedicatedServerId: plan.DedicatedServerId,
				Type:              types.StringValue(string(result.GetType())),
				Password:          types.StringValue(result.GetPassword()),
				Username:          types.StringValue(result.GetUsername()),
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

	request := c.client.GetServerCredential(
		ctx,
		state.DedicatedServerId.ValueString(),
		dedicatedserver.CredentialType(state.Type.ValueString()),
		state.Username.ValueString(),
	)
	result, response, err := request.Execute()
	if err != nil {
		summary := fmt.Sprintf(
			"Reading resource %s for username %q and dedicated_server_id %q",
			c.name,
			state.Username.ValueString(),
			state.DedicatedServerId.ValueString(),
		)
		utils.Error(ctx, &resp.Diagnostics, summary, err, response)
		return
	}

	resp.Diagnostics.Append(
		resp.State.Set(
			ctx,
			credentialResourceModel{
				DedicatedServerId: state.DedicatedServerId,
				Type:              types.StringValue(string(result.GetType())),
				Password:          types.StringValue(result.GetPassword()),
				Username:          types.StringValue(result.GetUsername()),
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

	opts := dedicatedserver.NewUpdateServerCredentialOpts(
		plan.Password.ValueString(),
	)
	request := c.client.UpdateServerCredential(
		ctx,
		plan.DedicatedServerId.ValueString(),
		dedicatedserver.CredentialType(plan.Type.ValueString()),
		plan.Username.ValueString(),
	).UpdateServerCredentialOpts(*opts)
	result, response, err := request.Execute()
	if err != nil {
		summary := fmt.Sprintf(
			"Updating resource %s for username %q and dedicated_server_id %q",
			c.name,
			plan.Username.ValueString(),
			plan.DedicatedServerId.ValueString(),
		)
		utils.Error(ctx, &resp.Diagnostics, summary, err, response)
		return
	}

	resp.Diagnostics.Append(
		resp.State.Set(
			ctx,
			credentialResourceModel{
				DedicatedServerId: plan.DedicatedServerId,
				Type:              types.StringValue(string(result.GetType())),
				Password:          types.StringValue(result.GetPassword()),
				Username:          types.StringValue(result.GetUsername()),
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

	request := c.client.DeleteServerCredential(
		ctx,
		state.DedicatedServerId.ValueString(),
		dedicatedserver.CredentialType(state.Type.ValueString()),
		state.Username.ValueString(),
	)
	response, err := request.Execute()
	if err != nil {
		summary := fmt.Sprintf(
			"Deleting resource %s for username %q and dedicated_server_id %q",
			c.name,
			state.Username.ValueString(),
			state.DedicatedServerId.ValueString(),
		)
		utils.Error(ctx, &resp.Diagnostics, summary, err, response)
	}
}
