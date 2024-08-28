package operating_systems

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/stretchr/testify/assert"
)

func Test_operatingSystemsDataSource_Metadata(t *testing.T) {
	resp := datasource.MetadataResponse{}
	operatingSystemsDataSource := New()

	operatingSystemsDataSource.Metadata(
		context.TODO(),
		datasource.MetadataRequest{ProviderTypeName: "test_provider"},
		&resp,
	)

	assert.Equal(
		t,
		"test_provider_dedicated_server_operating_systems",
		resp.TypeName,
		"Type name should be test_provider_dedicated_server_operating_systems",
	)
}
