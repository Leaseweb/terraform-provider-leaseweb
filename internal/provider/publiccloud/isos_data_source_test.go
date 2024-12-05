package publiccloud

import (
	"testing"

	"github.com/leaseweb/leaseweb-go-sdk/v2/publiccloud"
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

func Test_adaptIsosToISOsDataSource(t *testing.T) {
	sdkISOs := []publiccloud.Iso{
		{Id: "id"},
	}

	got := adaptIsosToISOsDataSource(sdkISOs)

	assert.Len(t, got.ISOs, 1)
	assert.Equal(t, "id", got.ISOs[0].ID.ValueString())
}
