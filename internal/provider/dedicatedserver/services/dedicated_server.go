package services

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	dedicatedServerSdk "github.com/leaseweb/leaseweb-go-sdk/dedicatedServer"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/dedicatedserver"
)

type DedicatedServerError struct {
	Summary string
	Message string
}

func (d DedicatedServerError) Error() string {
	return d.Message
}

func NewDedicatedServerError(summary string, message string) *DedicatedServerError {
	return &DedicatedServerError{
		Summary: summary,
		Message: message,
	}
}

// DedicatedServer has a single responsibility
type DedicatedServer struct {
	sdkUpdateReference          func(ctx context.Context, serverId string) dedicatedServerSdk.ApiUpdateServerReferenceRequest
	sdkServerPowerOn            func(ctx context.Context, id string) dedicatedServerSdk.ApiPowerServerOnRequest
	sdkServerPowerOff           func(ctx context.Context, id string) dedicatedServerSdk.ApiPowerServerOffRequest
	sdkUpdateIpProfile          func(ctx context.Context, id string, publicIp string) dedicatedServerSdk.ApiUpdateIpProfileRequest
	nullIpRoute                 func(ctx context.Context, id string, publicIp string) dedicatedServerSdk.ApiNullIpRouteRequest
	removeNullIpRoute           func(ctx context.Context, id string, publicIp string) dedicatedServerSdk.ApiRemoveNullIpRouteRequest
	createServerDhcpReservation func(ctx context.Context, id string) dedicatedServerSdk.ApiCreateServerDhcpReservationRequest
	deleteServerDhcpReservation func(ctx context.Context, id string) dedicatedServerSdk.ApiDeleteServerDhcpReservationRequest
	openNetworkInterface        func(ctx context.Context, id string, networkType dedicatedServerSdk.NetworkType) dedicatedServerSdk.ApiOpenNetworkInterfaceRequest
	closeNetworkInterface       func(ctx context.Context, id string, networkType dedicatedServerSdk.NetworkType) dedicatedServerSdk.ApiCloseNetworkInterfaceRequest
}

func (d DedicatedServer) updateReference(
	state *dedicatedserver.DedicatedServerModel,
	reference types.String,
	ctx context.Context,
) *DedicatedServerError {
	opts := dedicatedServerSdk.NewUpdateServerReferenceOpts(reference.ValueString())
	_, err := d.sdkUpdateReference(ctx, state.ID.ValueString()).UpdateServerReferenceOpts(*opts).Execute()
	if err != nil {
		return NewDedicatedServerError(
			fmt.Sprintf(
				"Error updating dedicated server reference with id: %q",
				state.ID.ValueString(),
			),
			err.Error(),
		)
	}

	state.Reference = reference
	return nil
}

func (d DedicatedServer) powerOn(id string, ctx context.Context) error {
	request := d.sdkServerPowerOn(ctx, id)
	_, err := request.Execute()
	if err != nil {
		return err
	}

	return nil
}

func (d DedicatedServer) powerOff(id string, ctx context.Context) error {
	request := d.sdkServerPowerOff(ctx, id)
	_, err := request.Execute()
	if err != nil {
		return err
	}

	return nil
}

func (d DedicatedServer) updatePower(
	state *dedicatedserver.DedicatedServerModel,
	poweredOn types.Bool,
	ctx context.Context,
) *DedicatedServerError {
	if poweredOn.ValueBool() {
		err := d.powerOn(state.ID.ValueString(), ctx)
		if err != nil {
			return NewDedicatedServerError(
				fmt.Sprintf(
					"Error seting power on for dedicated server: %q",
					state.ID.ValueString(),
				),
				err.Error(),
			)
		}

		state.PoweredOn = poweredOn
		return nil
	}

	err := d.powerOff(state.ID.ValueString(), ctx)
	if err != nil {
		return NewDedicatedServerError(
			fmt.Sprintf(
				"Error seting power off for dedicated server: %q",
				state.ID.ValueString(),
			),
			err.Error(),
		)
	}

	state.PoweredOn = poweredOn
	return nil
}

func (d DedicatedServer) updateIpProfile(
	state *dedicatedserver.DedicatedServerModel,
	reverseLookup types.String,
	ctx context.Context,
) *DedicatedServerError {
	opts := dedicatedServerSdk.NewUpdateIpProfileOpts()
	opts.ReverseLookup = reverseLookup.ValueStringPointer()
	_, _, err := d.sdkUpdateIpProfile(
		ctx,
		state.ID.ValueString(),
		state.PublicIP.ValueString(),
	).UpdateIpProfileOpts(*opts).Execute()
	if err != nil {
		return NewDedicatedServerError(
			fmt.Sprintf(
				"Error updating ip profile for dedicated server: %q",
				state.ID.ValueString(),
			),
			err.Error(),
		)
	}

	state.ReverseLookup = reverseLookup
	return nil
}

func (d DedicatedServer) updateIpNullRouting(
	state *dedicatedserver.DedicatedServerModel,
	isPublicIpNullRouted types.Bool,
	publicIp string,
	ctx context.Context,
) *DedicatedServerError {
	if isPublicIpNullRouted.ValueBool() {
		_, _, err := d.nullIpRoute(ctx, state.ID.ValueString(), publicIp).Execute()
		if err != nil {
			return NewDedicatedServerError(
				fmt.Sprintf(
					"Error setting ip null routing for dedicated server: %q",
					state.ID.ValueString(),
				),
				err.Error(),
			)
		}

		state.PublicIPNullRouted = isPublicIpNullRouted
		return nil
	}

	_, _, err := d.removeNullIpRoute(ctx, state.ID.ValueString(), publicIp).Execute()
	if err != nil {
		return NewDedicatedServerError(
			fmt.Sprintf(
				"Error removing ip null routing for dedicated server: %q",
				state.ID.ValueString(),
			),
			err.Error(),
		)
	}

	state.PublicIPNullRouted = isPublicIpNullRouted
	return nil
}

