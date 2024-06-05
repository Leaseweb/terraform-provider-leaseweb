package instance

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_instanceResource_Metadata(t *testing.T) {
	resp := resource.MetadataResponse{}
	instanceResource := NewInstanceResource()

	instanceResource.Metadata(context.TODO(), resource.MetadataRequest{ProviderTypeName: "tralala"}, &resp)

	assert.Equal(t, "tralala_instance", resp.TypeName, "Type name should be tralala_instances")
}
