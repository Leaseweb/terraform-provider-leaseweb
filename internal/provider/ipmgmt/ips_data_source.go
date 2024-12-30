package ipmgmt

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/ipmgmt"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

type ipsDataSourceModel struct {
	AssignedContractIDs []string     `tfsdk:"assigned_contract_ids"`
	EquipmentIDs        []string     `tfsdk:"equipment_ids"`
	FilteredIPs         []string     `tfsdk:"filtered_ips"`
	FromIP              types.String `tfsdk:"from_ip"`
	Primary             types.Bool   `tfsdk:"primary"`
	NullRouted          types.Bool   `tfsdk:"null_routed"`
	ReverseLookup       types.String `tfsdk:"reverse_lookup"`
	Sort                []string     `tfsdk:"sort"`
	SubnetID            types.String `tfsdk:"subnet_id"`
	ToIP                types.String `tfsdk:"to_ip"`
	Type                types.String `tfsdk:"type"`
	Version             types.Int32  `tfsdk:"version"`

	IPs []ipDataSourceModel `tfsdk:"ips"`
}

type ipDataSourceModel struct {
	AssignedContract *assignedContractDataSourceModel `tfsdk:"assigned_contract"`
	EquipmentID      types.String                     `tfsdk:"equipment_id"`
	IP               types.String                     `tfsdk:"ip"`
	NullLevel        types.Int32                      `tfsdk:"null_level"`
	NullRouted       types.Bool                       `tfsdk:"null_routed"`
	PrefixLength     types.Int32                      `tfsdk:"prefix_length"`
	Primary          types.Bool                       `tfsdk:"primary"`
	ReverseLookup    types.String                     `tfsdk:"reverse_lookup"`
	Subnet           subnetDataSourceModel            `tfsdk:"subnet"`
	Type             types.String                     `tfsdk:"type"`
	UnnullingAllowed types.Bool                       `tfsdk:"unnulling_allowed"`
	Version          types.Int32                      `tfsdk:"version"`
}

type assignedContractDataSourceModel struct {
	ID types.String `tfsdk:"id"`
}

type subnetDataSourceModel struct {
	Gateway      types.String `tfsdk:"gateway"`
	ID           types.String `tfsdk:"id"`
	NetworkIP    types.String `tfsdk:"network_ip"`
	PrefixLength types.Int32  `tfsdk:"prefix_length"`
}

var (
	_ datasource.DataSourceWithConfigure = &ipsDataSource{}
)

type ipsDataSource struct{ utils.DataSourceAPI }