func (d DedicatedServer) updateDhcpLease(
	state *dedicatedserver.DedicatedServerModel,
	dhcpLease types.String,
	ctx context.Context,
) *DedicatedServerError {
	if dhcpLease.ValueString() != "" {
		opts := dedicatedServerSdk.NewCreateServerDhcpReservationOpts(dhcpLease.ValueString())
		_, err := d.createServerDhcpReservation(
			ctx,
			state.ID.ValueString(),
		).CreateServerDhcpReservationOpts(*opts).Execute()
		if err != nil {
			return NewDedicatedServerError(
				fmt.Sprintf(
					"Error creating dhcp lease reservervation for dedicated server: %q",
					state.ID.ValueString(),
				),
				err.Error(),
			)
		}

		state.DHCPLease = dhcpLease
		return nil
	}

	_, err := d.deleteServerDhcpReservation(ctx, state.ID.ValueString()).Execute()
	if err != nil {
		return NewDedicatedServerError(
			fmt.Sprintf(
				"Error deleting dhcp lease reservervation for dedicated server: %q",
				state.ID.ValueString(),
			),
			err.Error(),
		)
	}

	state.DHCPLease = dhcpLease
	return nil
}

func (d DedicatedServer) updateNetworkInterfaceStatus(
	state *dedicatedserver.DedicatedServerModel,
	publicNetworkInterfaceOpened types.Bool,
	ctx context.Context,
) *DedicatedServerError {
	if publicNetworkInterfaceOpened.ValueBool() {
		_, err := d.openNetworkInterface(
			ctx,
			state.ID.ValueString(),
			dedicatedServerSdk.NETWORKTYPE_PUBLIC,
		).Execute()
		if err != nil {
			return NewDedicatedServerError(
				fmt.Sprintf(
					"Error opening network interface for dedicated server: %q",
					state.ID.ValueString(),
				),
				err.Error(),
			)
		}
	}

	_, err := d.closeNetworkInterface(
		ctx,
		state.ID.ValueString(),
		dedicatedServerSdk.NETWORKTYPE_PUBLIC,
	).Execute()
	if err != nil {
		return NewDedicatedServerError(
			fmt.Sprintf(
				"Error closing network interface for dedicated server: %q",
				state.ID.ValueString(),
			),
			err.Error(),
		)
	}

	state.PublicNetworkInterfaceOpened = publicNetworkInterfaceOpened
	return nil
}

// Update returns the new state if successful, otherwise an error.
func (d DedicatedServer) Update(
	plan dedicatedserver.DedicatedServerModel,
	state *dedicatedserver.DedicatedServerModel,
	ctx context.Context,
) *DedicatedServerError {
	if !plan.Reference.IsNull() && !plan.Reference.IsUnknown() {
		err := d.updateReference(state, plan.Reference, ctx)
		if err != nil {
			return err
		}
	}

	if !plan.PoweredOn.IsNull() && !plan.PoweredOn.IsUnknown() {
		err := d.updatePower(state, plan.PoweredOn, ctx)
		if err != nil {
			return err
		}
	}

	if !plan.ReverseLookup.IsNull() && !plan.ReverseLookup.IsUnknown() && d.doesPublicIpExist(state.PublicIP) {
		err := d.updateIpProfile(
			state,
			plan.ReverseLookup,
			ctx,
		)
		if err != nil {
			return err
		}
	}

	if !plan.PublicIPNullRouted.IsNull() && !plan.PublicIPNullRouted.IsUnknown() && plan.PublicIPNullRouted != state.PublicIPNullRouted && d.doesPublicIpExist(state.PublicIP) {
		err := d.updateIpNullRouting(
			state,
			plan.PublicIPNullRouted,
			state.PublicIP.ValueString(),
			ctx,
		)
		if err != nil {
			return err
		}
	}

	if !plan.DHCPLease.IsNull() && !plan.DHCPLease.IsUnknown() {
		err := d.updateDhcpLease(state, plan.DHCPLease, ctx)
		if err != nil {
			return err

		}
	}

	if !plan.PublicIPNullRouted.IsNull() && !plan.PublicIPNullRouted.IsUnknown() && plan.PublicNetworkInterfaceOpened != state.PublicNetworkInterfaceOpened {
		err := d.updateNetworkInterfaceStatus(
			state,
			plan.PublicNetworkInterfaceOpened,
			ctx,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d DedicatedServer) doesPublicIpExist(publicIp types.String) bool {
	return !publicIp.IsNull() && !publicIp.IsUnknown() && publicIp.ValueString() != ""
}

func NewDedicatedServer(
	sdkUpdateReference func(
		ctx context.Context,
		serverId string,
	) dedicatedServerSdk.ApiUpdateServerReferenceRequest,
	sdkPowerOn func(
		ctx context.Context,
		id string,
	) dedicatedServerSdk.ApiPowerServerOnRequest,
	sdkPowerOff func(
		ctx context.Context,
		id string,
	) dedicatedServerSdk.ApiPowerServerOffRequest,
	sdkUpdateIpProfile func(
		ctx context.Context,
		id string,
		publicIp string,
	) dedicatedServerSdk.ApiUpdateIpProfileRequest,
) DedicatedServer {
	return DedicatedServer{
		sdkUpdateReference: sdkUpdateReference,
		sdkServerPowerOn:   sdkPowerOn,
		sdkServerPowerOff:  sdkPowerOff,
		sdkUpdateIpProfile: sdkUpdateIpProfile,
	}
}
