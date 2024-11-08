package utils

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/stretchr/testify/assert"
)

func Test_generateName(t *testing.T) {
	got := generateName("a", "b")
	want := "a_b"

	assert.Equal(t, want, got)
}

func TestSetResourceTypeName(t *testing.T) {
	response := resource.MetadataResponse{}
	SetResourceTypeName(
		&response,
		resource.MetadataRequest{ProviderTypeName: "a"},
		"b",
	)

	assert.Equal(t, "a_b", response.TypeName)
}

func TestSetDataSourceTypeName(t *testing.T) {
	response := datasource.MetadataResponse{}
	SetDataSourceTypeName(
		&response,
		datasource.MetadataRequest{ProviderTypeName: "a"},
		"b",
	)

	assert.Equal(t, "a_b", response.TypeName)
}

func ExampleSetResourceTypeName() {
	response := resource.MetadataResponse{}
	SetResourceTypeName(
		&response,
		resource.MetadataRequest{ProviderTypeName: "leaseweb"},
		"public_cloud_instance",
	)

	fmt.Println(response.TypeName)
	// Output: leaseweb_public_cloud_instance
}

func ExampleSetDataSourceTypeName() {
	response := datasource.MetadataResponse{}
	SetDataSourceTypeName(
		&response,
		datasource.MetadataRequest{ProviderTypeName: "leaseweb"},
		"public_cloud_instances",
	)

	fmt.Println(response.TypeName)
	// Output: leaseweb_public_cloud_instances
}
