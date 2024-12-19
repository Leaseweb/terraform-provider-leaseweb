package publiccloud

import (
	"testing"

	"github.com/leaseweb/leaseweb-go-sdk/v3/publiccloud"
	"github.com/stretchr/testify/assert"
)

func Test_adaptIsoToISODataSource(t *testing.T) {
	sdkISO := publiccloud.Iso{
		Id:   "id",
		Name: "name",
	}

	got := adaptIsoToISODataSource(sdkISO)

	assert.Equal(t, "id", got.ID.ValueString())
	assert.Equal(t, "name", got.Name.ValueString())
}
