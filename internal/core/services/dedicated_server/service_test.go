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
	operatingSystems domain.OperatingSystems

	getAllDedicatedServerError  *sharedRepository.RepositoryError
	getAllOperatingSystemsError *sharedRepository.RepositoryError
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

	t.Run(
		"error from repository GetAllOperatingSystems bubbles up",
		func(t *testing.T) {
			want := "some error"
			generalError := sharedRepository.NewGeneralError("", errors.New(want))
			repoResponse := domain.OperatingSystems{{}, {}, {}}

			service := New(&repositorySpy{operatingSystems: repoResponse, getAllOperatingSystemsError: generalError})
			_, err := service.GetAllOperatingSystems(context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, want)
		},
	)
}
