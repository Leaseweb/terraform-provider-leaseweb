package provider

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
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
	ID                           types.String `tfsdk:"id"`
	Reference                    types.String `tfsdk:"reference"`
	ReverseLookup                types.String `tfsdk:"reverse_lookup"`
	DHCPLease                    types.String `tfsdk:"dhcp_lease"`
	PoweredOn                    types.Bool   `tfsdk:"powered_on"`
	PublicNetworkInterfaceOpened types.Bool   `tfsdk:"public_network_interface_opened"`
	PublicIpNullRouted           types.Bool   `tfsdk:"public_ip_null_routed"`
	PublicIp                     types.String `tfsdk:"public_ip"`
	RemoteManagementIp           types.String `tfsdk:"remote_management_ip"`
	InternalMac                  types.String `tfsdk:"internal_mac"`
	Location                     types.Object `tfsdk:"location"`
}

type dedicatedServerLocationResourceData struct {
	Rack  types.String `tfsdk:"rack"`
	Site  types.String `tfsdk:"site"`
	Suite types.String `tfsdk:"suite"`
	Unit  types.String `tfsdk:"unit"`
}

func (l dedicatedServerLocationResourceData) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"rack":  types.StringType,
		"site":  types.StringType,
		"suite": types.StringType,
		"unit":  types.StringType,
	}
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

func (d *dedicatedServerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
			},
			"reverse_lookup": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The reverse lookup associated with the dedicated server public IP.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"dhcp_lease": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The URL of PXE boot the dedicated server is booting from.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"powered_on": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether the dedicated server is powered on or not.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"public_network_interface_opened": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether the public network interface of the dedicated server is opened or not.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"public_ip_null_routed": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether the public IP of the dedicated server is null routed or not.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
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
}

func (d *dedicatedServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data dedicatedServerResourceData
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	dedicatedServer, err := d.getServer(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Reading dedicated server", err.Error())
		return
	}

	diags = resp.State.Set(ctx, &dedicatedServer)
	resp.Diagnostics.Append(diags...)
}

func (d *dedicatedServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	resource.ImportStatePassthroughID(
		ctx,
		path.Root("id"),
		req,
		resp,
	)

	dedicatedServer, err := d.getServer(ctx, req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Importing dedicated server", err.Error())
		return
	}

	diags := resp.State.Set(ctx, dedicatedServer)
	resp.Diagnostics.Append(diags...)
}

