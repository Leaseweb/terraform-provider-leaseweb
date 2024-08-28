package control_panels

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/stretchr/testify/assert"
)

func Test_controlPanelsDataSource_Metadata(t *testing.T) {
	resp := datasource.MetadataResponse{}
	controlPanelsDataSource := New()

	controlPanelsDataSource.Metadata(
		context.TODO(),
		datasource.MetadataRequest{ProviderTypeName: "tralala"},
		&resp,
	)

	assert.Equal(
		t,
		"tralala_dedicated_server_control_panels",
		resp.TypeName,
		"Type name should be tralala_dedicated_server_control_panels",
	)
}
