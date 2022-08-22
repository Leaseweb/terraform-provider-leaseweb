package leaseweb

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceDedicatedServerNotificationSettingDatatraffic() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDedicatedServerNotificationSettingDatatrafficCreate,
		ReadContext:   resourceDedicatedServerNotificationSettingDatatrafficRead,
		UpdateContext: resourceDedicatedServerNotificationSettingDatatrafficUpdate,
		DeleteContext: resourceDedicatedServerNotificationSettingDatatrafficDelete,
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
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"DAILY", "WEEKLY", "MONTHLY"}, false),
			},
			"threshold": {
				Type:         schema.TypeFloat,
				Required:     true,
				ValidateFunc: validation.FloatAtLeast(0),
			},
			"unit": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"MB", "GB", "TB"}, false),
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

func resourceDedicatedServerNotificationSettingDatatrafficCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	serverID := d.Get("dedicated_server_id").(string)

	var notificationSetting = NotificationSetting{
		Frequency: d.Get("frequency").(string),
		Threshold: d.Get("threshold").(float64),
		Unit:      d.Get("unit").(string),
	}

	createdNotificationSetting, err := createDedicatedServerNotificationSetting(serverID, "datatraffic", &notificationSetting)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdNotificationSetting.ID)

	return resourceDedicatedServerNotificationSettingDatatrafficRead(ctx, d, m)
}

func resourceDedicatedServerNotificationSettingDatatrafficRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	serverID := d.Get("dedicated_server_id").(string)
	notificationSettingID := d.Get("id").(string)

	var diags diag.Diagnostics

	notificationSetting, err := getDedicatedServerNotificationSetting(serverID, "datatraffic", notificationSettingID)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("frequency", notificationSetting.Frequency)
	d.Set("threshold", notificationSetting.Threshold)
	d.Set("unit", notificationSetting.Unit)

	return diags
}

func resourceDedicatedServerNotificationSettingDatatrafficUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	serverID := d.Get("dedicated_server_id").(string)
	notificationSettingID := d.Get("id").(string)

	var notificationSetting = NotificationSetting{
		Frequency: d.Get("frequency").(string),
		Threshold: d.Get("threshold").(float64),
		Unit:      d.Get("unit").(string),
	}

	if _, err := updateDedicatedServerNotificationSetting(serverID, "datatraffic", notificationSettingID, &notificationSetting); err != nil {
		return diag.FromErr(err)
	}

	return resourceDedicatedServerNotificationSettingDatatrafficRead(ctx, d, m)
}

func resourceDedicatedServerNotificationSettingDatatrafficDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	serverID := d.Get("dedicated_server_id").(string)
	notificationSettingID := d.Get("id").(string)

	if err := deleteDedicatedServerNotificationSetting(serverID, "datatraffic", notificationSettingID); err != nil {
		return diag.FromErr(err)
	}

	return diags
}
