package dedicated_servers

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/stretchr/testify/assert"
)

func Test_dedicatedServerDataSource_Metadata(t *testing.T) {
	resp := datasource.MetadataResponse{}
	dedicatedServerDataSource := NewDedicatedServerDataSource()

	dedicatedServerDataSource.Metadata(
		context.TODO(),
		datasource.MetadataRequest{ProviderTypeName: "tralala"},
		&resp,
	)

	assert.Equal(
		t,
		"tralala_dedicated_servers",
		resp.TypeName,
		"Type name should be tralala_dedicated_servers",
	)
}
