package dedicatedserver

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider"

	terraformResource "github.com/hashicorp/terraform-plugin-framework/resource"
	dedicatedServerSdk "github.com/leaseweb/leaseweb-go-sdk/dedicatedServer"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
)

var (
	_ terraformResource.Resource                = &dedicatedServerResource{}
	_ terraformResource.ResourceWithConfigure   = &dedicatedServerResource{}
	_ terraformResource.ResourceWithImportState = &dedicatedServerResource{}
)

type dedicatedServerResource struct {
	// TODO: Refactor this part, apiKey shouldn't be here.
	apiKey string
	client provider.Client
}

func NewDedicatedServerResource() terraformResource.Resource {
	return &dedicatedServerResource{}
}

func (d *dedicatedServerResource) Metadata(
	_ context.Context,
	req terraformResource.MetadataRequest,
	resp *terraformResource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_dedicated_server"
}

func (d *dedicatedServerResource) authContext(ctx context.Context) context.Context {
	return context.WithValue(
		ctx,
		dedicatedServerSdk.ContextAPIKeys,
		map[string]dedicatedServerSdk.APIKey{
			"X-LSW-Auth": {Key: d.apiKey, Prefix: ""},
		},
	)
}

func (d *dedicatedServerResource) Configure(
	ctx context.Context,
	req terraformResource.ConfigureRequest,
	resp *terraformResource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}
	configuration := dedicatedServerSdk.NewConfiguration()

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

	apiClient := dedicatedServerSdk.NewAPIClient(configuration)
	d.client = provider.NewClient(configuration.Host)
}

func (d *dedicatedServerResource) ImportState(
	ctx context.Context,
	req terraformResource.ImportStateRequest,
	resp *terraformResource.ImportStateResponse,
) {

	terraformResource.ImportStatePassthroughID(
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

func (d *dedicatedServerResource) Create(
	ctx context.Context,
	req terraformResource.CreateRequest,
	resp *terraformResource.CreateResponse,
) {
	panic("unimplemented")
}

func (d *dedicatedServerResource) Delete(
	ctx context.Context,
	req terraformResource.DeleteRequest,
	resp *terraformResource.DeleteResponse,
) {
	panic("unimplemented")
}

func (d *dedicatedServerResource) getServer(
	ctx context.Context,
	serverID string,
) (*DedicatedServerModel, error) {

	// Getting server info
	serverResult, _, err := d.client.GetServer(d.authContext(ctx), serverID).Execute()
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
	l := locationModel{
		Rack:  types.StringValue(serverLocation.GetRack()),
		Site:  types.StringValue(serverLocation.GetSite()),
		Suite: types.StringValue(serverLocation.GetSuite()),
		Unit:  types.StringValue(serverLocation.GetUnit()),
	}
	location, digs := types.ObjectValueFrom(ctx, l.AttributeTypes(), l)
	if digs.HasError() {
		return nil, fmt.Errorf("error reading dedicated server locationModel with id: %q", serverID)
	}

	// Getting server power info
	powerResult, _, err := d.client.GetServerPowerStatus(d.authContext(ctx), serverID).Execute()
	if err != nil {
		return nil, fmt.Errorf("error reading dedicated server power status with id: %q - %s", serverID, err.Error())
	}
	pdu := powerResult.GetPdu()
	ipmi := powerResult.GetIpmi()
	poweredOn := pdu.GetStatus() != "off" && ipmi.GetStatus() != "off"

	// Getting server public network interface info
	var publicNetworkOpened bool
	networkRequest := d.client.GetNetworkInterface(d.authContext(ctx), serverID, dedicatedServerSdk.NETWORKTYPE_PUBLIC)
	networkResult, networkResponse, err := networkRequest.Execute()
	if err != nil && networkResponse.StatusCode != http.StatusNotFound {
		return nil, fmt.Errorf("error reading dedicated server network interface with id: %q - %s", serverID, err.Error())
	} else {
		if _, ok := networkResult.GetStatusOk(); ok {
			publicNetworkOpened = networkResult.GetStatus() == "open"
		}
	}

	// Getting server DHCP info
	dhcpResult, _, err := d.client.GetServerDhcpReservationList(d.authContext(ctx), serverID).Execute()
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
		ipResult, _, err := d.client.GetServerIp(d.authContext(ctx), serverID, publicIP).Execute()
		if err != nil {
			return nil, fmt.Errorf("error reading dedicated server IP details with id: %q - %s", serverID, err.Error())
		}
		reverseLookup = ipResult.GetReverseLookup()
	}

	dedicatedServer := DedicatedServerModel{
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

// Update @TODO Wrote the comments below before refactoring, so they don't make sense in the refactored context
// Update @TODO I have no idea what's happening here or why or why this is happening in a certain order. This should be in a service, preferably one that's properly unit tested. This is like forgoing services in php and just throwing all the logic into an action. Refactoring this is going to be a nightmare in a year, especially with multiple teams. We've now replace `overkill structure` with no structure. There are also going to be a lot of cases where code in create & update is duplicated.
// Update @TODO What happens if reference is an empty string here?
// Update @TODO Abstract away the auth part
// Update @TODO All these edge cases now also have to be tested with acceptance tests, which are way more resource intensive than unit tests.
// Update @TODO Because this isn't in a service this is also difficult to read. $isPublicIpExists could be set anywhere above.
func (d *dedicatedServerResource) Update(ctx context.Context, req terraformResource.UpdateRequest, resp *terraformResource.UpdateResponse) {
	var plan DedicatedServerModel
	planDiags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(planDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state DedicatedServerModel
	stateDiags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(stateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := d.client.DedicatedServer.Update(plan, &state, d.authContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(err.Summary, err.Message)
		return
	}

	stateDiags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(stateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
