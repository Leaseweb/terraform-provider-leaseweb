package dedicatedserver

import (
	"context"
	"fmt"
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
	"github.com/leaseweb/leaseweb-go-sdk/v2/dedicatedserver"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ resource.Resource                = &serverResource{}
	_ resource.ResourceWithConfigure   = &serverResource{}
	_ resource.ResourceWithImportState = &serverResource{}
)

type serverResource struct {
	name   string
	client dedicatedserver.DedicatedserverAPI
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

func (l locationResourceModel) attributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"rack":  types.StringType,
		"site":  types.StringType,
		"suite": types.StringType,
		"unit":  types.StringType,
	}
}

func NewServerResource() resource.Resource {
	return &serverResource{
		name: "dedicated_server",
	}
}

func (s *serverResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, s.name)
}

func (s *serverResource) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}

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

	s.client = coreClient.DedicatedserverAPI
}

func (s *serverResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
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

	newState, response, err := s.getServer(ctx, state.ID.ValueString())
	if err != nil {
		summary := fmt.Sprintf(
			"Reading resource %s for id %q",
			s.name,
			state.ID.ValueString(),
		)
		utils.Error(ctx, &resp.Diagnostics, summary, err, response)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
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

	state, response, err := s.getServer(ctx, req.ID)
	if err != nil {
		summary := fmt.Sprintf(
			"Importing resource %s for id %q",
			s.name,
			req.ID,
		)
		utils.Error(ctx, &resp.Diagnostics, summary, err, response)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
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
		response, err := s.client.UpdateServerReference(
			ctx,
			state.ID.ValueString(),
		).UpdateServerReferenceOpts(*opts).Execute()
		if err != nil {
			summary := fmt.Sprintf(
				"Updating resource %s reference for id %q",
				s.name,
				plan.ID.ValueString(),
			)
			utils.Error(ctx, &resp.Diagnostics, summary, err, response)
			return
		}
		state.Reference = plan.Reference
	}

	// Updating Power status
	if !plan.PoweredOn.IsNull() && !plan.PoweredOn.IsUnknown() {
		if plan.PoweredOn.ValueBool() {
			request := s.client.PowerServerOn(ctx, state.ID.ValueString())
			response, err := request.Execute()
			if err != nil {
				summary := fmt.Sprintf(
					"Updating resource %s powering on for id %q",
					s.name,
					state.ID.ValueString(),
				)
				utils.Error(ctx, &resp.Diagnostics, summary, err, response)
				return
			}
		} else {
			request := s.client.PowerServerOff(ctx, state.ID.ValueString())
			response, err := request.Execute()
			if err != nil {
				summary := fmt.Sprintf(
					"Updating resource %s powering off for id %q",
					s.name,
					state.ID.ValueString(),
				)
				utils.Error(ctx, &resp.Diagnostics, summary, err, response)
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
		_, response, err := s.client.UpdateIpProfile(
			ctx,
			state.ID.ValueString(),
			state.PublicIP.ValueString(),
		).UpdateIpProfileOpts(*opts).Execute()
		if err != nil {
			summary := fmt.Sprintf(
				"Updating resource %s reverse lookup for id %q",
				s.name,
				state.ID.ValueString(),
			)
			utils.Error(ctx, &resp.Diagnostics, summary, err, response)
			return
		}
		state.ReverseLookup = plan.ReverseLookup
	}

	// Updating an IP null routing
	if !plan.PublicIPNullRouted.IsNull() && !plan.PublicIPNullRouted.IsUnknown() && plan.PublicIPNullRouted != state.PublicIPNullRouted && isPublicIPExists {
		if plan.PublicIPNullRouted.ValueBool() {
			_, response, err := s.client.NullIpRoute(
				ctx,
				state.ID.ValueString(),
				state.PublicIP.ValueString(),
			).Execute()
			if err != nil {
				summary := fmt.Sprintf(
					"Updating resource %s null routing an ip for id %q and ip %q",
					s.name,
					state.ID.ValueString(),
					state.PublicIP.ValueString(),
				)
				utils.Error(ctx, &resp.Diagnostics, summary, err, response)
				return
			}
		} else {
			_, response, err := s.client.RemoveNullIpRoute(
				ctx,
				state.ID.ValueString(),
				state.PublicIP.ValueString(),
			).Execute()
			if err != nil {
				summary := fmt.Sprintf(
					"Updating resource %s remove null routing an ip for id %q and ip %q",
					s.name,
					state.ID.ValueString(),
					state.PublicIP.ValueString(),
				)
				utils.Error(ctx, &resp.Diagnostics, summary, err, response)
				return
			}
		}
		state.PublicIPNullRouted = plan.PublicIPNullRouted
	}

	// Updating dhcp lease
	if !plan.DHCPLease.IsNull() && !plan.DHCPLease.IsUnknown() {
		if plan.DHCPLease.ValueString() != "" {
			opts := dedicatedserver.NewCreateServerDhcpReservationOpts(plan.DHCPLease.ValueString())
			response, err := s.client.CreateServerDhcpReservation(
				ctx,
				state.ID.ValueString(),
			).CreateServerDhcpReservationOpts(*opts).Execute()
			if err != nil {
				summary := fmt.Sprintf(
					"Updating resource %s creating a DHCP reservation for id %q",
					s.name,
					state.ID.ValueString(),
				)
				utils.Error(ctx, &resp.Diagnostics, summary, err, response)
				return
			}
		} else {
			response, err := s.client.DeleteServerDhcpReservation(
				ctx,
				state.ID.ValueString(),
			).Execute()
			if err != nil {
				summary := fmt.Sprintf(
					"Updating resource %s deleting DHCP reservation for id %q",
					s.name,
					state.ID.ValueString(),
				)
				utils.Error(ctx, &resp.Diagnostics, summary, err, response)
				return
			}
		}
		state.DHCPLease = plan.DHCPLease
	}

	// Updating network interface status
	if !plan.PublicIPNullRouted.IsNull() && !plan.PublicIPNullRouted.IsUnknown() && plan.PublicNetworkInterfaceOpened != state.PublicNetworkInterfaceOpened {
		if plan.PublicNetworkInterfaceOpened.ValueBool() {
			response, err := s.client.OpenNetworkInterface(
				ctx,
				state.ID.ValueString(),
				dedicatedserver.NETWORKTYPEURL_PUBLIC,
			).Execute()
			if err != nil {
				summary := fmt.Sprintf(
					"Updating resource %s opening public network interface for id %q",
					s.name,
					state.ID.ValueString(),
				)
				utils.Error(ctx, &resp.Diagnostics, summary, err, response)
				return
			}
		} else {
			response, err := s.client.CloseNetworkInterface(
				ctx,
				state.ID.ValueString(),
				dedicatedserver.NETWORKTYPEURL_PUBLIC,
			).Execute()
			if err != nil {
				summary := fmt.Sprintf(
					"Updating resource %s closing public network interface for id %q",
					s.name,
					state.ID.ValueString(),
				)
				utils.Error(ctx, &resp.Diagnostics, summary, err, response)
				return
			}
		}
		state.PublicNetworkInterfaceOpened = plan.PublicNetworkInterfaceOpened
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (s *serverResource) Create(
	ctx context.Context,
	_ resource.CreateRequest,
	response *resource.CreateResponse,
) {
	utils.Error(
		ctx,
		&response.Diagnostics,
		fmt.Sprintf(
			"Resource %s can only be imported, not created.",
			s.name,
		),
		nil,
		nil,
	)
}

func (s *serverResource) Delete(
	_ context.Context,
	_ resource.DeleteRequest,
	_ *resource.DeleteResponse,
) {
}

func (s *serverResource) getServer(
	ctx context.Context,
	serverID string,
) (*serverResourceModel, *http.Response, error) {
	// Getting server info
	serverResult, serverResponse, err := s.client.GetServer(ctx, serverID).Execute()
	if err != nil {
		return nil, serverResponse, err
	}

	var publicIP string
	var publicIPNullRouted bool
	if networkInterfaces, ok := serverResult.GetNetworkInterfacesOk(); ok {
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
	if contract, ok := serverResult.GetContractOk(); ok {
		reference = contract.GetReference()
	}

	var internalMAC string
	if networkInterfaces, ok := serverResult.GetNetworkInterfacesOk(); ok {
		if internalNetworkInterface, ok := networkInterfaces.GetInternalOk(); ok {
			internalMAC = internalNetworkInterface.GetMac()
		}
	}

	var remoteManagementIP string
	if networkInterfaces, ok := serverResult.GetNetworkInterfacesOk(); ok {
		if remoteNetworkInterface, ok := networkInterfaces.GetRemoteManagementOk(); ok {
			remoteManagementIPPart := strings.Split(remoteNetworkInterface.GetIp(), "/")
			ip := net.ParseIP(remoteManagementIPPart[0])
			if ip != nil {
				remoteManagementIP = ip.String()
			}
		}
	}

	serverLocation := serverResult.GetLocation()
	l := locationResourceModel{
		Rack:  types.StringValue(serverLocation.GetRack()),
		Site:  types.StringValue(serverLocation.GetSite()),
		Suite: types.StringValue(serverLocation.GetSuite()),
		Unit:  types.StringValue(serverLocation.GetUnit()),
	}
	location, digs := types.ObjectValueFrom(ctx, l.attributeTypes(), l)

	if digs.HasError() {
		return nil, nil, fmt.Errorf("error reading dedicated server location with id: %q", serverID)
	}

	// Getting server power info
	powerResult, powerResponse, err := s.client.GetServerPowerStatus(
		ctx,
		serverID,
	).Execute()
	if err != nil {
		return nil, powerResponse, err
	}

	pdu := powerResult.GetPdu()
	ipmi := powerResult.GetIpmi()
	poweredOn := pdu.GetStatus() != "off" && ipmi.GetStatus() != "off"

	// Getting server public network interface info
	var publicNetworkOpened bool
	networkRequest := s.client.GetNetworkInterface(
		ctx,
		serverID,
		dedicatedserver.NETWORKTYPEURL_PUBLIC,
	)
	networkResult, networkResponse, err := networkRequest.Execute()
	if err != nil && networkResponse != nil && networkResponse.StatusCode != http.StatusNotFound {
		return nil, networkResponse, err
	} else {
		if networkResult != nil {
			if _, ok := networkResult.GetStatusOk(); ok {
				publicNetworkOpened = networkResult.GetStatus() == "open"
			}
		}
	}

	// Getting server DHCP info
	dhcpResult, dhcpResponse, err := s.client.GetServerDhcpReservationList(
		ctx,
		serverID,
	).Execute()
	if err != nil {
		return nil, dhcpResponse, err
	}
	var dhcpLease string
	if len(dhcpResult.GetLeases()) != 0 {
		leases := dhcpResult.GetLeases()
		dhcpLease = leases[0].GetBootfile()
	}

	// Getting server public IP info
	var reverseLookup string
	if publicIP != "" {
		ipResult, ipResponse, err := s.client.GetServerIp(
			ctx,
			serverID,
			publicIP,
		).Execute()
		if err != nil {
			return nil, ipResponse, err
		}
		reverseLookup = ipResult.GetReverseLookup()
	}

	dedicatedserverResource := serverResourceModel{
		ID:                           types.StringValue(serverResult.GetId()),
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
	}

	return &dedicatedserverResource, nil, nil
}