func (i ipsDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	response *datasource.SchemaResponse,
) {
	versions := utils.NewIntMarkdownList(ipmgmt.AllowedProtocolVersionEnumValues)

	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"assigned_contract_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Return only IPs assigned to contracts with these IDs",
			},
			"equipment_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Return only IPs assigned to equipment items",
			},
			"filtered_ips": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Return only these IPs",
			},
			"from_ip": schema.StringAttribute{
				Optional:    true,
				Description: "Return only IPs greater or equal to the specified address",
			},
			"null_routed": schema.BoolAttribute{
				Optional:    true,
				Description: "Filter by whether the IP has an active null route",
			},
			"primary": schema.BoolAttribute{
				Optional:    true,
				Description: "Filter by whether or not the IP is primary",
			},
			"reverse_lookup": schema.StringAttribute{
				Optional:    true,
				Description: "Filter by reverse lookup",
			},
			"sort": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Sort field names. Prepend the field name with '-' for descending order. E.g. `ip,-nullrouted`. Sortable field names are `ip`, `nullRouted`, `reverseLookup`",
			},
			"subnet_id": schema.StringAttribute{
				Optional:    true,
				Description: "Filter by subnet",
			},
			"to_ip": schema.StringAttribute{
				Optional:    true,
				Description: "Return only IPs lower or equal to the specified address",
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Description: "Filter by IP type. Valid options are " + utils.StringTypeArrayToMarkdown(ipmgmt.AllowedIpTypeEnumValues),
				Validators: []validator.String{
					stringvalidator.OneOf(utils.AdaptStringTypeArrayToStringArray(ipmgmt.AllowedIpTypeEnumValues)...),
				},
			},
			"version": schema.Int32Attribute{
				Optional:    true,
				Description: "Filter by protocol version. Valid options are " + versions.Markdown(),
				Validators: []validator.Int32{
					int32validator.OneOf(versions.ToInt32()...),
				},
			},

			"ips": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"assigned_contract": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Computed:    true,
									Description: "ID of the contract connected to the IP",
								},
							},
						},
						"equipment_id": schema.StringAttribute{
							Computed:    true,
							Description: "ID of the equipment using the IP",
						},
						"ip": schema.StringAttribute{
							Computed:    true,
							Description: "IP address",
						},
						"null_level": schema.Int32Attribute{
							Computed:    true,
							Description: "Null route level",
						},
						"null_routed": schema.BoolAttribute{
							Computed:    true,
							Description: "Boolean to indicate if the IP is null-routed",
						},
						"prefix_length": schema.Int32Attribute{
							Computed:    true,
							Description: "Prefix length of the IP range represented by the record. Note: this is not the same as `subnet.prefixLength`",
						},
						"primary": schema.BoolAttribute{
							Computed:    true,
							Description: "Boolean indicating if this is the primary IP of the assigned equipment",
						},
						"reverse_lookup": schema.StringAttribute{
							Computed:    true,
							Description: "Reverse lookup set for the IP. This only applies to IPv4",
						},
						"subnet": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"gateway": schema.StringAttribute{
									Computed:    true,
									Description: "The gateway IP to be used in network settings",
								},
								"id": schema.StringAttribute{
									Computed:    true,
									Description: "Subnet identifier consisting of network IP and prefix length separated by underscore (e.g. 192.0.2.0_24)",
								},
								"network_ip": schema.StringAttribute{
									Computed:    true,
									Description: "Network IP of the subnet",
								},
								"prefix_length": schema.Int32Attribute{
									Computed:    true,
									Description: "Address prefix length",
								},
							},
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "IP type",
						},
						"unnulling_allowed": schema.BoolAttribute{
							Computed:    true,
							Description: "Boolean indicating if the null route can be removed",
						},
						"version": schema.Int32Attribute{
							Computed:    true,
							Description: "Protocol version",
						},
					},
				},
			},
		},
	}
}

