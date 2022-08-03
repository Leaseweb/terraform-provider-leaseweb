package leaseweb

import (
	"context"
	"strings"
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
			"hostname": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"timezone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"ssh_keys": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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

	if d.Get("hostname") != "" {
		payload["hostname"] = d.Get("hostname").(string)
	}

	if d.Get("timezone") != "" {
		payload["timezone"] = d.Get("timezone").(string)
	}

	sshKeysIf := d.Get("ssh_keys").([]interface{})
	if len(sshKeysIf) != 0 {
		sshKeys := make([]string, len(sshKeysIf))
		for i, sshKey := range sshKeysIf {
			sshKeys[i] = sshKey.(string)
		}
		payload["sshKeys"] = strings.Join(sshKeys, "\n")
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

	if hostname, ok := installationJob.Payload["hostname"]; ok {
		d.Set("hostname", hostname)
	}

	if timezone, ok := installationJob.Payload["timezone"]; ok {
		d.Set("timezone", timezone)
	}

	if sshKeys, ok := installationJob.Payload["sshKeys"]; ok {
		d.Set("ssh_keys", strings.Split(sshKeys.(string), "\n"))
	}

	return diags
}

func resourceDedicatedServerInstallationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	d.SetId("")

	return diags
}
