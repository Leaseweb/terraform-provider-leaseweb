package leaseweb

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceOperatingSystems() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOperatingSystemsRead,
		Schema: map[string]*schema.Schema{
			"names": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceOperatingSystemsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	operatingSystems, err := getOperatingSystems()
	if err != nil {
		return diag.FromErr(err)
	}

	operatingSystemsNames := make(map[string]string)

	for _, os := range operatingSystems {
		operatingSystemsNames[os.ID] = os.Name
	}

	if err := d.Set("names", operatingSystemsNames); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
