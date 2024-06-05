package instance

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_convertStringToInt32(t *testing.T) {
	got, _ := convertStringToInt32("", "32", &resource.CreateResponse{})

	assert.Equal(t, int32(32), got, "Should return 32")
}

func Test_convertStringToInt32Error(t *testing.T) {
	resp := resource.CreateResponse{}

	_, err := convertStringToInt32("tralala", "tralala", &resp)

	assert.NotNilf(t, err, "Should return error")
	assert.Equal(t,
		"Could not set \"tralala\", unexpected error:  \"strconv.Atoi: parsing \"tralala\": invalid syntax\"",
		resp.Diagnostics.Errors()[0].Detail(),
		"Should add diagnostics error",
	)
}
