package leaseweb

import (
	"context"
	"fmt"
	"strings"

	LSW "github.com/LeaseWeb/leaseweb-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceDedicatedServerNotificationSettingBandwidth() *schema.Resource {
	return &schema.Resource{
		Description: `
The ` + "`dedicated_server_notification_setting_bandwidth`" + ` resource manages a bandwidth
notification setting linked to a dedicated server.
`,
		CreateContext: resourceDedicatedServerNotificationSettingBandwidthCreate,
		ReadContext:   resourceDedicatedServerNotificationSettingBandwidthRead,
		UpdateContext: resourceDedicatedServerNotificationSettingBandwidthUpdate,
		DeleteContext: resourceDedicatedServerNotificationSettingBandwidthDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the notification setting.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"dedicated_server_id": {
				Description: "The ID of the dedicated server.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"frequency": {
				Description: `
The frequency of the notification.
Can be either ` + "`DAILY`" + `, ` + "`WEEKLY`" + `, or ` + "`MONTHLY`" + `.
`,
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"DAILY", "WEEKLY", "MONTHLY"}, false),
			},
			"threshold": {
				Description:  "The threshold of the notification.",
				Type:         schema.TypeFloat,
				Required:     true,
				ValidateFunc: validation.FloatAtLeast(0),
			},
			"unit": {
				Description: `
The unit of the notification.
Can be either ` + "`Mbps`" + ` or ` + "`Gbps`" + `.
`,
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"Mbps", "Gbps"}, false),
			},
		},
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				parts := strings.SplitN(d.Id(), ":", 2)

				if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
					return nil, fmt.Errorf("Invalid ID format (%s), expected dedicated_server_id:notification_setting_id", d.Id())
				}

				d.Set("dedicated_server_id", parts[0])
				d.SetId(parts[1])

				return []*schema.ResourceData{d}, nil
			},
		},
	}
}

func resourceDedicatedServerNotificationSettingBandwidthCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	serverID := d.Get("dedicated_server_id").(string)
	frequency := d.Get("frequency").(string)
	threshold := d.Get("threshold").(float64)
	unit := d.Get("unit").(string)

	createdNotificationSetting, err := LSW.DedicatedServerApi{}.CreateBandWidthNotificationSetting(ctx, serverID, frequency, threshold, unit)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdNotificationSetting.Id)

	return resourceDedicatedServerNotificationSettingBandwidthRead(ctx, d, m)
}

func resourceDedicatedServerNotificationSettingBandwidthRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	serverID := d.Get("dedicated_server_id").(string)
	notificationSettingID := d.Get("id").(string)

	var diags diag.Diagnostics

	notificationSetting, err := LSW.DedicatedServerApi{}.GetBandWidthNotificationSetting(ctx, serverID, notificationSettingID)
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

	threshold := d.Get("threshold").(float64)
	params := map[string]string{
		"frequency": d.Get("frequency").(string),
		"threshold": fmt.Sprintf("%.6g", threshold),
		"unit":      d.Get("unit").(string),
	}

	_, err := LSW.DedicatedServerApi{}.UpdateBandWidthNotificationSetting(ctx, serverID, notificationSettingID, params)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceDedicatedServerNotificationSettingBandwidthRead(ctx, d, m)
}

func resourceDedicatedServerNotificationSettingBandwidthDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	serverID := d.Get("dedicated_server_id").(string)
	notificationSettingID := d.Get("id").(string)

	err := LSW.DedicatedServerApi{}.DeleteBandWidthNotificationSetting(ctx, serverID, notificationSettingID)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
