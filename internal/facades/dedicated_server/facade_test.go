package dedicated_server

import (
	"context"
	"errors"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	serviceErrors "github.com/leaseweb/terraform-provider-leaseweb/internal/core/services/errors"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/data_sources/dedicated_server/model"

	domain "github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain/dedicated_server"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/ports"
	resourceModel "github.com/leaseweb/terraform-provider-leaseweb/internal/provider/resources/dedicated_server/model"
	"github.com/stretchr/testify/assert"
)

var (
	_ ports.DedicatedServerService = &serviceSpy{}
)

type serviceSpy struct {
	dedicatedServers             domain.DedicatedServers
	operatingSystems             domain.OperatingSystems
	controlPanels                domain.ControlPanels
	notificationSettingBandwidth *domain.NotificationSettingBandwidth

	getAllDedicatedServerError              *serviceErrors.ServiceError
	getAllOperatingSystemError              *serviceErrors.ServiceError
	getAllControlPanelError                 *serviceErrors.ServiceError
	createNotificationSettingBandwidthError *serviceErrors.ServiceError
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

func (s *serviceSpy) CreateNotificationSettingBandwidth(
	notificationSettingBandwidth domain.NotificationSettingBandwidth,
	ctx context.Context,
) (*domain.NotificationSettingBandwidth, *serviceErrors.ServiceError) {
	return s.notificationSettingBandwidth, s.createNotificationSettingBandwidthError
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

func TestCreateNotificationSettingBandwidth(t *testing.T) {
	t.Run(
		"facade passes back notification setting bandwidth from service while creating a notification setting bandwidth",
		func(t *testing.T) {
			date := "2024-05-04"
			action := map[string]attr.Value{
				"last_triggered_at": basetypes.NewStringValue(date),
				"type":              basetypes.NewStringValue("EMAIL"),
			}
			actionObject, _ := basetypes.NewObjectValue(resourceModel.Action{}.AttributeTypes(), action)
			actionsList, _ := basetypes.NewListValue(
				types.ObjectType{AttrTypes: resourceModel.Action{}.AttributeTypes()},
				[]attr.Value{actionObject},
			)

			want := resourceModel.NotificationSettingBandwidth{
				ServerId:            basetypes.NewStringValue("12345"),
				Id:                  basetypes.NewStringValue("123456"),
				LastCheckedAt:       basetypes.NewStringValue("2024-05-04"),
				ThresholdExceededAt: basetypes.NewStringValue("2024-05-04"),
				Frequency:           basetypes.NewStringValue("DAILY"),
				Threshold:           basetypes.NewStringValue("1"),
				Unit:                basetypes.NewStringValue("Gbps"),
				Actions:             actionsList,
			}

			domainNotificationSettingBandwidth := domain.NotificationSettingBandwidth{
				ServerId:            "12345",
				Id:                  "123456",
				LastCheckedAt:       &date,
				ThresholdExceededAt: &date,
				Frequency:           "DAILY",
				Threshold:           "1",
				Unit:                "Gbps",
				Actions:             domain.Actions{domain.Action{LastTriggeredAt: &date, Type: "EMAIL"}},
			}

			spy := serviceSpy{notificationSettingBandwidth: &domainNotificationSettingBandwidth}

			facade := New(&spy)

			got, err := facade.CreateNotificationSettingBandwidth(want, context.TODO())

			assert.Nil(t, err)
			assert.Equal(t, want, *got)
		},
	)

	t.Run(
		"error from service CreateNotificationSettingBandwidth bubbles up while creating a notification setting bandwidth",
		func(t *testing.T) {
			facade := New(
				&serviceSpy{
					createNotificationSettingBandwidthError: serviceErrors.NewError(
						"",
						errors.New("some error"),
					),
				},
			)

			resourceModelNotificationSettingBandwidth := resourceModel.NotificationSettingBandwidth{
				ServerId:  basetypes.NewStringValue("MONTHLY"),
				Frequency: basetypes.NewStringValue("DAILY"),
				Threshold: basetypes.NewStringValue("1"),
				Unit:      basetypes.NewStringValue("Gbps"),
			}

			_, err := facade.CreateNotificationSettingBandwidth(resourceModelNotificationSettingBandwidth, context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "some error")
		},
	)
}

func TestGetFrequencies(t *testing.T) {
	facade := DedicatedServerFacade{}
	want := []string{"DAILY", "WEEKLY", "MONTHLY"}
	got := facade.GetFrequencies()

	assert.Equal(t, want, got)
}

func TestGetUnits(t *testing.T) {
	facade := DedicatedServerFacade{}
	want := []string{"Mbps", "Gbps"}
	got := facade.GetUnits()

	assert.Equal(t, want, got)
}
