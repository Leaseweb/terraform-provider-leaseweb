package modify_plan

import (
	"testing"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func TestNewRegions(t *testing.T) {
	got := NewRegions([]publicCloud.Region{{Name: "region"}})

	assert.Equal(t, Regions{"region"}, got)
}

func TestRegions_Contains(t *testing.T) {
	regions := NewRegions([]publicCloud.Region{{Name: "region"}})

	assert.True(t, regions.Contains("region"))
	assert.False(t, regions.Contains("tralala"))
}
