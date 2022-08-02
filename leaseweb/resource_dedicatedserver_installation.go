package leaseweb

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
			"job_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				d.Set("dedicated_server_id", d.Id())

				return []*schema.ResourceData{d}, nil
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
		},
	}
}

func resourceDedicatedServerInstallationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	serverID := d.Get("dedicated_server_id").(string)

	var payload = Payload{
		"operatingSystemId": d.Get("operating_system_id").(string),
	}

	installationJob, err := launchInstallationJob(serverID, &payload)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("job_uuid", installationJob.UUID)
	d.SetId(serverID)

	createStateConf := &resource.StateChangeConf{
		Pending: []string{"ACTIVE"},
		Target:  []string{"FINISHED"},
		Refresh: func() (interface{}, string, error) {
			job, err := getJob(serverID, installationJob.UUID)
			if err != nil {
				return nil, "error", err
			}
			return job, job.Status, err
		},
		Timeout:      d.Timeout(schema.TimeoutCreate) - time.Minute,
		Delay:        5 * time.Minute,
		PollInterval: 30 * time.Second,
	}
	_, err = createStateConf.WaitForStateContext(ctx)

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceDedicatedServerInstallationRead(ctx, d, m)
}

func resourceDedicatedServerInstallationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	serverID := d.Get("dedicated_server_id").(string)

	var diags diag.Diagnostics

	installationJob, err := getLatestInstallationJob(serverID)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("job_uuid", installationJob.UUID)
	d.Set("operating_system_id", installationJob.Payload["operatingSystemId"])

	return diags
}

func resourceDedicatedServerInstallationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	d.SetId("")

	return diags
}
