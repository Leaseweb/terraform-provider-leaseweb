package leaseweb

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDedicatedServerControlPanels() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDedicatedServerControlPanelsRead,
		Schema: map[string]*schema.Schema{
			"names": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceDedicatedServerControlPanelsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	controlPanels, err := getControlPanels()
	if err != nil {
		return diag.FromErr(err)
	}

	controlPanelsNames := make(map[string]string)
	controlPanelsIds := make([]string, len(controlPanels))

	for i, cp := range controlPanels {
		controlPanelsNames[cp.ID] = cp.Name
		controlPanelsIds[i] = cp.ID
	}

	if err := d.Set("names", controlPanelsNames); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("ids", controlPanelsIds); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
