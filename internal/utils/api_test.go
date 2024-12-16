package utils

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/leaseweb/leaseweb-go-sdk/v2/dedicatedserver"
	"github.com/leaseweb/leaseweb-go-sdk/v2/publiccloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
	"github.com/stretchr/testify/assert"
)

func Test_generateTypeName(t *testing.T) {
	got := generateTypeName("provider", "name")
	want := "provider_name"

	assert.Equal(t, want, got)
}

func TestResourceAPI_Metadata(t *testing.T) {
	api := ResourceAPI{
		Name: "name",
	}
	response := resource.MetadataResponse{}
	request := resource.MetadataRequest{ProviderTypeName: "providerTypeName"}
	api.Metadata(context.TODO(), request, &response)

	assert.Equal(t, "providerTypeName_name", response.TypeName)
}

func TestDataSourceAPI_Metadata(t *testing.T) {
	api := DataSourceAPI{
		Name: "name",
	}
	response := datasource.MetadataResponse{}
	request := datasource.MetadataRequest{ProviderTypeName: "providerTypeName"}
	api.Metadata(context.TODO(), request, &response)

	assert.Equal(t, "providerTypeName_name", response.TypeName)
}

func Test_getCoreClient(t *testing.T) {
	t.Run(
		"return nil and don't set error if providerData is nil",
		func(t *testing.T) {
			diags := diag.Diagnostics{}
			coreClient := getCoreClient(nil, &diags)

			assert.False(t, diags.HasError())
			assert.Nil(t, coreClient)
		},
	)

	t.Run("return client if providerData is set", func(t *testing.T) {
		apiClient := publiccloud.NewAPIClient(publiccloud.NewConfiguration())
		want := client.Client{
			PubliccloudAPI: apiClient.PubliccloudAPI,
		}
		diags := diag.Diagnostics{}
		coreClient := getCoreClient(want, &diags)

		assert.False(t, diags.HasError())
		assert.Equal(t, want, *coreClient)
	})

	t.Run(
		"return nil and set error if resource configure type is unexpected",
		func(t *testing.T) {
			diags := diag.Diagnostics{}
			coreClient := getCoreClient("tralala", &diags)

			assert.Nil(t, coreClient)
			assert.Equal(t, 1, diags.ErrorsCount())
			assert.Equal(
				t,
				"Unexpected Resource Configure Type",
				diags[0].Summary(),
			)
			assert.Equal(
				t,
				"Expected client.Client, got: string. Please report this issue to the provider developers.",
				diags[0].Detail(),
			)
		},
	)
}

func TestPubliccloudResourceAPI_Configure(t *testing.T) {
	t.Run("nothing is set if providerData is nil", func(t *testing.T) {
		api := PubliccloudResourceAPI{}
		response := resource.ConfigureResponse{}
		api.Configure(context.TODO(), resource.ConfigureRequest{}, &response)

		assert.Nil(t, api.Client)
	})

	t.Run("client is set from ProviderData", func(t *testing.T) {
		api := PubliccloudResourceAPI{}
		response := resource.ConfigureResponse{}
		apiClient := publiccloud.NewAPIClient(publiccloud.NewConfiguration())
		api.Configure(
			context.TODO(),
			resource.ConfigureRequest{
				ProviderData: client.Client{
					PubliccloudAPI: apiClient.PubliccloudAPI,
				},
			},
			&response,
		)

		assert.Equal(t, apiClient.PubliccloudAPI, api.Client)
	})
}

func TestNewPubliccloudResourceAPI(t *testing.T) {
	got := NewPubliccloudResourceAPI("name")
	want := PubliccloudResourceAPI{
		ResourceAPI: ResourceAPI{
			Name: "name",
		},
	}

	assert.Equal(t, want, got)
}

func ExampleNewPubliccloudResourceAPI() {
	api := NewPubliccloudResourceAPI("name")
	fmt.Println(api)
	// Output: {{name} <nil>}
}