func (d *dedicatedServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan dedicatedServerResourceData
	planDiags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(planDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state dedicatedServerResourceData
	stateDiags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(stateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Updating reference
	if !plan.Reference.IsNull() && !plan.Reference.IsUnknown() {
		ropts := dedicatedServer.NewUpdateServerReferenceOpts(plan.Reference.ValueString())
		referenceRequest := d.client.UpdateServerReference(d.authContext(ctx), state.ID.ValueString()).UpdateServerReferenceOpts(*ropts)
		_, err := referenceRequest.Execute()
		if err != nil {
			summary := fmt.Sprintf("Error updating dedicated server reference with id: %q", plan.ID.ValueString())
			resp.Diagnostics.AddError(summary, err.Error())
			return
		}
		state.Reference = plan.Reference
	}

	// Updating Power status
	if !plan.PoweredOn.IsNull() && !plan.PoweredOn.IsUnknown() {
		if plan.PoweredOn.ValueBool() {
			request := d.client.PowerServerOn(d.authContext(ctx), state.ID.ValueString())
			_, err := request.Execute()
			if err != nil {
				summary := fmt.Sprintf("Error powering on for dedicated server: %q", state.ID.ValueString())
				resp.Diagnostics.AddError(summary, err.Error())
				return
			}
		} else {
			request := d.client.PowerServerOff(d.authContext(ctx), state.ID.ValueString())
			_, err := request.Execute()
			if err != nil {
				summary := fmt.Sprintf("Error powering off for dedicated server: %q", state.ID.ValueString())
				resp.Diagnostics.AddError(summary, err.Error())
				return
			}
		}
		state.PoweredOn = plan.PoweredOn
	}

	// Updateing Reverse Lookup
	if !plan.ReverseLookup.IsNull() && !plan.ReverseLookup.IsUnknown() {
		iopts := dedicatedServer.NewUpdateIpProfileOpts()
		iopts.ReverseLookup = plan.ReverseLookup.ValueStringPointer()
		ipRequest := d.client.UpdateIpProfile(d.authContext(ctx), state.ID.ValueString(), state.PublicIp.ValueString()).UpdateIpProfileOpts(*iopts)
		_, _, err := ipRequest.Execute()
		if err != nil {
			summary := fmt.Sprintf("Error updating dedicated server reverse lookup with id: %q", state.ID.ValueString())
			resp.Diagnostics.AddError(summary, err.Error())
			return
		}
		state.ReverseLookup = plan.ReverseLookup
	}

	// Updating an IP null routing
	if !plan.PublicIpNullRouted.IsNull() && !plan.PublicIpNullRouted.IsUnknown() {
		if plan.PublicIpNullRouted.ValueBool() {
			nullRoutedRequest := d.client.NullIpRoute(d.authContext(ctx), state.ID.ValueString(), state.PublicIp.ValueString())
			_, _, err := nullRoutedRequest.Execute()
			if err != nil {
				summary := fmt.Sprintf("Error null routing an IP for dedicated server: %q", state.ID.ValueString())
				resp.Diagnostics.AddError(summary, err.Error())
				return
			}
		} else {
			nullRoutedRequest := d.client.RemoveNullIpRoute(d.authContext(ctx), state.ID.ValueString(), state.PublicIp.ValueString())
			_, _, err := nullRoutedRequest.Execute()
			if err != nil {
				summary := fmt.Sprintf("Error remove null routing an IP for dedicated server: %q", state.ID.ValueString())
				resp.Diagnostics.AddError(summary, err.Error())
				return
			}
		}
		state.PublicIpNullRouted = plan.PublicIpNullRouted
	}

	// Updating dhcp lease
	if !plan.DHCPLease.IsNull() && !plan.DHCPLease.IsUnknown() {
		if plan.DHCPLease.ValueString() != "" {
			opts := dedicatedServer.NewCreateServerDhcpReservationOpts(plan.DHCPLease.ValueString())
			dhcpRequest := d.client.CreateServerDhcpReservation(d.authContext(ctx), state.ID.ValueString()).CreateServerDhcpReservationOpts(*opts)
			_, err := dhcpRequest.Execute()
			if err != nil {
				summary := fmt.Sprintf("Error creating a DHCP reservation for dedicated server: %q", state.ID.ValueString())
				resp.Diagnostics.AddError(summary, err.Error())
				return
			}
		} else {
			dhcpRequest := d.client.DeleteServerDhcpReservation(d.authContext(ctx), state.ID.ValueString())
			_, err := dhcpRequest.Execute()
			if err != nil {
				summary := fmt.Sprintf("Error deleting DHCP reservation for dedicated server: %q", state.ID.ValueString())
				resp.Diagnostics.AddError(summary, err.Error())
				return
			}
		}
		state.DHCPLease = plan.DHCPLease
	}

	// Updating network interface status
	if !plan.PublicNetworkInterfaceOpened.IsNull() && !plan.PublicNetworkInterfaceOpened.IsUnknown() {
		if plan.PublicNetworkInterfaceOpened.ValueBool() {
			request := d.client.OpenNetworkInterface(d.authContext(ctx), state.ID.ValueString(), dedicatedServer.NETWORKTYPE_PUBLIC)
			_, err := request.Execute()
			if err != nil {
				summary := fmt.Sprintf("Error opening public network interface for dedicated server: %q", state.ID.ValueString())
				resp.Diagnostics.AddError(summary, err.Error())
				return
			}
		} else {
			request := d.client.CloseNetworkInterface(d.authContext(ctx), state.ID.ValueString(), dedicatedServer.NETWORKTYPE_PUBLIC)
			_, err := request.Execute()
			if err != nil {
				summary := fmt.Sprintf("Error closing public network interface for dedicated server: %q", state.ID.ValueString())
				resp.Diagnostics.AddError(summary, err.Error())
				return
			}
		}
		state.PublicNetworkInterfaceOpened = plan.PublicNetworkInterfaceOpened
	}

	stateDiags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(stateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *dedicatedServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	panic("unimplemented")
}

func (d *dedicatedServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	panic("unimplemented")
}

func (d *dedicatedServerResource) getServer(ctx context.Context, serverID string) (*dedicatedServerResourceData, error) {

	// Getting server info
	serverRequest := d.client.GetServer(d.authContext(ctx), serverID)
	serverResult, _, err := serverRequest.Execute()
	if err != nil {
		return nil, fmt.Errorf("error reading dedicated server with id: %q - %s", serverID, err.Error())
	}

	var publicIp string
	var publicIpNullRouted *bool
	if networkInterfaces, ok := serverResult.GetNetworkInterfacesOk(); ok {
		if publicNetworkInterface, ok := networkInterfaces.GetPublicOk(); ok {
			publicIpPart := strings.Split(publicNetworkInterface.GetIp(), "/")
			ip := net.ParseIP(publicIpPart[0])
			if ip != nil {
				publicIp = ip.String()
			}
			publicIpNullRouted, _ = publicNetworkInterface.GetNullRoutedOk()
		}
	}

	var reference string
	if contract, ok := serverResult.GetContractOk(); ok {
		reference = contract.GetReference()
	}

	var internalMac string
	if networkInterfaces, ok := serverResult.GetNetworkInterfacesOk(); ok {
		if internalNetworkInterface, ok := networkInterfaces.GetInternalOk(); ok {
			internalMac = internalNetworkInterface.GetMac()
		}
	}

	var remoteManagementIp string
	if networkInterfaces, ok := serverResult.GetNetworkInterfacesOk(); ok {
		if remoteNetworkInterface, ok := networkInterfaces.GetRemoteManagementOk(); ok {
			remoteManagementIpPart := strings.Split(remoteNetworkInterface.GetIp(), "/")
			ip := net.ParseIP(remoteManagementIpPart[0])
			if ip != nil {
				remoteManagementIp = ip.String()
			}
		}
	}

	serverLocation := serverResult.GetLocation()
	l := dedicatedServerLocationResourceData{
		Rack:  types.StringValue(serverLocation.GetRack()),
		Site:  types.StringValue(serverLocation.GetSite()),
		Suite: types.StringValue(serverLocation.GetSuite()),
		Unit:  types.StringValue(serverLocation.GetUnit()),
	}
	location, digs := types.ObjectValueFrom(ctx, l.AttributeTypes(), l)
	if digs.HasError() {
		return nil, fmt.Errorf("error reading dedicated server location with id: %q", serverID)
	}

	// Getting server power info
	powerRequest := d.client.GetServerPowerStatus(d.authContext(ctx), serverID)
	powerResult, _, err := powerRequest.Execute()
	if err != nil {
		return nil, fmt.Errorf("error reading dedicated server power status with id: %q - %s", serverID, err.Error())
	}
	pdu := powerResult.GetPdu()
	ipmi := powerResult.GetIpmi()
	poweredOn := pdu.GetStatus() != "off" && ipmi.GetStatus() != "off"

	// Getting server public network interface info
	var publicNetworkOpened *bool
	networkRequest := d.client.GetNetworkInterface(d.authContext(ctx), serverID, dedicatedServer.NETWORKTYPE_PUBLIC)
	networkResult, networkResponse, err := networkRequest.Execute()
	if err != nil && networkResponse.StatusCode != http.StatusNotFound {
		return nil, fmt.Errorf("error reading dedicated server network interface with id: %q - %s", serverID, err.Error())
	} else {
		if _, ok := networkResult.GetStatusOk(); ok {
			isStatusOpen := networkResult.GetStatus() == "open"
			publicNetworkOpened = &isStatusOpen
		}
	}

	// Getting server DHCP info
	dhcpRequest := d.client.GetServerDhcpReservationList(d.authContext(ctx), serverID)
	dhcpResult, _, err := dhcpRequest.Execute()
	if err != nil {
		return nil, fmt.Errorf("error reading dedicated server DHCP with id: %q - %s", serverID, err.Error())
	}
	var dhcpLease string
	if len(dhcpResult.GetLeases()) != 0 {
		leases := dhcpResult.GetLeases()
		dhcpLease = leases[0].GetBootfile()
	}

	// Getting server public IP info
	var reverseLookup string
	if publicIp != "" {
		ipRequest := d.client.GetServerIp(d.authContext(ctx), serverID, publicIp)
		ipResult, _, err := ipRequest.Execute()
		if err != nil {
			return nil, fmt.Errorf("error reading dedicated server IP details with id: %q - %s", serverID, err.Error())
		}
		reverseLookup = ipResult.GetReverseLookup()
	}

	dedicatedServer := dedicatedServerResourceData{
		ID:                           types.StringValue(serverResult.GetId()),
		Reference:                    types.StringValue(reference),
		ReverseLookup:                types.StringValue(reverseLookup),
		DHCPLease:                    types.StringValue(dhcpLease),
		PoweredOn:                    types.BoolValue(poweredOn),
		PublicNetworkInterfaceOpened: types.BoolPointerValue(publicNetworkOpened),
		PublicIpNullRouted:           types.BoolPointerValue(publicIpNullRouted),
		PublicIp:                     types.StringValue(publicIp),
		RemoteManagementIp:           types.StringValue(remoteManagementIp),
		InternalMac:                  types.StringValue(internalMac),
		Location:                     location,
	}

	return &dedicatedServer, nil
}
