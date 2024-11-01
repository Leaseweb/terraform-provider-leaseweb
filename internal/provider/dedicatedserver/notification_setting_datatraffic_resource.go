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
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ resource.Resource              = &notificationSettingDatatrafficResource{}
	_ resource.ResourceWithConfigure = &notificationSettingDatatrafficResource{}
)

type notificationSettingDatatrafficResource struct {
	name   string
	client dedicatedServer.DedicatedServerAPI
}

type notificationSettingDatatrafficResourceModel struct {
	Id                types.String `tfsdk:"id"`
	DedicatedServerId types.String `tfsdk:"dedicated_server_id"`
	Frequency         types.String `tfsdk:"frequency"`
	Threshold         types.String `tfsdk:"threshold"`
	Unit              types.String `tfsdk:"unit"`
}

func NewNotificationSettingDatatrafficResource() resource.Resource {
	return &notificationSettingDatatrafficResource{
		name: "dedicated_server_notification_setting_datatraffic",
	}
}

func (n *notificationSettingDatatrafficResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, n.name)
}

func (n *notificationSettingDatatrafficResource) Configure(
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

func (n *notificationSettingDatatrafficResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the notification setting.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"dedicated_server_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the dedicated server.",
			},
			"frequency": schema.StringAttribute{
				Required:    true,
				Description: `The frequency of the notification. Can be either "DAILY", "WEEKLY" or "MONTHLY".`,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"DAILY", "WEEKLY", "MONTHLY"}...),
				},
			},
			"threshold": schema.StringAttribute{
				Required:    true,
				Description: "The threshold of the notification.",
				Validators: []validator.String{
					greaterThanZero(),
				},
			},
			"unit": schema.StringAttribute{
				Required:    true,
				Description: `The unit of the notification. Can be either "MB", "GB" or "TB".`,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"MB", "GB", "TB"}...),
				},
			},
		},
	}
}

func (n *notificationSettingDatatrafficResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var data notificationSettingDatatrafficResourceModel
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	opts := dedicatedServer.NewDataTrafficNotificationSettingOpts(
		data.Frequency.ValueString(),
		data.Threshold.ValueString(),
		data.Unit.ValueString(),
	)
	request := n.client.CreateServerDataTrafficNotificationSetting(
		ctx,
		data.DedicatedServerId.ValueString(),
	).DataTrafficNotificationSettingOpts(*opts)
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

	dataTrafficNotificationSetting := notificationSettingDatatrafficResourceModel{
		DedicatedServerId: data.DedicatedServerId,
		Id:                types.StringValue(result.GetId()),
		Frequency:         types.StringValue(result.GetFrequency()),
		Threshold:         types.StringValue(result.GetThreshold()),
		Unit:              types.StringValue(result.GetUnit()),
	}
	diags = resp.State.Set(ctx, dataTrafficNotificationSetting)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (n *notificationSettingDatatrafficResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var data notificationSettingDatatrafficResourceModel
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := n.client.GetServerDataTrafficNotificationSetting(
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

	dataTrafficNotificationSetting := notificationSettingDatatrafficResourceModel{
		DedicatedServerId: data.DedicatedServerId,
		Id:                types.StringValue(result.GetId()),
		Frequency:         types.StringValue(result.GetFrequency()),
		Threshold:         types.StringValue(result.GetThreshold()),
		Unit:              types.StringValue(result.GetUnit()),
	}
	diags = resp.State.Set(ctx, dataTrafficNotificationSetting)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (n *notificationSettingDatatrafficResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var data notificationSettingDatatrafficResourceModel
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	opts := dedicatedServer.NewDataTrafficNotificationSettingOpts(
		data.Frequency.ValueString(),
		data.Threshold.ValueString(),
		data.Unit.ValueString(),
	)
	request := n.client.UpdateServerDataTrafficNotificationSetting(
		ctx,
		data.DedicatedServerId.ValueString(),
		data.Id.ValueString(),
	).DataTrafficNotificationSettingOpts(*opts)
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

	dataTrafficNotificationSetting := notificationSettingDatatrafficResourceModel{
		Id:                data.Id,
		DedicatedServerId: data.DedicatedServerId,
		Frequency:         types.StringValue(result.GetFrequency()),
		Threshold:         types.StringValue(result.GetThreshold()),
		Unit:              types.StringValue(result.GetUnit()),
	}
	diags = resp.State.Set(ctx, dataTrafficNotificationSetting)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (n *notificationSettingDatatrafficResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var data notificationSettingDatatrafficResourceModel
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := n.client.DeleteServerDataTrafficNotificationSetting(
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
