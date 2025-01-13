package dedicatedserver

import (
	"context"
	"net"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/dedicatedserver"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ resource.Resource                = &serverResource{}
	_ resource.ResourceWithConfigure   = &serverResource{}
	_ resource.ResourceWithImportState = &serverResource{}
)

type serverResource struct {
	utils.ResourceAPI
}

type serverResourceModel struct {
	ID                           types.String `tfsdk:"id"`
	Reference                    types.String `tfsdk:"reference"`
	ReverseLookup                types.String `tfsdk:"reverse_lookup"`
	DHCPLease                    types.String `tfsdk:"dhcp_lease"`
	PoweredOn                    types.Bool   `tfsdk:"powered_on"`
	PublicNetworkInterfaceOpened types.Bool   `tfsdk:"public_network_interface_opened"`
	PublicIPNullRouted           types.Bool   `tfsdk:"public_ip_null_routed"`
	PublicIP                     types.String `tfsdk:"public_ip"`
	RemoteManagementIP           types.String `tfsdk:"remote_management_ip"`
	InternalMAC                  types.String `tfsdk:"internal_mac"`
	Location                     types.Object `tfsdk:"location"`
}

type locationResourceModel struct {
	Rack  types.String `tfsdk:"rack"`
	Site  types.String `tfsdk:"site"`
	Suite types.String `tfsdk:"suite"`
	Unit  types.String `tfsdk:"unit"`
}

func NewServerResource() resource.Resource {
	return &serverResource{
		ResourceAPI: utils.ResourceAPI{
			Name: "dedicated_server",
		},
	}
}

func (s *serverResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier of the server.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"reference": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Reference of server.",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
			},
			"reverse_lookup": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The reverse lookup associated with the dedicated server public IP.",
			},
			"dhcp_lease": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The URL of PXE boot the dedicated server is booting from.",
			},
			"powered_on": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether the dedicated server is powered on or not.",
			},
			"public_network_interface_opened": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether the public network interface of the dedicated server is opened or not.",
			},
			"public_ip_null_routed": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether the public IP of the dedicated server is null routed or not.",
			},
			"public_ip": schema.StringAttribute{
				Computed:    true,
				Description: "The public IP of the dedicated server.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"remote_management_ip": schema.StringAttribute{
				Computed:    true,
				Description: "The remote management IP of the dedicated server.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"internal_mac": schema.StringAttribute{
				Computed:    true,
				Description: "The MAC address of the interface connected to internal private network.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"location": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"rack": schema.StringAttribute{
						Computed:    true,
						Description: "the location rack",
					},
					"site": schema.StringAttribute{
						Computed:    true,
						Description: "the location site",
					},
					"suite": schema.StringAttribute{
						Computed:    true,
						Description: "the location suite",
					},
					"unit": schema.StringAttribute{
						Computed:    true,
						Description: "the location unit",
					},
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}

	utils.AddUnsupportedActionsNotation(
		resp,
		[]utils.Action{utils.CreateAction, utils.DeleteAction},
	)
}

