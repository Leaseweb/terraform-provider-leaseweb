package dedicated_server

import (
	"context"
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	serviceErrors "github.com/leaseweb/terraform-provider-leaseweb/internal/core/services/errors"
	model "github.com/leaseweb/terraform-provider-leaseweb/internal/provider/data_sources/dedicated_server/model"
	resourceModel "github.com/leaseweb/terraform-provider-leaseweb/internal/provider/resources/dedicated_server/model"

	domain "github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain/dedicated_server"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/ports"
	"github.com/stretchr/testify/assert"
)

var (
	_ ports.DedicatedServerService = &serviceSpy{}
)

type serviceSpy struct {
	dedicatedServers               domain.DedicatedServers
	dedicatedServer                domain.DedicatedServer
	operatingSystems               domain.OperatingSystems
	controlPanels                  domain.ControlPanels
	dataTrafficNotificationSetting domain.DataTrafficNotificationSetting

	getAllDedicatedServerError                *serviceErrors.ServiceError
	getAllOperatingSystemError                *serviceErrors.ServiceError
	getAllControlPanelError                   *serviceErrors.ServiceError
	getDedicatedServerError                   *serviceErrors.ServiceError
	getDataTrafficNotificationSettingError    *serviceErrors.ServiceError
	createDataTrafficNotificationSettingError *serviceErrors.ServiceError
	deleteDataTrafficNotificationSettingError *serviceErrors.ServiceError
	updateDataTrafficNotificationSettingError *serviceErrors.ServiceError
}

func (s *serviceSpy) UpdateDataTrafficNotificationSetting(ctx context.Context, serverId string, dataTrafficNotificationSettingId string, dataTrafficNotificationSetting domain.DataTrafficNotificationSetting) (*domain.DataTrafficNotificationSetting, *serviceErrors.ServiceError) {
	return &s.dataTrafficNotificationSetting, s.updateDataTrafficNotificationSettingError
}

func (s *serviceSpy) GetAllDedicatedServers(ctx context.Context) (domain.DedicatedServers, *serviceErrors.ServiceError) {
	return s.dedicatedServers, s.getAllDedicatedServerError
}

func (s *serviceSpy) GetAllOperatingSystems(ctx context.Context) (domain.OperatingSystems, *serviceErrors.ServiceError) {
	return s.operatingSystems, s.getAllOperatingSystemError
}

func (s *serviceSpy) GetAllControlPanels(ctx context.Context) (domain.ControlPanels, *serviceErrors.ServiceError) {
	return s.controlPanels, s.getAllControlPanelError
}

func (s *serviceSpy) GetDedicatedServer(ctx context.Context, id string) (*domain.DedicatedServer, *serviceErrors.ServiceError) {
	return &s.dedicatedServer, s.getDedicatedServerError
}

func (s *serviceSpy) GetDataTrafficNotificationSetting(
	ctx context.Context,
	serverId string,
	dataTrafficNotificationSettingId string,
) (
	*domain.DataTrafficNotificationSetting,
	*serviceErrors.ServiceError,
) {
	return &s.dataTrafficNotificationSetting, s.getDataTrafficNotificationSettingError
}

func (s *serviceSpy) CreateDataTrafficNotificationSetting(ctx context.Context, serverId string, dataTrafficNotificationSetting domain.DataTrafficNotificationSetting) (*domain.DataTrafficNotificationSetting, *serviceErrors.ServiceError) {
	return &s.dataTrafficNotificationSetting, s.createDataTrafficNotificationSettingError
}

func (s *serviceSpy) DeleteDataTrafficNotificationSetting(
	ctx context.Context,
	serverId string,
	dataTrafficNotificationSettingId string,
) *serviceErrors.ServiceError {
	return s.deleteDataTrafficNotificationSettingError
}

