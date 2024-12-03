package publiccloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/v2/publiccloud"
	"github.com/stretchr/testify/assert"
)

func Test_adaptIsoToInstanceISOResource(t *testing.T) {
	t.Run("expected value is returned if ISO is not set", func(t *testing.T) {
		desiredISOId := "desiredIsoId"
		got := adaptIsoToInstanceISOResource(
			&desiredISOId,
			"instanceId",
			nil,
		)

		want := instanceISOResourceModel{
			DesiredID:  basetypes.NewStringPointerValue(&desiredISOId),
			ID:         basetypes.NewStringPointerValue(nil),
			InstanceID: basetypes.NewStringValue("instanceId"),
			Name:       basetypes.NewStringPointerValue(nil),
		}

		assert.Equal(t, want, got)
	})

	t.Run("expected value is returned if ISO is set", func(t *testing.T) {
		sdkISO := publiccloud.Iso{
			Id:   "id",
			Name: "name",
		}
		desiredISOId := "desiredIsoId"
		got := adaptIsoToInstanceISOResource(
			&desiredISOId,
			"instanceId",
			&sdkISO,
		)

		want := instanceISOResourceModel{
			DesiredID:  basetypes.NewStringPointerValue(&desiredISOId),
			ID:         basetypes.NewStringPointerValue(&sdkISO.Id),
			InstanceID: basetypes.NewStringValue("instanceId"),
			Name:       basetypes.NewStringPointerValue(&sdkISO.Name),
		}

		assert.Equal(t, want, got)
	})
}
