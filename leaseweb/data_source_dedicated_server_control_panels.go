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
		Description: `
The ` + "`dedicated_server_control_panels`" + ` data source allows access to the list of
control panels available for installation on a dedicated server.
`,
		ReadContext: dataSourceDedicatedServerControlPanelsRead,
		Schema: map[string]*schema.Schema{
			"names": {
				Description: "List of the control panel names.",
				Type:        schema.TypeMap,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ids": {
				Description: "List of the control panel IDs.",
				Type:        schema.TypeSet,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"operating_system_id": {
				Description: "Filter the list of control panels to return only the ones available to an operating system.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func dataSourceDedicatedServerControlPanelsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	operatingSystemID := d.Get("operating_system_id").(string)
	controlPanels, err := getControlPanels(ctx, operatingSystemID)
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
