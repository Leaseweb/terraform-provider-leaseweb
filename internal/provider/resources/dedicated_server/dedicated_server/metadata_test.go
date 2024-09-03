package dedicated_server

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/stretchr/testify/assert"
)

func Test_dedicatedServerResource_Metadata(t *testing.T) {
	resp := resource.MetadataResponse{}
	dedicatedServerResource := NewDedicatedServerResource()

	dedicatedServerResource.Metadata(
		context.TODO(),
		resource.MetadataRequest{ProviderTypeName: "tralala"},
		&resp,
	)

	assert.Equal(t,
		"tralala_dedicated_servers",
		resp.TypeName,
		"Type name should be tralala_dedicated_servers",
	)
}
