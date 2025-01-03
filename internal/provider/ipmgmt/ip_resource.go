package ipmgmt

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/ipmgmt"
	"github.com/leaseweb/leaseweb-go-sdk/publiccloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ resource.ResourceWithConfigure   = &ipResource{}
	_ resource.ResourceWithImportState = &ipResource{}
)

type ipResourceModel struct {
	AssignedContract types.Object `tfsdk:"assigned_contract"`
	EquipmentID      types.String `tfsdk:"equipment_id"`
	IP               types.String `tfsdk:"ip"`
	NullLevel        types.Int32  `tfsdk:"null_level"`
	NullRouted       types.Bool   `tfsdk:"null_routed"`
	PrefixLength     types.Int32  `tfsdk:"prefix_length"`
	Primary          types.Bool   `tfsdk:"primary"`
	ReverseLookup    types.String `tfsdk:"reverse_lookup"`
	Subnet           types.Object `tfsdk:"subnet"`
	Type             types.String `tfsdk:"type"`
	UnnullingAllowed types.Bool   `tfsdk:"unnulling_allowed"`
	Version          types.Int32  `tfsdk:"version"`
}

func adaptIPToIPResourceModel(
	ip ipmgmt.Ip,
	ctx context.Context,
	diags *diag.Diagnostics,
) *ipResourceModel {
	sdkAssignedContract, _ := ip.GetAssignedContractOk()
	assignedContract := utils.AdaptNullableSdkModelToResourceObject(
		sdkAssignedContract,
		map[string]attr.Type{"id": types.StringType},
		ctx,
		func(contract ipmgmt.AssignedContract) assignedContractResourceModel {
			return assignedContractResourceModel{
				ID: basetypes.NewStringValue(contract.GetId()),
			}
		},
		diags,
	)
	if diags.HasError() {
		return nil
	}

	subnet := utils.AdaptSdkModelToResourceObject(
		ip.GetSubnet(),
		map[string]attr.Type{
			"gateway":       types.StringType,
			"id":            types.StringType,
			"network_ip":    types.StringType,
			"prefix_length": types.Int32Type,
		},
		ctx,
		func(subnet ipmgmt.Subnet) subnetResourceModel {
			return subnetResourceModel{
				Gateway:      basetypes.NewStringValue(subnet.GetGateway()),
				ID:           basetypes.NewStringValue(subnet.GetId()),
				NetworkIP:    basetypes.NewStringValue(subnet.GetNetworkIp()),
				PrefixLength: basetypes.NewInt32Value(subnet.GetPrefixLength()),
			}
		},
		diags,
	)
	if diags.HasError() {
		return nil
	}

	return &ipResourceModel{
		AssignedContract: assignedContract,
		EquipmentID:      basetypes.NewStringValue(ip.GetEquipmentId()),
		IP:               basetypes.NewStringValue(ip.GetIp()),
		NullLevel:        basetypes.NewInt32Value(ip.GetNullLevel()),
		NullRouted:       basetypes.NewBoolValue(ip.GetNullRouted()),
		PrefixLength:     basetypes.NewInt32Value(ip.GetPrefixLength()),
		Primary:          basetypes.NewBoolValue(ip.GetPrimary()),
		ReverseLookup:    basetypes.NewStringValue(ip.GetReverseLookup()),
		Subnet:           subnet,
		Type:             basetypes.NewStringValue(string(ip.GetType())),
		UnnullingAllowed: basetypes.NewBoolValue(ip.GetUnnullingAllowed()),
		Version:          basetypes.NewInt32Value(int32(ip.GetVersion())),
	}
}

type subnetResourceModel struct {
	Gateway      types.String `tfsdk:"gateway"`
	ID           types.String `tfsdk:"id"`
	NetworkIP    types.String `tfsdk:"network_ip"`
	PrefixLength types.Int32  `tfsdk:"prefix_length"`
}

type ipResource struct {
	utils.ResourceAPI
}

func (i ipResource) ImportState(
	ctx context.Context,
	request resource.ImportStateRequest,
	response *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(
		ctx,
		path.Root("ip"),
		request,
		response,
	)
}

func (i ipResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	response *resource.SchemaResponse,
) {
	response.Schema = schema.Schema{
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
				Required:    true,
				Description: "IP address. Changing this value updates the resource state with data related to the new IP address.",
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
				Optional:    true,
				Computed:    true,
				Description: "Set reverse lookup for the IP",
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
	}
}

func (i ipResource) Create(
	ctx context.Context,
	request resource.CreateRequest,
	response *resource.CreateResponse,
) {
	utils.ImportOnlyError(&response.Diagnostics)
}

func (i ipResource) Read(
	ctx context.Context,
	request resource.ReadRequest,
	response *resource.ReadResponse,
) {
	var originalState ipResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &originalState)...)
	if response.Diagnostics.HasError() {
		return
	}

	ip, httpResponse, err := i.IPmgmtAPI.InspectIP(
		ctx,
		originalState.IP.ValueString(),
	).Execute()
	if err != nil {
		utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
		return
	}

	state := adaptIPToIPResourceModel(*ip, ctx, &response.Diagnostics)
	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

func (i ipResource) Update(
	ctx context.Context,
	request resource.UpdateRequest,
	response *resource.UpdateResponse,
) {
	var plan ipResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	opts := publiccloud.NewUpdateIPOpts(plan.ReverseLookup.ValueString())
	ip, httpResponse, err := i.IPmgmtAPI.UpdateIP(ctx, plan.IP.ValueString()).
		UpdateIPOpts(ipmgmt.UpdateIPOpts(*opts)).
		Execute()
	if err != nil {
		utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
		return
	}

	state := adaptIPToIPResourceModel(*ip, ctx, &response.Diagnostics)
	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

func (i ipResource) Delete(
	_ context.Context,
	_ resource.DeleteRequest,
	_ *resource.DeleteResponse,
) {
}

func NewIPResource() resource.Resource {
	return &ipResource{
		ResourceAPI: utils.ResourceAPI{Name: "ipmgmt_ip"},
	}
}
