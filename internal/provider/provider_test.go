package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMetadata(t *testing.T) {
	leasewebProvider := New("dev")
	metadataResponse := provider.MetadataResponse{}
	leasewebProvider().Metadata(
		context.TODO(),
		provider.MetadataRequest{},
		&metadataResponse,
	)

	want := "dev"
	got := metadataResponse.Version

	assert.Equal(t, want, got, "version should be passed to provider")
}

func TestProviderSchema(t *testing.T) {
	leasewebProvider := New("dev")
	schemaResponse := provider.SchemaResponse{}
	leasewebProvider().Schema(context.TODO(), provider.SchemaRequest{}, &schemaResponse)

	assert.True(t, schemaResponse.Schema.Attributes["host"].IsOptional(), "host is optional")
	assert.True(t, schemaResponse.Schema.Attributes["scheme"].IsOptional(), "scheme is optional")
	assert.True(t, schemaResponse.Schema.Attributes["token"].IsSensitive(), "token is sensitive")
}

func TestDatSources(t *testing.T) {
	leasewebProvider := New("dev")

	assert.True(
		t,
		implementsDataSource(
			leasewebProvider().DataSources(
				context.TODO()), "_instances",
		),
		"data sources should implement InstancesDataSource",
	)
}

func TestResources(t *testing.T) {
	leasewebProvider := New("dev")

	assert.Nil(t, leasewebProvider().Resources(context.TODO()))
}

func implementsDataSource(dataSources []func() datasource.DataSource, expectedTypeName string) bool {
	for _, dataSource := range dataSources {
		resp := datasource.MetadataResponse{}
		dataSource().Metadata(context.TODO(), datasource.MetadataRequest{}, &resp)

		if resp.TypeName == expectedTypeName {
			return true
		}
	}

	return false
}
