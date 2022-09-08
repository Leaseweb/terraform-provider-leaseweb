package leaseweb

import (
	"context"
	"time"

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

	// get basic data
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

	// get IP data
	ip, err := getServerIP(serverID, server.NetworkInterfaces.Public.IP)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("reverse_lookup", ip.ReverseLookup)
	d.Set("public_ip_null_routed", ip.NullRouted)

	// get lease data
	lease, err := getServerLease(serverID)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("dhcp_lease", lease.GetBootfile())

	// get power data
	powerInfo, err := getPowerInfo(serverID)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("powered_on", powerInfo.IsPoweredOn())

	// get public network interface data
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

		// Wait a bit for the change to be made available in the contract before reading the resource again
		time.Sleep(5 * time.Second)
	}

	if d.HasChange("reverse_lookup") {
		publicIP := d.Get("public_ip").(string)
		reverseLookup := d.Get("reverse_lookup").(string)
		if err := updateReverseLookup(serverID, publicIP, reverseLookup); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("dhcp_lease") {
		bootFile := d.Get("dhcp_lease").(string)
		if bootFile != "" {
			if err := addDHCPLease(serverID, bootFile); err != nil {
				return diag.FromErr(err)
			}
		} else {
			if err := removeDHCPLease(serverID); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("powered_on") {
		if d.Get("powered_on").(bool) {
			if err := powerOnServer(serverID); err != nil {
				return diag.FromErr(err)
			}
		} else {
			if err := powerOffServer(serverID); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("public_network_interface_opened") {
		if d.Get("public_network_interface_opened").(bool) {
			if err := openNetworkInterface(serverID, "public"); err != nil {
				return diag.FromErr(err)
			}
		} else {
			if err := closeNetworkInterface(serverID, "public"); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("public_ip_null_routed") {
		publicIP := d.Get("public_ip").(string)
		if d.Get("public_ip_null_routed").(bool) {
			if err := nullIP(serverID, publicIP); err != nil {
				return diag.FromErr(err)
			}
		} else {
			if err := unnullIP(serverID, publicIP); err != nil {
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
