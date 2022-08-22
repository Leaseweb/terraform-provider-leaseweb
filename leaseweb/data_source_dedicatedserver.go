package leaseweb

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDedicatedServer() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDedicatedServerRead,
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

func dataSourceDedicatedServerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	site := d.Get("site").(string)
	serverList, err := getServerList(site)
	if err != nil {
		return diag.FromErr(err)
	}

	serverListIds := make([]string, len(serverList))

	for i, server := range serverList {
		serverListIds[i] = server.ID
	}

	if err := d.Set("ids", serverListIds); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}