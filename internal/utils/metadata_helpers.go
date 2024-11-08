package utils

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// SetDataSourceTypeName generates & sets the data source type name.
func SetDataSourceTypeName(
	response *datasource.MetadataResponse,
	request datasource.MetadataRequest,
	name string,
) {
	response.TypeName = generateName(request.ProviderTypeName, name)
}

// SetResourceTypeName generates & sets the resource type name.
func SetResourceTypeName(
	response *resource.MetadataResponse,
	request resource.MetadataRequest,
	name string,
) {
	response.TypeName = generateName(request.ProviderTypeName, name)
}

func generateName(providerTypeName string, name string) string {
	return fmt.Sprintf("%s_%s", providerTypeName, name)
}
