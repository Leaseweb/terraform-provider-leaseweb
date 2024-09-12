package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/leaseweb/leaseweb-go-sdk/dedicatedServer"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
	customValidators "github.com/leaseweb/terraform-provider-leaseweb/internal/provider/validators"
)

var (
	_ resource.Resource              = &dataTrafficNotificationSettingResource{}
	_ resource.ResourceWithConfigure = &dataTrafficNotificationSettingResource{}
)

type dataTrafficNotificationSettingResource struct {
	// TODO: Refactor this part, apiKey shouldn't be here.
	apiKey string
	client dedicatedServer.DedicatedServerAPI
}

type dataTrafficNotificationSettingResourceData struct {
	Id                types.String `tfsdk:"id"`
	DedicatedServerId types.String `tfsdk:"dedicated_server_id"`
	Frequency         types.String `tfsdk:"frequency"`
	Threshold         types.String `tfsdk:"threshold"`
	Unit              types.String `tfsdk:"unit"`
}

func NewDataTrafficNotificationSettingResource() resource.Resource {
	return &dataTrafficNotificationSettingResource{}
}

func (d *dataTrafficNotificationSettingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dedicated_server_data_traffic_notification_setting"
}

func (d *dataTrafficNotificationSettingResource) authContext(ctx context.Context) context.Context {
	return context.WithValue(
		ctx,
		dedicatedServer.ContextAPIKeys,
		map[string]dedicatedServer.APIKey{
			"X-LSW-Auth": {Key: d.apiKey, Prefix: ""},
		},
	)
}

func (d *dataTrafficNotificationSettingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (d *dataTrafficNotificationSettingResource) Schema(
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
					customValidators.GreaterThanZero(),
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

func (d *dataTrafficNotificationSettingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data dataTrafficNotificationSettingResourceData
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
	request := d.client.CreateServerDataTrafficNotificationSetting(d.authContext(ctx), data.DedicatedServerId.ValueString()).DataTrafficNotificationSettingOpts(*opts)
	result, _, err := request.Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf(
				"Error creating data traffic notification setting with dedicated_server_id: %q",
				data.DedicatedServerId.ValueString(),
			),
			err.Error(),
		)
		return
	}

	dataTrafficNotificationSetting := dataTrafficNotificationSettingResourceData{
		DedicatedServerId: data.DedicatedServerId,
		Id:                types.StringValue(result.GetId()),
		Frequency:         types.StringValue(result.GetFrequency()),
		Threshold:         types.StringValue(result.GetThreshold()),
		Unit:              types.StringValue(result.GetUnit()),
	}
	diags = resp.State.Set(ctx, dataTrafficNotificationSetting)
	resp.Diagnostics.Append(diags...)
}

func (d *dataTrafficNotificationSettingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data dataTrafficNotificationSettingResourceData
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := d.client.GetServerDataTrafficNotificationSetting(d.authContext(ctx), data.DedicatedServerId.ValueString(), data.Id.ValueString())
	result, _, err := request.Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf(
				"Error reading data traffic notification setting with id: %q and dedicated_server_id: %q",
				data.Id.ValueString(),
				data.DedicatedServerId.ValueString(),
			),
			err.Error(),
		)
		return
	}

	dataTrafficNotificationSetting := dataTrafficNotificationSettingResourceData{
		DedicatedServerId: data.DedicatedServerId,
		Id:                types.StringValue(result.GetId()),
		Frequency:         types.StringValue(result.GetFrequency()),
		Threshold:         types.StringValue(result.GetThreshold()),
		Unit:              types.StringValue(result.GetUnit()),
	}
	diags = resp.State.Set(ctx, dataTrafficNotificationSetting)
	resp.Diagnostics.Append(diags...)
}

func (d *dataTrafficNotificationSettingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data dataTrafficNotificationSettingResourceData
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
	request := d.client.UpdateServerDataTrafficNotificationSetting(d.authContext(ctx), data.DedicatedServerId.ValueString(), data.Id.ValueString()).DataTrafficNotificationSettingOpts(*opts)
	result, _, err := request.Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf(
				"Error updating data traffic notification setting with id: %q and dedicated_server_id: %q",
				data.Id.ValueString(),
				data.DedicatedServerId.ValueString(),
			),
			err.Error(),
		)
		return
	}

	dataTrafficNotificationSetting := dataTrafficNotificationSettingResourceData{
		Id:                data.Id,
		DedicatedServerId: data.DedicatedServerId,
		Frequency:         types.StringValue(result.GetFrequency()),
		Threshold:         types.StringValue(result.GetThreshold()),
		Unit:              types.StringValue(result.GetUnit()),
	}
	diags = resp.State.Set(ctx, dataTrafficNotificationSetting)
	resp.Diagnostics.Append(diags...)
}

func (d *dataTrafficNotificationSettingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data dataTrafficNotificationSettingResourceData
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := d.client.DeleteServerDataTrafficNotificationSetting(d.authContext(ctx), data.DedicatedServerId.ValueString(), data.Id.ValueString())
	_, err := request.Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf(
				"Error deleting data traffic notification setting with id: %q and dedicated_server_id: %q",
				data.Id.ValueString(),
				data.DedicatedServerId.ValueString(),
			),
			err.Error(),
		)
		return
	}
}
