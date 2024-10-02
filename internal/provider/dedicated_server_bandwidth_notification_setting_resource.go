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
	_ resource.Resource              = &bandwidthNotificationSettingResource{}
	_ resource.ResourceWithConfigure = &bandwidthNotificationSettingResource{}
)

type bandwidthNotificationSettingResource struct {
	// TODO: Refactor this part, apiKey shouldn't be here.
	apiKey string
	client dedicatedServer.DedicatedServerAPI
}

type bandwidthNotificationSettingResourceData struct {
	Id                types.String `tfsdk:"id"`
	DedicatedServerId types.String `tfsdk:"dedicated_server_id"`
	Frequency         types.String `tfsdk:"frequency"`
	Threshold         types.String `tfsdk:"threshold"`
	Unit              types.String `tfsdk:"unit"`
}

func NewBandwidthNotificationSettingResource() resource.Resource {
	return &bandwidthNotificationSettingResource{}
}

func (b *bandwidthNotificationSettingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dedicated_server_bandwidth_notification_setting"
}

func (b *bandwidthNotificationSettingResource) authContext(ctx context.Context) context.Context {
	return context.WithValue(
		ctx,
		dedicatedServer.ContextAPIKeys,
		map[string]dedicatedServer.APIKey{
			"X-LSW-Auth": {Key: b.apiKey, Prefix: ""},
		},
	)
}

func (b *bandwidthNotificationSettingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	b.apiKey = coreClient.ProviderData.ApiKey
	if coreClient.ProviderData.Host != nil {
		configuration.Host = *coreClient.ProviderData.Host
	}
	if coreClient.ProviderData.Scheme != nil {
		configuration.Scheme = *coreClient.ProviderData.Scheme
	}

	apiClient := dedicatedServer.NewAPIClient(configuration)
	b.client = apiClient.DedicatedServerAPI
}

func (b *bandwidthNotificationSettingResource) Schema(
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
					customValidators.GreaterThanZero(),
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

func (b *bandwidthNotificationSettingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data bandwidthNotificationSettingResourceData
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
	request := b.client.CreateServerBandwidthNotificationSetting(b.authContext(ctx), data.DedicatedServerId.ValueString()).BandwidthNotificationSettingOpts(*opts)
	result, response, err := request.Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf(
				"Error creating bandwidth notification setting with dedicated_server_id: %q",
				data.DedicatedServerId.ValueString(),
			),
			getHttpErrorMessage(response, err),
		)
		return
	}

	newData := bandwidthNotificationSettingResourceData{
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

func (b *bandwidthNotificationSettingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data bandwidthNotificationSettingResourceData
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := b.client.GetServerBandwidthNotificationSetting(b.authContext(ctx), data.DedicatedServerId.ValueString(), data.Id.ValueString())
	result, response, err := request.Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf(
				"Error reading bandwidth notification setting with id: %q and dedicated_server_id: %q",
				data.Id.ValueString(),
				data.DedicatedServerId.ValueString(),
			),
			getHttpErrorMessage(response, err),
		)
		return
	}

	newData := bandwidthNotificationSettingResourceData{
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

func (b *bandwidthNotificationSettingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data bandwidthNotificationSettingResourceData
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
	request := b.client.UpdateServerBandwidthNotificationSetting(b.authContext(ctx), data.DedicatedServerId.ValueString(), data.Id.ValueString()).BandwidthNotificationSettingOpts(*opts)
	result, response, err := request.Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf(
				"Error updating bandwidth notification setting with id: %q and dedicated_server_id: %q",
				data.Id.ValueString(),
				data.DedicatedServerId.ValueString(),
			),
			getHttpErrorMessage(response, err),
		)
		return
	}

	newData := bandwidthNotificationSettingResourceData{
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

func (b *bandwidthNotificationSettingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data bandwidthNotificationSettingResourceData
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := b.client.DeleteServerBandwidthNotificationSetting(b.authContext(ctx), data.DedicatedServerId.ValueString(), data.Id.ValueString())
	response, err := request.Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf(
				"Error deleting bandwidth notification setting with id: %q and dedicated_server_id: %q",
				data.Id.ValueString(),
				data.DedicatedServerId.ValueString(),
			),
			getHttpErrorMessage(response, err),
		)
		return
	}
}
