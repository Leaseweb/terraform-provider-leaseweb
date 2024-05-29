package instances

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSchema(t *testing.T) {
	instancesDataSource := NewInstancesDataSource()

	schemaResponse := datasource.SchemaResponse{}
	instancesDataSource.Schema(context.TODO(), datasource.SchemaRequest{}, &schemaResponse)
	_, instancesExists := schemaResponse.Schema.GetAttributes()["instances"]

	assert.True(t, instancesExists, "schema should contain instances")
}
