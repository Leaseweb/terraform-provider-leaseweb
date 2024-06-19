package modify_plan

import (
	"context"
	"testing"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
	providerClient "terraform-provider-leaseweb/internal/client"
)

func TestInstanceTypes_AccGetAllowedInstanceTypes(t *testing.T) {
	client := providerClient.NewClient("token", &providerClient.Options{
		Host:   "localhost:8080",
		Scheme: "http",
	})
	instanceTypes := NewInstanceTypes(*client, context.TODO())
	allowedInstanceTypes, _, _ := instanceTypes.GetAllowedInstanceTypes("e8b1b70b-5c42-4b12-9c21-5c2f95f6b476")

	assert.NotEmptyf(t, allowedInstanceTypes, "")
}

func TestInstanceTypes_AccGetAllowedInstanceTypes_Error(t *testing.T) {
	client := providerClient.NewClient("token", &providerClient.Options{
		Host:   "localhost:8080",
		Scheme: "http",
	})
	instanceTypes := NewInstanceTypes(*client, context.TODO())
	allowedInstanceTypes, response, err := instanceTypes.GetAllowedInstanceTypes("tralala")

	assert.Empty(t, allowedInstanceTypes, "expected types to be empty")
	assert.NotNilf(t, response, "expected response to not be nil")
	assert.NotNil(t, err, "expected error to not be nil")
}

func Test_convertAllowedInstancesTypesToString(t *testing.T) {
	instanceType := publicCloud.NewInstanceType()
	instanceType.SetName("tralala")

	updateInstanceTypes := []publicCloud.InstanceType{*instanceType}

	got := convertSdkInstanceTypesToString(updateInstanceTypes)
	want := []string{"tralala"}

	assert.Equal(t, want, got)
}
