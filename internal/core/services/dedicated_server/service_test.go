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
	dedicatedServers domain.DedicatedServers
	controlPanels    domain.ControlPanels

	getAllDedicatedServerError *sharedRepository.RepositoryError
	getAllControlPanelError    *sharedRepository.RepositoryError
}

func (r *repositorySpy) GetAllDedicatedServers(ctx context.Context) (
	domain.DedicatedServers,
	*sharedRepository.RepositoryError,
) {
	return r.dedicatedServers, r.getAllDedicatedServerError
}

func (r *repositorySpy) GetAllControlPanels(ctx context.Context) (
	domain.ControlPanels,
	*sharedRepository.RepositoryError,
) {
	return r.controlPanels, r.getAllControlPanelError
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
