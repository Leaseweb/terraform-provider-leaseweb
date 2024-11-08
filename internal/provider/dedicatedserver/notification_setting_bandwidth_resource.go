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
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/leaseweb/leaseweb-go-sdk/dedicatedServer"
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
	coreClient, ok := utils.GetResourceClient(req, resp)
	if !ok {
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
	var data notificationSettingBandwidthResourceModel
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	opts := dedicatedServer.NewBandwidthNotificationSettingOpts(
		data.Frequency.ValueString(),
		data.Threshold.ValueString(),
		data.Unit.ValueString(),
	)
	request := n.client.CreateServerBandwidthNotificationSetting(
		ctx,
		data.DedicatedServerId.ValueString(),
	).BandwidthNotificationSettingOpts(*opts)
	result, response, err := request.Execute()
	if err != nil {
		summary := fmt.Sprintf(
			"Creating resource %s for dedicated_server_id %q",
			n.name,
			data.DedicatedServerId.ValueString(),
		)
		resp.Diagnostics.AddError(summary, utils.NewError(response, err).Error())
		tflog.Error(ctx, fmt.Sprintf("%s %s", summary, utils.NewError(response, err).Error()))
		return
	}

	newData := notificationSettingBandwidthResourceModel{
		Id:        types.StringValue(result.GetId()),
		Frequency: types.StringValue(result.GetFrequency()),
		Threshold: types.StringValue(result.GetThreshold()),
		Unit:      types.StringValue(result.GetUnit()),
	}
	newData.DedicatedServerId = data.DedicatedServerId
	diags = resp.State.Set(ctx, newData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (n *notificationSettingBandwidthResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var data notificationSettingBandwidthResourceModel
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := n.client.GetServerBandwidthNotificationSetting(
		ctx,
		data.DedicatedServerId.ValueString(),
		data.Id.ValueString(),
	)
	result, response, err := request.Execute()
	if err != nil {
		summary := fmt.Sprintf(
			"Reading resource %s for id %q and dedicated_server_id %q",
			n.name,
			data.Id.ValueString(),
			data.DedicatedServerId.ValueString(),
		)
		resp.Diagnostics.AddError(summary, utils.NewError(response, err).Error())
		tflog.Error(ctx, fmt.Sprintf("%s %s", summary, utils.NewError(response, err).Error()))
		return
	}

	newData := notificationSettingBandwidthResourceModel{
		Id:        types.StringValue(result.GetId()),
		Frequency: types.StringValue(result.GetFrequency()),
		Threshold: types.StringValue(result.GetThreshold()),
		Unit:      types.StringValue(result.GetUnit()),
	}
	newData.DedicatedServerId = data.DedicatedServerId
	diags = resp.State.Set(ctx, newData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (n *notificationSettingBandwidthResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var data notificationSettingBandwidthResourceModel
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	opts := dedicatedServer.NewBandwidthNotificationSettingOpts(
		data.Frequency.ValueString(),
		data.Threshold.ValueString(),
		data.Unit.ValueString(),
	)
	request := n.client.UpdateServerBandwidthNotificationSetting(
		ctx,
		data.DedicatedServerId.ValueString(),
		data.Id.ValueString(),
	).BandwidthNotificationSettingOpts(*opts)
	result, response, err := request.Execute()
	if err != nil {
		summary := fmt.Sprintf(
			"Updating resource %s for id %q and dedicated_server_id %q",
			n.name,
			data.Id.ValueString(),
			data.DedicatedServerId.ValueString(),
		)
		resp.Diagnostics.AddError(summary, utils.NewError(response, err).Error())
		tflog.Error(ctx, fmt.Sprintf("%s %s", summary, utils.NewError(response, err).Error()))
		return
	}

	newData := notificationSettingBandwidthResourceModel{
		Id:                data.Id,
		DedicatedServerId: data.DedicatedServerId,
		Frequency:         types.StringValue(result.GetFrequency()),
		Threshold:         types.StringValue(result.GetThreshold()),
		Unit:              types.StringValue(result.GetUnit()),
	}
	diags = resp.State.Set(ctx, newData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (n *notificationSettingBandwidthResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var data notificationSettingBandwidthResourceModel
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := n.client.DeleteServerBandwidthNotificationSetting(
		ctx,
		data.DedicatedServerId.ValueString(),
		data.Id.ValueString(),
	)
	response, err := request.Execute()
	if err != nil {
		summary := fmt.Sprintf(
			"Deleting resource %s for id %q and dedicated_server_id %q",
			n.name,
			data.Id.ValueString(),
			data.DedicatedServerId.ValueString(),
		)
		resp.Diagnostics.AddError(summary, utils.NewError(response, err).Error())
		tflog.Error(ctx, fmt.Sprintf("%s %s", summary, utils.NewError(response, err).Error()))
		return
	}
}