func TestFacadeGetAllDedicatedServers(t *testing.T) {
	t.Run(
		"facade passes back dedicated servers from service",
		func(t *testing.T) {

			id := "123456"

			want := model.DedicatedServers{
				DedicatedServers: []model.DedicatedServer{
					{Id: basetypes.NewStringValue(id)},
				},
			}

			facade := New(&serviceSpy{
				dedicatedServers: domain.DedicatedServers{
					domain.DedicatedServer{Id: id},
				},
			})
			facade.adaptDedicatedServersToDatasourceModel = func(dedicatedServers domain.DedicatedServers) model.DedicatedServers {
				return model.DedicatedServers{
					DedicatedServers: []model.DedicatedServer{{Id: basetypes.NewStringValue(id)}},
				}
			}
			got, err := facade.GetAllDedicatedServers(context.TODO())

			assert.Nil(t, err)
			assert.Equal(t, want, *got)
		},
	)

	t.Run(
		"error from service GetAllDedicatedServers bubbles up",
		func(t *testing.T) {

			want := "some error"
			serviceError := serviceErrors.NewError("", errors.New(want))

			facade := New(&serviceSpy{getAllDedicatedServerError: serviceError})
			_, err := facade.GetAllDedicatedServers(context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "some error")
		},
	)
}

func TestFacadeGetAllOperatingSystems(t *testing.T) {
	t.Run(
		"facade passes back operating sysetms from service",
		func(t *testing.T) {

			id := "123456"
			name := "name"

			want := model.OperatingSystems{
				OperatingSystems: []model.OperatingSystem{
					{Id: basetypes.NewStringValue(id), Name: basetypes.NewStringValue(name)},
				},
			}

			facade := New(&serviceSpy{
				operatingSystems: domain.OperatingSystems{
					domain.NewOperatingSystem(id, name),
				},
			})
			got, err := facade.GetAllOperatingSystems(context.TODO())

			assert.Nil(t, err)
			assert.Equal(t, want, *got)
		},
	)

	t.Run(
		"error from service getAllOperatingSystems bubbles up",
		func(t *testing.T) {

			want := "some error"
			serviceError := serviceErrors.NewError("", errors.New(want))

			facade := New(&serviceSpy{getAllOperatingSystemError: serviceError})
			_, err := facade.GetAllOperatingSystems(context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "some error")
		},
	)
}

func TestGetAllControlPanels(t *testing.T) {
	t.Run(
		"facade passes back control panels from service",
		func(t *testing.T) {

			id := "123456"
			name := "name"

			want := model.ControlPanels{
				ControlPanels: []model.ControlPanel{
					{Id: basetypes.NewStringValue(id), Name: basetypes.NewStringValue(name)},
				},
			}

			spy := serviceSpy{}
			spy.controlPanels = domain.ControlPanels{
				domain.NewControlPanel(id, name),
			}

			facade := New(&spy)

			got, err := facade.GetAllControlPanels(context.TODO())

			assert.Nil(t, err)
			assert.Equal(t, want, *got)
		},
	)

	t.Run(
		"error from service getAllControlPanels bubbles up",
		func(t *testing.T) {
			facade := New(
				&serviceSpy{
					getAllControlPanelError: serviceErrors.NewError(
						"",
						errors.New("some error"),
					),
				},
			)

			_, err := facade.GetAllControlPanels(context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "some error")
		},
	)
}

func TestFacadeGetDedicatedServer(t *testing.T) {
	t.Run(
		"facade passes back dedicated server from service",
		func(t *testing.T) {

			id := "123456"

			want := resourceModel.DedicatedServer{
				Id: basetypes.NewStringValue(id),
			}

			facade := New(&serviceSpy{
				dedicatedServer: domain.DedicatedServer{Id: id},
			})
			facade.adaptDedicatedServerToResourceModel = func(dedicatedServer domain.DedicatedServer) resourceModel.DedicatedServer {
				return resourceModel.DedicatedServer{Id: basetypes.NewStringValue(id)}
			}

			got, err := facade.GetDedicatedServer(context.TODO(), id)

			assert.Nil(t, err)
			assert.Equal(t, want, *got)
		},
	)

	t.Run(
		"error from service GetDedicatedServer bubbles up",
		func(t *testing.T) {

			want := "some error"
			serviceError := serviceErrors.NewError("", errors.New(want))

			facade := New(&serviceSpy{getDedicatedServerError: serviceError})
			_, err := facade.GetDedicatedServer(context.TODO(), "id")

			assert.Error(t, err)
			assert.ErrorContains(t, err, "some error")
		},
	)
}
