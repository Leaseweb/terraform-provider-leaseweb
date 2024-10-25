package publiccloud

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/stretchr/testify/assert"
)

func Test_credentialDataSource_Metadata(t *testing.T) {
	resp := datasource.MetadataResponse{}
	credentialDataSource := NewCredentialDataSource()

	credentialDataSource.Metadata(
		context.TODO(),
		datasource.MetadataRequest{ProviderTypeName: "tralala"},
		&resp,
	)

	assert.Equal(
		t,
		"tralala_public_cloud_credential",
		resp.TypeName,
		"Type name should be tralala_public_cloud_credential",
	)
}
