package leaseweb

import (
	"context"
	"encoding/base64"
	"strings"
	"time"

	LSW "github.com/LeaseWeb/leaseweb-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceDedicatedServerInstallation() *schema.Resource {
	return &schema.Resource{
		Description: `
The ` + "`dedicated_server_installation`" + ` resource is used to define an installation to a dedicated server.
The resource cannot be updated in place, modifying any data will launch a new installation.
`,
		CreateContext: resourceDedicatedServerInstallationCreate,
		ReadContext:   resourceDedicatedServerInstallationRead,
		DeleteContext: resourceDedicatedServerInstallationDelete,
		Schema: map[string]*schema.Schema{
			"dedicated_server_id": {
				Description: "The ID of the dedicated server.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"operating_system_id": {
				Description: "The ID of the operating system to install.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"callback_url": {
				Description:  "The URL which will receive callbacks when the installation is finished or failed.",
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsURLWithScheme([]string{"http", "https"}),
			},
			"control_panel_id": {
				Description: "The ID of the control panel to install.",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"job_uuid": {
				Description: "The UUID of the installation job.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"hostname": {
				Description: "The hostname to configure on the dedicated server.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
			},
			"timezone": {
				Description: "The timezone to configure on the dedicated server.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
			},
			"ssh_keys": {
				Description: "List of public SSH keys to authorize on the dedicated server.",
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"post_install_script": {
				Description: "Script to run right after the installation.",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"password": {
				Description: "The root password to configure on the dedicated server.",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"raid": {
				Description: "The RAID configuration to use on the dedicated server.",
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				ForceNew:    true,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Description: `
The RAID type to apply.
Valid types are ` + "`HW`" + `, ` + "`SW`" + ` and ` + "`NONE`" + `(pass-through).
`,
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice([]string{"HW", "SW", "NONE"}, false),
						},
						"level": {
							Description: `
The RAID level to apply (only valid with HW and SW types).
Valid levels are ` + "`0`" + `, ` + "`1`" + `, ` + "`5`" + ` and ` + "`10`" + `.
`,
							Type:         schema.TypeInt,
							Optional:     true,
							ForceNew:     true,
							Computed:     true,
							ValidateFunc: validation.IntInSlice([]int{0, 1, 5, 10}),
						},
						"number_of_disks": {
							Description: `
The number of disks to apply RAID on (only valid with HW and SW types).
All disks are used if unspecified.`,
							Type:         schema.TypeInt,
							Optional:     true,
							ForceNew:     true,
							Computed:     true,
							ValidateFunc: validation.IntAtLeast(1),
						},
					},
				},
			},
			"device": {
				Description: `
Block devices in a disk set in which the partitions will be installed.
Supported values are any disk set id, ` + "`SATA_SAS`" + ` or ` + "`NVME`" + `.
`,
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"partition": {
				Description: "The partition configuration to use on the dedicated server.",
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"filesystem": {
							Description: "Filesystem of the partition.",
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
						},
						"mountpoint": {
							Description: `
Mountpoint of the partition.
Mandatory for root partition, unnecessary for swap partition.
`,
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"size": {
							Description: "Size of the partition (Normally in MB, but this is OS-specific).",
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
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

	var payload = map[string]interface{}{
		"operatingSystemId":   d.Get("operating_system_id").(string),
		"doEmailNotification": false,
	}

	if d.Get("control_panel_id") != "" {
		payload["controlPanelId"] = d.Get("control_panel_id").(string)
	}

	if d.Get("callback_url") != "" {
		payload["callbackUrl"] = d.Get("callback_url").(string)
	}

	if d.Get("hostname") != "" {
		payload["hostname"] = d.Get("hostname").(string)
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

	installationJob, err := LSW.DedicatedServerApi{}.LaunchInstallation(ctx, serverID, payload)
	if err != nil {
		logAPIError(ctx, err)
		return diag.FromErr(err)
	}

	d.Set("job_uuid", installationJob.Uuid)
	d.SetId(serverID)

	createStateConf := &resource.StateChangeConf{
		Pending: []string{"ACTIVE"},
		Target:  []string{"FINISHED"},
		Refresh: func() (interface{}, string, error) {
			job, err := LSW.DedicatedServerApi{}.GetJob(ctx, serverID, installationJob.Uuid)
			if err != nil {
				logAPIError(ctx, err)
				return nil, "error", err
			}
			return job, job.Status, err
		},
		Timeout:      d.Timeout(schema.TimeoutCreate) - time.Minute,
		PollInterval: 30 * time.Second,
	}
	_, err = createStateConf.WaitForStateContext(ctx)

	if err != nil {
		logAPIError(ctx, err)
		return diag.FromErr(err)
	}
	return resourceDedicatedServerInstallationRead(ctx, d, m)
}

func resourceDedicatedServerInstallationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	serverID := d.Get("dedicated_server_id").(string)

	var diags diag.Diagnostics

	opts := LSW.DedicatedServerListJobOptions{
		PaginationOptions: LSW.PaginationOptions{
			Offset: LSW.Int(0),
			Limit:  LSW.Int(1),
		},
		Type: LSW.String("install"),
	}

	installationJobs, err := LSW.DedicatedServerApi{}.ListJobs(ctx, serverID, opts)
	if err != nil {
		logAPIError(ctx, err)
		return diag.FromErr(err)
	}
	if len(installationJobs.Jobs) == 0 {
		return diag.Errorf("no installation jobs found for server %s", serverID)
	}
	installationJob := installationJobs.Jobs[0]

	d.Set("job_uuid", installationJob.Uuid)

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

	if raid, ok := installationJob.Payload["raid"]; ok {
		raid := raid.(map[string]interface{})
		var raidConfig = map[string]interface{}{
			"type": raid["type"].(string),
		}

		if raidConfig["type"] != "NONE" {
			raidConfig["level"] = raid["level"]
			if numberOfDisks, ok := raid["numberOfDisks"]; ok {
				raidConfig["number_of_disks"] = numberOfDisks
			}
		}
		d.Set("raid", []interface{}{
			raidConfig,
		})
	}

	return diags
}

func resourceDedicatedServerInstallationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	d.SetId("")

	return diags
}