func (i ipsDataSource) Read(
	ctx context.Context,
	request datasource.ReadRequest,
	response *datasource.ReadResponse,
) {
	var config ipsDataSourceModel
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}

	var ips []ipmgmt.Ip
	var offset *int32
	var state ipsDataSourceModel

	ipListRequest := i.IPmgmtAPI.GetIPList(ctx)
	if len(config.AssignedContractIDs) > 0 {
		ipListRequest = ipListRequest.AssignedContractIds(strings.Join(config.AssignedContractIDs[:], ","))
		state.AssignedContractIDs = config.AssignedContractIDs
	}
	if len(config.EquipmentIDs) > 0 {
		ipListRequest = ipListRequest.EquipmentIds(strings.Join(config.EquipmentIDs[:], ","))
		state.EquipmentIDs = config.EquipmentIDs
	}
	if len(config.FilteredIPs) > 0 {
		ipListRequest = ipListRequest.Ips(strings.Join(config.FilteredIPs[:], ","))
		state.IPs = config.IPs
	}
	if !config.FromIP.IsNull() {
		ipListRequest = ipListRequest.FromIp(config.FromIP.ValueString())
		state.FromIP = config.FromIP
	}
	if !config.NullRouted.IsNull() {
		ipListRequest = ipListRequest.NullRouted(config.NullRouted.ValueBool())
		state.NullRouted = config.NullRouted
	}
	if !config.Primary.IsNull() {
		ipListRequest = ipListRequest.Primary(config.NullRouted.ValueBool())
		state.Primary = config.Primary
	}
	if !config.ReverseLookup.IsNull() {
		ipListRequest = ipListRequest.ReverseLookup(config.ReverseLookup.ValueString())
		state.ReverseLookup = config.ReverseLookup
	}
	if len(config.Sort) > 0 {
		ipListRequest = ipListRequest.Sort(strings.Join(config.Sort[:], ","))
		state.Sort = config.Sort
	}
	if !config.SubnetID.IsNull() {
		ipListRequest = ipListRequest.SubnetId(config.SubnetID.ValueString())
		state.SubnetID = config.SubnetID
	}
	if !config.ToIP.IsNull() {
		ipListRequest = ipListRequest.ToIp(config.ToIP.ValueString())
		state.ToIP = config.ToIP
	}
	if !config.Type.IsNull() {
		ipListRequest = ipListRequest.Type_(ipmgmt.IpType(config.Type.ValueString()))
		state.Type = config.Type
	}
	if !config.Version.IsNull() {
		ipListRequest = ipListRequest.Version(ipmgmt.ProtocolVersion(config.Version.ValueInt32()))
		state.Version = config.Version
	}

	for {
		result, httpResponse, err := ipListRequest.Execute()
		if err != nil {
			utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
			return
		}

		ips = append(ips, result.GetIps()...)

		metadata := result.GetMetadata()

		offset = utils.NewOffset(
			metadata.GetLimit(),
			metadata.GetOffset(),
			metadata.GetTotalCount(),
		)

		if offset == nil {
			break
		}

		ipListRequest = ipListRequest.Offset(*offset)
	}

	for _, sdkIP := range ips {
		var assignedContract *assignedContractDataSourceModel
		sdkAssignedContract, _ := sdkIP.GetAssignedContractOk()
		if sdkAssignedContract != nil {
			assignedContract = &assignedContractDataSourceModel{
				ID: basetypes.NewStringValue(sdkAssignedContract.GetId()),
			}
		}

		reverseLookup, _ := sdkIP.GetReverseLookupOk()
		nullLevel, _ := sdkIP.GetNullLevelOk()
		subnet := sdkIP.GetSubnet()

		ip := ipDataSourceModel{
			AssignedContract: assignedContract,
			EquipmentID:      basetypes.NewStringValue(sdkIP.GetEquipmentId()),
			IP:               basetypes.NewStringValue(sdkIP.GetIp()),
			NullLevel:        basetypes.NewInt32PointerValue(nullLevel),
			NullRouted:       basetypes.NewBoolValue(sdkIP.GetNullRouted()),
			PrefixLength:     basetypes.NewInt32Value(sdkIP.GetPrefixLength()),
			Primary:          basetypes.NewBoolValue(sdkIP.GetPrimary()),
			ReverseLookup:    basetypes.NewStringPointerValue(reverseLookup),
			Subnet: subnetDataSourceModel{
				Gateway:      basetypes.NewStringValue(subnet.GetGateway()),
				ID:           basetypes.NewStringValue(subnet.GetId()),
				NetworkIP:    basetypes.NewStringValue(subnet.GetNetworkIp()),
				PrefixLength: basetypes.NewInt32Value(subnet.GetPrefixLength()),
			},
			Type:             basetypes.NewStringValue(string(sdkIP.GetType())),
			UnnullingAllowed: basetypes.NewBoolValue(sdkIP.GetUnnullingAllowed()),
			Version:          basetypes.NewInt32Value(int32(sdkIP.GetVersion())),
		}
		state.IPs = append(state.IPs, ip)
	}

	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

func NewIPsDataSource() datasource.DataSource {
	return &ipsDataSource{
		DataSourceAPI: utils.DataSourceAPI{
			Name: "ipmgmt_ips",
		},
	}
}
