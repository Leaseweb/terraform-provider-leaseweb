package provider

import (
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewLeasewebProviderClient(t *testing.T) {
	client := NewLeasewebProviderClient("token", &LeasewebProviderClientOptions{
		Host: "tralala.com",
	})

	got := client.Client.GetConfig().Host
	want := "tralala.com"

	assert.Equal(t, want, got, "want should be passed to client")
}

func TestAuthContext(t *testing.T) {
	client := NewLeasewebProviderClient("token", &LeasewebProviderClientOptions{})

	got := client.AuthContext().Value(publicCloud.ContextAPIKeys)
	want := map[string]publicCloud.APIKey{
		"X-LSW-Auth": {Key: "token", Prefix: ""},
	}

	assert.Equal(t, want, got, "token should be passed to context")
}
