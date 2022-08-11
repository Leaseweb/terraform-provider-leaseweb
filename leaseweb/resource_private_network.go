package leaseweb

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourcePrivateNetwork() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourcePrivateNetworkRead,
		UpdateContext: resourcePrivateNetworkUpdate,
		DeleteContext: resourcePrivateNetworkDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "default",
			},
			"dhcp": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"ENABLED", "DISABLED"}, false),
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourcePrivateNetworkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	privateNetworkID := d.Get("id").(string)

	var diags diag.Diagnostics

	privateNetwork, err := getPrivateNetwork(privateNetworkID)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("name", privateNetwork.Name)
	d.Set("dhcp", privateNetwork.Dhcp)

	return diags
}

func resourcePrivateNetworkUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	privateNetworkID := d.Get("id").(string)

	var privateNetwork = PrivateNetwork{
		Name: d.Get("name").(string),
		Dhcp: d.Get("dhcp").(string),
	}

	if _, err := updatePrivateNetwork(privateNetworkID, &privateNetwork); err != nil {
		return diag.FromErr(err)
	}

	return resourcePrivateNetworkRead(ctx, d, m)
}

func resourcePrivateNetworkDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	privateNetworkID := d.Get("id").(string)

	if err := deletePrivateNetwork(privateNetworkID); err != nil {
		return diag.FromErr(err)
	}

	return diags
}
