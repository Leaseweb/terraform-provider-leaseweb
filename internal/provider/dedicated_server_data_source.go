package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/leaseweb/leaseweb-go-sdk/dedicatedServer"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
)

var (
	_ datasource.DataSource              = &dedicatedServerDataSource{}
	_ datasource.DataSourceWithConfigure = &dedicatedServerDataSource{}
)

type dedicatedServerDataSource struct {
	// TODO: Refactor this part, apiKey shouldn't be here.
	apiKey string
	client dedicatedServer.DedicatedServerAPI
}

type dedicatedServerDataSourceData struct {
	Id                                 types.String `tfsdk:"id"`
	AssetId                            types.String `tfsdk:"asset_id"`
	SerialNumber                       types.String `tfsdk:"serial_number"`
	ContractId                         types.String `tfsdk:"contract_id"`
	RackId                             types.String `tfsdk:"rack_id"`
	RackCapacity                       types.String `tfsdk:"rack_capacity"`
	RackType                           types.String `tfsdk:"rack_type"`
	IsAutomationFeatureAvailable       types.Bool   `tfsdk:"is_automation_feature_available"`
	IsIpmiRebootFeatureAvailable       types.Bool   `tfsdk:"is_ipmi_reboot_feature_available"`
	IsPowerCycleFeatureAvailable       types.Bool   `tfsdk:"is_power_cycle_feature_available"`
	IsPrivateNetworkFeatureAvailable   types.Bool   `tfsdk:"is_private_network_feature_available"`
	IsRemoteManagementFeatureAvailable types.Bool   `tfsdk:"is_remote_management_feature_available"`
	LocationRack                       types.String `tfsdk:"location_rack"`
	LocationSite                       types.String `tfsdk:"location_site"`
	LocationSuite                      types.String `tfsdk:"location_suite"`
	LocationUnit                       types.String `tfsdk:"location_unit"`
	PublicMac                          types.String `tfsdk:"public_mac"`
	PublicIp                           types.String `tfsdk:"public_ip"`
	PublicGateway                      types.String `tfsdk:"public_gateway"`
	InternalMac                        types.String `tfsdk:"internal_mac"`
	InternalIp                         types.String `tfsdk:"internal_ip"`
	InternalGateway                    types.String `tfsdk:"internal_gateway"`
	RemoteMac                          types.String `tfsdk:"remote_mac"`
	RemoteIp                           types.String `tfsdk:"remote_ip"`
	RemoteGateway                      types.String `tfsdk:"remote_gateway"`
	RamSize                            types.Int32  `tfsdk:"ram_size"`
	RamUnit                            types.String `tfsdk:"ram_unit"`
	CpuQuantity                        types.Int32  `tfsdk:"cpu_quantity"`
	CpuType                            types.String `tfsdk:"cpu_type"`
}

func (d *dedicatedServerDataSource) authContext(ctx context.Context) context.Context {
	return context.WithValue(
		ctx,
		dedicatedServer.ContextAPIKeys,
		map[string]dedicatedServer.APIKey{
			"X-LSW-Auth": {Key: d.apiKey, Prefix: ""},
		},
	)
}

func (d *dedicatedServerDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *dedicatedServerDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dedicated_server"
}

