package resource

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/leaseweb/leaseweb-go-sdk/dedicatedServer"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/customerror"
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
	PublicIPNullRouted           types.Bool   `tfsdk:"public_ip_null_routed"`
	PublicIP                     types.String `tfsdk:"public_ip"`
	RemoteManagementIP           types.String `tfsdk:"remote_management_ip"`
	InternalMAC                  types.String `tfsdk:"internal_mac"`
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

func New() resource.Resource {
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

func (d *dedicatedServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data dedicatedServerResourceData
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	dedicatedServer, err := d.getServer(ctx, data.ID.ValueString())
	if err != nil {
		summary := "Reading dedicated server"
		resp.Diagnostics.AddError(summary, customerror.NewError(nil, err).Error())
		tflog.Error(ctx, fmt.Sprintf("%s %s", summary, customerror.NewError(nil, err).Error()))
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
		summary := "Importing dedicated server"
		resp.Diagnostics.AddError(summary, customerror.NewError(nil, err).Error())
		tflog.Error(ctx, fmt.Sprintf("%s %s", summary, customerror.NewError(nil, err).Error()))
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
		response, err := d.client.UpdateServerReference(d.authContext(ctx), state.ID.ValueString()).UpdateServerReferenceOpts(*ropts).Execute()
		if err != nil {
			summary := fmt.Sprintf("Error updating dedicated server reference with id: %q", plan.ID.ValueString())
			resp.Diagnostics.AddError(summary, customerror.NewError(response, err).Error())
			tflog.Error(ctx, fmt.Sprintf("%s %s", summary, customerror.NewError(response, err).Error()))
			return
		}
		state.Reference = plan.Reference
	}

	// Updating Power status
	if !plan.PoweredOn.IsNull() && !plan.PoweredOn.IsUnknown() {
		if plan.PoweredOn.ValueBool() {
			request := d.client.PowerServerOn(d.authContext(ctx), state.ID.ValueString())
			response, err := request.Execute()
			if err != nil {
				summary := fmt.Sprintf("Error powering on for dedicated server: %q", state.ID.ValueString())
				resp.Diagnostics.AddError(summary, customerror.NewError(response, err).Error())
				tflog.Error(ctx, fmt.Sprintf("%s %s", summary, customerror.NewError(response, err).Error()))
				return
			}
		} else {
			request := d.client.PowerServerOff(d.authContext(ctx), state.ID.ValueString())
			response, err := request.Execute()
			if err != nil {
				summary := fmt.Sprintf("Error powering off for dedicated server: %q", state.ID.ValueString())
				resp.Diagnostics.AddError(summary, customerror.NewError(response, err).Error())
				tflog.Error(ctx, fmt.Sprintf("%s %s", summary, customerror.NewError(response, err).Error()))
				return
			}
		}
		state.PoweredOn = plan.PoweredOn
	}

	// Updateing Reverse Lookup
	isPublicIPExists := !state.PublicIP.IsNull() && !state.PublicIP.IsUnknown() && state.PublicIP.ValueString() != ""
	if !plan.ReverseLookup.IsNull() && !plan.ReverseLookup.IsUnknown() && isPublicIPExists {
		iopts := dedicatedServer.NewUpdateIpProfileOpts()
		iopts.ReverseLookup = plan.ReverseLookup.ValueStringPointer()
		_, response, err := d.client.UpdateIpProfile(d.authContext(ctx), state.ID.ValueString(), state.PublicIP.ValueString()).UpdateIpProfileOpts(*iopts).Execute()
		if err != nil {
			summary := fmt.Sprintf("Error updating dedicated server reverse lookup with id: %q", state.ID.ValueString())
			resp.Diagnostics.AddError(summary, customerror.NewError(response, err).Error())
			tflog.Error(ctx, fmt.Sprintf("%s %s", summary, customerror.NewError(response, err).Error()))
			return
		}
		state.ReverseLookup = plan.ReverseLookup
	}

	// Updating an IP null routing
	if !plan.PublicIPNullRouted.IsNull() && !plan.PublicIPNullRouted.IsUnknown() && plan.PublicIPNullRouted != state.PublicIPNullRouted && isPublicIPExists {
		if plan.PublicIPNullRouted.ValueBool() {
			_, response, err := d.client.NullIpRoute(d.authContext(ctx), state.ID.ValueString(), state.PublicIP.ValueString()).Execute()
			if err != nil {
				summary := fmt.Sprintf("Error null routing an IP for dedicated server: %q and IP: %q", state.ID.ValueString(), state.PublicIP.ValueString())
				resp.Diagnostics.AddError(summary, customerror.NewError(response, err).Error())
				tflog.Error(ctx, fmt.Sprintf("%s %s", summary, customerror.NewError(response, err).Error()))
				return
			}
		} else {
			_, response, err := d.client.RemoveNullIpRoute(d.authContext(ctx), state.ID.ValueString(), state.PublicIP.ValueString()).Execute()
			if err != nil {
				summary := fmt.Sprintf("Error remove null routing an IP for dedicated server: %q and IP: %q", state.ID.ValueString(), state.PublicIP.ValueString())
				resp.Diagnostics.AddError(summary, customerror.NewError(response, err).Error())
				tflog.Error(ctx, fmt.Sprintf("%s %s", summary, customerror.NewError(response, err).Error()))
				return
			}
		}
		state.PublicIPNullRouted = plan.PublicIPNullRouted
	}

	// Updating dhcp lease
	if !plan.DHCPLease.IsNull() && !plan.DHCPLease.IsUnknown() {
		if plan.DHCPLease.ValueString() != "" {
			opts := dedicatedServer.NewCreateServerDhcpReservationOpts(plan.DHCPLease.ValueString())
			response, err := d.client.CreateServerDhcpReservation(d.authContext(ctx), state.ID.ValueString()).CreateServerDhcpReservationOpts(*opts).Execute()
			if err != nil {
				summary := fmt.Sprintf("Error creating a DHCP reservation for dedicated server: %q", state.ID.ValueString())
				resp.Diagnostics.AddError(summary, customerror.NewError(response, err).Error())
				tflog.Error(ctx, fmt.Sprintf("%s %s", summary, customerror.NewError(response, err).Error()))
				return
			}
		} else {
			response, err := d.client.DeleteServerDhcpReservation(d.authContext(ctx), state.ID.ValueString()).Execute()
			if err != nil {
				summary := fmt.Sprintf("Error deleting DHCP reservation for dedicated server: %q", state.ID.ValueString())
				resp.Diagnostics.AddError(summary, customerror.NewError(response, err).Error())
				tflog.Error(ctx, fmt.Sprintf("%s %s", summary, customerror.NewError(response, err).Error()))
				return
			}
		}
		state.DHCPLease = plan.DHCPLease
	}

	// Updating network interface status
	if !plan.PublicIPNullRouted.IsNull() && !plan.PublicIPNullRouted.IsUnknown() && plan.PublicNetworkInterfaceOpened != state.PublicNetworkInterfaceOpened {
		if plan.PublicNetworkInterfaceOpened.ValueBool() {
			response, err := d.client.OpenNetworkInterface(d.authContext(ctx), state.ID.ValueString(), dedicatedServer.NETWORKTYPE_PUBLIC).Execute()
			if err != nil {
				summary := fmt.Sprintf("Error opening public network interface for dedicated server: %q", state.ID.ValueString())
				resp.Diagnostics.AddError(summary, customerror.NewError(response, err).Error())
				tflog.Error(ctx, fmt.Sprintf("%s %s", summary, customerror.NewError(response, err).Error()))
				return
			}
		} else {
			response, err := d.client.CloseNetworkInterface(d.authContext(ctx), state.ID.ValueString(), dedicatedServer.NETWORKTYPE_PUBLIC).Execute()
			if err != nil {
				summary := fmt.Sprintf("Error closing public network interface for dedicated server: %q", state.ID.ValueString())
				resp.Diagnostics.AddError(summary, customerror.NewError(response, err).Error())
				tflog.Error(ctx, fmt.Sprintf("%s %s", summary, customerror.NewError(response, err).Error()))
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