func TestPubliccloudDataSourceAPI_Configure(t *testing.T) {
	t.Run("nothing is set if providerData is nil", func(t *testing.T) {
		api := PubliccloudDataSourceAPI{}
		response := datasource.ConfigureResponse{}
		api.Configure(context.TODO(), datasource.ConfigureRequest{}, &response)

		assert.Nil(t, api.Client)
	})

	t.Run("client is set from ProviderData", func(t *testing.T) {
		api := PubliccloudDataSourceAPI{}
		response := datasource.ConfigureResponse{}
		apiClient := publiccloud.NewAPIClient(publiccloud.NewConfiguration())
		api.Configure(
			context.TODO(),
			datasource.ConfigureRequest{
				ProviderData: client.Client{
					PubliccloudAPI: apiClient.PubliccloudAPI,
				},
			},
			&response,
		)

		assert.Equal(t, apiClient.PubliccloudAPI, api.Client)
	})
}

func TestNewPubliccloudDataSourceAPI(t *testing.T) {
	got := NewPubliccloudDataSourceAPI("name")
	want := PubliccloudDataSourceAPI{
		DataSourceAPI: DataSourceAPI{
			Name: "name",
		},
	}

	assert.Equal(t, want, got)
}

func TestDedicatedserverResourceAPI_Configure(t *testing.T) {
	t.Run("nothing is set if providerData is nil", func(t *testing.T) {
		api := DedicatedserverResourceAPI{}
		response := resource.ConfigureResponse{}
		api.Configure(context.TODO(), resource.ConfigureRequest{}, &response)

		assert.Nil(t, api.Client)
	})

	t.Run("client is set from ProviderData", func(t *testing.T) {
		api := DedicatedserverResourceAPI{}
		response := resource.ConfigureResponse{}
		apiClient := dedicatedserver.NewAPIClient(dedicatedserver.NewConfiguration())
		api.Configure(
			context.TODO(),
			resource.ConfigureRequest{
				ProviderData: client.Client{
					DedicatedserverAPI: apiClient.DedicatedserverAPI,
				},
			},
			&response,
		)

		assert.Equal(t, apiClient.DedicatedserverAPI, api.Client)
	})
}

func TestNewDedicatedserverResourceAPI(t *testing.T) {
	got := NewDedicatedserverResourceAPI("name")
	want := DedicatedserverResourceAPI{
		ResourceAPI: ResourceAPI{
			Name: "name",
		},
	}

	assert.Equal(t, want, got)
}

func ExampleNewDedicatedserverResourceAPI() {
	api := NewDedicatedserverResourceAPI("name")
	fmt.Println(api)
	// Output: {{name} <nil>}
}

func TestDedicatedserverDataSourceAPI_Configure(t *testing.T) {
	t.Run("nothing is set if providerData is nil", func(t *testing.T) {
		api := DedicatedserverDataSourceAPI{}
		response := datasource.ConfigureResponse{}
		api.Configure(context.TODO(), datasource.ConfigureRequest{}, &response)

		assert.Nil(t, api.Client)
	})

	t.Run("client is set from ProviderData", func(t *testing.T) {
		api := DedicatedserverDataSourceAPI{}
		response := datasource.ConfigureResponse{}
		apiClient := dedicatedserver.NewAPIClient(dedicatedserver.NewConfiguration())
		api.Configure(
			context.TODO(),
			datasource.ConfigureRequest{
				ProviderData: client.Client{
					DedicatedserverAPI: apiClient.DedicatedserverAPI,
				},
			},
			&response,
		)

		assert.Equal(t, apiClient.DedicatedserverAPI, api.Client)
	})
}

func TestNewDedicatedserverDataSourceAPI(t *testing.T) {
	got := NewDedicatedserverDataSourceAPI("name")
	want := DedicatedserverDataSourceAPI{
		DataSourceAPI: DataSourceAPI{
			Name: "name",
		},
	}

	assert.Equal(t, want, got)
}

func ExampleNewDedicatedserverDataSourceAPI() {
	api := NewDedicatedserverDataSourceAPI("name")
	fmt.Println(api)
	// Output: {{name} <nil>}
}
