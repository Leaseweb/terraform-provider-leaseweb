package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/leaseweb/leaseweb-go-sdk/dedicatedServer"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
)

var (
	_ resource.Resource                = &dedicatedServerResource{}
	_ resource.ResourceWithConfigure   = &dedicatedServerResource{}
	_ resource.ResourceWithImportState = &dedicatedServerResource{}
)

type dedicatedServerResource struct {
	// TODO: Refactor this part, apiKey shouldn't be here.
	apiKey string
	client dedicatedServer.DedicatedServerAPI
}

type dedicatedServerResourceData struct {
	Id                                  types.String `tfsdk:"id"`
	AssetId                             types.String `tfsdk:"asset_id"`
	SerialNumber                        types.String `tfsdk:"serial_number"`
	RackId                              types.String `tfsdk:"rack_id"`
	RackCapacity                        types.String `tfsdk:"rack_capacity"`
	RackType                            types.String `tfsdk:"rack_type"`
	FeatureAvailabilityAutomation       types.Bool   `tfsdk:"feature_availability_automation"`
	FeatureAvailabilityIpmiReboot       types.Bool   `tfsdk:"feature_availability_ipmi_reboot"`
	FeatureAvailabilityPowerCycle       types.Bool   `tfsdk:"feature_availability_power_cycle"`
	FeatureAvailabilityPrivateNetwork   types.Bool   `tfsdk:"feature_availability_private_network"`
	FeatureAvailabilityRemoteManagement types.Bool   `tfsdk:"feature_availability_remote_management"`
	LocationRack                        types.String `tfsdk:"location_rack"`
	LocationSite                        types.String `tfsdk:"location_site"`
	LocationSuite                       types.String `tfsdk:"location_suite"`
	LocationUnit                        types.String `tfsdk:"location_unit"`
	PublicNetworkInterfaceMac           types.String `tfsdk:"public_network_interface_mac"`
	PublicNetworkInterfaceIp            types.String `tfsdk:"public_network_interface_ip"`
	PublicNetworkInterfaceGateway       types.String `tfsdk:"public_network_interface_gateway"`
	PublicNetworkInterfaceLocationId    types.String `tfsdk:"public_network_interface_location_id"`
	PublicNetworkInterfaceNullRouted    types.Bool   `tfsdk:"public_network_interface_null_routed"`
	InternalNetworkInterfaceMac         types.String `tfsdk:"internal_network_interface_mac"`
	InternalNetworkInterfaceIp          types.String `tfsdk:"internal_network_interface_ip"`
	InternalNetworkInterfaceGateway     types.String `tfsdk:"internal_network_interface_gateway"`
	InternalNetworkInterfaceLocationId  types.String `tfsdk:"internal_network_interface_location_id"`
	InternalNetworkInterfaceNullRouted  types.Bool   `tfsdk:"internal_network_interface_null_routed"`
	RemoteNetworkInterfaceMac           types.String `tfsdk:"remote_network_interface_mac"`
	RemoteNetworkInterfaceIp            types.String `tfsdk:"remote_network_interface_ip"`
	RemoteNetworkInterfaceGateway       types.String `tfsdk:"remote_network_interface_gateway"`
	RemoteNetworkInterfaceLocationId    types.String `tfsdk:"remote_network_interface_location_id"`
	RemoteNetworkInterfaceNullRouted    types.Bool   `tfsdk:"remote_network_interface_null_routed"`
	RamSize                             types.Int32  `tfsdk:"ram_size"`
	RamUnit                             types.String `tfsdk:"ram_unit"`
	CpuQuantity                         types.Int32  `tfsdk:"cpu_quantity"`
	CpuType                             types.String `tfsdk:"cpu_type"`
}

func NewDedicatedServerResource() resource.Resource {
	return &dedicatedServerResource{}
}

func (d *dedicatedServerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dedicated_server"
}

func (d *dedicatedServerResource) authContext(ctx context.Context) context.Context {
	return context.WithValue(
		ctx,
		dedicatedServer.ContextAPIKeys,
		map[string]dedicatedServer.APIKey{
			"X-LSW-Auth": {Key: d.apiKey, Prefix: ""},
		},
	)
}

