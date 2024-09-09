package datatrafficnotificationsetting

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/leaseweb/leaseweb-go-sdk/dedicatedServer"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
)

var (
	_ resource.Resource              = &dataTrafficNotificationSettingResource{}
	_ resource.ResourceWithConfigure = &dataTrafficNotificationSettingResource{}
)

type dataTrafficNotificationSettingResource struct {
	apiKey string
	client dedicatedServer.DedicatedServerAPI
}

type resourceData struct {
	Id        types.String `tfsdk:"id"`
	ServerId  types.String `tfsdk:"server_id"`
	Frequency types.String `tfsdk:"frequency"`
	Threshold types.String `tfsdk:"threshold"`
	Unit      types.String `tfsdk:"unit"`
}

func NewDataTrafficNotificationSettingResource() resource.Resource {
	return &dataTrafficNotificationSettingResource{}
}

func (d *dataTrafficNotificationSettingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data_traffic_notification_setting"
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
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"server_id": schema.StringAttribute{
				Required: true,
			},
			"frequency": schema.StringAttribute{
				Required: true,
			},
			"threshold": schema.StringAttribute{
				Required: true,
			},
			"unit": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (d *dataTrafficNotificationSettingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data resourceData
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
	request := d.client.CreateServerDataTrafficNotificationSetting(d.authContext(ctx), data.ServerId.ValueString()).DataTrafficNotificationSettingOpts(*opts)
	result, _, err := request.Execute()
	if err != nil {
		return
	}

	newData := resourceData{
		Id:        types.StringValue(result.GetId()),
		Frequency: types.StringValue(result.GetFrequency()),
		Threshold: types.StringValue(result.GetThreshold()),
		Unit:      types.StringValue(result.GetUnit()),
	}
	newData.ServerId = data.ServerId
	diags = resp.State.Set(ctx, newData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *dataTrafficNotificationSettingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resourceData
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := d.client.GetServerDataTrafficNotificationSetting(d.authContext(ctx), data.ServerId.ValueString(), data.Id.ValueString())
	result, _, err := request.Execute()
	if err != nil {
		return
	}

	newData := resourceData{
		Id:        types.StringValue(result.GetId()),
		Frequency: types.StringValue(result.GetFrequency()),
		Threshold: types.StringValue(result.GetThreshold()),
		Unit:      types.StringValue(result.GetUnit()),
	}
	newData.ServerId = data.ServerId
	diags = resp.State.Set(ctx, newData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *dataTrafficNotificationSettingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data resourceData
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
	request := d.client.UpdateServerDataTrafficNotificationSetting(d.authContext(ctx), data.ServerId.ValueString(), data.Id.ValueString()).DataTrafficNotificationSettingOpts(*opts)
	result, _, err := request.Execute()
	if err != nil {
		return
	}

	newData := resourceData{
		Id:        data.Id,
		ServerId:  data.ServerId,
		Frequency: types.StringValue(result.GetFrequency()),
		Threshold: types.StringValue(result.GetThreshold()),
		Unit:      types.StringValue(result.GetUnit()),
	}
	diags = resp.State.Set(ctx, newData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *dataTrafficNotificationSettingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resourceData
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := d.client.DeleteServerDataTrafficNotificationSetting(d.authContext(ctx), data.ServerId.ValueString(), data.Id.ValueString())
	_, err := request.Execute()
	if err != nil {
		return
	}
}
