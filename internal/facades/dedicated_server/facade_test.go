package dedicated_server

import (
	"context"
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	serviceErrors "github.com/leaseweb/terraform-provider-leaseweb/internal/core/services/errors"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/data_sources/dedicated_server/model"

	domain "github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain/dedicated_server"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/ports"
	"github.com/stretchr/testify/assert"
)

var (
	_ ports.DedicatedServerService = &serviceSpy{}
)

type serviceSpy struct {
	dedicatedServers domain.DedicatedServers
	operatingSystems domain.OperatingSystems

	getAllDedicatedServerError *serviceErrors.ServiceError
	getAllOperatingSystemError *serviceErrors.ServiceError
}

func (s *serviceSpy) GetAllDedicatedServers(ctx context.Context) (*domain.DedicatedServers, *serviceErrors.ServiceError) {
	return &s.dedicatedServers, s.getAllDedicatedServerError
}

func (s *serviceSpy) GetAllOperatingSystems(ctx context.Context) (domain.OperatingSystems, *serviceErrors.ServiceError) {
	return s.operatingSystems, s.getAllOperatingSystemError
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
