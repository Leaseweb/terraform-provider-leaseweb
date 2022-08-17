package leaseweb

import (
	"context"
	"encoding/base64"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
			"callback_url": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"control_panel_id": {
				Type:     schema.TypeString,
				Optional: true,
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
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"post_install_script": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"raid": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice([]string{"HW", "SW", "NONE"}, false),
						},
						"level": {
							Type:         schema.TypeInt,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.IntInSlice([]int{0, 1, 5, 10}),
						},
						"number_of_disks": {
							Type:         schema.TypeInt,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.IntAtLeast(1),
						},
					},
				},
			},
			"device": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"partition": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bootable": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"filesystem": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"mountpoint": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"primary": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"size": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
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
		"operatingSystemId":   d.Get("operating_system_id").(string),
		"doEmailNotification": false,
	}

	if d.Get("control_panel_id") != "" {
		payload["controlPanelId"] = d.Get("control_panel_id").(string)
	}

	if d.Get("timezone") != "" {
		payload["timezone"] = d.Get("timezone").(string)
	}

	raid := d.Get("raid").([]interface{})

	if len(raid) != 0 {
		raidDetails := raid[0].(map[string]interface{})
		var raidConfig = map[string]interface{}{
			"type": raidDetails["type"].(string),
		}

		if raidConfig["type"] != "NONE" {
			raidConfig["level"] = raidDetails["level"]
			if raidDetails["number_of_disks"].(int) != 0 {
				raidConfig["numberOfDisks"] = raidDetails["number_of_disks"]
			}
		}
		payload["raid"] = raidConfig
	}

	sshKeysSet := d.Get("ssh_keys").(*schema.Set)
	if sshKeysSet.Len() != 0 {
		sshKeys := make([]string, sshKeysSet.Len())
		for i, sshKey := range sshKeysSet.List() {
			sshKeys[i] = sshKey.(string)
		}
		payload["sshKeys"] = strings.Join(sshKeys, "\n")
	}

	if d.Get("post_install_script") != "" {
		payload["postInstallScript"] = base64.StdEncoding.EncodeToString([]byte(d.Get("post_install_script").(string)))
	}

	if d.Get("password") != "" {
		payload["password"] = d.Get("password").(string)
	}

	if d.Get("device") != "" {
		payload["device"] = d.Get("device").(string)
	}

	partitions := d.Get("partition").([]interface{})
	if len(partitions) != 0 {
		payload["partitions"] = partitions
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

	if controlPanelID, ok := installationJob.Payload["controlPanelId"]; ok {
		d.Set("control_panel_id", controlPanelID)
	}

	if callbackURL, ok := installationJob.Payload["callbackUrl"]; ok {
		d.Set("callback_url", callbackURL)
	}

	if hostname, ok := installationJob.Payload["hostname"]; ok {
		d.Set("hostname", hostname)
	}

	if timezone, ok := installationJob.Payload["timezone"]; ok {
		d.Set("timezone", timezone)
	}

	if sshKeys, ok := installationJob.Payload["sshKeys"]; ok {
		sshKeysList := strings.Split(sshKeys.(string), "\n")
		sshKeysIf := make([]interface{}, len(sshKeysList))
		for i, sshKey := range sshKeysList {
			sshKeysIf[i] = sshKey
		}
		d.Set("ssh_keys", schema.NewSet(schema.HashString, sshKeysIf))
	}

	if device, ok := installationJob.Payload["device"]; ok {
		d.Set("device", device)
	}

	if partitions, ok := installationJob.Payload["partitions"]; ok {
		d.Set("partition", partitions)
	}

	return diags
}

func resourceDedicatedServerInstallationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	d.SetId("")

	return diags
}