func (s *serverResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state serverResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Getting server info
	server, httpResponse, err := s.DedicatedserverAPI.GetServer(
		ctx,
		state.ID.ValueString(),
	).Execute()
	if err != nil {
		utils.SdkError(ctx, &resp.Diagnostics, err, httpResponse)
		return
	}

	var publicIP string
	var publicIPNullRouted bool
	if networkInterfaces, ok := server.GetNetworkInterfacesOk(); ok {
		if publicNetworkInterface, ok := networkInterfaces.GetPublicOk(); ok {
			publicIPPart := strings.Split(publicNetworkInterface.GetIp(), "/")
			ip := net.ParseIP(publicIPPart[0])
			if ip != nil {
				publicIP = ip.String()
			}
			publicIPNullRouted = publicNetworkInterface.GetNullRouted()
		}
	}

	var reference string
	if contract, ok := server.GetContractOk(); ok {
		reference = contract.GetReference()
	}

	var internalMAC string
	if networkInterfaces, ok := server.GetNetworkInterfacesOk(); ok {
		if internalNetworkInterface, ok := networkInterfaces.GetInternalOk(); ok {
			internalMAC = internalNetworkInterface.GetMac()
		}
	}

	var remoteManagementIP string
	if networkInterfaces, ok := server.GetNetworkInterfacesOk(); ok {
		if remoteNetworkInterface, ok := networkInterfaces.GetRemoteManagementOk(); ok {
			remoteManagementIPPart := strings.Split(remoteNetworkInterface.GetIp(), "/")
			ip := net.ParseIP(remoteManagementIPPart[0])
			if ip != nil {
				remoteManagementIP = ip.String()
			}
		}
	}

	serverLocation := server.GetLocation()
	location, diags := types.ObjectValueFrom(
		ctx,
		map[string]attr.Type{
			"rack":  types.StringType,
			"site":  types.StringType,
			"suite": types.StringType,
			"unit":  types.StringType,
		},
		locationResourceModel{
			Rack:  types.StringValue(serverLocation.GetRack()),
			Site:  types.StringValue(serverLocation.GetSite()),
			Suite: types.StringValue(serverLocation.GetSuite()),
			Unit:  types.StringValue(serverLocation.GetUnit()),
		},
	)

	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Getting server power info
	getServerPowerStatusResult, httpResponse, err := s.DedicatedserverAPI.GetServerPowerStatus(
		ctx,
		state.ID.ValueString(),
	).Execute()
	if err != nil {
		utils.SdkError(ctx, &resp.Diagnostics, err, httpResponse)
		return
	}

	pdu := getServerPowerStatusResult.GetPdu()
	ipmi := getServerPowerStatusResult.GetIpmi()
	poweredOn := pdu.GetStatus() != "off" && ipmi.GetStatus() != "off"

	// Getting server public network interface info
	var publicNetworkOpened bool
	operationNetworkInterface, httpResponse, err := s.DedicatedserverAPI.GetNetworkInterface(
		ctx,
		state.ID.ValueString(),
		dedicatedserver.NETWORKTYPEURL_PUBLIC,
	).Execute()
	if err != nil && httpResponse != nil && httpResponse.StatusCode != http.StatusNotFound {
		utils.SdkError(ctx, &resp.Diagnostics, err, httpResponse)
		return
	} else {
		if operationNetworkInterface != nil {
			if _, ok := operationNetworkInterface.GetStatusOk(); ok {
				publicNetworkOpened = operationNetworkInterface.GetStatus() == "open"
			}
		}
	}

	// Getting server DHCP info
	getServerDhcpReservationListResult, httpResponse, err := s.DedicatedserverAPI.GetServerDhcpReservationList(
		ctx,
		state.ID.ValueString(),
	).Execute()
	if err != nil {
		utils.SdkError(ctx, &resp.Diagnostics, err, httpResponse)
		return
	}
	var dhcpLease string
	if len(getServerDhcpReservationListResult.GetLeases()) != 0 {
		leases := getServerDhcpReservationListResult.GetLeases()
		dhcpLease = leases[0].GetBootfile()
	}

	// Getting server public IP info
	var reverseLookup string
	if publicIP != "" {
		ip, httpResponse, err := s.DedicatedserverAPI.GetServerIp(
			ctx,
			state.ID.ValueString(),
			publicIP,
		).Execute()
		if err != nil {
			utils.SdkError(ctx, &resp.Diagnostics, err, httpResponse)
			return
		}
		reverseLookup = ip.GetReverseLookup()
	}

	resp.Diagnostics.Append(
		resp.State.Set(
			ctx,
			serverResourceModel{
				ID:                           types.StringValue(server.GetId()),
				Reference:                    types.StringValue(reference),
				ReverseLookup:                types.StringValue(reverseLookup),
				DHCPLease:                    types.StringValue(dhcpLease),
				PoweredOn:                    types.BoolValue(poweredOn),
				PublicNetworkInterfaceOpened: types.BoolValue(publicNetworkOpened),
				PublicIPNullRouted:           types.BoolValue(publicIPNullRouted),
				PublicIP:                     types.StringValue(publicIP),
				RemoteManagementIP:           types.StringValue(remoteManagementIP),
				InternalMAC:                  types.StringValue(internalMAC),
				Location:                     location,
			},
		)...,
	)
}

func (s *serverResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(
		ctx,
		path.Root("id"),
		req,
		resp,
	)
}

