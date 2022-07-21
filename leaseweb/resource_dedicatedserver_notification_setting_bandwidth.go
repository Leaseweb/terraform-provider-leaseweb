package leaseweb

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-cty/cty"
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
				ValidateDiagFunc: func(v interface{}, p cty.Path) diag.Diagnostics {
					value := v.(string)
					var diags diag.Diagnostics
					if value != "DAILY" && value != "WEEKLY" && value != "MONTHLY" {
						diag := diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Incorrect attribute value",
							Detail:   `Inappropriate value for attribute "frequency": valid values are "DAILY", "WEEKLY", "MONTHLY"`,
						}
						diags = append(diags, diag)
					}
					return diags
				},
			},
			"threshold": {
				Type:     schema.TypeFloat,
				Required: true,
				ValidateDiagFunc: func(v interface{}, p cty.Path) diag.Diagnostics {
					value := v.(float64)
					var diags diag.Diagnostics
					if value <= 0 {
						diag := diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Incorrect attribute value",
							Detail:   `Inappropriate value for attribute "threshold": strictly positive float required`,
						}
						diags = append(diags, diag)
					}
					return diags
				},
			},
			"unit": {
				Type:     schema.TypeString,
				Required: true,
				ValidateDiagFunc: func(v interface{}, p cty.Path) diag.Diagnostics {
					value := v.(string)
					var diags diag.Diagnostics
					if value != "Gbps" && value != "Mbps" {
						diag := diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Incorrect attribute value",
							Detail:   `Inappropriate value for attribute "unit": valid values are "Gbps", "Mbps"`,
						}
						diags = append(diags, diag)
					}
					return diags
				},
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

	var notificationSetting = NotificationSetting{
		Frequency: d.Get("frequency").(string),
		Threshold: d.Get("threshold").(float64),
		Unit:      d.Get("unit").(string),
	}

	createdNotificationSetting, err := createDedicatedServerNotificationSetting(serverID, "bandwidth", &notificationSetting)
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

	notificationSetting, err := getDedicatedServerNotificationSetting(serverID, "bandwidth", notificationSettingID)
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
			Threshold: d.Get("threshold").(float64),
			Unit:      d.Get("unit").(string),
		}

		if _, err := updateDedicatedServerNotificationSetting(serverID, "bandwidth", notificationSettingID, &notificationSetting); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDedicatedServerNotificationSettingBandwidthRead(ctx, d, m)
}

func resourceDedicatedServerNotificationSettingBandwidthDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	serverID := d.Get("dedicated_server_id").(string)
	notificationSettingID := d.Get("id").(string)

	if err := deleteDedicatedServerNotificationSetting(serverID, "bandwidth", notificationSettingID); err != nil {
		return diag.FromErr(err)
	}

	return diags
}
