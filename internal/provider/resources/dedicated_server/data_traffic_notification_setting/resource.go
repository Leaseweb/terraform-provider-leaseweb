package data_traffic_notification_setting

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
)

var (
	_ resource.Resource                = &dataTrafficNotificationSettingResource{}
	_ resource.ResourceWithConfigure   = &dataTrafficNotificationSettingResource{}
	_ resource.ResourceWithModifyPlan  = &dataTrafficNotificationSettingResource{}
	_ resource.ResourceWithImportState = &dataTrafficNotificationSettingResource{}
)

func NewDataTrafficNotificationSettingResource() resource.Resource {
	return &dataTrafficNotificationSettingResource{}
}

type dataTrafficNotificationSettingResource struct {
	client client.Client
}
