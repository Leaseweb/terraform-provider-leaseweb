package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLeasewebProvider_Metadata(t *testing.T) {
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

func TestLeasewebProvider_Schema(t *testing.T) {
	leasewebProvider := New("dev")
	schemaResponse := provider.SchemaResponse{}
	leasewebProvider().Schema(context.TODO(), provider.SchemaRequest{}, &schemaResponse)

	assert.True(t, schemaResponse.Schema.Attributes["host"].IsOptional(), "host is optional")
	assert.True(t, schemaResponse.Schema.Attributes["scheme"].IsOptional(), "scheme is optional")
	assert.True(t, schemaResponse.Schema.Attributes["token"].IsSensitive(), "token is sensitive")
}

func TestLeasewebProvider_DataSources(t *testing.T) {
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

func TestLeasewebProvider_Resources(t *testing.T) {
	leasewebProvider := New("dev")

	assert.True(
		t,
		implementsResource(
			leasewebProvider().Resources(
				context.TODO()), "_instance",
		),
		"resources should implement InstanceResource",
	)
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

func implementsResource(resources []func() resource.Resource, expectedTypeName string) bool {
	for _, element := range resources {
		resp := resource.MetadataResponse{}
		element().Metadata(context.TODO(), resource.MetadataRequest{}, &resp)

		if resp.TypeName == expectedTypeName {
			return true
		}
	}

	return false
}
