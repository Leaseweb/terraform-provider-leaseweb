package client

import (
	"context"
	"testing"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func TestClientSupportsHost(t *testing.T) {
	client := NewClient("token", &Options{
		Host: "tralala.com",
	})

	got := client.PublicCloudClient.GetConfig().Host
	want := "tralala.com"

	assert.Equal(t, want, got, "client supports host")
}

func TestClientSupportsScheme(t *testing.T) {
	client := NewClient("token", &Options{
		Scheme: "http",
	})

	got := client.PublicCloudClient.GetConfig().Scheme
	want := "http"

	assert.Equal(t, want, got, "client supports scheme")
}

func TestAuthContext(t *testing.T) {
	client := NewClient("token", &Options{})

	got := client.AuthContext(context.TODO()).Value(publicCloud.ContextAPIKeys)
	want := map[string]publicCloud.APIKey{
		"X-LSW-Auth": {Key: "token", Prefix: ""},
	}

	assert.Equal(t, want, got, "token should be passed to context")
}
