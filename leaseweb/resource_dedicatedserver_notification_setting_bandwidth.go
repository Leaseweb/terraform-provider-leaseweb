package leaseweb

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDedicatedServerNotificationSettingBandwidth() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDedicatedServerNotificationSettingBandwidthCreate,
		ReadContext:   resourceDedicatedServerNotificationSettingBandwidthRead,
		UpdateContext: resourceDedicatedServerNotificationSettingBandwidthUpdate,
		DeleteContext: resourceDedicatedServerNotificationSettingBandwidthDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dedicated_server_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"frequency": {
				Type:     schema.TypeString,
				Required: true,
			},
			"threshold": {
				Type:     schema.TypeString,
				Required: true,
			},
			"unit": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceDedicatedServerNotificationSettingBandwidthCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	serverID := d.Get("dedicated_server_id").(string)

	var notificationSetting = NotificationSetting{
		Frequency: d.Get("frequency").(string),
		Threshold: d.Get("threshold").(string),
		Unit:      d.Get("unit").(string),
	}

	createdNotificationSetting, err := createDedicatedServerNotificationSettingBandwidth(serverID, &notificationSetting)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdNotificationSetting.ID)

	return resourceDedicatedServerNotificationSettingBandwidthRead(ctx, d, m)
}

func resourceDedicatedServerNotificationSettingBandwidthRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	serverID := d.Get("dedicated_server_id").(string)
	notificationSettingID := d.Get("id").(string)

	var diags diag.Diagnostics

	notificationSetting, err := getDedicatedServerNotificationSettingBandwidth(serverID, notificationSettingID)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("frequency", notificationSetting.Frequency)
	d.Set("threshold", notificationSetting.Threshold)
	d.Set("unit", notificationSetting.Unit)

	return diags
}

func resourceDedicatedServerNotificationSettingBandwidthUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	serverID := d.Get("dedicated_server_id").(string)
	notificationSettingID := d.Get("id").(string)

	if d.HasChange("frequency") || d.HasChange("threshold") || d.HasChange("unit") {
		var notificationSetting = NotificationSetting{
			Frequency: d.Get("frequency").(string),
			Threshold: d.Get("threshold").(string),
			Unit:      d.Get("unit").(string),
		}

		if _, err := updateDedicatedServerNotificationSettingBandwidth(serverID, notificationSettingID, &notificationSetting); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDedicatedServerNotificationSettingBandwidthRead(ctx, d, m)
}

func resourceDedicatedServerNotificationSettingBandwidthDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	serverID := d.Get("dedicated_server_id").(string)
	notificationSettingID := d.Get("id").(string)

	if err := deleteDedicatedServerNotificationSettingBandwidth(serverID, notificationSettingID); err != nil {
		return diag.FromErr(err)
	}

	return diags
}
