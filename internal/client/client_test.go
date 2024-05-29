package client

import (
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClientSupportsHost(t *testing.T) {
	client := NewClient("token", &Options{
		Host: "tralala.com",
	})

	got := client.SdkClient.GetConfig().Host
	want := "tralala.com"

	assert.Equal(t, want, got, "client supports host")
}

func TestClientSupportsScheme(t *testing.T) {
	client := NewClient("token", &Options{
		Scheme: "http",
	})

	got := client.SdkClient.GetConfig().Scheme
	want := "http"

	assert.Equal(t, want, got, "client supports scheme")
}

func TestAuthContext(t *testing.T) {
	client := NewClient("token", &Options{})

	got := client.AuthContext().Value(publicCloud.ContextAPIKeys)
	want := map[string]publicCloud.APIKey{
		"X-LSW-Auth": {Key: "token", Prefix: ""},
	}

	assert.Equal(t, want, got, "token should be passed to context")
}
