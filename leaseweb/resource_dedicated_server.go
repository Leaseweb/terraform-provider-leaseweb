package leaseweb

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDedicatedServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDedicatedServerCreate,
		ReadContext:   resourceDedicatedServerRead,
		UpdateContext: resourceDedicatedServerUpdate,
		DeleteContext: resourceDedicatedServerDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"reference": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"reverse_lookup": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"dhcp_lease": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"powered_on": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"public_network_interface_opened": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"public_ip_null_routed": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"location": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"remote_management_ip": {
				Type:     schema.TypeString,
				Computed: true,
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

	// 1) get basic data from /v2/servers/{id}
	server, err := getServer(serverID)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("reference", server.Contract.Reference)
	d.Set("public_ip", server.NetworkInterfaces.Public.IP)
	d.Set("remote_management_ip", server.NetworkInterfaces.RemoteManagement.IP)

	d.Set("location", map[string]string{
		"rack":  server.Location.Rack,
		"site":  server.Location.Site,
		"suite": server.Location.Suite,
		"unit":  server.Location.Unit,
	})

	// 2) get reverse lookup from /v2/servers/{id}/ips/{ip}
	ip, err := getServerIP(serverID, server.NetworkInterfaces.Public.IP)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("reverse_lookup", ip.ReverseLookup)
	d.Set("public_ip_null_routed", ip.NullRouted)

	// 3) get leases info from /v2/servers/{id}/leases
	lease, err := getServerLease(serverID)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("dhcp_lease", lease.GetBootfile())

	// 4) get power info from /v2/servers/{id}/powerInfo
	powerInfo, err := getPowerInfo(serverID)
	d.Set("powered_on", powerInfo.IsPoweredOn())

	// 5) get public network interface info from /v2/servers/{serverId}/networkInterfaces/{networkType}
	publicNetworkInterfaceInfo, err := getNetworkInterfaceInfo(serverID, "public")
	d.Set("public_network_interface_opened", publicNetworkInterfaceInfo.IsOpened())

	return diags
}

func resourceDedicatedServerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	serverID := d.Get("id").(string)

	if d.HasChange("reference") {
		reference := d.Get("reference").(string)
		if err := updateReference(serverID, reference); err != nil {
			return diag.FromErr(err)
		}
		d.Set("reference", reference)
	}

	if d.HasChange("reverse_lookup") {
		publicIP := d.Get("public_ip").(string)
		reverseLookup := d.Get("reverse_lookup").(string)
		if err := updateReverseLookup(serverID, publicIP, reverseLookup); err != nil {
			return diag.FromErr(err)
		}
		d.Set("reverse_lookup", reverseLookup)
	}

	if d.HasChange("dhcp_lease") {
		bootFile := d.Get("dhcp_lease").(string)
		if bootFile != "" {
			if err := addDHCPLease(serverID, bootFile); err != nil {
				return diag.FromErr(err)
			}
			d.Set("dhcp_lease", bootFile)
		} else {
			if err := removeDHCPLease(serverID); err != nil {
				return diag.FromErr(err)
			}
			d.Set("dhcp_lease", "")
		}
	}

	if d.HasChange("powered_on") {
		if d.Get("powered_on").(bool) {
			if err := powerOnServer(serverID); err != nil {
				return diag.FromErr(err)
			}
			d.Set("powered_on", true)
		} else {
			if err := powerOffServer(serverID); err != nil {
				return diag.FromErr(err)
			}
			d.Set("powered_on", false)
		}
	}

	if d.HasChange("public_network_interface_opened") {
		if d.Get("public_network_interface_opened").(bool) {
			if err := openNetworkInterface(serverID, "public"); err != nil {
				return diag.FromErr(err)
			}
			d.Set("public_network_interface_opened", true)
		} else {
			if err := closeNetworkInterface(serverID, "public"); err != nil {
				return diag.FromErr(err)
			}
			d.Set("public_network_interface_opened", false)
		}
	}

	if d.HasChange("public_ip_null_routed") {
		publicIP := d.Get("public_ip").(string)
		if d.Get("public_ip_null_routed").(bool) {
			if err := nullIP(serverID, publicIP); err != nil {
				return diag.FromErr(err)
			}
			d.Set("public_ip_null_routed", true)
		} else {
			if err := unnullIP(serverID, publicIP); err != nil {
				return diag.FromErr(err)
			}
			d.Set("public_ip_null_routed", false)
		}
	}

	var diags diag.Diagnostics

	return diags

	// NOTE: return resourceDedicatedServerRead(ctx, d, m)
}

func resourceDedicatedServerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// NOTE: we do not support destroying resources at this moment
	d.SetId("")

	return diags
}
