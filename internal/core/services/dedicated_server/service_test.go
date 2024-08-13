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

	getAllDedicatedServerError *sharedRepository.RepositoryError
}

func (r *repositorySpy) GetAllDedicatedServers(ctx context.Context) (
	domain.DedicatedServers,
	*sharedRepository.RepositoryError,
) {
	return r.dedicatedServers, r.getAllDedicatedServerError
}

func newRepositorySpy() repositorySpy {
	return repositorySpy{}
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

			spy := newRepositorySpy()
			spy.dedicatedServers = want

			service := New(&spy)

			got, err := service.GetAllDedicatedServers(context.TODO())

			assert.Nil(t, err)
			assert.Equal(t, want, *got)
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
