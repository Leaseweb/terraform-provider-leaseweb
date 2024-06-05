package instances

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_instancesDataSource_Metadata(t *testing.T) {
	resp := datasource.MetadataResponse{}
	instancesDataSource := NewInstancesDataSource()

	instancesDataSource.Metadata(context.TODO(), datasource.MetadataRequest{ProviderTypeName: "tralala"}, &resp)

	assert.Equal(t, "tralala_instances", resp.TypeName, "Type name should be tralala_instances")
}