func (d *dedicatedServerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (d *dedicatedServerResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {

	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"asset_id": schema.StringAttribute{
				Computed: true,
			},
			"serial_number": schema.StringAttribute{
				Computed: true,
			},
			"rack_id": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"rack_capacity": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"rack_type": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"feature_availability_automation": schema.BoolAttribute{
				Computed: true,
				Optional: true,
			},
			"feature_availability_ipmi_reboot": schema.BoolAttribute{
				Computed: true,
				Optional: true,
			},
			"feature_availability_power_cycle": schema.BoolAttribute{
				Computed: true,
				Optional: true,
			},
			"feature_availability_private_network": schema.BoolAttribute{
				Computed: true,
				Optional: true,
			},
			"feature_availability_remote_management": schema.BoolAttribute{
				Computed: true,
				Optional: true,
			},
			"location_rack": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"location_site": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"location_suite": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"location_unit": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"public_network_interface_mac": schema.StringAttribute{
				Computed: true,
			},
			"public_network_interface_ip": schema.StringAttribute{
				Computed: true,
			},
			"public_network_interface_gateway": schema.StringAttribute{
				Computed: true,
			},
			"public_network_interface_location_id": schema.StringAttribute{
				Computed: true,
			},
			"public_network_interface_null_routed": schema.BoolAttribute{
				Computed: true,
			},
			"internal_network_interface_mac": schema.StringAttribute{
				Computed: true,
			},
			"internal_network_interface_ip": schema.StringAttribute{
				Computed: true,
			},
			"internal_network_interface_gateway": schema.StringAttribute{
				Computed: true,
			},
			"internal_network_interface_location_id": schema.StringAttribute{
				Computed: true,
			},
			"internal_network_interface_null_routed": schema.BoolAttribute{
				Computed: true,
			},
			"remote_network_interface_mac": schema.StringAttribute{
				Computed: true,
			},
			"remote_network_interface_ip": schema.StringAttribute{
				Computed: true,
			},
			"remote_network_interface_gateway": schema.StringAttribute{
				Computed: true,
			},
			"remote_network_interface_location_id": schema.StringAttribute{
				Computed: true,
			},
			"remote_network_interface_null_routed": schema.BoolAttribute{
				Computed: true,
			},
			"ram_size": schema.Int32Attribute{
				Computed: true,
				Optional: true,
			},
			"ram_unit": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"cpu_quantity": schema.Int32Attribute{
				Computed: true,
				Optional: true,
			},
			"cpu_type": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
		},
	}
}

func (d *dedicatedServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	panic("unimplemented")
}

func (d *dedicatedServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data dedicatedServerResourceData
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := d.client.GetServer(d.authContext(ctx), data.Id.ValueString())
	result, _, err := request.Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error reading dedicated server", err.Error())
		return
	}

	newData := dedicatedServerResourceData{
		Id:                                  types.StringValue(result.GetId()),
		AssetId:                             types.StringValue(result.GetAssetId()),
		SerialNumber:                        types.StringValue(result.GetSerialNumber()),
		RackId:                              types.StringValue(*result.GetRack().Id),
		RackCapacity:                        types.StringValue(*result.GetRack().Capacity),
		RackType:                            types.StringValue(*result.GetRack().Type),
		FeatureAvailabilityAutomation:       types.BoolValue(*result.GetFeatureAvailability().Automation),
		FeatureAvailabilityIpmiReboot:       types.BoolValue(*result.GetFeatureAvailability().IpmiReboot),
		FeatureAvailabilityPowerCycle:       types.BoolValue(*result.GetFeatureAvailability().PowerCycle),
		FeatureAvailabilityPrivateNetwork:   types.BoolValue(*result.GetFeatureAvailability().PrivateNetwork),
		FeatureAvailabilityRemoteManagement: types.BoolValue(*result.GetFeatureAvailability().RemoteManagement),
		LocationRack:                        types.StringValue(*result.GetLocation().Rack),
		LocationSite:                        types.StringValue(*result.GetLocation().Site),
		LocationSuite:                       types.StringValue(*result.GetLocation().Suite),
		LocationUnit:                        types.StringValue(*result.GetLocation().Unit),
		PublicNetworkInterfaceMac:           types.StringValue(result.GetNetworkInterfaces().Public.GetMac()),
		PublicNetworkInterfaceIp:            types.StringValue(result.GetNetworkInterfaces().Public.GetIp()),
		PublicNetworkInterfaceGateway:       types.StringValue(result.GetNetworkInterfaces().Public.GetGateway()),
		PublicNetworkInterfaceLocationId:    types.StringValue(result.GetNetworkInterfaces().Public.GetLocationId()),
		PublicNetworkInterfaceNullRouted:    types.BoolValue(result.GetNetworkInterfaces().Public.GetNullRouted()),
		InternalNetworkInterfaceMac:         types.StringValue(result.GetNetworkInterfaces().Internal.GetMac()),
		InternalNetworkInterfaceIp:          types.StringValue(result.GetNetworkInterfaces().Internal.GetIp()),
		InternalNetworkInterfaceGateway:     types.StringValue(result.GetNetworkInterfaces().Internal.GetGateway()),
		InternalNetworkInterfaceLocationId:  types.StringValue(result.GetNetworkInterfaces().Internal.GetLocationId()),
		InternalNetworkInterfaceNullRouted:  types.BoolValue(result.GetNetworkInterfaces().Internal.GetNullRouted()),
		RemoteNetworkInterfaceMac:           types.StringValue(result.GetNetworkInterfaces().RemoteManagement.GetMac()),
		RemoteNetworkInterfaceIp:            types.StringValue(result.GetNetworkInterfaces().RemoteManagement.GetIp()),
		RemoteNetworkInterfaceGateway:       types.StringValue(result.GetNetworkInterfaces().RemoteManagement.GetGateway()),
		RemoteNetworkInterfaceLocationId:    types.StringValue(result.GetNetworkInterfaces().RemoteManagement.GetLocationId()),
		RemoteNetworkInterfaceNullRouted:    types.BoolValue(result.GetNetworkInterfaces().RemoteManagement.GetNullRouted()),
		RamSize:                             types.Int32Value(*result.GetSpecs().Ram.Size),
		RamUnit:                             types.StringValue(*result.GetSpecs().Ram.Unit),
		CpuQuantity:                         types.Int32Value(*result.GetSpecs().Cpu.Quantity),
		CpuType:                             types.StringValue(*result.GetSpecs().Cpu.Type),
	}
	diags = resp.State.Set(ctx, newData)
	resp.Diagnostics.Append(diags...)
}

