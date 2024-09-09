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
	dedicatedServers             domain.DedicatedServers
	operatingSystems             domain.OperatingSystems
	controlPanels                domain.ControlPanels
	notificationSettingBandwidth *domain.NotificationSettingBandwidth

	getAllDedicatedServerError              *sharedRepository.RepositoryError
	getAllOperatingSystemsError             *sharedRepository.RepositoryError
	getAllControlPanelError                 *sharedRepository.RepositoryError
	createNotificationSettingBandwidthError *sharedRepository.RepositoryError
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

func (r *repositorySpy) CreateNotificationSettingBandwidth(
	notificationSettingBandwidth domain.NotificationSettingBandwidth,
	ctx context.Context,
) (*domain.NotificationSettingBandwidth, *sharedRepository.RepositoryError) {
	return r.notificationSettingBandwidth, r.createNotificationSettingBandwidthError
}

func TestService_GetAllDedicatedServers(t *testing.T) {
	t.Run(
		"service passes back dedicated server from repository",
		func(t *testing.T) {

			id := "123456"

			want := domain.DedicatedServers{
				domain.DedicatedServer{
					Id: id,
				},
			}

			spy := repositorySpy{dedicatedServers: want}

			service := New(&spy)

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
		"service passes back dedicated server from repository",
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

func TestService_CreateNotificationSettingBandwidth(t *testing.T) {
	t.Run(
		"service passed back notificationSettingBandwidth from repository while creating a notification setting bandwidth",
		func(t *testing.T) {

			want := domain.NotificationSettingBandwidth{
				ServerId:  "12345",
				Frequency: "DAILY",
				Threshold: "1",
				Unit:      "Gbps",
			}
			spy := repositorySpy{notificationSettingBandwidth: &want}
			service := New(&spy)
			got, err := service.CreateNotificationSettingBandwidth(want, context.TODO())

			assert.Equal(t, want, *got)
			assert.Nil(t, err)
		},
	)

	t.Run(
		"service passed back error while creating a notification setting bandwidth",
		func(t *testing.T) {
			service := New(
				&repositorySpy{
					createNotificationSettingBandwidthError: sharedRepository.NewGeneralError(
						"",
						errors.New("some error"),
					),
				},
			)

			notificationSettingBandwidth := domain.NotificationSettingBandwidth{
				ServerId:  "12345",
				Frequency: "DAILY",
				Threshold: "1",
				Unit:      "Gbps",
			}
			_, err := service.CreateNotificationSettingBandwidth(notificationSettingBandwidth, context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "some error")
		},
	)
}
