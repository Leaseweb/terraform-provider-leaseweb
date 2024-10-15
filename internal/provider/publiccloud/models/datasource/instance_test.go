package datasource

import (
	"testing"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func TestNewInstance(t *testing.T) {
	reference := "reference"
	marketAppId := "marketAppId"

	sdkInstance := publicCloud.Instance{
		Id:        "id",
		Region:    "region",
		Reference: *publicCloud.NewNullableString(&reference),
		Image: publicCloud.Image{
			Id: "imageId",
		},
		State:               publicCloud.STATE_CREATING,
		Type:                publicCloud.TYPENAME_C3_2XLARGE,
		RootDiskSize:        50,
		RootDiskStorageType: publicCloud.STORAGETYPE_CENTRAL,
		Ips: []publicCloud.Ip{
			{Ip: "127.0.0.1"},
		},
		Contract: publicCloud.Contract{
			Term: publicCloud.CONTRACTTERM__1,
		},
		MarketAppId: *publicCloud.NewNullableString(&marketAppId),
	}

	got := NewInstance(sdkInstance)

	assert.Equal(t, "id", got.Id.ValueString())
	assert.Equal(t, "region", got.Region.ValueString())
	assert.Equal(t, "reference", got.Reference.ValueString())
	assert.Equal(t, "imageId", got.Image.Id.ValueString())
	assert.Equal(t, "CREATING", got.State.ValueString())
	assert.Equal(t, "lsw.c3.2xlarge", got.Type.ValueString())
	assert.Equal(t, int64(50), got.RootDiskSize.ValueInt64())
	assert.Equal(t, "CENTRAL", got.RootDiskStorageType.ValueString())
	assert.Len(t, got.Ips, 1)
	assert.Equal(t, "127.0.0.1", got.Ips[0].Ip.ValueString())
	assert.Equal(t, int64(1), got.Contract.Term.ValueInt64())
	assert.Equal(t, "marketAppId", got.MarketAppId.ValueString())

}
