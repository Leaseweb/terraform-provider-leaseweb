package dedicated_server

import (
	"context"
	"errors"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	serviceErrors "github.com/leaseweb/terraform-provider-leaseweb/internal/core/services/errors"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/data_sources/dedicated_server/model"
	"testing"

	domain "github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain/dedicated_server"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/ports"
	"github.com/stretchr/testify/assert"
)

var (
	_ ports.DedicatedServerService = &serviceSpy{}
)

type serviceSpy struct {
	dedicatedServers domain.DedicatedServers
	controlPanels    domain.ControlPanels

	getAllDedicatedServerError *serviceErrors.ServiceError
	getAllControlPanelError    *serviceErrors.ServiceError
}

func (s *serviceSpy) GetAllDedicatedServers(ctx context.Context) (domain.DedicatedServers, *serviceErrors.ServiceError) {
	return s.dedicatedServers, s.getAllDedicatedServerError
}

func (s *serviceSpy) GetAllControlPanels(ctx context.Context) (domain.ControlPanels, *serviceErrors.ServiceError) {
	return s.controlPanels, s.getAllControlPanelError
}

func newServiceSpy() serviceSpy {
	return serviceSpy{}
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

			spy := newServiceSpy()
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
