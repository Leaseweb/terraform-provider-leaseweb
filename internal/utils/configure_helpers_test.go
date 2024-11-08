package utils

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
	"github.com/stretchr/testify/assert"
)

func TestGetResourceClient(t *testing.T) {
	t.Run("returns false if ProviderData is nil", func(t *testing.T) {
		request := resource.ConfigureRequest{}
		response := resource.ConfigureResponse{}

		_, ok := GetResourceClient(request, &response)

		assert.False(t, ok)
		assert.False(t, response.Diagnostics.HasError())
	})

	t.Run(
		"sets error if ProviderData is an unexpected type",
		func(t *testing.T) {
			type tralala struct{}
			request := resource.ConfigureRequest{
				ProviderData: tralala{},
			}
			response := resource.ConfigureResponse{}

			_, ok := GetResourceClient(request, &response)

			assert.False(t, ok)
			assert.True(t, response.Diagnostics.HasError())
			assert.Equal(
				t,
				response.Diagnostics[0].Summary(),
				"Unexpected Resource Configure Type",
			)
			assert.Equal(
				t,
				response.Diagnostics[0].Detail(),
				"Expected client.Client, got: utils.tralala. Please report this issue to the provider developers.",
			)
		},
	)

	t.Run("client is returned if type is expected", func(t *testing.T) {
		want := client.Client{}
		request := resource.ConfigureRequest{
			ProviderData: want,
		}
		response := resource.ConfigureResponse{}

		apiClient, ok := GetResourceClient(request, &response)

		assert.True(t, ok)
		assert.False(t, response.Diagnostics.HasError())
		assert.Equal(t, want, *apiClient)
	})
}

func TestGetDataSourceClient(t *testing.T) {
	t.Run("sets false if ProviderData is nil", func(t *testing.T) {
		request := datasource.ConfigureRequest{}
		response := datasource.ConfigureResponse{}

		_, ok := GetDataSourceClient(request, &response)

		assert.False(t, ok)
		assert.False(t, response.Diagnostics.HasError())
	})

	t.Run(
		"sets error if ProviderData is an unexpected type",
		func(t *testing.T) {
			type tralala struct{}
			request := datasource.ConfigureRequest{
				ProviderData: tralala{},
			}
			response := datasource.ConfigureResponse{}

			_, ok := GetDataSourceClient(request, &response)

			assert.False(t, ok)
			assert.True(t, response.Diagnostics.HasError())
			assert.Equal(
				t,
				response.Diagnostics[0].Summary(),
				"Unexpected Resource Configure Type",
			)
			assert.Equal(
				t,
				response.Diagnostics[0].Detail(),
				"Expected client.Client, got: utils.tralala. Please report this issue to the provider developers.",
			)
		},
	)

	t.Run("client is returned if type is expected", func(t *testing.T) {
		want := client.Client{}
		request := datasource.ConfigureRequest{
			ProviderData: want,
		}
		response := datasource.ConfigureResponse{}

		apiClient, ok := GetDataSourceClient(request, &response)

		assert.True(t, ok)
		assert.False(t, response.Diagnostics.HasError())
		assert.Equal(t, want, *apiClient)
	})
}

func ExampleGetResourceClient() {
	request := resource.ConfigureRequest{
		ProviderData: client.Client{},
	}
	response := resource.ConfigureResponse{}
	apiClient, ok := GetResourceClient(request, &response)

	fmt.Println(apiClient, ok)
	// Output: &{<nil> <nil>} true
}

func ExampleGetDataSourceClient() {
	request := datasource.ConfigureRequest{
		ProviderData: client.Client{},
	}
	response := datasource.ConfigureResponse{}
	apiClient, ok := GetDataSourceClient(request, &response)

	fmt.Println(apiClient, ok)
	// Output: &{<nil> <nil>} true
}
