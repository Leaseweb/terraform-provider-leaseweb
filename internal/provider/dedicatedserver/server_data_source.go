package dedicatedserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ datasource.DataSource              = &serverDataSource{}
	_ datasource.DataSourceWithConfigure = &serverDataSource{}
)

type serverDataSource struct {
	utils.DataSourceAPI
}

type serverDataSourceModel struct {
	ID                                 types.String `tfsdk:"id"`
	AssetID                            types.String `tfsdk:"asset_id"`
	ContractID                         types.String `tfsdk:"contract_id"`
	CPUQuantity                        types.Int32  `tfsdk:"cpu_quantity"`
	CPUType                            types.String `tfsdk:"cpu_type"`
	InternalGateway                    types.String `tfsdk:"internal_gateway"`
	InternalIP                         types.String `tfsdk:"internal_ip"`
	InternalMAC                        types.String `tfsdk:"internal_mac"`
	IsAutomationFeatureAvailable       types.Bool   `tfsdk:"is_automation_feature_available"`
	IsIPMIRebootFeatureAvailable       types.Bool   `tfsdk:"is_ipmi_reboot_feature_available"`
	IsPowerCycleFeatureAvailable       types.Bool   `tfsdk:"is_power_cycle_feature_available"`
	IsPrivateNetworkFeatureAvailable   types.Bool   `tfsdk:"is_private_network_feature_available"`
	IsRemoteManagementFeatureAvailable types.Bool   `tfsdk:"is_remote_management_feature_available"`
	LocationRack                       types.String `tfsdk:"location_rack"`
	LocationSite                       types.String `tfsdk:"location_site"`
	LocationSuite                      types.String `tfsdk:"location_suite"`
	LocationUnit                       types.String `tfsdk:"location_unit"`
	PublicGateway                      types.String `tfsdk:"public_gateway"`
	PublicIP                           types.String `tfsdk:"public_ip"`
	PublicMAC                          types.String `tfsdk:"public_mac"`
	RackCapacity                       types.String `tfsdk:"rack_capacity"`
	RackID                             types.String `tfsdk:"rack_id"`
	RackType                           types.String `tfsdk:"rack_type"`
	RAMSize                            types.Int32  `tfsdk:"ram_size"`
	RAMUnit                            types.String `tfsdk:"ram_unit"`
	RemoteGateway                      types.String `tfsdk:"remote_gateway"`
	RemoteIP                           types.String `tfsdk:"remote_ip"`
	RemoteMAC                          types.String `tfsdk:"remote_mac"`
	SerialNumber                       types.String `tfsdk:"serial_number"`
}

func (s *serverDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config serverDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	request := s.DedicatedserverAPI.GetServer(ctx, config.ID.ValueString())
	result, response, err := request.Execute()
	if err != nil {
		utils.SdkError(ctx, &resp.Diagnostics, err, response)
		return
	}

	var contractID *string
	if contract, ok := result.GetContractOk(); ok {
		contractID, _ = contract.GetIdOk()
	}

	var rackID, rackCapacity, rackType *string
	if rack, ok := result.GetRackOk(); ok {
		rackID, _ = rack.GetIdOk()
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

	var publicMAC, publicIP, publicGateway *string
	if networkInterfaces, ok := result.GetNetworkInterfacesOk(); ok {
		if publicNetworkInterface, ok := networkInterfaces.GetPublicOk(); ok {
			publicMAC, _ = publicNetworkInterface.GetMacOk()
			publicIP, _ = publicNetworkInterface.GetIpOk()
			publicGateway, _ = publicNetworkInterface.GetGatewayOk()
		}
	}

	var internalMAC, internalIP, internalGateway *string
	if networkInterfaces, ok := result.GetNetworkInterfacesOk(); ok {
		if internalNetworkInterface, ok := networkInterfaces.GetInternalOk(); ok {
			internalMAC, _ = internalNetworkInterface.GetMacOk()
			internalIP, _ = internalNetworkInterface.GetIpOk()
			internalGateway, _ = internalNetworkInterface.GetGatewayOk()
		}
	}

	var remoteMAC, remoteIP, remoteGateway *string
	if networkInterfaces, ok := result.GetNetworkInterfacesOk(); ok {
		if remoteNetworkInterface, ok := networkInterfaces.GetRemoteManagementOk(); ok {
			remoteMAC, _ = remoteNetworkInterface.GetMacOk()
			remoteIP, _ = remoteNetworkInterface.GetIpOk()
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

	resp.Diagnostics.Append(
		resp.State.Set(
			ctx,
			serverDataSourceModel{
				ID:                                 types.StringValue(result.GetId()),
				AssetID:                            types.StringValue(result.GetAssetId()),
				ContractID:                         types.StringPointerValue(contractID),
				CPUQuantity:                        types.Int32PointerValue(cpuQuantity),
				CPUType:                            types.StringPointerValue(cpuType),
				InternalGateway:                    types.StringPointerValue(internalGateway),
				InternalIP:                         types.StringPointerValue(internalIP),
				InternalMAC:                        types.StringPointerValue(internalMAC),
				SerialNumber:                       types.StringValue(result.GetSerialNumber()),
				IsAutomationFeatureAvailable:       types.BoolPointerValue(automation),
				IsIPMIRebootFeatureAvailable:       types.BoolPointerValue(ipmiReboot),
				IsPowerCycleFeatureAvailable:       types.BoolPointerValue(powerCycle),
				IsPrivateNetworkFeatureAvailable:   types.BoolPointerValue(privateNetwork),
				IsRemoteManagementFeatureAvailable: types.BoolPointerValue(remoteManagement),
				LocationRack:                       types.StringPointerValue(locationRack),
				LocationSite:                       types.StringPointerValue(locationSite),
				LocationSuite:                      types.StringPointerValue(locationSuite),
				LocationUnit:                       types.StringPointerValue(locationUnit),
				PublicGateway:                      types.StringPointerValue(publicGateway),
				PublicIP:                           types.StringPointerValue(publicIP),
				PublicMAC:                          types.StringPointerValue(publicMAC),
				RackCapacity:                       types.StringPointerValue(rackCapacity),
				RackID:                             types.StringPointerValue(rackID),
				RackType:                           types.StringPointerValue(rackType),
				RAMSize:                            types.Int32PointerValue(ramSize),
				RAMUnit:                            types.StringPointerValue(ramUnit),
				RemoteGateway:                      types.StringPointerValue(remoteGateway),
				RemoteIP:                           types.StringPointerValue(remoteIP),
				RemoteMAC:                          types.StringPointerValue(remoteMAC),
			},
		)...,
	)
}

func (s *serverDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The unique identifier of the server.",
			},
			"asset_id": schema.StringAttribute{
				Computed:    true,
				Description: "The Asset ID of the server.",
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
				Description: "The ID of the rack.",
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

func NewServerDataSource() datasource.DataSource {
	return &serverDataSource{
		DataSourceAPI: utils.DataSourceAPI{
			Name: "dedicated_server",
		},
	}
}