func (s *serverResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan serverResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state serverResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Updating reference
	if !plan.Reference.IsNull() && !plan.Reference.IsUnknown() {
		opts := dedicatedserver.NewUpdateServerReferenceOpts(plan.Reference.ValueString())
		response, err := s.DedicatedserverAPI.UpdateServerReference(
			ctx,
			state.ID.ValueString(),
		).UpdateServerReferenceOpts(*opts).Execute()
		if err != nil {
			utils.SdkError(ctx, &resp.Diagnostics, err, response)
			return
		}
		state.Reference = plan.Reference
	}

	// Updating Power status
	if !plan.PoweredOn.IsNull() && !plan.PoweredOn.IsUnknown() {
		if plan.PoweredOn.ValueBool() {
			request := s.DedicatedserverAPI.PowerServerOn(ctx, state.ID.ValueString())
			response, err := request.Execute()
			if err != nil {
				utils.SdkError(ctx, &resp.Diagnostics, err, response)
				return
			}
		} else {
			request := s.DedicatedserverAPI.PowerServerOff(ctx, state.ID.ValueString())
			response, err := request.Execute()
			if err != nil {
				utils.SdkError(ctx, &resp.Diagnostics, err, response)
				return
			}
		}
		state.PoweredOn = plan.PoweredOn
	}

	// Updating Reverse Lookup
	isPublicIPExists := !state.PublicIP.IsNull() && !state.PublicIP.IsUnknown() && state.PublicIP.ValueString() != ""
	if !plan.ReverseLookup.IsNull() && !plan.ReverseLookup.IsUnknown() && isPublicIPExists {
		opts := dedicatedserver.NewUpdateIpProfileOpts()
		opts.ReverseLookup = plan.ReverseLookup.ValueStringPointer()
		_, response, err := s.DedicatedserverAPI.UpdateIpProfile(
			ctx,
			state.ID.ValueString(),
			state.PublicIP.ValueString(),
		).UpdateIpProfileOpts(*opts).Execute()
		if err != nil {
			utils.SdkError(ctx, &resp.Diagnostics, err, response)
			return
		}
		state.ReverseLookup = plan.ReverseLookup
	}

	// Updating an IP null routing
	if !plan.PublicIPNullRouted.IsNull() && !plan.PublicIPNullRouted.IsUnknown() && plan.PublicIPNullRouted != state.PublicIPNullRouted && isPublicIPExists {
		if plan.PublicIPNullRouted.ValueBool() {
			_, response, err := s.DedicatedserverAPI.NullIpRoute(
				ctx,
				state.ID.ValueString(),
				state.PublicIP.ValueString(),
			).Execute()
			if err != nil {
				utils.SdkError(ctx, &resp.Diagnostics, err, response)
				return
			}
		} else {
			_, response, err := s.DedicatedserverAPI.RemoveNullIpRoute(
				ctx,
				state.ID.ValueString(),
				state.PublicIP.ValueString(),
			).Execute()
			if err != nil {
				utils.SdkError(ctx, &resp.Diagnostics, err, response)
				return
			}
		}
		state.PublicIPNullRouted = plan.PublicIPNullRouted
	}

	// Updating dhcp lease
	if !plan.DHCPLease.IsNull() && !plan.DHCPLease.IsUnknown() {
		if plan.DHCPLease.ValueString() != "" {
			opts := dedicatedserver.NewCreateServerDhcpReservationOpts(plan.DHCPLease.ValueString())
			response, err := s.DedicatedserverAPI.CreateServerDhcpReservation(
				ctx,
				state.ID.ValueString(),
			).CreateServerDhcpReservationOpts(*opts).Execute()
			if err != nil {
				utils.SdkError(ctx, &resp.Diagnostics, err, response)
				return
			}
		} else {
			response, err := s.DedicatedserverAPI.DeleteServerDhcpReservation(
				ctx,
				state.ID.ValueString(),
			).Execute()
			if err != nil {
				utils.SdkError(ctx, &resp.Diagnostics, err, response)
				return
			}
		}
		state.DHCPLease = plan.DHCPLease
	}

	// Updating network interface status
	if !plan.PublicIPNullRouted.IsNull() && !plan.PublicIPNullRouted.IsUnknown() && plan.PublicNetworkInterfaceOpened != state.PublicNetworkInterfaceOpened {
		if plan.PublicNetworkInterfaceOpened.ValueBool() {
			response, err := s.DedicatedserverAPI.OpenNetworkInterface(
				ctx,
				state.ID.ValueString(),
				dedicatedserver.NETWORKTYPEURL_PUBLIC,
			).Execute()
			if err != nil {
				utils.SdkError(ctx, &resp.Diagnostics, err, response)
				return
			}
		} else {
			response, err := s.DedicatedserverAPI.CloseNetworkInterface(
				ctx,
				state.ID.ValueString(),
				dedicatedserver.NETWORKTYPEURL_PUBLIC,
			).Execute()
			if err != nil {
				utils.SdkError(ctx, &resp.Diagnostics, err, response)
				return
			}
		}
		state.PublicNetworkInterfaceOpened = plan.PublicNetworkInterfaceOpened
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (s *serverResource) Create(
	_ context.Context,
	_ resource.CreateRequest,
	response *resource.CreateResponse,
) {
	utils.ImportOnlyError(&response.Diagnostics)
}

func (s *serverResource) Delete(
	_ context.Context,
	_ resource.DeleteRequest,
	_ *resource.DeleteResponse,
) {
}
