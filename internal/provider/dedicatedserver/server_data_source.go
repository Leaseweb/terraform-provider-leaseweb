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
	Id                                 types.String `tfsdk:"id"`
	AssetId                            types.String `tfsdk:"asset_id"`
	ContractId                         types.String `tfsdk:"contract_id"`
	CpuQuantity                        types.Int32  `tfsdk:"cpu_quantity"`
	CpuType                            types.String `tfsdk:"cpu_type"`
	InternalGateway                    types.String `tfsdk:"internal_gateway"`
	InternalIp                         types.String `tfsdk:"internal_ip"`
	InternalMac                        types.String `tfsdk:"internal_mac"`
	IsAutomationFeatureAvailable       types.Bool   `tfsdk:"is_automation_feature_available"`
	IsIpmiRebootFeatureAvailable       types.Bool   `tfsdk:"is_ipmi_reboot_feature_available"`
	IsPowerCycleFeatureAvailable       types.Bool   `tfsdk:"is_power_cycle_feature_available"`
	IsPrivateNetworkFeatureAvailable   types.Bool   `tfsdk:"is_private_network_feature_available"`
	IsRemoteManagementFeatureAvailable types.Bool   `tfsdk:"is_remote_management_feature_available"`
	LocationRack                       types.String `tfsdk:"location_rack"`
	LocationSite                       types.String `tfsdk:"location_site"`
	LocationSuite                      types.String `tfsdk:"location_suite"`
	LocationUnit                       types.String `tfsdk:"location_unit"`
	PublicGateway                      types.String `tfsdk:"public_gateway"`
	PublicIp                           types.String `tfsdk:"public_ip"`
	PublicMac                          types.String `tfsdk:"public_mac"`
	RackCapacity                       types.String `tfsdk:"rack_capacity"`
	RackId                             types.String `tfsdk:"rack_id"`
	RackType                           types.String `tfsdk:"rack_type"`
	RamSize                            types.Int32  `tfsdk:"ram_size"`
	RamUnit                            types.String `tfsdk:"ram_unit"`
	RemoteGateway                      types.String `tfsdk:"remote_gateway"`
	RemoteIp                           types.String `tfsdk:"remote_ip"`
	RemoteMac                          types.String `tfsdk:"remote_mac"`
	SerialNumber                       types.String `tfsdk:"serial_number"`
}

func (s *serverDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config serverDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	request := s.DedicatedserverAPI.GetServer(ctx, config.Id.ValueString())
	result, response, err := request.Execute()
	if err != nil {
		utils.SdkError(ctx, &resp.Diagnostics, err, response)
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

	resp.Diagnostics.Append(
		resp.State.Set(
			ctx,
			serverDataSourceModel{
				Id:                                 types.StringValue(result.GetId()),
				AssetId:                            types.StringValue(result.GetAssetId()),
				ContractId:                         types.StringPointerValue(contractId),
				CpuQuantity:                        types.Int32PointerValue(cpuQuantity),
				CpuType:                            types.StringPointerValue(cpuType),
				InternalGateway:                    types.StringPointerValue(internalGateway),
				InternalIp:                         types.StringPointerValue(internalIp),
				InternalMac:                        types.StringPointerValue(internalMac),
				SerialNumber:                       types.StringValue(result.GetSerialNumber()),
				IsAutomationFeatureAvailable:       types.BoolPointerValue(automation),
				IsIpmiRebootFeatureAvailable:       types.BoolPointerValue(ipmiReboot),
				IsPowerCycleFeatureAvailable:       types.BoolPointerValue(powerCycle),
				IsPrivateNetworkFeatureAvailable:   types.BoolPointerValue(privateNetwork),
				IsRemoteManagementFeatureAvailable: types.BoolPointerValue(remoteManagement),
				LocationRack:                       types.StringPointerValue(locationRack),
				LocationSite:                       types.StringPointerValue(locationSite),
				LocationSuite:                      types.StringPointerValue(locationSuite),
				LocationUnit:                       types.StringPointerValue(locationUnit),
				PublicGateway:                      types.StringPointerValue(publicGateway),
				PublicIp:                           types.StringPointerValue(publicIp),
				PublicMac:                          types.StringPointerValue(publicMac),
				RackCapacity:                       types.StringPointerValue(rackCapacity),
				RackId:                             types.StringPointerValue(rackId),
				RackType:                           types.StringPointerValue(rackType),
				RamSize:                            types.Int32PointerValue(ramSize),
				RamUnit:                            types.StringPointerValue(ramUnit),
				RemoteGateway:                      types.StringPointerValue(remoteGateway),
				RemoteIp:                           types.StringPointerValue(remoteIp),
				RemoteMac:                          types.StringPointerValue(remoteMac),
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

func NewServerDataSource() datasource.DataSource {
	return &serverDataSource{
		DataSourceAPI: utils.DataSourceAPI{
			Name: "dedicated_server",
		},
	}
}
