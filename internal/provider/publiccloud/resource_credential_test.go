package publiccloud

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/stretchr/testify/assert"
)

func Test_credentialResource_Metadata(t *testing.T) {
	resp := resource.MetadataResponse{}
	credentialResource := NewCredentialResource()

	credentialResource.Metadata(
		context.TODO(),
		resource.MetadataRequest{ProviderTypeName: "tralala"},
		&resp,
	)

	assert.Equal(
		t,
		"tralala_public_cloud_resource",
		resp.TypeName,
		"Type name should be tralala_public_cloud_resource",
	)
}
