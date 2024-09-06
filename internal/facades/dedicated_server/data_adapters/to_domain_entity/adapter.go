package to_domain_entity

import (
	domain "github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain/dedicated_server"
	resourceModel "github.com/leaseweb/terraform-provider-leaseweb/internal/provider/resources/dedicated_server/model"
)

func AdaptToCreateDataTrafficNotificationSettingOpts(
	dataTrafficNotificationSetting resourceModel.DataTrafficNotificationSetting,
) domain.DataTrafficNotificationSetting {
	return domain.NewCreateDataTrafficNotificationSetting(
		dataTrafficNotificationSetting.Frequency.ValueString(),
		dataTrafficNotificationSetting.Threshold.ValueString(),
		dataTrafficNotificationSetting.Unit.ValueString(),
	)
}
