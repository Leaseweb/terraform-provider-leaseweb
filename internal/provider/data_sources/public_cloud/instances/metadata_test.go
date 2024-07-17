package instances

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/stretchr/testify/assert"
)

func Test_instancesDataSource_Metadata(t *testing.T) {
	resp := datasource.MetadataResponse{}
	instancesDataSource := NewInstancesDataSource()

	instancesDataSource.Metadata(
		context.TODO(),
		datasource.MetadataRequest{ProviderTypeName: "tralala"},
		&resp,
	)

	assert.Equal(
		t,
		"tralala_public_cloud_instances",
		resp.TypeName,
		"Type name should be tralala_public_cloud_instances",
	)
}
