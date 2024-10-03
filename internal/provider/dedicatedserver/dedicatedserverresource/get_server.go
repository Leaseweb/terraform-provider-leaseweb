package dedicatedserverresource

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/dedicatedServer"
)

func (d *resource) getServer(ctx context.Context, serverID string) (*resourceData, error) {

	// Getting server info
	serverResult, _, err := d.Client.GetServer(d.AuthContext(ctx), serverID).Execute()
	if err != nil {
		return nil, fmt.Errorf("error reading dedicated server with id: %q - %s", serverID, err.Error())
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
	l := locationResourceData{
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
	powerResult, _, err := d.Client.GetServerPowerStatus(d.AuthContext(ctx), serverID).Execute()
	if err != nil {
		return nil, fmt.Errorf("error reading dedicated server power status with id: %q - %s", serverID, err.Error())
	}
	pdu := powerResult.GetPdu()
	ipmi := powerResult.GetIpmi()
	poweredOn := pdu.GetStatus() != "off" && ipmi.GetStatus() != "off"

	// Getting server public network interface info
	var publicNetworkOpened bool
	networkRequest := d.Client.GetNetworkInterface(d.AuthContext(ctx), serverID, dedicatedServer.NETWORKTYPE_PUBLIC)
	networkResult, networkResponse, err := networkRequest.Execute()
	if err != nil && networkResponse.StatusCode != http.StatusNotFound {
		return nil, fmt.Errorf("error reading dedicated server network interface with id: %q - %s", serverID, err.Error())
	} else {
		if _, ok := networkResult.GetStatusOk(); ok {
			publicNetworkOpened = networkResult.GetStatus() == "open"
		}
	}

	// Getting server DHCP info
	dhcpResult, _, err := d.Client.GetServerDhcpReservationList(d.AuthContext(ctx), serverID).Execute()
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
	if publicIP != "" {
		ipResult, _, err := d.Client.GetServerIp(d.AuthContext(ctx), serverID, publicIP).Execute()
		if err != nil {
			return nil, fmt.Errorf("error reading dedicated server IP details with id: %q - %s", serverID, err.Error())
		}
		reverseLookup = ipResult.GetReverseLookup()
	}

	dedicatedServer := resourceData{
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

	return &dedicatedServer, nil
}
