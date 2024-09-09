package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/stretchr/testify/assert"
)

const (
	providerConfig = `
provider "leaseweb" {
  host     = "localhost:8080"
  scheme = "http"
  token = "tralala"
}
`
)

var (
	testAccProtoV6ProviderFactories = map[string]func() (
		tfprotov6.ProviderServer,
		error,
	){
		"leaseweb": providerserver.NewProtocol6WithError(NewProvider("test")()),
	}
)

func TestLeasewebProvider_Metadata(t *testing.T) {
	leasewebProvider := NewProvider("dev")
	metadataResponse := provider.MetadataResponse{}
	leasewebProvider().Metadata(
		context.TODO(),
		provider.MetadataRequest{},
		&metadataResponse,
	)

	want := "dev"
	got := metadataResponse.Version

	assert.Equal(
		t,
		want,
		got,
		"version should be passed to provider",
	)
}

func TestLeasewebProvider_Schema(t *testing.T) {
	leasewebProvider := NewProvider("dev")
	schemaResponse := provider.SchemaResponse{}
	leasewebProvider().Schema(
		context.TODO(),
		provider.SchemaRequest{},
		&schemaResponse,
	)

	assert.True(
		t,
		schemaResponse.Schema.Attributes["host"].IsOptional(),
		"host is optional",
	)
	assert.True(
		t,
		schemaResponse.Schema.Attributes["scheme"].IsOptional(),
		"scheme is optional",
	)
	assert.True(
		t,
		schemaResponse.Schema.Attributes["token"].IsSensitive(),
		"token is sensitive",
	)
}

func TestLeasewebProvider_DataSources(t *testing.T) {
	leasewebProvider := NewProvider("dev")

	assert.True(
		t,
		implementsDataSource(
			leasewebProvider().DataSources(
				context.TODO()), "_public_cloud_instances",
		),
		"data sources should implement Public Cloud InstancesDataSource",
	)

	assert.True(
		t,
		implementsDataSource(
			leasewebProvider().DataSources(
				context.TODO()), "_dedicated_servers",
		),
		"data sources should implement DedicatedServersDataSource",
	)
}

func TestLeasewebProvider_Resources(t *testing.T) {
	leasewebProvider := NewProvider("dev")

	assert.True(
		t,
		implementsResource(
			leasewebProvider().Resources(
				context.TODO()), "_public_cloud_instance",
		),
		"resource should implement Public Cloud InstanceResource",
	)

	assert.True(
		t,
		implementsResource(
			leasewebProvider().Resources(
				context.TODO()), "_dedicated_server",
		),
		"resource should implement Dedicated Server dedicatedServerResource",
	)
}

func implementsDataSource(
	dataSources []func() datasource.DataSource,
	expectedTypeName string,
) bool {
	for _, dataSource := range dataSources {
		resp := datasource.MetadataResponse{}
		dataSource().Metadata(context.TODO(), datasource.MetadataRequest{}, &resp)

		if resp.TypeName == expectedTypeName {
			return true
		}
	}

	return false
}

func implementsResource(
	resources []func() resource.Resource,
	expectedTypeName string,
) bool {
	for _, element := range resources {
		resp := resource.MetadataResponse{}
		element().Metadata(context.TODO(), resource.MetadataRequest{}, &resp)

		if resp.TypeName == expectedTypeName {
			return true
		}
	}

	return false
}
