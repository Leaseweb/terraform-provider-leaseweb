package dedicated_server

import (
	"context"
	"errors"
	"testing"

	domain "github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain/dedicated_server"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/ports"
	sharedRepository "github.com/leaseweb/terraform-provider-leaseweb/internal/repositories/shared"
	"github.com/stretchr/testify/assert"
)

var (
	_ ports.DedicatedServerRepository = &repositorySpy{}
)

type repositorySpy struct {
	dedicatedServers               domain.DedicatedServers
	operatingSystems               domain.OperatingSystems
	controlPanels                  domain.ControlPanels
	dedicatedServer                domain.DedicatedServer
	dataTrafficNotificationSetting domain.DataTrafficNotificationSetting

	getAllDedicatedServerError                *sharedRepository.RepositoryError
	getAllOperatingSystemsError               *sharedRepository.RepositoryError
	getAllControlPanelError                   *sharedRepository.RepositoryError
	getDedicatedServerError                   *sharedRepository.RepositoryError
	getDataTrafficNotificationSettingError    *sharedRepository.RepositoryError
	createDataTrafficNotificationSettingError *sharedRepository.RepositoryError
	updateDataTrafficNotificationSettingError *sharedRepository.RepositoryError
	deleteDataTrafficNotificationSettingError *sharedRepository.RepositoryError
}

func (r *repositorySpy) UpdateDataTrafficNotificationSetting(ctx context.Context, serverId string, dataTrafficNotificationSettingId string, dataTrafficNotificationSetting domain.DataTrafficNotificationSetting) (*domain.DataTrafficNotificationSetting, *sharedRepository.RepositoryError) {
	return &r.dataTrafficNotificationSetting, r.updateDataTrafficNotificationSettingError
}

func (r *repositorySpy) GetAllDedicatedServers(ctx context.Context) (
	domain.DedicatedServers,
	*sharedRepository.RepositoryError,
) {
	return r.dedicatedServers, r.getAllDedicatedServerError
}

func (r *repositorySpy) GetAllOperatingSystems(ctx context.Context) (
	domain.OperatingSystems,
	*sharedRepository.RepositoryError,
) {
	return r.operatingSystems, r.getAllOperatingSystemsError
}

func (r *repositorySpy) GetAllControlPanels(ctx context.Context) (
	domain.ControlPanels,
	*sharedRepository.RepositoryError,
) {
	return r.controlPanels, r.getAllControlPanelError
}

func (r *repositorySpy) GetDedicatedServer(ctx context.Context, id string) (
	*domain.DedicatedServer,
	*sharedRepository.RepositoryError,
) {
	return &r.dedicatedServer, r.getDedicatedServerError
}

func (r *repositorySpy) GetDataTrafficNotificationSetting(ctx context.Context, serverId string, dataTrafficNotificationSettingId string) (
	*domain.DataTrafficNotificationSetting,
	*sharedRepository.RepositoryError,
) {
	return &r.dataTrafficNotificationSetting, r.getDataTrafficNotificationSettingError
}

func (r *repositorySpy) CreateDataTrafficNotificationSetting(
	ctx context.Context,
	serverId string,
	dataTrafficNotificationSetting domain.DataTrafficNotificationSetting,
) (
	*domain.DataTrafficNotificationSetting,
	*sharedRepository.RepositoryError,
) {
	return &r.dataTrafficNotificationSetting, r.createDataTrafficNotificationSettingError
}

func (r *repositorySpy) DeleteDataTrafficNotificationSetting(ctx context.Context, serverId string, dataTrafficNotificationSettingId string) *sharedRepository.RepositoryError {
	return r.deleteDataTrafficNotificationSettingError
}

func TestService_GetAllDedicatedServers(t *testing.T) {
	t.Run(
		"service passes back dedicated servers from repository",
		func(t *testing.T) {
			id := "123456"
			want := domain.DedicatedServers{domain.DedicatedServer{Id: id}}
			service := New(&repositorySpy{dedicatedServers: want})

			got, err := service.GetAllDedicatedServers(context.TODO())

			assert.Nil(t, err)
			assert.Equal(t, want, got)
		},
	)

	t.Run(
		"error from repository getAllDedicatedServers bubbles up",
		func(t *testing.T) {
			service := New(
				&repositorySpy{
					getAllDedicatedServerError: sharedRepository.NewGeneralError(
						"",
						errors.New("some error"),
					),
				},
			)

			_, err := service.GetAllDedicatedServers(context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "some error")
		},
	)

	t.Run(
		"error from repository getAllDedicatedServers bubbles up", func(t *testing.T) {
			service := New(
				&repositorySpy{
					dedicatedServers: domain.DedicatedServers{
						{Id: "1"},
						{Id: "2"},
						{Id: "3"},
					},
					getAllDedicatedServerError: sharedRepository.NewGeneralError(
						"",
						errors.New("some error"),
					),
				},
			)

			_, err := service.GetAllDedicatedServers(context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "some error")
		},
	)
}

func TestService_GetAllOperatingSystems(t *testing.T) {
	t.Run(
		"service passes back operating systems from repository",
		func(t *testing.T) {
			want := domain.OperatingSystems{{Id: "123456"}}
			service := New(&repositorySpy{operatingSystems: want})

			got, err := service.GetAllOperatingSystems(context.TODO())

			assert.Nil(t, err)
			assert.Equal(t, want, got)
		},
	)

	t.Run(
		"error from repository getAllOperatingSystems bubbles up",
		func(t *testing.T) {
			want := "some error"
			generalError := sharedRepository.NewGeneralError("", errors.New(want))
			service := New(&repositorySpy{getAllOperatingSystemsError: generalError})

			_, err := service.GetAllOperatingSystems(context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, want)
		},
	)
}

func TestService_GetAllControlPanels(t *testing.T) {
	t.Run(
		"service passes back control panels from repository",
		func(t *testing.T) {
			want := domain.ControlPanels{
				domain.ControlPanel{
					Id:   "id",
					Name: "name",
				},
			}
			spy := repositorySpy{controlPanels: want}
			service := New(&spy)

			got, err := service.GetAllControlPanels(context.TODO())

			assert.Equal(t, want, got)
			assert.Nil(t, err)
		},
	)

	t.Run(
		"error from repository getAllControlPanels bubbles up",
		func(t *testing.T) {
			service := New(
				&repositorySpy{
					getAllControlPanelError: sharedRepository.NewGeneralError(
						"",
						errors.New("some error"),
					),
				},
			)

			_, err := service.GetAllControlPanels(context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "some error")
		},
	)
}

func TestService_GetDedicatedServer(t *testing.T) {
	t.Run(
		"service passes back dedicated server from repository",
		func(t *testing.T) {
			id := "123456"
			want := domain.DedicatedServer{Id: id}
			service := New(&repositorySpy{dedicatedServer: want})

			got, err := service.GetDedicatedServer(context.TODO(), id)

			assert.Nil(t, err)
			assert.Equal(t, want, *got)
		},
	)

	t.Run(
		"error from repository GetDedicatedServer bubbles up",
		func(t *testing.T) {
			service := New(
				&repositorySpy{
					getDedicatedServerError: sharedRepository.NewGeneralError(
						"",
						errors.New("some error"),
					),
				},
			)

			_, err := service.GetDedicatedServer(context.TODO(), "id")

			assert.Error(t, err)
			assert.ErrorContains(t, err, "some error")
		},
	)
}
