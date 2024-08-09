package to_data_source_model

import (
	"testing"

	domain "github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain/dedicated_server"
	"github.com/stretchr/testify/assert"
)

func TestAdaptDedicatedServers(t *testing.T) {
	id := "123456"
	dedicatedServers := domain.DedicatedServers{
		domain.DedicatedServer{Id: id},
	}

	got := AdaptDedicatedServers(dedicatedServers)

	assert.Len(t, got.DedicatedServers, 1)
	assert.Equal(t, id, got.DedicatedServers[0].Id.ValueString())
}
