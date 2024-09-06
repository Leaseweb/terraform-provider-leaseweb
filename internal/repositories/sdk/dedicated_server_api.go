package sdk

import (
	"context"

	"github.com/leaseweb/leaseweb-go-sdk/dedicatedServer"
)

// DedicatedServerApi contains all methods that the sdk must support.
type DedicatedServerApi interface {
	GetServerList(ctx context.Context) dedicatedServer.ApiGetServerListRequest
	GetOperatingSystemList(ctx context.Context) dedicatedServer.ApiGetOperatingSystemListRequest
	GetControlPanelList(ctx context.Context) dedicatedServer.ApiGetControlPanelListRequest
	GetServer(ctx context.Context, id string) dedicatedServer.ApiGetServerRequest
	GetServerDataTrafficNotificationSetting(
		ctx context.Context,
		serverId string,
		dataTrafficNotificationSettingId string,
	) dedicatedServer.ApiGetServerDataTrafficNotificationSettingRequest
	CreateServerDataTrafficNotificationSetting(
		ctx context.Context,
		serverId string,
	) dedicatedServer.ApiCreateServerDataTrafficNotificationSettingRequest
	UpdateServerDataTrafficNotificationSetting(
		ctx context.Context,
		serverId string,
		dataTrafficNotificationSettingId string,
	) dedicatedServer.ApiUpdateServerDataTrafficNotificationSettingRequest
	DeleteServerDataTrafficNotificationSetting(
		ctx context.Context,
		serverId string,
		dataTrafficNotificationSettingId string,
	) dedicatedServer.ApiDeleteServerDataTrafficNotificationSettingRequest
}
