package leaseweb

import (
	"context"
	"strconv"
	"time"

	LSW "github.com/LeaseWeb/leaseweb-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDedicatedServers() *schema.Resource {
	return &schema.Resource{
		Description: `
The ` + "`dedicated_servers`" + ` data source allows access to the list of
dedicated servers available in your account.
`,
		ReadContext: dataSourceDedicatedServersRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Description: "List of the dedicated server IDs available to the account.",
				Type:        schema.TypeSet,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ip": {
				Description: "Filter the list of servers by ip address.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"mac_address": {
				Description: "Filter the list of servers by mac address.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"private_network_capable": {
				Description: "Filter the list for private network capable servers.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"private_network_enabled": {
				Description: "Filter the list for private network enabled servers.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"private_rack_id": {
				Description: "Filter the list of servers by dedicated rack id.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"reference": {
				Description: "Filter the list of servers by reference.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"site": {
				Description: "Filter the list of servers by location site.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func dataSourceDedicatedServersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var servers []LSW.DedicatedServer
	var opts LSW.DedicatedServerListOptions

	// For each possible DedicatedServerListOption we check if is defined and add it to the filter when needed
	if ip := d.Get("ip").(string); ip != "" {
		opts.IP = &ip
	}

	if macAddress := d.Get("mac_address").(string); macAddress != "" {
		opts.MacAddress = &macAddress
	}

	// The IsOkExists method will trigger a deprecated warning on runtime[1] hence we use the GetRawConfigmethod to determine if a boolean typed field is defined.
	// [1] https://discuss.hashicorp.com/t/terraform-sdk-usage-which-out-of-get-getok-getokexists-with-boolean/41815
	if !d.GetRawConfig().AsValueMap()["private_network_capable"].IsNull() {
		privateNetworkCapable := d.Get("private_network_capable").(bool)
		opts.PrivateNetworkCapable = &privateNetworkCapable
	}
	if !d.GetRawConfig().AsValueMap()["private_network_enabled"].IsNull() {
		privateNetworkEnabled := d.Get("private_network_enabled").(bool)
		opts.PrivateNetworkEnabled = &privateNetworkEnabled
	}

	if privateRackId := d.Get("private_rack_id").(string); privateRackId != "" {
		opts.PrivateRackID = &privateRackId
	}

	if reference := d.Get("reference").(string); reference != "" {
		opts.Reference = &reference
	}

	if site := d.Get("site").(string); site != "" {
		opts.Site = &site
	}

	servers, err := LSW.DedicatedServerApi{}.AllWithOpts(ctx, opts)
	if err != nil {
		return diag.FromErr(err)
	}

	serverIds := make([]string, len(servers))
	for i, server := range servers {
		serverIds[i] = server.Id
	}

	if err := d.Set("ids", serverIds); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
