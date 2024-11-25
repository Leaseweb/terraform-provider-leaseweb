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
	"github.com/leaseweb/leaseweb-go-sdk/dedicatedServer"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ resource.Resource              = &notificationSettingBandwidthResource{}
	_ resource.ResourceWithConfigure = &notificationSettingBandwidthResource{}
)

type notificationSettingBandwidthResource struct {
	name   string
	client dedicatedServer.DedicatedServerAPI
}

type notificationSettingBandwidthResourceModel struct {
	Id                types.String `tfsdk:"id"`
	DedicatedServerId types.String `tfsdk:"dedicated_server_id"`
	Frequency         types.String `tfsdk:"frequency"`
	Threshold         types.String `tfsdk:"threshold"`
	Unit              types.String `tfsdk:"unit"`
}

func NewNotificationSettingBandwidthResource() resource.Resource {
	return &notificationSettingBandwidthResource{
		name: "dedicated_server_notification_setting_bandwidth",
	}
}

func (n *notificationSettingBandwidthResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, n.name)
}

func (n *notificationSettingBandwidthResource) Configure(
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

	n.client = coreClient.DedicatedServerAPI
}

func (n *notificationSettingBandwidthResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The notification setting bandwidth unique identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"dedicated_server_id": schema.StringAttribute{
				Required:    true,
				Description: "The server unique identifier",
			},
			"frequency": schema.StringAttribute{
				Required:    true,
				Description: "The notification frequency. Valid options can be *DAILY* or *WEEKLY* or *MONTHLY*.",
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"DAILY", "WEEKLY", "MONTHLY"}...),
				},
			},
			"threshold": schema.StringAttribute{
				Required:    true,
				Description: "Threshold Value. Value can be a number greater than 0.",
				Validators: []validator.String{
					greaterThanZero(),
				},
			},
			"unit": schema.StringAttribute{
				Required:    true,
				Description: "The notification unit. Valid options can be *Mbps* or *Gbps*.",
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"Mbps", "Gbps"}...),
				},
			},
		},
	}
}

func (n *notificationSettingBandwidthResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan notificationSettingBandwidthResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	opts := dedicatedServer.NewBandwidthNotificationSettingOpts(
		plan.Frequency.ValueString(),
		plan.Threshold.ValueString(),
		plan.Unit.ValueString(),
	)
	request := n.client.CreateServerBandwidthNotificationSetting(
		ctx,
		plan.DedicatedServerId.ValueString(),
	).BandwidthNotificationSettingOpts(*opts)
	result, response, err := request.Execute()
	if err != nil {
		summary := fmt.Sprintf(
			"Creating resource %s for dedicated_server_id %q",
			n.name,
			plan.DedicatedServerId.ValueString(),
		)
		utils.Error(ctx, &resp.Diagnostics, summary, err, response)
		return
	}

	resp.Diagnostics.Append(
		resp.State.Set(
			ctx,
			notificationSettingBandwidthResourceModel{
				Id:                types.StringValue(result.GetId()),
				Frequency:         types.StringValue(result.GetFrequency()),
				Threshold:         types.StringValue(result.GetThreshold()),
				Unit:              types.StringValue(result.GetUnit()),
				DedicatedServerId: plan.DedicatedServerId,
			},
		)...,
	)
}

func (n *notificationSettingBandwidthResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state notificationSettingBandwidthResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := n.client.GetServerBandwidthNotificationSetting(
		ctx,
		state.DedicatedServerId.ValueString(),
		state.Id.ValueString(),
	)
	result, response, err := request.Execute()
	if err != nil {
		summary := fmt.Sprintf(
			"Reading resource %s for id %q and dedicated_server_id %q",
			n.name,
			state.Id.ValueString(),
			state.DedicatedServerId.ValueString(),
		)
		utils.Error(ctx, &resp.Diagnostics, summary, err, response)
		return
	}

	resp.Diagnostics.Append(
		resp.State.Set(
			ctx,
			notificationSettingBandwidthResourceModel{
				Id:                types.StringValue(result.GetId()),
				Frequency:         types.StringValue(result.GetFrequency()),
				Threshold:         types.StringValue(result.GetThreshold()),
				Unit:              types.StringValue(result.GetUnit()),
				DedicatedServerId: state.DedicatedServerId,
			},
		)...,
	)
}

func (n *notificationSettingBandwidthResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan notificationSettingBandwidthResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	opts := dedicatedServer.NewBandwidthNotificationSettingOpts(
		plan.Frequency.ValueString(),
		plan.Threshold.ValueString(),
		plan.Unit.ValueString(),
	)
	request := n.client.UpdateServerBandwidthNotificationSetting(
		ctx,
		plan.DedicatedServerId.ValueString(),
		plan.Id.ValueString(),
	).BandwidthNotificationSettingOpts(*opts)
	result, response, err := request.Execute()
	if err != nil {
		summary := fmt.Sprintf(
			"Updating resource %s for id %q and dedicated_server_id %q",
			n.name,
			plan.Id.ValueString(),
			plan.DedicatedServerId.ValueString(),
		)
		utils.Error(ctx, &resp.Diagnostics, summary, err, response)
		return
	}

	resp.Diagnostics.Append(
		resp.State.Set(
			ctx,
			notificationSettingBandwidthResourceModel{
				Id:                plan.Id,
				DedicatedServerId: plan.DedicatedServerId,
				Frequency:         types.StringValue(result.GetFrequency()),
				Threshold:         types.StringValue(result.GetThreshold()),
				Unit:              types.StringValue(result.GetUnit()),
			},
		)...,
	)
}

func (n *notificationSettingBandwidthResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state notificationSettingBandwidthResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := n.client.DeleteServerBandwidthNotificationSetting(
		ctx,
		state.DedicatedServerId.ValueString(),
		state.Id.ValueString(),
	)
	response, err := request.Execute()
	if err != nil {
		summary := fmt.Sprintf(
			"Deleting resource %s for id %q and dedicated_server_id %q",
			n.name,
			state.Id.ValueString(),
			state.DedicatedServerId.ValueString(),
		)
		utils.Error(ctx, &resp.Diagnostics, summary, err, response)
	}
}
