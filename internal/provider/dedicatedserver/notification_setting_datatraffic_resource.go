package dedicatedserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/v2/dedicatedserver"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ resource.Resource              = &notificationSettingDatatrafficResource{}
	_ resource.ResourceWithConfigure = &notificationSettingDatatrafficResource{}
)

type notificationSettingDatatrafficResource struct {
	utils.DedicatedserverResourceAPI
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
		DedicatedserverResourceAPI: utils.NewDedicatedserverResourceAPI("dedicated_server_notification_setting_datatraffic"),
	}
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
	var plan notificationSettingDatatrafficResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	opts := dedicatedserver.NewDataTrafficNotificationSettingOpts(
		plan.Frequency.ValueString(),
		plan.Threshold.ValueString(),
		plan.Unit.ValueString(),
	)
	request := n.Client.CreateServerDataTrafficNotificationSetting(
		ctx,
		plan.DedicatedServerId.ValueString(),
	).DataTrafficNotificationSettingOpts(*opts)
	result, response, err := request.Execute()
	if err != nil {
		utils.SdkError(ctx, &resp.Diagnostics, err, response)
		return
	}

	resp.Diagnostics.Append(
		resp.State.Set(
			ctx,
			notificationSettingDatatrafficResourceModel{
				DedicatedServerId: plan.DedicatedServerId,
				Id:                types.StringValue(result.GetId()),
				Frequency:         types.StringValue(result.GetFrequency()),
				Threshold:         types.StringValue(result.GetThreshold()),
				Unit:              types.StringValue(result.GetUnit()),
			},
		)...,
	)
}

func (n *notificationSettingDatatrafficResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state notificationSettingDatatrafficResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := n.Client.GetServerDataTrafficNotificationSetting(
		ctx,
		state.DedicatedServerId.ValueString(),
		state.Id.ValueString(),
	)
	result, response, err := request.Execute()
	if err != nil {
		utils.SdkError(ctx, &resp.Diagnostics, err, response)
		return
	}

	resp.Diagnostics.Append(
		resp.State.Set(
			ctx,
			notificationSettingDatatrafficResourceModel{
				DedicatedServerId: state.DedicatedServerId,
				Id:                types.StringValue(result.GetId()),
				Frequency:         types.StringValue(result.GetFrequency()),
				Threshold:         types.StringValue(result.GetThreshold()),
				Unit:              types.StringValue(result.GetUnit()),
			},
		)...,
	)
}

func (n *notificationSettingDatatrafficResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan notificationSettingDatatrafficResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	opts := dedicatedserver.NewDataTrafficNotificationSettingOpts(
		plan.Frequency.ValueString(),
		plan.Threshold.ValueString(),
		plan.Unit.ValueString(),
	)
	request := n.Client.UpdateServerDataTrafficNotificationSetting(
		ctx,
		plan.DedicatedServerId.ValueString(),
		plan.Id.ValueString(),
	).DataTrafficNotificationSettingOpts(*opts)
	result, response, err := request.Execute()
	if err != nil {
		utils.SdkError(ctx, &resp.Diagnostics, err, response)
		return
	}

	resp.Diagnostics.Append(
		resp.State.Set(
			ctx,
			notificationSettingDatatrafficResourceModel{
				Id:                plan.Id,
				DedicatedServerId: plan.DedicatedServerId,
				Frequency:         types.StringValue(result.GetFrequency()),
				Threshold:         types.StringValue(result.GetThreshold()),
				Unit:              types.StringValue(result.GetUnit()),
			},
		)...,
	)
}

func (n *notificationSettingDatatrafficResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state notificationSettingDatatrafficResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := n.Client.DeleteServerDataTrafficNotificationSetting(
		ctx,
		state.DedicatedServerId.ValueString(),
		state.Id.ValueString(),
	)
	response, err := request.Execute()
	if err != nil {
		utils.SdkError(ctx, &resp.Diagnostics, err, response)
	}
}
