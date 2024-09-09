// Package to_resource_model implements adapters to convert domain entities to resource models.
package to_resource_model

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain/dedicated_server"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/facades/shared"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/resources/dedicated_server/model"
)

// AdaptNotificationSettingBandwidth adapts dedicated_server.NotificationSettingBandwidth to model.NotificationSettingBandwidth.
func AdaptNotificationSettingBandwidth(
	ctx context.Context,
	notificationSettingBandwidth dedicated_server.NotificationSettingBandwidth,
) (*model.NotificationSettingBandwidth, error) {
	plan := model.NotificationSettingBandwidth{}

	plan.ServerId = basetypes.NewStringValue(notificationSettingBandwidth.ServerId)
	plan.Id = basetypes.NewStringValue(notificationSettingBandwidth.Id)
	plan.Frequency = basetypes.NewStringValue(notificationSettingBandwidth.Frequency)
	plan.LastCheckedAt = shared.AdaptNullableStringToStringValue(notificationSettingBandwidth.LastCheckedAt)
	plan.Threshold = basetypes.NewStringValue(notificationSettingBandwidth.Threshold)
	plan.ThresholdExceededAt = shared.AdaptNullableStringToStringValue(notificationSettingBandwidth.ThresholdExceededAt)
	plan.Unit = basetypes.NewStringValue(notificationSettingBandwidth.Unit)

	actions, err := shared.AdaptEntitiesToListValue(
		notificationSettingBandwidth.Actions,
		model.Action{}.AttributeTypes(),
		ctx,
		adaptAction,
	)
	if err != nil {
		return nil, fmt.Errorf("AdaptNotificationSettingBandwidth: %w", err)
	}
	plan.Actions = actions

	return &plan, nil
}

func adaptAction(
	ctx context.Context,
	action dedicated_server.Action,
) (*model.Action, error) {

	return &model.Action{
		LastTriggeredAt: shared.AdaptNullableStringToStringValue(action.LastTriggeredAt),
		Type:            basetypes.NewStringValue(action.Type),
	}, nil
}