func (d *dedicatedServerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var data dedicatedServerDataSourceData
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	request := d.client.GetServer(d.authContext(ctx), data.Id.ValueString())
	result, response, err := request.Execute()
	if err != nil {
		summary := fmt.Sprintf("Error reading dedicated server with id: %q", data.Id.ValueString())
		resp.Diagnostics.AddError(summary, NewError(response, err).Error())
		tflog.Error(ctx, fmt.Sprintf("%s %s", summary, NewError(response, err).Error()))
		return
	}

	var contractId *string
	if contract, ok := result.GetContractOk(); ok {
		contractId, _ = contract.GetIdOk()
	}

	var rackId, rackCapacity, rackType *string
	if rack, ok := result.GetRackOk(); ok {
		rackId, _ = rack.GetIdOk()
		rackCapacity, _ = rack.GetCapacityOk()

		if rt, ok := rack.GetTypeOk(); ok && rt != nil {
			rtStr := string(*rt)
			rackType = &rtStr
		}
	}

	var automation, ipmiReboot, powerCycle, privateNetwork, remoteManagement *bool
	if featureAvailability, ok := result.GetFeatureAvailabilityOk(); ok {
		automation, _ = featureAvailability.GetAutomationOk()
		ipmiReboot, _ = featureAvailability.GetIpmiRebootOk()
		powerCycle, _ = featureAvailability.GetPowerCycleOk()
		privateNetwork, _ = featureAvailability.GetPrivateNetworkOk()
		remoteManagement, _ = featureAvailability.GetRemoteManagementOk()
	}

	var locationRack, locationSite, locationSuite, locationUnit *string
	if location, ok := result.GetLocationOk(); ok {
		locationRack, _ = location.GetRackOk()
		locationSite, _ = location.GetSiteOk()
		locationSuite, _ = location.GetSuiteOk()
		locationUnit, _ = location.GetUnitOk()
	}

	var publicMac, publicIp, publicGateway *string
	if networkInterfaces, ok := result.GetNetworkInterfacesOk(); ok {
		if publicNetworkInterface, ok := networkInterfaces.GetPublicOk(); ok {
			publicMac, _ = publicNetworkInterface.GetMacOk()
			publicIp, _ = publicNetworkInterface.GetIpOk()
			publicGateway, _ = publicNetworkInterface.GetGatewayOk()
		}
	}

	var internalMac, internalIp, internalGateway *string
	if networkInterfaces, ok := result.GetNetworkInterfacesOk(); ok {
		if internalNetworkInterface, ok := networkInterfaces.GetInternalOk(); ok {
			internalMac, _ = internalNetworkInterface.GetMacOk()
			internalIp, _ = internalNetworkInterface.GetIpOk()
			internalGateway, _ = internalNetworkInterface.GetGatewayOk()
		}
	}

	var remoteMac, remoteIp, remoteGateway *string
	if networkInterfaces, ok := result.GetNetworkInterfacesOk(); ok {
		if remoteNetworkInterface, ok := networkInterfaces.GetRemoteManagementOk(); ok {
			remoteMac, _ = remoteNetworkInterface.GetMacOk()
			remoteIp, _ = remoteNetworkInterface.GetIpOk()
			remoteGateway, _ = remoteNetworkInterface.GetGatewayOk()
		}
	}

	var ramSize *int32
	var ramUnit *string
	if specs, ok := result.GetSpecsOk(); ok {
		if ram, ok := specs.GetRamOk(); ok {
			ramSize, _ = ram.GetSizeOk()
			ramUnit, _ = ram.GetUnitOk()
		}
	}

	var cpuQuantity *int32
	var cpuType *string
	if specs, ok := result.GetSpecsOk(); ok {
		if cpu, ok := specs.GetCpuOk(); ok {
			cpuQuantity, _ = cpu.GetQuantityOk()
			cpuType, _ = cpu.GetTypeOk()
		}
	}

	data = dedicatedServerDataSourceData{
		Id:                                 types.StringValue(result.GetId()),
		AssetId:                            types.StringValue(result.GetAssetId()),
		SerialNumber:                       types.StringValue(result.GetSerialNumber()),
		ContractId:                         types.StringPointerValue(contractId),
		RackId:                             types.StringPointerValue(rackId),
		RackCapacity:                       types.StringPointerValue(rackCapacity),
		RackType:                           types.StringPointerValue(rackType),
		IsAutomationFeatureAvailable:       types.BoolPointerValue(automation),
		IsIpmiRebootFeatureAvailable:       types.BoolPointerValue(ipmiReboot),
		IsPowerCycleFeatureAvailable:       types.BoolPointerValue(powerCycle),
		IsPrivateNetworkFeatureAvailable:   types.BoolPointerValue(privateNetwork),
		IsRemoteManagementFeatureAvailable: types.BoolPointerValue(remoteManagement),
		LocationRack:                       types.StringPointerValue(locationRack),
		LocationSite:                       types.StringPointerValue(locationSite),
		LocationSuite:                      types.StringPointerValue(locationSuite),
		LocationUnit:                       types.StringPointerValue(locationUnit),
		PublicMac:                          types.StringPointerValue(publicMac),
		PublicIp:                           types.StringPointerValue(publicIp),
		PublicGateway:                      types.StringPointerValue(publicGateway),
		InternalMac:                        types.StringPointerValue(internalMac),
		InternalIp:                         types.StringPointerValue(internalIp),
		InternalGateway:                    types.StringPointerValue(internalGateway),
		RemoteMac:                          types.StringPointerValue(remoteMac),
		RemoteIp:                           types.StringPointerValue(remoteIp),
		RemoteGateway:                      types.StringPointerValue(remoteGateway),
		RamSize:                            types.Int32PointerValue(ramSize),
		RamUnit:                            types.StringPointerValue(ramUnit),
		CpuQuantity:                        types.Int32PointerValue(cpuQuantity),
		CpuType:                            types.StringPointerValue(cpuType),
	}

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *dedicatedServerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The unique identifier of the server.",
			},
			"asset_id": schema.StringAttribute{
				Computed:    true,
				Description: "The Asset Id of the server.",
			},
			"serial_number": schema.StringAttribute{
				Computed:    true,
				Description: "Serial number of server.",
			},
			"contract_id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier of the contract.",
			},
			"rack_id": schema.StringAttribute{
				Computed:    true,
				Description: "The Id of the rack.",
			},
			"rack_capacity": schema.StringAttribute{
				Computed:    true,
				Description: "The capacity of the rack.",
			},
			"rack_type": schema.StringAttribute{
				Computed:    true,
				Description: "The type of the rack.",
			},
			"is_automation_feature_available": schema.BoolAttribute{
				Computed:    true,
				Description: "To check if automation feature is available for the server.",
			},
			"is_ipmi_reboot_feature_available": schema.BoolAttribute{
				Computed:    true,
				Description: "To check if ipmi_reboot feature is available for the server.",
			},
			"is_power_cycle_feature_available": schema.BoolAttribute{
				Computed:    true,
				Description: "To check if power_cycle feature is available for the server.",
			},
			"is_private_network_feature_available": schema.BoolAttribute{
				Computed:    true,
				Description: "To check if private network feature is available for the server.",
			},
			"is_remote_management_feature_available": schema.BoolAttribute{
				Computed:    true,
				Description: "To check if remote management feature is available for the server.",
			},
			"location_rack": schema.StringAttribute{
				Computed: true,
			},
			"location_site": schema.StringAttribute{
				Computed:    true,
				Description: "The site of the location.",
			},
			"location_suite": schema.StringAttribute{
				Computed:    true,
				Description: "The suite of the location.",
			},
			"location_unit": schema.StringAttribute{
				Computed:    true,
				Description: "The unit of the location.",
			},
			"public_mac": schema.StringAttribute{
				Computed:    true,
				Description: "Public mac address.",
			},
			"public_ip": schema.StringAttribute{
				Computed:    true,
				Description: "Public ip address.",
			},
			"public_gateway": schema.StringAttribute{
				Computed:    true,
				Description: "Public gateway.",
			},
			"internal_mac": schema.StringAttribute{
				Computed:    true,
				Description: "Internal mac address.",
			},
			"internal_ip": schema.StringAttribute{
				Computed:    true,
				Description: "Internal ip address.",
			},
			"internal_gateway": schema.StringAttribute{
				Computed:    true,
				Description: "Internal gateway.",
			},
			"remote_mac": schema.StringAttribute{
				Computed:    true,
				Description: "Remote mac address.",
			},
			"remote_ip": schema.StringAttribute{
				Computed:    true,
				Description: "Remote ip address.",
			},
			"remote_gateway": schema.StringAttribute{
				Computed:    true,
				Description: "Remote gateway.",
			},
			"ram_size": schema.Int32Attribute{
				Computed:    true,
				Description: "The size of the ram.",
			},
			"ram_unit": schema.StringAttribute{
				Computed:    true,
				Description: "The unit of the ram.",
			},
			"cpu_quantity": schema.Int32Attribute{
				Computed:    true,
				Description: "The quantity of the cpu.",
			},
			"cpu_type": schema.StringAttribute{
				Computed:    true,
				Description: "The type of the cpu.",
			},
		},
	}
}

func NewDedicatedServerDataSource() datasource.DataSource {
	return &dedicatedServerDataSource{}
}
