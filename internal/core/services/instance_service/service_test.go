package instance_service

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/domain/entity"
)

func TestService_GetAllInstances(t *testing.T) {
	instanceService := New()

	got, err := instanceService.GetAllInstances()

	assert.Nil(t, err)
	assert.Equal(t, entity.Instances{}, got)
}

func TestService_GetInstance(t *testing.T) {
	instanceService := New()
	id, _ := uuid.NewUUID()

	got, err := instanceService.GetInstance(id)

	assert.Nil(t, err)
	assert.Equal(t, entity.Instance{}, got)
}

func TestService_CreateInstance(t *testing.T) {
	instanceService := New()

	got, err := instanceService.CreateInstance(entity.Instance{})

	assert.Nil(t, err)
	assert.Equal(t, entity.Instance{}, got)
}

func TestService_UpdateInstance(t *testing.T) {
	instanceService := New()

	got, err := instanceService.UpdateInstance(entity.Instance{})

	assert.Nil(t, err)
	assert.Equal(t, entity.Instance{}, got)
}

func TestService_DeleteInstance(t *testing.T) {
	instanceService := New()
	id, _ := uuid.NewUUID()

	err := instanceService.DeleteInstance(id)

	assert.Nil(t, err)
}
