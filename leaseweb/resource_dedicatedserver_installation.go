package leaseweb

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDedicatedServerInstallation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDedicatedServerInstallationCreate,
		ReadContext:   resourceDedicatedServerInstallationRead,
		DeleteContext: resourceDedicatedServerInstallationDelete,
		Schema: map[string]*schema.Schema{
			"dedicated_server_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"operating_system_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceDedicatedServerInstallationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	serverID := d.Get("dedicated_server_id").(string)

	var payload = InstallationJobPayload{
		OperatingSystemID: d.Get("operating_system_id").(string),
	}

	installationJob, err := launchInstallationJob(serverID, &payload)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(installationJob.UUID)

	return resourceDedicatedServerInstallationRead(ctx, d, m)
}

func resourceDedicatedServerInstallationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	serverID := d.Get("dedicated_server_id").(string)

	var diags diag.Diagnostics

	installationJob, err := getLatestInstallationJob(serverID)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("operating_system_id", installationJob.Payload.OperatingSystemID)

	return diags
}

func resourceDedicatedServerInstallationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	d.SetId("")

	return diags
}
