package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/domain/entity"
	"terraform-provider-leaseweb/internal/core/shared/value_object"
)

func TestInstances_Populate(t *testing.T) {
	t.Run("instance is set properly", func(t *testing.T) {
		instanceId := value_object.NewGeneratedUuid()

		instance := entity.Instance{
			Id: instanceId,
		}

		instances := Instances{}
		instances.Populate(entity.Instances{instance})

		assert.Equal(
			t,
			instanceId.String(),
			instances.Instances[0].Id.ValueString(),
			"instance should be set",
		)
	})
}
