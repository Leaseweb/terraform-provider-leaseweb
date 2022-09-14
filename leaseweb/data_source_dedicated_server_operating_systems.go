package leaseweb

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDedicatedServerOperatingSystems() *schema.Resource {
	return &schema.Resource{
		Description: `
The ` + "`dedicated_server_operating_systems`" + ` data source allows access to the list of
operating systems available for installation on a dedicated server.
`,
		ReadContext: dataSourceDedicatedServerOperatingSystemsRead,
		Schema: map[string]*schema.Schema{
			"names": {
				Description: "List of the operating system names.",
				Type:        schema.TypeMap,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ids": {
				Description: "List of the operating system IDs.",
				Type:        schema.TypeSet,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceDedicatedServerOperatingSystemsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	operatingSystems, err := getOperatingSystems()
	if err != nil {
		return diag.FromErr(err)
	}

	operatingSystemsNames := make(map[string]string)
	operatingSystemsIds := make([]string, len(operatingSystems))

	for i, os := range operatingSystems {
		operatingSystemsNames[os.ID] = os.Name
		operatingSystemsIds[i] = os.ID
	}

	if err := d.Set("names", operatingSystemsNames); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("ids", operatingSystemsIds); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
