package leaseweb

import (
	"context"
	"strings"
	"time"

	LSW "github.com/LeaseWeb/leaseweb-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceDedicatedServer() *schema.Resource {
	return &schema.Resource{
		Description: `
The ` + "`dedicated_server`" + ` resource manages several items linked to a dedicated server.
The resource cannot currently be created automatically, it needs to be imported first.
`,
		CreateContext: resourceDedicatedServerCreate,
		ReadContext:   resourceDedicatedServerRead,
		UpdateContext: resourceDedicatedServerUpdate,
		DeleteContext: resourceDedicatedServerDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the dedicated server.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"reference": {
				Description: "The reference of the dedicated server.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"reverse_lookup": {
				Description: "The reverse lookup associated with the dedicated server public IP.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"dhcp_lease": {
				Description:  "The URL of PXE boot the dedicated server is booting from.",
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IsURLWithScheme([]string{"http", "https"}),
			},
			"powered_on": {
				Description: "Whether the dedicated server is powered on or not.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"public_network_interface_opened": {
				Description: "Whether the public network interface of the dedicated server is opened or not.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"public_ip_null_routed": {
				Description: "Whether the public IP of the dedicated server is null routed or not.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"location": {
				Description: `
The location of the server.
Available fields are ` + "`rack`" + `, ` + "`site`" + `, ` + "`suite`" + ` and ` + "`unit`" + `.
`,
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"public_ip": {
				Description: "The public IP of the dedicated server.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"remote_management_ip": {
				Description: "The remote management IP of the dedicated server.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceDedicatedServerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// NOTE: we do not support creating resources at this moment
	return resourceDedicatedServerRead(ctx, d, m)
}

func resourceDedicatedServerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	serverID := d.Get("id").(string)

	var diags diag.Diagnostics

	// get basic data
	server, err := LSW.DedicatedServerApi{}.Get(ctx, serverID)
	if err != nil {
		logApiError(ctx, err)
		return diag.FromErr(err)
	}

	server.NetworkInterfaces.Public.Ip = strings.SplitN(server.NetworkInterfaces.Public.Ip, "/", 2)[0]
	server.NetworkInterfaces.RemoteManagement.Ip = strings.SplitN(server.NetworkInterfaces.RemoteManagement.Ip, "/", 2)[0]

	d.Set("reference", server.Contract.Reference)
	d.Set("public_ip", server.NetworkInterfaces.Public.Ip)
	d.Set("remote_management_ip", server.NetworkInterfaces.RemoteManagement.Ip)

	d.Set("location", map[string]string{
		"rack":  server.Location.Rack,
		"site":  server.Location.Site,
		"suite": server.Location.Suite,
		"unit":  server.Location.Unit,
	})

	// get IP data
	ip, err := LSW.DedicatedServerApi{}.GetIp(ctx, serverID, server.NetworkInterfaces.Public.Ip)
	if err != nil {
		logApiError(ctx, err)
		return diag.FromErr(err)
	}
	d.Set("reverse_lookup", ip.ReverseLookup)
	d.Set("public_ip_null_routed", ip.NullRouted)

	// get lease data
	lease, err := LSW.DedicatedServerApi{}.ListDhcpReservation(ctx, serverID, LSW.PaginationOptions{})
	if err != nil {
		logApiError(ctx, err)
		return diag.FromErr(err)
	}

	if len(lease.Leases) == 0 {
		d.Set("dhcp_lease", "")
	} else {
		d.Set("dhcp_lease", lease.Leases[0].BootFile)
	}

	// get power data
	powerStatus, err := LSW.DedicatedServerApi{}.GetPowerStatus(ctx, serverID)
	if err != nil {
		logApiError(ctx, err)
		return diag.FromErr(err)
	}

	d.Set("powered_on", powerStatus.Pdu.Status != "off" && powerStatus.Ipmi.Status != "off")

	// get public network interface data
	publicNetworkInterfaceInfo, err := LSW.DedicatedServerApi{}.GetNetworkInterface(ctx, serverID, "public")
	if err != nil {
		logApiError(ctx, err)
		return diag.FromErr(err)
	}
	d.Set("public_network_interface_opened", publicNetworkInterfaceInfo.Status == "OPEN")

	return diags
}

func resourceDedicatedServerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	serverID := d.Get("id").(string)

	if d.HasChange("reference") {
		reference := d.Get("reference").(string)
		err := LSW.DedicatedServerApi{}.Update(ctx, serverID, map[string]interface{}{"reference": reference})
		if err != nil {
			logApiError(ctx, err)
			return diag.FromErr(err)
		}

		// Wait a bit for the change to be made available in the contract before reading the resource again
		time.Sleep(5 * time.Second)
	}

	if d.HasChange("reverse_lookup") {
		publicIP := d.Get("public_ip").(string)
		reverseLookup := d.Get("reverse_lookup").(string)
		_, err := LSW.DedicatedServerApi{}.UpdateIp(ctx, serverID, publicIP, map[string]string{"reverseLookup": reverseLookup})
		if err != nil {
			logApiError(ctx, err)
			return diag.FromErr(err)
		}
	}

	if d.HasChange("dhcp_lease") {
		bootFile := d.Get("dhcp_lease").(string)
		if bootFile != "" {
			err := LSW.DedicatedServerApi{}.CreateDhcpReservation(ctx, serverID, map[string]string{"bootfile": bootFile})
			if err != nil {
				logApiError(ctx, err)
				return diag.FromErr(err)
			}
		} else {
			err := LSW.DedicatedServerApi{}.DeleteDhcpReservation(ctx, serverID)
			if err != nil {
				logApiError(ctx, err)
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("powered_on") {
		if d.Get("powered_on").(bool) {
			err := LSW.DedicatedServerApi{}.PowerOnServer(ctx, serverID)
			if err != nil {
				logApiError(ctx, err)
				return diag.FromErr(err)
			}
		} else {
			err := LSW.DedicatedServerApi{}.PowerOffServer(ctx, serverID)
			if err != nil {
				logApiError(ctx, err)
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("public_network_interface_opened") {
		if d.Get("public_network_interface_opened").(bool) {
			err := LSW.DedicatedServerApi{}.OpenNetworkInterface(ctx, serverID, "public")
			if err != nil {
				logApiError(ctx, err)
				return diag.FromErr(err)
			}
		} else {
			err := LSW.DedicatedServerApi{}.CloseNetworkInterface(ctx, serverID, "public")
			if err != nil {
				logApiError(ctx, err)
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("public_ip_null_routed") {
		publicIP := d.Get("public_ip").(string)
		if d.Get("public_ip_null_routed").(bool) {
			_, err := LSW.DedicatedServerApi{}.NullRouteAnIp(ctx, serverID, publicIP)
			if err != nil {
				logApiError(ctx, err)
				return diag.FromErr(err)
			}
		} else {
			_, err := LSW.DedicatedServerApi{}.NullRouteAnIp(ctx, serverID, publicIP)
			if err != nil {
				logApiError(ctx, err)
				return diag.FromErr(err)
			}
		}
	}

	return resourceDedicatedServerRead(ctx, d, m)
}

func resourceDedicatedServerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// NOTE: we do not support destroying resources at this moment
	d.SetId("")

	return diags
}