func (d *dedicatedServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	resource.ImportStatePassthroughID(
		ctx,
		path.Root("id"),
		req,
		resp,
	)

	request := d.client.GetServer(d.authContext(ctx), req.ID)
	result, _, err := request.Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error importing dedicated server", err.Error())
		return
	}

	newData := dedicatedServerResourceData{
		Id:                                  types.StringValue(result.GetId()),
		AssetId:                             types.StringValue(result.GetAssetId()),
		SerialNumber:                        types.StringValue(result.GetSerialNumber()),
		RackId:                              types.StringValue(*result.GetRack().Id),
		RackCapacity:                        types.StringValue(*result.GetRack().Capacity),
		RackType:                            types.StringValue(*result.GetRack().Type),
		FeatureAvailabilityAutomation:       types.BoolValue(*result.GetFeatureAvailability().Automation),
		FeatureAvailabilityIpmiReboot:       types.BoolValue(*result.GetFeatureAvailability().IpmiReboot),
		FeatureAvailabilityPowerCycle:       types.BoolValue(*result.GetFeatureAvailability().PowerCycle),
		FeatureAvailabilityPrivateNetwork:   types.BoolValue(*result.GetFeatureAvailability().PrivateNetwork),
		FeatureAvailabilityRemoteManagement: types.BoolValue(*result.GetFeatureAvailability().RemoteManagement),
		LocationRack:                        types.StringValue(*result.GetLocation().Rack),
		LocationSite:                        types.StringValue(*result.GetLocation().Site),
		LocationSuite:                       types.StringValue(*result.GetLocation().Suite),
		LocationUnit:                        types.StringValue(*result.GetLocation().Unit),
		PublicNetworkInterfaceMac:           types.StringValue(result.GetNetworkInterfaces().Public.GetMac()),
		PublicNetworkInterfaceIp:            types.StringValue(result.GetNetworkInterfaces().Public.GetIp()),
		PublicNetworkInterfaceGateway:       types.StringValue(result.GetNetworkInterfaces().Public.GetGateway()),
		PublicNetworkInterfaceLocationId:    types.StringValue(result.GetNetworkInterfaces().Public.GetLocationId()),
		PublicNetworkInterfaceNullRouted:    types.BoolValue(result.GetNetworkInterfaces().Public.GetNullRouted()),
		InternalNetworkInterfaceMac:         types.StringValue(result.GetNetworkInterfaces().Internal.GetMac()),
		InternalNetworkInterfaceIp:          types.StringValue(result.GetNetworkInterfaces().Internal.GetIp()),
		InternalNetworkInterfaceGateway:     types.StringValue(result.GetNetworkInterfaces().Internal.GetGateway()),
		InternalNetworkInterfaceLocationId:  types.StringValue(result.GetNetworkInterfaces().Internal.GetLocationId()),
		InternalNetworkInterfaceNullRouted:  types.BoolValue(result.GetNetworkInterfaces().Internal.GetNullRouted()),
		RemoteNetworkInterfaceMac:           types.StringValue(result.GetNetworkInterfaces().RemoteManagement.GetMac()),
		RemoteNetworkInterfaceIp:            types.StringValue(result.GetNetworkInterfaces().RemoteManagement.GetIp()),
		RemoteNetworkInterfaceGateway:       types.StringValue(result.GetNetworkInterfaces().RemoteManagement.GetGateway()),
		RemoteNetworkInterfaceLocationId:    types.StringValue(result.GetNetworkInterfaces().RemoteManagement.GetLocationId()),
		RemoteNetworkInterfaceNullRouted:    types.BoolValue(result.GetNetworkInterfaces().RemoteManagement.GetNullRouted()),
		RamSize:                             types.Int32Value(*result.GetSpecs().Ram.Size),
		RamUnit:                             types.StringValue(*result.GetSpecs().Ram.Unit),
		CpuQuantity:                         types.Int32Value(*result.GetSpecs().Cpu.Quantity),
		CpuType:                             types.StringValue(*result.GetSpecs().Cpu.Type),
	}
	diags := resp.State.Set(ctx, newData)
	resp.Diagnostics.Append(diags...)
}

func (d *dedicatedServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	panic("unimplemented")
}

func (d *dedicatedServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	panic("unimplemented")
}
