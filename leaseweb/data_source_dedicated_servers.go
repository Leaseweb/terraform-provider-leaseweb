package leaseweb

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDedicatedServers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDedicatedServersRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"site": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceDedicatedServersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	site := d.Get("site").(string)
	servers, err := getServersBatch(0, 5, site)
	if err != nil {
		return diag.FromErr(err)
	}

	serverIds := make([]string, len(servers))

	for i, server := range servers {
		serverIds[i] = server.ID
	}

	if err := d.Set("ids", serverIds); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
