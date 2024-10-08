package dedicatedserverresource

import (
	"context"
	"fmt"

	tfResource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/leaseweb/leaseweb-go-sdk/dedicatedServer"
)

func (d *resource) Update(ctx context.Context, req tfResource.UpdateRequest, resp *tfResource.UpdateResponse) {
	var plan resourceData
	planDiags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(planDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state resourceData
	stateDiags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(stateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Updating reference
	if !plan.Reference.IsNull() && !plan.Reference.IsUnknown() {
		ropts := dedicatedServer.NewUpdateServerReferenceOpts(plan.Reference.ValueString())
		_, err := d.Client.UpdateServerReference(d.AuthContext(ctx), state.ID.ValueString()).UpdateServerReferenceOpts(*ropts).Execute()
		if err != nil {
			summary := fmt.Sprintf("Error updating dedicated server reference with id: %q", plan.ID.ValueString())
			resp.Diagnostics.AddError(summary, err.Error())
			tflog.Error(ctx, fmt.Sprintf("%s %s", summary, err.Error()))
			return
		}
		state.Reference = plan.Reference
	}

	// Updating Power status
	if !plan.PoweredOn.IsNull() && !plan.PoweredOn.IsUnknown() {
		if plan.PoweredOn.ValueBool() {
			request := d.Client.PowerServerOn(d.AuthContext(ctx), state.ID.ValueString())
			_, err := request.Execute()
			if err != nil {
				summary := fmt.Sprintf("Error powering on for dedicated server: %q", state.ID.ValueString())
				resp.Diagnostics.AddError(summary, err.Error())
				tflog.Error(ctx, fmt.Sprintf("%s %s", summary, err.Error()))
				return
			}
		} else {
			request := d.Client.PowerServerOff(d.AuthContext(ctx), state.ID.ValueString())
			_, err := request.Execute()
			if err != nil {
				summary := fmt.Sprintf("Error powering off for dedicated server: %q", state.ID.ValueString())
				resp.Diagnostics.AddError(summary, err.Error())
				tflog.Error(ctx, fmt.Sprintf("%s %s", summary, err.Error()))
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
		_, _, err := d.Client.UpdateIpProfile(d.AuthContext(ctx), state.ID.ValueString(), state.PublicIP.ValueString()).UpdateIpProfileOpts(*iopts).Execute()
		if err != nil {
			summary := fmt.Sprintf("Error updating dedicated server reverse lookup with id: %q", state.ID.ValueString())
			resp.Diagnostics.AddError(summary, err.Error())
			tflog.Error(ctx, fmt.Sprintf("%s %s", summary, err.Error()))
			return
		}
		state.ReverseLookup = plan.ReverseLookup
	}

	// Updating an IP null routing
	if !plan.PublicIPNullRouted.IsNull() && !plan.PublicIPNullRouted.IsUnknown() && plan.PublicIPNullRouted != state.PublicIPNullRouted && isPublicIPExists {
		if plan.PublicIPNullRouted.ValueBool() {
			_, _, err := d.Client.NullIpRoute(d.AuthContext(ctx), state.ID.ValueString(), state.PublicIP.ValueString()).Execute()
			if err != nil {
				summary := fmt.Sprintf("Error null routing an IP for dedicated server: %q and IP: %q", state.ID.ValueString(), state.PublicIP.ValueString())
				resp.Diagnostics.AddError(summary, err.Error())
				tflog.Error(ctx, fmt.Sprintf("%s %s", summary, err.Error()))
				return
			}
		} else {
			_, _, err := d.Client.RemoveNullIpRoute(d.AuthContext(ctx), state.ID.ValueString(), state.PublicIP.ValueString()).Execute()
			if err != nil {
				summary := fmt.Sprintf("Error remove null routing an IP for dedicated server: %q and IP: %q", state.ID.ValueString(), state.PublicIP.ValueString())
				resp.Diagnostics.AddError(summary, err.Error())
				tflog.Error(ctx, fmt.Sprintf("%s %s", summary, err.Error()))
				return
			}
		}
		state.PublicIPNullRouted = plan.PublicIPNullRouted
	}

	// Updating dhcp lease
	if !plan.DHCPLease.IsNull() && !plan.DHCPLease.IsUnknown() {
		if plan.DHCPLease.ValueString() != "" {
			opts := dedicatedServer.NewCreateServerDhcpReservationOpts(plan.DHCPLease.ValueString())
			_, err := d.Client.CreateServerDhcpReservation(d.AuthContext(ctx), state.ID.ValueString()).CreateServerDhcpReservationOpts(*opts).Execute()
			if err != nil {
				summary := fmt.Sprintf("Error creating a DHCP reservation for dedicated server: %q", state.ID.ValueString())
				resp.Diagnostics.AddError(summary, err.Error())
				tflog.Error(ctx, fmt.Sprintf("%s %s", summary, err.Error()))
				return
			}
		} else {
			_, err := d.Client.DeleteServerDhcpReservation(d.AuthContext(ctx), state.ID.ValueString()).Execute()
			if err != nil {
				summary := fmt.Sprintf("Error deleting DHCP reservation for dedicated server: %q", state.ID.ValueString())
				resp.Diagnostics.AddError(summary, err.Error())
				tflog.Error(ctx, fmt.Sprintf("%s %s", summary, err.Error()))
				return
			}
		}
		state.DHCPLease = plan.DHCPLease
	}

	// Updating network interface status
	if !plan.PublicIPNullRouted.IsNull() && !plan.PublicIPNullRouted.IsUnknown() && plan.PublicNetworkInterfaceOpened != state.PublicNetworkInterfaceOpened {
		if plan.PublicNetworkInterfaceOpened.ValueBool() {
			_, err := d.Client.OpenNetworkInterface(d.AuthContext(ctx), state.ID.ValueString(), dedicatedServer.NETWORKTYPE_PUBLIC).Execute()
			if err != nil {
				summary := fmt.Sprintf("Error opening public network interface for dedicated server: %q", state.ID.ValueString())
				resp.Diagnostics.AddError(summary, err.Error())
				tflog.Error(ctx, fmt.Sprintf("%s %s", summary, err.Error()))
				return
			}
		} else {
			_, err := d.Client.CloseNetworkInterface(d.AuthContext(ctx), state.ID.ValueString(), dedicatedServer.NETWORKTYPE_PUBLIC).Execute()
			if err != nil {
				summary := fmt.Sprintf("Error closing public network interface for dedicated server: %q", state.ID.ValueString())
				resp.Diagnostics.AddError(summary, err.Error())
				tflog.Error(ctx, fmt.Sprintf("%s %s", summary, err.Error()))
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
