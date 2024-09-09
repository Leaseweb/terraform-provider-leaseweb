package notification_setting_bandwidth

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
)

var (
	_ resource.Resource              = &notificationSettingBandwidthResource{}
	_ resource.ResourceWithConfigure = &notificationSettingBandwidthResource{}
)

func New() resource.Resource {
	return &notificationSettingBandwidthResource{}
}

type notificationSettingBandwidthResource struct {
	client client.Client
}
