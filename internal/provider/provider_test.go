package provider

import (
	"context"
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

func TestSchema(t *testing.T) {
	leasewebProvider := New("dev")
	schemaResponse := provider.SchemaResponse{}
	leasewebProvider().Schema(context.TODO(), provider.SchemaRequest{}, &schemaResponse)

	assert.True(t, schemaResponse.Schema.Attributes["host"].IsOptional(), "host is optional")
	assert.True(t, schemaResponse.Schema.Attributes["token"].IsSensitive(), "token is sensitive")
}

func TestDatSources(t *testing.T) {
	leasewebProvider := New("dev")

	assert.Nil(t, leasewebProvider().DataSources(context.TODO()))
}

func TestResources(t *testing.T) {
	leasewebProvider := New("dev")

	assert.Nil(t, leasewebProvider().Resources(context.TODO()))
}
